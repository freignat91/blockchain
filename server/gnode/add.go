package gnode

import (
	"fmt"
	"strings"
	"time"

	proto "github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

type EntryManager struct {
	nbId  int64
	gnode *GNode
	idMap map[string]*CheckRequestCounter
}

type CheckRequestCounter struct {
	nbOk   int
	nbSame int
	done   bool
}

func (m *EntryManager) init(g *GNode) {
	m.gnode = g
	m.idMap = make(map[string]*CheckRequestCounter)
}

func (m *EntryManager) getNewId() string {
	m.nbId++
	return fmt.Sprintf("%s-%d", m.gnode.host, m.nbId)
}

func (m *EntryManager) getRequestCounter(id string) *CheckRequestCounter {
	if counter, ok := m.idMap[id]; ok {
		return counter
	}
	counter := &CheckRequestCounter{}
	m.idMap[id] = counter
	return counter
}

func (m *EntryManager) addBranch(mes *AntMes) error {
	return m.addItem(mes, true)
}

func (m *EntryManager) addEntry(mes *AntMes) error {
	return m.addItem(mes, false)
}

func (m *EntryManager) addItem(mes *AntMes, isBranch bool) error {
	//Verify user signature
	payload := mes.Data
	logf.info("Received add item: %s\n", string(payload))
	labels := m.convertLabels(mes.Args)
	signedData := m.getDataToSign(mes.Data, labels)
	if err := m.gnode.key.verifyUserSignature(mes.UserName, mes.Key, signedData); err != nil {
		return fmt.Errorf("User %s not authenticated", mes.UserName)
	}
	//Verify request validity
	if _, err := m.gnode.treeManager.getLastExistingBranchBlock(labels, isBranch); err != nil {
		return err
	}
	//Return answer to client
	answer := m.gnode.createAnswer(mes, false)
	m.gnode.senderManager.sendMessage(answer)
	signMap := make(map[string][]byte)
	timeNow, _ := time.Now().MarshalBinary()
	entry := &BCEntry{
		Date:          timeNow,
		Labels:        labels,
		Payload:       payload,
		UserName:      mes.UserName,
		UserSignature: mes.Key,
	}
	entryData, errm := proto.Marshal(entry)
	if errm != nil {
		return fmt.Errorf("marshaling entry error: ", errm)
	}
	if signature, err := m.gnode.key.sign(entryData); err != nil {
		return fmt.Errorf("node signature error: ", err)
	} else {
		signMap[m.gnode.name] = signature
	}
	entryHash := getNewHash()
	entryHash.setHash(entryData)
	rootHash := m.computeRootHash(entryHash.hash)
	req := &CheckEntryRequest{
		Id:          m.getNewId(),
		OriginNode:  m.gnode.name,
		NodeSignMap: signMap,
		Entry:       entryData,
		EntryHash:   entryHash.hash,
		RootHash:    rootHash,
		IsBranch:    isBranch,
	}
	if err := m.sendCheckEntry(req, ""); err != nil {
		return fmt.Errorf("sendCheckEntry error: %v\n", err)
	}
	return nil
}

func (m *EntryManager) convertLabels(args []string) []*TreeLabel {
	labels := []*TreeLabel{}
	for _, arg := range args {
		list := strings.Split(arg, "=")
		if len(list) == 1 {
			labels = append(labels, &TreeLabel{Name: list[0]})
		} else {
			labels = append(labels, &TreeLabel{Name: list[0], Value: list[1]})
		}
	}
	return labels
}

func (m *EntryManager) computeRootHash(entry []byte) []byte {
	return m.gnode.treeManager.root.FullHash
}

// Send message using all targets
func (m *EntryManager) sendCheckEntry(req *CheckEntryRequest, originNodeName string) error {
	for _, target := range m.gnode.targetMap {
		if target.name != originNodeName {
			m.sendCheckEntryToTarget(target, req)
		}
	}
	return nil
}

func (m *EntryManager) sendCheckEntryToTarget(target *gnodeTarget, req *CheckEntryRequest) {
	_, err := target.client.CheckEntry(context.Background(), req)
	if err != nil {
		logf.info("CheckEntry error: %v\n", err)
	}
}

func (m *EntryManager) checkEntry(req *CheckEntryRequest) error {
	//Vertify if not already added in blockchain
	counter := m.getRequestCounter(req.Id)
	if counter.done {
		//fmt.Printf("Received alreday treated CheckEntry: %s\n", req.Id)
		return nil
	}
	logf.info("Received CheckEntry: %s\n", req.Id)

	//Verify the other nodes signatures
	nb := 0
	for nodeName, signature := range req.NodeSignMap {
		if err := m.gnode.key.verifyNodeSignature(nodeName, signature, req.Entry); err != nil {
			logf.warn("The node %s is not authenticated in entry: %v\n", nodeName, err)
		} else {
			nb++
		}
	}

	//If not already verify then verify and sign it
	if _, exist := req.NodeSignMap[m.gnode.name]; !exist {
		if !m.verifyEntry(req) {
			logf.warn("Entry request %s not validated\n", req.Id)
			return nil
		}
	}

	//if majority of nodes accepted it then add it in blockchain
	logf.info("Entry nbNode ok=%d\n", nb)
	if nb > m.gnode.nbNode/2 {
		m.addItemInBlockchain(req.Entry, req.EntryHash, req.IsBranch)
		counter.done = true
	}

	//if received too much time the same unchanged request then set as done
	if counter.nbOk == nb {
		counter.nbSame++
		if counter.nbSame > m.gnode.nbNode*2 {
			logf.warn("Received too many time the same request %s with no more node validation: cancel request\n", req.Id)
			counter.done = true
			return nil
		}
	}
	counter.nbOk = nb

	//send back the updated request to all connections expect the one it comes from.
	originNode := req.OriginNode
	req.OriginNode = m.gnode.name
	if err := m.sendCheckEntry(req, originNode); err != nil {
		return fmt.Errorf("sendCheckEntry error: %v\n", err)
	}
	return nil
}

func (m *EntryManager) verifyEntry(req *CheckEntryRequest) bool {
	logf.info("Verify entry request: %s\n", req.Id)

	//verify root hash
	if string(m.gnode.treeManager.root.FullHash) != string(req.RootHash) {
		logf.error("Roothash not authenticate\n")
		return false
	}
	//Verify entry hash
	hash := getNewHash()
	hash.setHash(req.Entry)
	if string(hash.hash) != string(req.EntryHash) {
		logf.error("Entry hash not authenticate\n")
		return false
	}
	//unmarchal entry
	entry := &BCEntry{}
	if err := proto.Unmarshal(req.Entry, entry); err != nil {
		logf.error("unmarshaling entry error: %v\n", err)
		return false
	}
	logf.info("Entry: %s\n", m.entryString(entry))

	//Verify the entry user signature
	signedData := m.getDataToSign(entry.Payload, entry.Labels)
	if err := m.gnode.key.verifyUserSignature(entry.UserName, entry.UserSignature, signedData); err != nil {
		logf.error("User %s not authenticated in entry\n", entry.UserName)
		return false
	}

	//add the node signature
	signature, err := m.gnode.key.sign(req.Entry)
	if err != nil {
		logf.warn("node signature error: ", err)
		return false
	}
	req.NodeSignMap[m.gnode.name] = signature

	//verify root hashs
	//TODO
	return true
}

func (m *EntryManager) entryString(e *BCEntry) string {
	return fmt.Sprintf("user:%s payload:%s", e.UserName, e.Payload)
}

func (m *EntryManager) getDataToSign(payload []byte, labels []*TreeLabel) []byte {
	size := len(payload)
	for _, label := range labels {
		size += len(label.Name) + len(label.Value) + 1
	}
	dataToSign := make([]byte, size, size)
	nn := m.appendData(dataToSign, 0, payload)
	for _, label := range labels {
		nn = m.appendData(dataToSign, nn, []byte(fmt.Sprintf("%s=%s", label.Name, label.Value)))
	}
	return dataToSign
}

func (m *EntryManager) appendData(buffer []byte, nn int, item []byte) int {
	for i := 0; i < len(item); i++ {
		buffer[nn+i] = item[i]
	}
	return nn + len(item)
}

func (m *EntryManager) addItemInBlockchain(entryData []byte, hash []byte, isBranch bool) error {
	entry := &BCEntry{}
	if err := proto.Unmarshal(entryData, entry); err != nil {
		return fmt.Errorf("unmarshaling entry error: %v\n", err)
	}
	entry.Hash = hash
	if err := m.gnode.treeManager.addItem(entry, isBranch); err != nil {
		return err
	}
	logf.info("Entry added in blockchain %s\n", m.entryString(entry))
	return nil
}
