package gnode

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

type RequestStatusManager struct {
	gnode            *GNode
	requestStatusMap map[string]*RequestStatus
	requestList      []*RequestStatus
}

func (m *RequestStatusManager) init(g *GNode) {
	m.gnode = g
	m.requestStatusMap = make(map[string]*RequestStatus)
}

func (m RequestStatusManager) createReqStatus(req *CheckEntryRequest, entry *BCEntry) *RequestStatus {
	s := &RequestStatus{
		Id:       req.Id,
		Date:     entry.Date,
		UserName: entry.UserName,
		Status:   ReqStatusAccepted,
	}
	m.requestStatusMap[req.Id] = s
	return s
}

func (m RequestStatusManager) setReqStatus(id string, status string) *RequestStatus {
	s, exist := m.requestStatusMap[id]
	if !exist {
		s = &RequestStatus{Id: id}
		m.requestStatusMap[id] = s
	}
	s.Date = time.Now().Unix()
	s.Status = status
	return s
}

func (m *RequestStatusManager) setReqStatusError(id string, errorMes string, reqError error) {
	s := m.setReqStatus(id, ReqStatusOnError)
	if reqError != nil {
		s.ReqError = fmt.Sprintf("%s %v", errorMes, reqError)
	} else {
		s.ReqError = errorMes
	}
	logf.error("request %s error: %s\n", id, s.ReqError)
}

func (m *RequestStatusManager) getCreateReqStatus(id string, userName string) *RequestStatus {
	if status, ok := m.requestStatusMap[id]; ok {
		return status
	}
	s := m.setReqStatus(id, ReqStatusOnGoing)
	s.UserName = userName
	return s
}

func (m *RequestStatusManager) sort() {
	m.requestList = []*RequestStatus{}
	for _, status := range m.requestStatusMap {
		m.requestList = append(m.requestList, status)
	}
	sort.Sort(m)
}

func (m RequestStatusManager) Len() int {
	return len(m.requestList)
}

func (m RequestStatusManager) Swap(i, j int) {
	m.requestList[i], m.requestList[j] = m.requestList[j], m.requestList[i]
}

func (m RequestStatusManager) Less(i, j int) bool {
	if m.requestList[i].Date > m.requestList[j].Date {
		return true
	}
	return false
}

func (m *RequestStatusManager) addRequestStatus(mes *AntMes) error {
	if len(mes.Args) == 0 {
		return fmt.Errorf("Invalid argument, needs request id")
	}
	id := mes.Args[0]
	status, exist := m.requestStatusMap[id]
	if !exist {
		return fmt.Errorf("Request not found")
	}
	answer := m.gnode.createAnswer(mes, false)
	answer.Status = status
	m.gnode.senderManager.sendMessage(answer)
	return nil
}

func (m *RequestStatusManager) lastAddRequestStatus(mes *AntMes) error {
	if len(mes.Args) < 3 {
		return fmt.Errorf("Not enough argument, need [nb] [userName] [errorOnly]")
	}
	nb, err := strconv.Atoi(mes.Args[0])
	if err != nil {
		return fmt.Errorf("Invalid argument, should be a number")
	}
	userName := mes.Args[1]
	errorOnly := false
	if mes.Args[2] == "true" {
		errorOnly = true
	}
	m.sort()
	for _, status := range m.requestList {
		if userName == "" || userName == status.UserName {
			if !errorOnly || status.Status == ReqStatusOnError {
				answer := m.gnode.createAnswer(mes, false)
				answer.Status = status
				m.gnode.sendBackClient(answer.FromClient, answer)
				nb--
				if nb == 0 {
					break
				}
			}
		}
	}
	answer := m.gnode.createAnswer(mes, false)
	answer.Args = []string{"end"}
	m.gnode.sendBackClient(answer.FromClient, answer)
	return nil
}
