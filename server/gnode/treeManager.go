package gnode

import (
	"fmt"
	"os"
	"path"
	"time"
)

const (
	loadedTreeDepth    = 100
	treeBlockMaxSizeMB = 1000
)

var (
	nbEntryByBlock = 3
)

type TreeManager struct {
	gnode    *GNode
	root     *TreeBlock
	blockMap map[string]*TreeBlockItem
}

type TreeBlockItem struct {
	timestamp time.Time
	block     *TreeBlock
}

func newTreeBlockItem(block *TreeBlock) *TreeBlockItem {
	return &TreeBlockItem{
		timestamp: time.Now(),
		block:     block,
	}
}

func (t *TreeManager) init(g *GNode) error {
	t.gnode = g
	nbEntryByBlock = config.maxEntryNumberPerBlock
	t.blockMap = make(map[string]*TreeBlockItem)
	if err := os.MkdirAll(path.Join(config.rootDataPath, "tree"), 0700); err != nil {
		logf.error("Imposible to create the node tree directory: %v\n", err)
		os.Exit(1)
	}
	logf.info("Load blockchain tree: %s", config.rootDataPath)
	root, err := t.getBlock("root")
	if err != nil {
		logf.warn("root blockchain tree doesn't exist (%v)\n", err)
		logf.warn("initialize blockchain\n")
		r, errn := newRoot()
		if errn != nil {
			logf.error("%v\n", errn)
			os.Exit(1)
		}
		root = r
	}
	t.root = root
	t.blockMap["root"] = newTreeBlockItem(root)
	if err := t.loadBranchesTree("root"); err != nil {
		logf.error("loadTree error: %v\n", err)
	}
	logf.info("blockchain tree ready")
	return nil
}

func (t *TreeManager) getBlockItem(id string) *TreeBlockItem {
	blockItem, ok := t.blockMap[id]
	if ok {
		return blockItem
	}
	return nil
}

func (t *TreeManager) getBlock(id string) (*TreeBlock, error) {
	if blockItem, ok := t.blockMap[id]; ok {
		return blockItem.block, nil
	}
	block := newBlock(id)
	if err := block.load(); err != nil {
		return nil, fmt.Errorf("Load block %s error: %v\n", id, err)
	}
	t.blockMap[id] = newTreeBlockItem(block)
	return block, nil
}

func (t *TreeManager) loadBranchesTree(nodeId string) error {
	block, err := t.getBlock(nodeId)
	if err != nil {
		return err
	}
	if block.Depth > loadedTreeDepth {
		return nil
	}
	for _, branchId := range block.BranchMap {
		t.loadBranchesTree(branchId)
	}
	return nil
}

func (t *TreeManager) addItem(entry *BCEntry, isBranch bool) error {
	logf.info("----------------------------------------------------------------------------------------------------------------\n")
	if isBranch {
		return t.addBranch(entry)
	}
	return t.addEntry(entry)
}

func (t *TreeManager) addBranch(entry *BCEntry) error {
	if len(entry.Labels) == 0 {
		return fmt.Errorf("a branch can't be created without label")
	}
	logf.info("add branch: %v\n", entry.Labels)
	branch, err := t.getLastExistingBranchBlock(entry.Labels, true)
	if err != nil {
		return err
	}
	logf.info("lastExistingBranch: depth=%d label: %s=%s\n", branch.Depth, branch.LabelName, branch.LabelValue)
	depth := int(branch.Depth)
	if _, exist := branch.BranchMap[entry.Labels[depth].Value]; exist {
		return fmt.Errorf("The branch alreday exists")
	}
	block := newBlock("")
	block.Entries = []*BCEntry{entry}
	block.Hash = block.computeBlockHash()
	block.Id = fmt.Sprintf("%x", block.Hash)
	block.setParentBranch(branch, entry.Labels[depth].Value)
	block.Updated = true
	logf.info("add branch: depth=%d label: %s=%s\n", block.Depth, block.LabelName, block.LabelValue)
	return t.updateBlockBranch(block)
}

func (t *TreeManager) addEntry(entry *BCEntry) error {
	logf.info("add entry %s: %v\n", string(entry.Payload), entry.Labels)
	block, errbl := t.getEntryBlock(entry)
	if errbl != nil {
		return errbl
	}
	return t.updateBlockBranch(block)
}

func (t *TreeManager) updateBlockBranch(block *TreeBlock) error {
	if err := t.saveBranch(block); err != nil {
		logf.error("Error saving block %s branch: %v\n", block.Id, err)
		return err
	}
	if err := t.updateAndSaveBranch(block); err != nil {
		logf.error("Error adding block %s: %v\n", block.Id, err)
		if err := t.rollbackBranch(block); err != nil {
			logf.error("Error rollbacking block %s branch: %v\n", block.Id, err)
			return err
		}
	}
	logf.info("block %s saved\n", block.Id)
	return nil
}

func (t *TreeManager) getLastExistingBranchBlock(labels []*TreeLabel, isBranch bool) (*TreeBlock, error) {
	branch, err := t.getLastExistingBranchBlockEff("root", labels, isBranch)
	if err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, fmt.Errorf("The branch doesn't exist")
	}
	depth := int(branch.Depth)
	if isBranch {
		if len(labels) > depth+1 {
			return nil, fmt.Errorf("Impossible to create several branches in the same request")
		}
	} else {
		if len(labels) != depth {
			return nil, fmt.Errorf("the branch doesn't exist")
		}
	}
	return branch, nil
}

