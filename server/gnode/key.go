package gnode

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"hash"
	"math/big"
)

type gnodeKey struct {
	privateKey   *rsa.PrivateKey
	publicKey    *rsa.PublicKey
	label        []byte
	shaHash      hash.Hash
	publicKeyMap map[string]*rsa.PublicKey
}

func (g *GNode) initKey() error {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	g.key = &gnodeKey{
		privateKey:   key,
		publicKey:    &key.PublicKey,
		label:        []byte(""),
		shaHash:      sha256.New(),
		publicKeyMap: make(map[string]*rsa.PublicKey),
	}
	return nil
}

func (k *gnodeKey) sign(block []byte) ([]byte, error) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto //TODO: to be updated
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(block)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, k.privateKey, newhash, hashed, &opts)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (k *gnodeKey) getPublicKey() *DataKey {
	data := &DataKey{
		Expo: int64(k.publicKey.E),
		Data: k.publicKey.N.Bytes(),
	}
	return data
}

func (k *gnodeKey) addPublicKey(nodeName string, dataKey *DataKey) {
	n := big.NewInt(0)
	n.SetBytes(dataKey.Data)
	publicKey := &rsa.PublicKey{
		E: int(dataKey.Expo),
		N: n,
	}
	k.publicKeyMap[nodeName] = publicKey
}

func (k *gnodeKey) verifySignature(nodeName string, signature []byte, block []byte) error {
	publicKey, ok := k.publicKeyMap[nodeName]
	if !ok {
		return fmt.Errorf("unknown node %s", nodeName)
	}
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto //TODO: to be updated
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(block)
	hashed := pssh.Sum(nil)
	if err := rsa.VerifyPSS(publicKey, newhash, hashed, signature, &opts); err != nil {
		return err
	}
	return nil
}

func (g *GNode) sendPublicKey() {
	g.senderManager.sendMessage(&AntMes{
		Target:   "*",
		Function: "setPublicKey",
		Key:      g.key.getPublicKey(),
	})
}

func (g *GNode) setPublicKey(mes *AntMes) error {
	g.key.addPublicKey(mes.Origin, mes.Key)
	logf.info("Received public key from %s\n", mes.Origin)
	return nil
}

func (g *GNode) testKey() {
	key := g.key
	pKey := key.getPublicKey()
	key.addPublicKey("test", pKey)
	block := []byte("essai")
	signature, err := key.sign(block)
	if err != nil {
		fmt.Printf("Error auto sign: %v\n", err)
		return
	}
	if err := key.verifySignature("test", signature, block); err != nil {
		fmt.Printf("Error auto sign verifie: %v\n", err)
		return
	}
	fmt.Println("Signature ok")
}
