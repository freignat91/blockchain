package gnode

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"encoding/pem"
)

func (n *nodeFunctions) createUser(mes *AntMes) error {
	if len(mes.Args) < 1 {
		return fmt.Errorf("Number of argument error, need userName")
	}
	userName := mes.Args[0]
	logf.info("Received create user %s\n", userName)
	key, err := n.gnode.newKey(false)
	if err != nil {
		return err
	}
	data, errg := key.getPublicKey()
	if errg != nil {
		return errg
	}
	if err := n.gnode.createUser(userName, data); err != nil {
		return err
	}
	keyData, errg := key.getPublicKey()
	if errg != nil {
		return errg
	}
	n.gnode.senderManager.sendMessage(&AntMes{
		Target:   "*",
		Origin:   n.gnode.name,
		Function: "createNodeUser",
		UserName: userName,
		Key:      keyData,
	})
	answer := n.gnode.createAnswer(mes, true)
	answer.Key = key.getPrivateKey()
	answer.UserName = userName
	n.gnode.senderManager.sendMessage(answer)
	return nil
}

func (n *nodeFunctions) createNodeUser(mes *AntMes) error {
	userName := mes.UserName
	key := mes.Key
	logf.info("Received create node user %s\n", userName)
	err := n.gnode.createUser(userName, key)
	if err != nil {
		logf.error("createUser error: user=%s: %v\n", userName, err)
		return err
	}
	return nil
}

func (g *GNode) createUser(userName string, publicKey []byte) error {
	logf.info("Create user %s\n", userName)
	_, err := ioutil.ReadDir(path.Join(config.rootDataPath, "users", userName))
	if err == nil {
		g.loadOneUser(userName)
		return fmt.Errorf("User %s : already exist", userName)
	}
	os.MkdirAll(path.Join(config.rootDataPath, "users"), os.ModeDir)
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKey,
	})
	ioutil.WriteFile(path.Join(config.rootDataPath, "users", userName), pubBytes, 0644)
	g.loadOneUser(userName)
	logf.info("Save user: %s\n", userName)
	return nil
}

func (n *nodeFunctions) removeUser(mes *AntMes) error {
	userName := mes.UserName
	err := n.gnode.removeUser(userName)
	if err != nil {
		return err
	}
	n.gnode.senderManager.sendMessage(&AntMes{
		Target:   "*",
		Origin:   n.gnode.name,
		Function: "removeNodeUser",
		UserName: userName,
	})
	answer := n.gnode.createAnswer(mes, true)
	n.gnode.senderManager.sendMessage(answer)
	return nil
}

func (n *nodeFunctions) removeNodeUser(mes *AntMes) error {
	userName := mes.UserName
	err := n.gnode.removeUser(userName)
	if err != nil {
		return err
	}
	return nil
}

func (g *GNode) removeUser(userName string) error {
	logf.warn("Remove user %s\n", userName)
	os.Remove(path.Join(config.rootDataPath, "users", userName, "key"))
	return nil
}

func (g *GNode) loadUser() error {
	logf.printf("Load users:\n")
	fileList, err := ioutil.ReadDir(path.Join(config.rootDataPath, "users"))
	if err != nil {
		return err
	}
	for _, fd := range fileList {
		g.loadOneUser(fd.Name())
	}
	logf.printf("End load users\n")
	return nil
}

func (g *GNode) loadOneUser(name string) {
	data, err := ioutil.ReadFile(path.Join(config.rootDataPath, "users", name))
	if err != nil {
		logf.error("loadOneUser error user=%s: %v\n", name, err)
		return
	}
	block, _ := pem.Decode(data)
	if block == nil {
		logf.error("Error public of user %s key file is not a pem format", name)
		return
	}
	if err := g.key.addUserKey(name, block.Bytes); err != nil {
		logf.error("parse user %s publicKey error: %v\n", name, err)
		return
	}
	logf.info("Add user %s\n", name)
}

func (g *GNode) checkUser(user string, token string) bool {
	//TODO
	return true
}
