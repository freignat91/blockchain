package tests

import (
	"github.com/freignat91/blockchain/api"
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

//time.Sleep(1*time.Second) need on v0.0.1, will be deleted on v0.0.2

func TestReady(t *testing.T) {
	if err := bcapi.UserSignup("test", "/tmp/blockchain/private.key"); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	if err := bcapi.SetUser("test", "/tmp/blockchain/private.key"); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	if err := bcapi.AddBranch([]string{"organization=company1"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddBranch([]string{"organization=company1", "member=member11"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddBranch([]string{"organization=company1", "member=member12"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddBranch([]string{"organization=company2"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddBranch([]string{"organization=company2", "member=member21"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddBranch([]string{"organization=company2", "member=member22"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddEntry([]byte("doc111"), []string{"organization=company1", "member=member11"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddEntry([]byte("doc112"), []string{"organization=company1", "member=member11"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddEntry([]byte("doc113"), []string{"organization=company1", "member=member11"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddEntry([]byte("doc114"), []string{"organization=company1", "member=member11"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddEntry([]byte("doc115"), []string{"organization=company1", "member=member11"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddEntry([]byte("doc116"), []string{"organization=company1", "member=member11"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddEntry([]byte("doc117"), []string{"organization=company1", "member=member11"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddEntry([]byte("doc118"), []string{"organization=company1", "member=member11"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	if err := bcapi.AddEntry([]byte("doc119"), []string{"organization=company1", "member=member11"}); err != nil {
		t.Fatalf("error: %v\n", err)
	}
}
