package gnode

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"hash"
	"io/ioutil"
)

type GNodeKey struct {
	gnode      *GNode
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	label      []byte
	shaHash    hash.Hash
	nodeKeyMap map[string]*rsa.PublicKey
	userKeyMap map[string]*rsa.PublicKey
}

func (g *GNode) newKey(init bool) (*GNodeKey, error) {
	keyp, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	key := &GNodeKey{
		gnode:      g,
		privateKey: keyp,
		publicKey:  &keyp.PublicKey,
		label:      []byte(""),
		shaHash:    sha256.New(),
	}
	if !init {
		return key, nil
	}
	key.nodeKeyMap = make(map[string]*rsa.PublicKey)
	key.userKeyMap = make(map[string]*rsa.PublicKey)
	return key, nil
}

func (k *GNodeKey) sign(block []byte) ([]byte, error) {
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

func (k *GNodeKey) getPublicKey() ([]byte, error) {
	data, err := x509.MarshalPKIXPublicKey(k.publicKey)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (k *GNodeKey) parsePublicKey(dataKey []byte) (*rsa.PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(dataKey)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PublicKey), nil
}

func (k *GNodeKey) addNodeKey(nodeName string, dataKey []byte) error {
	key, err := k.parsePublicKey(dataKey)
	if err != nil {
		return err
	}
	k.nodeKeyMap[nodeName] = key
	return nil
}

func (k *GNodeKey) addUserKey(userName string, dataKey []byte) error {
	key, err := k.parsePublicKey(dataKey)
	if err != nil {
		return err
	}
	k.userKeyMap[userName] = key
	return nil
}

func (k *GNodeKey) verifyNodeSignature(nodeName string, signature []byte, block []byte) error {
	publicKey, ok := k.nodeKeyMap[nodeName]
	if !ok {
		if nodeName == k.gnode.name {
			publicKey = k.publicKey
		} else {
			return fmt.Errorf("unknown node %s", nodeName)
		}
	}
	return k.verifySignature(publicKey, signature, block)
}

func (k *GNodeKey) verifyUserSignature(userName string, signature []byte, block []byte) error {
	publicKey, ok := k.userKeyMap[userName]
	if !ok {
		return fmt.Errorf("unknown user %s", userName)
	}
	return k.verifySignature(publicKey, signature, block)
}

func (k *GNodeKey) verifySignature(publicKey *rsa.PublicKey, signature []byte, block []byte) error {
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

func (k *GNodeKey) getPrivateKey() []byte {
	data := x509.MarshalPKCS1PrivateKey(k.privateKey)
	_, err := x509.ParsePKCS1PrivateKey(data)
	if err != nil {
		logf.error("error getPrivateKey: %v\n", err)
	}
	ioutil.WriteFile("/data/tmp/key", data, 0777)
	return data
}

func (g *GNode) testKey() {
	key := g.key
	pKey, errg := key.getPublicKey()
	if errg != nil {
		fmt.Printf("Error getPublicKey: %v\n", errg)
	}
	key.addNodeKey("test", pKey)
	block := []byte("essai")
	signature, err := key.sign(block)
	if err != nil {
		fmt.Printf("Error auto sign: %v\n", err)
		return
	}
	if err := key.verifyNodeSignature("test", signature, block); err != nil {
		fmt.Printf("Error auto sign verifie: %v\n", err)
		return
	}
	fmt.Println("Signature ok")
}
