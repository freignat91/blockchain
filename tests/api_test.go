package tests

import (
	"github.com/freignat91/blockchain/api"
	"os"
	"testing"
)

const (
	server = "127.0.0.1:30103"
)

var bcApi *api.BchainAPI

func TestMain(m *testing.M) {
	bcApi = api.New(server)
	bcApi.SetLogLevel("info")
	os.Exit(m.Run())
}
