package gnode

import (
	"fmt"
	"strings"
	"time"

	proto "github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const (
	ReqStatusAccepted = "accepted"
	ReqStatusOnGoing  = "on going"
	ReqStatusOnError  = "on error"
	ReqStatusExecuted = "executed"
)

type EntryManager struct {
	nbId                 int64
	gnode                *GNode
	requestQueuer        requestQueuer
	requestStatusManager RequestStatusManager
}

func (m *EntryManager) init(g *GNode) {
	fmt.Println("init EntryManager")
	m.gnode = g
	m.requestStatusManager.init(g)
	m.requestQueuer.init(g, m)
}

func (m *EntryManager) getNewId() string {
	m.nbId++
	return fmt.Sprintf("%s%x%x", m.gnode.host, time.Now().UnixNano(), m.nbId)
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
	labels, errl := m.convertLabels(mes.Args)
	if errl != nil {
		return errl
	}
	signedData := m.getDataToSign(mes.Data, labels)
	if err := m.gnode.key.verifyUserSignature(mes.UserName, mes.Key, signedData); err != nil {
		return fmt.Errorf("User %s signature not authenticated", mes.UserName)
	}
	signMap := make(map[string][]byte)
	timeNow := time.Now().Unix()
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
	req := &CheckEntryRequest{
		Id:          m.getNewId(),
		OriginNode:  m.gnode.name,
		NodeSignMap: signMap,
		Entry:       entryData,
		EntryHash:   entryHash.hash,
		IsBranch:    isBranch,
		Labels:      labels,
		UserName:    entry.UserName,
	}
	if err := m.requestQueuer.push(req); err != nil {
		return err
	}
	//Return answer to client
	m.requestStatusManager.createReqStatus(req, entry)
	answer := m.gnode.createAnswer(mes, false)
	answer.Args = []string{req.Id}
	m.gnode.senderManager.sendMessage(answer)
	return nil
}

func (m *EntryManager) convertLabels(args []string) ([]*TreeLabel, error) {
	labels := []*TreeLabel{}
	for _, arg := range args {
		list := strings.Split(arg, "=")
		if len(list) == 1 {
			if strings.Index(arg, "=") < 0 {
				return nil, fmt.Errorf("Label format error, needs '=': %s", arg)
			}
			labels = append(labels, &TreeLabel{Name: list[0]})
		} else if len(list) == 2 {
			labels = append(labels, &TreeLabel{Name: list[0], Value: list[1]})
		} else {
			return nil, fmt.Errorf("Label format error: %s", arg)
		}
	}
	return labels, nil
}

func (m *EntryManager) computeRootHash() []byte {
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
	status := m.requestStatusManager.getCreateReqStatus(req.Id, req.UserName)
	if status.Done {
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
		if err := m.verifyEntry(req); err != nil {
			m.requestStatusManager.setReqStatusError(req.Id, "not valid:", err)
			return nil
		}
	}

	//if majority of nodes accepted it then add it in blockchain
	logf.info("Entry nbNode ok=%d\n", nb)
	if nb > m.gnode.nbNode/2 {
		if err := m.addItemInBlockchain(req.Entry, req.EntryHash, req.IsBranch); err != nil {
			m.requestStatusManager.setReqStatusError(req.Id, "add error:", err)
		} else {
			m.requestStatusManager.setReqStatus(req.Id, ReqStatusExecuted)
		}
		status.Done = true
	}

	//if received too much time the same unchanged request then set as done
	if int(status.NbOk) == nb {
		status.NbSame++
		if int(status.NbSame) > m.gnode.nbNode*2 {
			m.requestStatusManager.setReqStatusError(req.Id, "not validated at majority", nil)
			status.Done = true
			return nil
		}
	}
	status.NbOk = int32(nb)

	//send back the updated request to all connections expect the one it comes from.
	originNode := req.OriginNode
	req.OriginNode = m.gnode.name
	if err := m.sendCheckEntry(req, originNode); err != nil {
		return err
	}
	return nil
}

func (m *EntryManager) verifyEntry(req *CheckEntryRequest) error {
	logf.info("Verify entry request: %s\n", req.Id)

	//verify root hash
	if string(m.gnode.treeManager.root.FullHash) != string(req.RootHash) {
		return fmt.Errorf("root hash not authenticate")
	}
	//Verify entry hash
	hash := getNewHash()
	hash.setHash(req.Entry)
	if string(hash.hash) != string(req.EntryHash) {
		return fmt.Errorf("entry hash not authenticate")
	}
	//unmarchal entry
	entry := &BCEntry{}
	if err := proto.Unmarshal(req.Entry, entry); err != nil {
		return fmt.Errorf("unmarshaling entry error: %v", err)
	}
	logf.info("Entry: %s\n", m.entryString(entry))

	//Verify the entry user signature
	signedData := m.getDataToSign(entry.Payload, entry.Labels)
	if err := m.gnode.key.verifyUserSignature(entry.UserName, entry.UserSignature, signedData); err != nil {
		return fmt.Errorf("user %s not authenticated in entry\n", entry.UserName)
	}

	//add the node signature
	signature, err := m.gnode.key.sign(req.Entry)
	if err != nil {
		return fmt.Errorf("node signature error: ", err)
	}
	req.NodeSignMap[m.gnode.name] = signature
	return nil
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
