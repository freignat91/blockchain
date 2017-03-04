package gnode

import (
	"fmt"
	"io/ioutil"
	"path"

	proto "github.com/golang/protobuf/proto"
)

func newRoot() (*TreeBlock, error) {
	//create new bock id=root
	root := newBlock("root")
	//create dedicated first entry for root
	entry := &BCEntry{
		Payload: []byte("branch:root"),
	}
	//sign entry and added it in root block
	data, err := proto.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("marshaling root error: ", err)
	}
	hash := getNewHash()
	hash.setHash(data)
	entry.Hash = hash.hash
	root.Entries = append(root.Entries, entry)
	//sign the root block
	root.Hash = root.computeBlockHash()
	//save the root block
	if err := root.save(); err != nil {
		return nil, fmt.Errorf("Impossible to save the blockchain root: %v\n", err)
	}
	return root, nil
}

func newBlock(id string) *TreeBlock {
	return &TreeBlock{
		Id:        id,
		Size:      1,
		BranchMap: make(map[string]string),
		Entries:   []*BCEntry{},
	}
}

func (n *TreeBlock) load() error {
	if n.Loaded {
		return nil
	}
	data, err := ioutil.ReadFile(path.Join(config.rootDataPath, "tree", n.Id))
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(data, n); err != nil {
		return fmt.Errorf("unmarshaling block %s error: %v\n", n.Id, err)
	}
	if n.BranchMap == nil {
		n.BranchMap = make(map[string]string)
	}
	n.Loaded = true
	return nil
}

func (n *TreeBlock) save() error {
	data, err := proto.Marshal(n)
	if err != nil {
		return fmt.Errorf("marshaling block error: ", err)
	}
	if err := ioutil.WriteFile(path.Join(config.rootDataPath, "tree", n.Id), data, 0600); err != nil {
		return err
	}
	n.Updated = false
	return nil
}

func (n *TreeBlock) setParent(parent *TreeBlock) {
	n.ParentId = parent.Id
	n.Depth = parent.Depth
	parent.ChildId = n.Id
}

func (n *TreeBlock) setParentBranch(parent *TreeBlock, labelValue string) {
	parent.BranchMap[labelValue] = n.Id
	n.LabelName = parent.SubBranchLabelName
	n.LabelValue = labelValue
	n.ParentId = parent.Id
	n.Depth = parent.Depth + 1
	logf.info("add branch %s at parent branch value %s: %s\n", n.Id, labelValue, parent.Id)
}

func (n *TreeBlock) computeBlockHash() []byte {
	hash := getNewHash()
	hash.setBlockHash(n)
	return hash.hash
}

func (n *TreeBlock) saveBlockFullHash(t *TreeManager) error {
	hash := getNewHash()
	if err := hash.setBlockFullHash(t, n); err != nil {
		return err
	}
	n.FullHash = hash.hash
	logf.info("save block %s (%x:%d)\n", n.Id, n.FullHash, n.Size)
	return n.save()
}
