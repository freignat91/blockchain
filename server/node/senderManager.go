package node

import (
	"fmt"
)

type SenderManager struct {
	usage      int
	gnode      *GNode
	buffer     MessageBuffer
	senderList []*MessageSender
	ioChan     chan *AntMes
	nbSender   int
	sender     MessageSender
}

func (m *SenderManager) start(gnode *GNode, bufferSize int, maxGoRoutine int) {
	maxGoRoutine = 0
	m.gnode = gnode
	m.nbSender = maxGoRoutine
	m.buffer.init(bufferSize)
	m.ioChan = make(chan *AntMes)
	if maxGoRoutine <= 0 {
		m.sender.gnode = gnode
		return
	}
	m.senderList = []*MessageSender{}
	for i := 0; i < maxGoRoutine; i++ {
		routine := &MessageSender{
			id:            i,
			gnode:         m.gnode,
			senderManager: m,
		}
		m.senderList = append(m.senderList, routine)
		routine.start()
	}
	go func() {
		for {
			mes, ok := m.buffer.get(true)
			if ok && mes != nil {
				//log.Printf("Message ack %s\n", mes.Id)
				m.ioChan <- mes
			}
		}
	}()
}

func (m *SenderManager) sendMessage(mes *AntMes) error {
	m.usage++
	if mes.Id == "" {
		mes.Id = m.gnode.getNewId(true)
	}
	if mes.Origin == "" {
		mes.Origin = m.gnode.name
	}
	if err := m.sender.sendMessage(mes); err != nil {
		logf.error("send message error: %s: %v\n", mes.toString(), err)
	}
	//logf.info("send message: %s\n", mes.toString())
	return nil
}

func (m *SenderManager) sendMessageReturnAnswer(mes *AntMes, timeoutSecond int) (*AntMes, error) {
	mes.Id = m.gnode.getNewId(true)
	mes.Origin = m.gnode.name
	mes.AnswerWait = true
	m.sendMessage(mes)
	ret, err := m.gnode.receiverManager.waitForAnswer(mes.Id, timeoutSecond)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (m *SenderManager) stats() {
	fmt.Printf("Sender: nb=%d maxbuf=%d\n", m.usage, m.buffer.max)
	execVal := ""
	for _, exec := range m.senderList {
		execVal = fmt.Sprintf("%s %d", execVal, exec.usage)
	}
	fmt.Printf("senders: %s\n", execVal)
}
