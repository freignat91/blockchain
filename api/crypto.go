package api

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

func (api *BchainAPI) sign(payload []byte) ([]byte, error) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto //TODO: to be updated
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(payload)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, api.key, newhash, hashed, &opts)
	if err != nil {
		return nil, err
	}
	return signature, nil
}
