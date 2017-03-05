.PHONY: all clean install fmt check version build run test

SHELL := /bin/sh
BASEDIR := $(shell echo $${PWD})

# build variables (provided to binaries by linker LDFLAGS below)
VERSION := 0.0.2

LDFLAGS=-ldflags "-X=main.Version=$(VERSION)"

# ignore vendor directory for go files
SRC := $(shell find . -type f -name '*.go' -not -path './vendor/*' -not -path './.git/*')

# for walking directory tree (like for proto rule)
DIRS = $(shell find . -type d -not -path '.' -not -path './vendor' -not -path './vendor/*' -not -path './.git' -not -path './.git/*')

# generated files that can be cleaned
GENERATED := $(shell find . -type f -name '*.pb.go' -not -path './vendor/*' -not -path './.git/*')

# ignore generated files when formatting/linting/vetting
CHECKSRC := $(shell find . -type f -name '*.go' -not -name '*.pb.go' -not -path './vendor/*' -not -path './.git/*')

OWNER := freignat91
NAME :=  blockchain
TAG := latest

IMAGE := $(OWNER)/$(NAME):$(TAG)
IMAGETEST := $(OWNER)/$(NAME):test
REPO := github.com/$(OWNER)/blockchain

CLIENT := bchain
SERVER := server
TESTS := tests

all: version check install

version: 
	@echo "version: $(VERSION)"

clean: 
	@rm -rf $(GENERATED)

install-client: 
	@go install $(LDFLAGS) $(REPO)/$(CLIENT)

install-server: 
	@go install $(LDFLAGS) $(REPO)/$(SERVER)  

install: install-server install-client

proto: 
	@protoc server/gnode/gnode.proto --go_out=plugins=grpc:.

# format and simplify if possible (https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
fmt:
	@gofmt -s -l -w $(CHECKSRC)

check:
	@test -z $(shell gofmt -l ${CHECKSRC} | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d} | sed '/pb\.go/d'; done
	@go tool vet ${CHECKSRC}

build:	install-client
	@docker build -t $(IMAGE) .

buildtest: install-client
	@docker build -t $(IMAGETEST) .


run: 	build
        @CID=$(shell docker run --net=host -d --name $(NAME) $(IMAGE)) && echo $${CID}

test:   
	@go test ./tests -v

testp:
	@go test ./testParallel -v

install-deps:
	@glide install --strip-vcs --strip-vendor --update-vendored

update-deps:
	@glide update --strip-vcs --strip-vendor --update-vendored

start:
	@docker node inspect self > /dev/null 2>&1 || docker swarm inspect > /dev/null 2>&1 || (echo "> Initializing swarm" && docker swarm init --advertise-addr 127.0.0.1)
	@docker network ls | grep aNetwork || (echo "> Creating overlay network 'aNetwork'" && docker network create -d overlay aNetwork)
	@mkdir -p /tmp/blockchain/data
	@docker service create --network aNetwork --name antblockchain \
	--publish 30103:30103 \
	--mount type=bind,source=/tmp/blockchain/data,target=/data \
	--replicas=5 \
	$(IMAGE)


starttest:
	@docker node inspect self > /dev/null 2>&1 || docker swarm inspect > /dev/null 2>&1 || (echo "> Initializing swarm" && docker swarm init --advertise-addr 127.0.0.1)
	@docker network ls | grep aNetwork || (echo "> Creating overlay network 'aNetwork'" && docker network create -d overlay aNetwork)
	@mkdir -p /tmp/blockchain/data
	@docker service create --network aNetwork --name antblockchain \
	--publish 30103:30103 \
	--mount type=bind,source=/tmp/blockchain/data,target=/data \
	--replicas=5 \
	$(IMAGETEST)

stop:
	@docker service rm antblockchain || true

init:
	@docker service rm antblockchain || true
	@rm -f ./logs/*
