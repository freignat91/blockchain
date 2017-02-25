package api

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const (
	LOG_ERROR = 0
	LOG_WARN  = 1
	LOG_INFO  = 2
	LOG_DEBUG = 3
)

type BchainAPI struct {
	serverList  []string
	serverIndex int
	logLevel    int
	userName    string
	key         *rsa.PrivateKey
}

// New create an blockchain api instance
func New(servers string) *BchainAPI {
	serverList := strings.Split(servers, ",")
	for i, serv := range serverList {
		serverList[i] = strings.Trim(serv, " ")
	}
	api := &BchainAPI{
		serverList: serverList,
		logLevel:   LOG_WARN,
	}
	return api
}

func (api *BchainAPI) getNextServerAddr() string {
	addr := api.serverList[api.serverIndex]
	api.serverIndex++
	if api.serverIndex >= len(api.serverList) {
		api.serverIndex = 0
	}
	return addr
}

func (api *BchainAPI) getClient() (*gnodeClient, error) {
	client := gnodeClient{}
	err := client.init(api)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (api *BchainAPI) SetLogLevel(level string) {
	if strings.ToLower(level) == "error" {
		api.logLevel = LOG_ERROR
	} else if strings.ToLower(level) == "warn" {
		api.logLevel = LOG_WARN
	} else if strings.ToLower(level) == "info" {
		api.logLevel = LOG_INFO
	} else if strings.ToLower(level) == "debug" {
		api.logLevel = LOG_DEBUG
	}
}

func (api *BchainAPI) LogLevelString() string {
	switch api.logLevel {
	case LOG_ERROR:
		return "error"
	case LOG_WARN:
		return "warn"
	case LOG_INFO:
		return "info"
	case LOG_DEBUG:
		return "debug"
	default:
		return "?"
	}
}

func (api *BchainAPI) error(format string, args ...interface{}) {
	if api.logLevel >= LOG_ERROR {
		log.Printf(format, args...)
	}
}

func (api *BchainAPI) warn(format string, args ...interface{}) {
	if api.logLevel >= LOG_WARN {
		log.Printf(format, args...)
	}
}

func (api *BchainAPI) info(format string, args ...interface{}) {
	if api.logLevel >= LOG_INFO {
		log.Printf(format, args...)
	}
}

func (api *BchainAPI) debug(format string, args ...interface{}) {
	if api.logLevel >= LOG_DEBUG {
		log.Printf(format, args...)
	}
}

func (api *BchainAPI) isDebug() bool {
	if api.logLevel >= LOG_DEBUG {
		return true
	}
	return false
}

func (api *BchainAPI) printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// SetUser define the current user
func (api *BchainAPI) SetUser(user string, keyPath string) error {
	//fmt.Printf("setUser: %s path:%s\n", user, keyPath)
	api.userName = user
	data, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("Read private key error: %v\n", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return fmt.Errorf("Error private key file is not a pem format")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("Error parsing private key: %v\n", err)
	}
	api.key = key
	return nil
}

// UserCreate create an user and return a token
func (api *BchainAPI) UserSignup(name string) error {
	if err := api.verifyUserName(name); err != nil {
		return fmt.Errorf("Invalide user name: %v", err)
	}
	client, err := api.getClient()
	if err != nil {
		return err
	}
	ret, errs := client.createSendMessage("", true, "createUser", name)
	if errs != nil {
		return errs
	}
	data := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: ret.Key,
		},
	)
	ioutil.WriteFile(fmt.Sprintf("./%s.key", name), data, 0644)
	return nil
}

func (api *BchainAPI) verifyUserName(name string) error {
	if name == "" || name == "common" {
		return fmt.Errorf("Invalid user name")
	}
	if strings.IndexAny(name, " /\\") >= 0 {
		return fmt.Errorf("Invalid character")
	}
	return nil
}

// UserRemove create an user
func (api *BchainAPI) UserRemove(name string) error {
	defer func() {
		api.userName = ""
		api.key = nil
	}()
	client, err := api.getClient()
	if err != nil {
		return err
	}
	_, errs := client.createSendMessage("*", true, "removeUser")
	if errs != nil {
		return errs
	}
	return nil
}

func (api *BchainAPI) AddEntry(entry []byte, args []string) error {
	client, err := api.getClient()
	if err != nil {
		return err
	}
	mes, errc := client.createSignedMessage("", true, "addEntry", entry, args...)
	if errc != nil {
		return errc
	}
	ret, errs := client.sendMessage(mes, true)
	if errs != nil {
		return errs
	}
	if ret.ErrorMes != "" {
		return fmt.Errorf(ret.ErrorMes)
	}
	return nil
}
