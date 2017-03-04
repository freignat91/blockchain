package tests

import (
	"github.com/freignat91/blockchain/api"
	"os"
	"testing"
)

const (
	server = "127.0.0.1:30103"
)

var bcapi *api.BchainAPI

func TestMain(m *testing.M) {
	bcapi = api.New(server)
	bcapi.SetLogLevel("info")
	os.Exit(m.Run())
}

func TestReady(t *testing.T) {
	bcapi.GetNbReady(3)
}
