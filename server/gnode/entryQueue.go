package gnode

import (
	"fmt"
	"sync"
	"time"
)

type requestQueuer struct {
	gnode         *GNode
	manager       *EntryManager
	stackIndexIn  int
	stackIndexOut int
	stackSize     int
	requestStack  []*CheckEntryRequest
	lock          sync.RWMutex
}

func (q *requestQueuer) init(g *GNode, m *EntryManager) {
	q.lock = sync.RWMutex{}
	q.gnode = g
	q.manager = m
	q.stackIndexIn = -1
	q.stackIndexOut = -1
	q.stackSize = 0
	q.requestStack = make([]*CheckEntryRequest, config.entryStackSize, config.entryStackSize)
	q.startReader()
}

func (q *requestQueuer) push(req *CheckEntryRequest) error {
	q.lock.Lock()
	if q.incrIn() {
		q.requestStack[q.stackIndexIn] = req
	} else {
		return fmt.Errorf("Max number of waiting add request reached, re-try later")
	}
	logf.info("in index=%d size=%d\n", q.stackIndexIn, q.stackSize)
	q.lock.Unlock()
	return nil
}

func (q *requestQueuer) incrIn() bool {
	if q.stackSize >= config.entryStackSize {
		return false
	}
	q.stackIndexIn++
	q.stackSize++
	if q.stackIndexIn >= config.entryStackSize {
		q.stackIndexIn = 0
	}
	return true
}

func (q *requestQueuer) startReader() {
	go func() {
		for {
			if q.incrOut() {
				req := q.requestStack[q.stackIndexOut]
				req.RootHash = q.manager.computeRootHash()
				time.Sleep(1 * time.Second)
				//Verify request validity
				logf.info("Send req id=%s labels=%v\n", req.Id, req.Labels)
				if _, errc := q.gnode.treeManager.getLastExistingBranchBlock(req.Labels, req.IsBranch); errc == nil {
					if err := q.manager.sendCheckEntry(req, ""); err != nil {
						logf.error("sendCheckEntry error: %v\n", err)
					}
				} else {
					logf.error("entry control error: %v\n", errc)
				}
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func (q *requestQueuer) incrOut() bool {
	q.lock.Lock()
	if q.stackSize <= 0 {
		q.lock.Unlock()
		return false
	}
	q.stackIndexOut++
	q.stackSize--
	if q.stackIndexOut >= config.entryStackSize {
		q.stackIndexOut = 0
	}
	logf.info("out index=%d size=%d\n", q.stackIndexOut, q.stackSize)
	q.lock.Unlock()
	return true
}
