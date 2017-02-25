package gnode

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type ReceiverManager struct {
	usage        int
	gnode        *GNode
	buffer       MessageBuffer
	receiverList []*MessageReceiver
	ioChan       chan *AntMes
	nbReceiver   int
	receiver     MessageReceiver
	answerMap    map[string]*AntMes
	getChan      chan string
	lockClient   sync.RWMutex
	functionMap  map[string]interface{}
}

func (m *ReceiverManager) loadFunctions() {
	m.functionMap = make(map[string]interface{})
	//node Functions
	m.functionMap["ping"] = m.gnode.nodeFunctions.ping
	m.functionMap["pingFromTo"] = m.gnode.nodeFunctions.pingFromTo
	m.functionMap["setLogLevel"] = m.gnode.nodeFunctions.setLogLevel
	m.functionMap["killNode"] = m.gnode.nodeFunctions.killNode
	m.functionMap["updateGrid"] = m.gnode.nodeFunctions.updateGrid
	m.functionMap["writeStatsInLog"] = m.gnode.nodeFunctions.writeStatsInLog
	m.functionMap["clear"] = m.gnode.nodeFunctions.clear
	m.functionMap["forceGC"] = m.gnode.nodeFunctions.forceGCMes
	m.functionMap["getConnections"] = m.gnode.nodeFunctions.getConnections
	m.functionMap["createUser"] = m.gnode.nodeFunctions.createUser
	m.functionMap["createNodeUser"] = m.gnode.nodeFunctions.createNodeUser
	m.functionMap["removeUser"] = m.gnode.nodeFunctions.removeUser
	m.functionMap["removeNodeUser"] = m.gnode.nodeFunctions.removeNodeUser
	//gnode Function
	m.functionMap["sendBackEvent"] = m.gnode.sendBackEvent
	m.functionMap["setPublicKey"] = m.gnode.setPublicKey
}

func (m *ReceiverManager) start(gnode *GNode, bufferSize int, maxGoRoutine int) {
	m.gnode = gnode
	m.loadFunctions()
	m.lockClient = sync.RWMutex{}
	m.nbReceiver = maxGoRoutine
	m.buffer.init(bufferSize)
	m.ioChan = make(chan *AntMes)
	m.getChan = make(chan string)
	m.answerMap = make(map[string]*AntMes)
	m.receiverList = []*MessageReceiver{}
	if maxGoRoutine <= 0 {
		m.receiver.gnode = gnode
		return
	}
	for i := 0; i < maxGoRoutine; i++ {
		routine := &MessageReceiver{
			id:              i,
			gnode:           m.gnode,
			receiverManager: m,
		}
		m.receiverList = append(m.receiverList, routine)
		routine.start()
	}
	go func() {
		for {
			mes, ok := m.buffer.get(true)
			//logf.info("Receive message ok=%t %v\n", ok, mes.toString())
			if ok && mes != nil {
				m.ioChan <- mes
			}
		}
	}()

}

func (m *ReceiverManager) waitForAnswer(id string, timeoutSecond int) (*AntMes, error) {
	if mes, ok := m.answerMap[id]; ok {
		return mes, nil
	}
	timer := time.AfterFunc(time.Second*time.Duration(timeoutSecond), func() {
		m.getChan <- "timeout"
	})
	logf.info("Waiting for answer originId=%s\n", id)
	for {
		retId := <-m.getChan
		if retId == "timeout" {
			return nil, fmt.Errorf("Timeout wiating for message answer id=%s", id)
		}
		if mes, ok := m.answerMap[id]; ok {
			logf.info("Found answer originId=%s\n", id)
			timer.Stop()
			return mes, nil
		}
	}
}

func (m *ReceiverManager) receiveMessage(mes *AntMes) bool {
	m.usage++
	logf.debugMes(mes, "recceive message: %s\n", mes.toString())
	if m.nbReceiver <= 0 {
		m.receiver.executeMessage(mes)
		return true
	}
	if m.buffer.put(mes) {
		//logf.info("receive message function=%s duplicate=%d order=%d ok\n", mes.Function, mes.Duplicate, mes.Order)
		return true
	}
	return false
}

func (m *ReceiverManager) stats() {
	fmt.Printf("Receiver: nb=%d maxbuf=%d\n", m.usage, m.buffer.max)
	execVal := ""
	for _, exec := range m.receiverList {
		execVal = fmt.Sprintf("%s %d", execVal, exec.usage)
	}
	fmt.Printf("Receivers: %s\n", execVal)
}

func (m *ReceiverManager) startClientReader(stream GNodeService_GetClientStreamServer) {
	m.lockClient.Lock()
	clientName := fmt.Sprintf("client-%d-%d", time.Now().UnixNano(), m.gnode.clientMap.len()+1)
	m.gnode.clientMap.set(clientName, &gnodeClient{
		name:   clientName,
		stream: stream,
	})
	stream.Send(&AntMes{
		Function:   "ClientAck",
		FromClient: clientName,
	})
	logf.info("Client stream open: %s\n", clientName)
	m.lockClient.Unlock() //unlock far to be sure to have several nano
	for {
		mes, err := stream.Recv()
		if err == io.EOF {
			logf.error("Client reader %s: EOF\n", clientName)
			m.gnode.clientMap.del(clientName)
			m.gnode.removeEventListener(clientName)
			m.gnode.nodeFunctions.forceGC()
			return
		}
		if err != nil {
			logf.error("Client reader %s: Failed to receive message: %v\n", clientName, err)
			m.gnode.clientMap.del(clientName)
			m.gnode.removeEventListener(clientName)
			m.gnode.nodeFunctions.forceGC()
			return
		}
		if mes.Function == "setEventListener" {
			m.gnode.setEventListener(mes.Args[0], mes.Args[1], mes.UserName, clientName)
		} else {
			mes.Id = m.gnode.getNewId(false)
			mes.Origin = m.gnode.name
			mes.FromClient = clientName
			m.gnode.idMap.Add(mes.Id)
			if mes.Debug {
				logf.debugMes(mes, "-------------------------------------------------------------------------------------------------------------\n")
				logf.debugMes(mes, "Receive mes from client %s : %v\n", clientName, mes)
			}
			for {
				if m.gnode.receiverManager.receiveMessage(mes) {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}
