package gnode

import (
	"crypto/sha256"
	"sort"
)

const hashLength = 32

type TreeHash struct {
	hash []byte
}

func getNewHash() *TreeHash {
	return &TreeHash{hash: make([]byte, hashLength, hashLength)}
}

func (h *TreeHash) setHash(data []byte) {
	hash := sha256.Sum256(data)
	for i, b := range hash {
		h.hash[i] = b
	}
}

func (h *TreeHash) setBlockHash(block *TreeBlock) {
	data := make([]byte, hashLength*len(block.Entries), hashLength*len(block.Entries))
	nn := 0
	for _, entry := range block.Entries {
		nn = h.copyHash(data, nn, entry.Hash)
	}
	h.setHash(data)
}

func (h *TreeHash) setBlockFullHash(t *TreeManager, block *TreeBlock) error {
	data := make([]byte, hashLength*(len(block.BranchMap)+3), hashLength*(len(block.BranchMap)+3))
	var size int64 = 1
	nn := h.copyHash(data, 0, block.Hash)
	if block.ChildId == "" {
		nn += hashLength
	} else {
		child, err := t.getBlock(block.ChildId)
		if err != nil {
			return err
		}
		size += child.Size
		nn = h.copyHash(data, nn, child.FullHash)
	}
	list := []string{}
	for _, branchId := range block.BranchMap {
		list = append(list, branchId)
	}
	sort.Strings(list)
	for _, branchId := range list {
		branch, err := t.getBlock(branchId)
		if err != nil {
			return err
		}
		size += branch.Size
		nn = h.copyHash(data, nn, branch.FullHash)
	}
	block.Size = size
	h.setHash(data)
	return nil
}

func (h *TreeHash) copyHash(data []byte, nn int, item []byte) int {
	for i, b := range item {
		data[nn+i] = b
	}
	return len(item) + nn
}