func (t *TreeManager) getLastExistingBranchBlockEff(branchId string, labels []*TreeLabel, isBranch bool) (*TreeBlock, error) {
	//load branch block
	branch, err := t.getBlock(branchId)
	if err != nil {
		return nil, err
	}
	depth := int(branch.Depth)
	if branch.SubBranchLabelName == "" {
		if isBranch {
			branch.SubBranchLabelName = labels[depth].Name
		}
		return branch, nil
	}
	// no more label criteria, return the current branch
	if depth >= len(labels) {
		return branch, nil
	}
	if branch.SubBranchLabelName != labels[depth].Name {
		return nil, nil
	}
	childId, exist := branch.BranchMap[labels[depth].Value]
	if !exist {
		if isBranch {
			return branch, nil
		}
		return nil, fmt.Errorf("the branch %s=%s doesn't exist", labels[depth].Name, labels[depth].Value)
	}
	return t.getLastExistingBranchBlockEff(childId, labels, isBranch)
}

func (t *TreeManager) getEntryBlock(entry *BCEntry) (*TreeBlock, error) {
	branch := t.root
	if len(entry.Labels) > 0 {
		b, err := t.getLastExistingBranchBlock(entry.Labels, false)
		if err != nil {
			return nil, err
		}
		branch = b
	}
	logf.info("lastExistingBranch: depth=%d label: %s=%s\n", branch.Depth, branch.LabelName, branch.LabelValue)
	parent, errp := t.getLastBranchBlock(branch)
	if errp != nil {
		return nil, errp
	}
	logf.info("lastExistingBlock: %s having %d entry(ies)\n", parent.Id, len(parent.Entries))
	if len(parent.Entries) < nbEntryByBlock {
		parent.Entries = append(parent.Entries, entry)
		parent.Size = int64(len(parent.Entries))
		parent.Hash = parent.computeBlockHash()
		logf.info("add entry in existing block id=%s\n", parent.Id)
		return parent, nil
	}
	block := newBlock("")
	block.Entries = []*BCEntry{entry}
	block.Hash = block.computeBlockHash()
	block.Id = fmt.Sprintf("%x", block.Hash)
	block.setParent(parent)
	block.Updated = true
	logf.info("add entry in new entry block id=%s\n", block.Id)
	return block, nil
}

func (t *TreeManager) getLastBranchBlock(block *TreeBlock) (*TreeBlock, error) {
	childId := block.ChildId
	for childId != "" {
		bl, err := t.getBlock(childId)
		if err != nil {
			return nil, err
		}
		block = bl
		childId = block.ChildId
	}
	return block, nil
}

func (t *TreeManager) saveBranch(block *TreeBlock) error {
	//TODO
	logf.info("push branch %s\n", block.Id)
	return nil
}

func (t *TreeManager) updateAndSaveBranch(block *TreeBlock) error {
	if err := block.saveBlockFullHash(t); err != nil {
		return err
	}
	parentId := block.ParentId
	for parentId != "" {
		block, err := t.getBlock(parentId)
		if err != nil {
			return err
		}
		if err := block.saveBlockFullHash(t); err != nil {
			return err
		}
		parentId = block.ParentId
	}
	return nil
}

func (t *TreeManager) rollbackBranch(block *TreeBlock) error {
	//TODO
	logf.info("Rollback branch %s\n", block.Id)
	return nil
}

func (t *TreeManager) getTree(mes *AntMes) error {
	fmt.Printf("getTree: %v\n", mes.Args)
	isBlocks := false
	if mes.Args[0] == "true" {
		isBlocks = true
	}
	isEntries := false
	if mes.Args[1] == "true" {
		isEntries = true
	}
	t.getTreeBranches(mes, "root", isBlocks, isEntries, "branch")
	answer := t.gnode.createAnswer(mes, false)
	answer.Args = []string{"end"}
	t.gnode.sendBackClient(answer.FromClient, answer)
	return nil
}

func (t *TreeManager) getTreeBranches(mes *AntMes, blockId string, isBlocks bool, isEntries bool, blockType string) {
	block := t.getTreeBlock(mes, blockId, isBlocks, isEntries, blockType)
	if block != nil {
		if isBlocks {
			childId := block.ChildId
			if childId != "" {
				t.getTreeBranches(mes, childId, isBlocks, isEntries, "block")
			}
		}
		for _, branchId := range block.BranchMap {
			t.getTreeBranches(mes, branchId, isBlocks, isEntries, "branch")
		}
	}
}

func (t *TreeManager) getTreeBlock(mes *AntMes, blockId string, isBlocks bool, isEntries bool, blockType string) *TreeBlock {
	block, err := t.getBlock(blockId)
	if err != nil {
		answer := t.gnode.createAnswer(mes, false)
		answer.ErrorMes = fmt.Sprintf("%v", err)
		answer.Args = []string{blockId, blockType}
		t.gnode.sendBackClient(answer.FromClient, answer)
		return nil
	}
	answer := t.gnode.createAnswer(mes, false)
	answer.Block = block
	answer.Args = []string{blockId, blockType}
	t.gnode.sendBackClient(answer.FromClient, answer)
	return block
}
