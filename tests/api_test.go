package tests

import (
	"github.com/freignat91/blockchain/api"
	"os"
	"testing"
)

const (
	server = "127.0.0.1:30103"
)

var api *api.BchainAPI

func TestMain(m *testing.M) {
	tapi = api.New(server)
	api.SetLogLevel("info")
	os.Exit(m.Run())
}
