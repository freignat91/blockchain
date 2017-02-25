package main

import (
	"fmt"
	"github.com/freignat91/blockchain/server/node"
	"log"
	"net/http"
	"os"
)

// build vars
var (
	Version string
	Build   string
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "healthcheck" {
		fmt.Println("health")
		if !healthcheck() {
			os.Exit(1)
		}
		os.Exit(0)
	}
	gnodeServer := node.GNode{}
	err := gnodeServer.Start(Version, Build)
	if err != nil {
		log.Printf("Exit on init error: %v\n", err)
		os.Exit(1)
	}
}

func healthcheck() bool {
	log.Println("main healthcheck")
	response, err := http.Get("http://127.0.0.1:3000/api/v1/health")
	if err != nil {
		return false
	}
	if response.StatusCode == 200 {
		return true
	}
	return false
}
