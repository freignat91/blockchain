package tests

import (
	"fmt"
	"github.com/freignat91/blockchain/api"
	"math/rand"
	"os"
	"testing"
	"time"
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

var branches = [][]string{}

//time.Sleep(1*time.Second) needed on v0.0.1, will be deleted on v0.0.2

func TestReady(t *testing.T) {
	branches = append(branches, []string{"organization=company1", "member=member11"})
	branches = append(branches, []string{"organization=company1", "member=member12"})
	branches = append(branches, []string{"organization=company2", "member=member21"})
	branches = append(branches, []string{"organization=company2", "member=member22"})
	user := fmt.Sprintf("user%d", time.Now().UnixNano())
	key := fmt.Sprintf("/tmp/blockchain/%s.key", user)
	if err := bcapi.UserSignup(user, key); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	if err := bcapi.SetUser(user, key); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	for nn := 0; nn < 100; nn++ {
		labels := branches[rand.Int31n(int32(len(branches)))]
		entry := fmt.Sprintf("doc-%s-%d", user, nn)
		if err := bcapi.AddEntry([]byte(entry), labels); err != nil {
			t.Fatalf("error: %v\n", err)
		}
	}
}
