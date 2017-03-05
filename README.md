# AntBlockChain

AntBlockChain release 0.0.1 (version 0.0.2 on going)

# Purpose

This project is an experimental project aiming to build a full secured block chain service implementing the following ideas:

- lower the needed node interconnections using ants behavior to handle nodes communication (as Agrid project https://github.com/freignat91/agrid), keeping the full block chain security and integrity.

- be able to do not have the full block-chain stored on each node, but have blockchain branches creation/append on sub node network keeping their security and integrity. These branches stays related to block chain root, but stored on a random part of the network. a node can this way handle several branches, but not all (all should stay possible). It's a kind of sharding.

- let open the payload stored in the blockchain and have capability to treat it as an extensible object with its own logical behavior. For instance if the payload is a financial transaction, one of its logical behavior could be to do not have the same user spending two times the same amount of money. The goal is to have a clear separation between the block-chain way of working and the object stored own logical constraints.

- technically developed in Go, using Docker services and able to dynamically interconnected several AntBlocChain services to extend the network with open services interconnection topologies.


This project will incrementally move to v1.0.0 if it is prove to work as expected, stop if not.


# Install

It needsfirst Docker 1.13 installed, then

- clone the git project: https://github.com/freignat91/blockchain
- execute `make install` to build the antblockchain CLI executable
- execute `make build` to create a image freignat91/blockchain:latest
- exec
- execute the command `make start` to initialize swarm, create an overlay network and start the antblockchain service on this network
- it starts a 30 nodes blockchain services locally

antblockchain can't be used as a single container, it needs to be started as a service on a swarm machine (manager or worker).

an docker image is available on docker hub: freignat91/blockchain:latest

antblockchain take more time to start than the "ready" docker status. Every CLI command executed before will be rejected with a message "Node not yet ready"
The anblockchain is really ready when the CLI command `bchain node info` display the 30 nodes at each call


# Configuration using System Variables:


- GRPCPORT:               grpc server port used by nodes (default 30103)
- NB_LINE_CONNECT:        number of "line" type connection in grid: default 0 means computed automatically
- NB_CROSS_CONNECT:       number of "cross" type connection in grid: default 0 means computed automatically
- DATA_PATH:              path in container where file data is stored: default: /data (should be mapped on host using mount docker argument (--mount 
type=bind,source=/[hostpath],target=/data)
- MAX_ENTRIES_NB_PER_BLOCK: max entries number per block, for debug reason default is 3



# Resilience

For resilience reason, it's mandatory to have a separated disk file system for each node (each node on its own VM), but for test reason it's possible to have several nodes on the same file system or have architecture with nodes spreaded on a less number of several VMs.


## Grid simulation

To simulate nodes connections using different parameters as, node number, line connections, cross connections,and see the grid topology, use the cli command:

`bchain grid simul [nodes] <--line xx > <--cross yy>`
- [nodes] the number of nodes
- <--line xx> optionally: xx the number of line connections 
- <--cross yy> optionally: yy the number of cross connections 

this command as not effect on the real cluster grid connections, see:
- ./docs/grid-building.pptx
- ./docs/ant-net.pptx

# CLI

AntBlockChain command lines is implemented using the AntBlockchain Go API

### common options

- --verbose: display more informations messages
- --server: format addr1:port,addr2:port, ...   list of the cluster servers (can list less servers than really in the cluster, just one for instance).
- --user: user name, default in ~/.config/antblockchain/blockchain.yaml config file, key: userName
- --key: privateKey file path, default in ~/.config/antblockchain/blockchain.yaml config file, key: keyPath

### create a user <--key-path>
`bchain user signup [username]`

Create a user and return a file containing his private key

- user: the user name to create
- --key-path: optionnal file name path in witch the key is stored, default is: ./private.key

The private key file path and user name can be set in the conffile file: ~/.config/antblockchain/blockchain.yaml, to do have to add it for all commands, as for instance:

```
username: aUserName
keypath: ~/.config/antblockchain/private.key

```

So the private key can be anywhere including in a temporary mounted USB key for more security

### remove a user
`bchain user remove [user]`

Remove a user.

- [user] the user name to remove


### list the cluster nodes and display information

`bchain node info`

Display the full list of nodes with their names, root hash, number of users, number of blockchain entries

### ping a cluster node

`bchain node ping |node]`
- [node] the node name to ping

### add a new branch in the blockchain 

`bchain add branch <label1=value1> <label2=value2> ...`

Validate and add add a branch in the blockchain tree

- labeln=valuen: the labels defined the new branch to add. Labels should be set in order, see doc: ./docs/blochchain-tree.pptx

### add a new entry in the blockchain (on going, not yet working)

`bchain add entry payload <label1=value1> <label2=value2> ...`

Validate and add in the blockchain tree a new blockchain entry containing the payload

- payload: on the version 0.0.1 the payload is a string
- labeln=valuen: the labels defined the branch on swhich the payload as to be added. Label should be set in order, see doc: ./docs/blochchain-tree.pptx


### display the blockchain tree

`bchain display label1:val1 ... labeln:valn <--blocks> <--entries> <--debug> <--hash>`

display the blockchain tree accordinling to the follwing options:
- default: display blockchain tree banches starting from the branch corresponding to the labels: label1:val1 ... labeln:valn
- label1:val1 ... labeln:valn, the branch to display, default: root
- --blocks: display the blocks under the branches
- --entries: display the entries in the blocks
- --debug: display blocks debug information (child id, paerent id, ...)
- --hash: display blocls hash instead of id


### display a "add request" status 

`bchain add status [id]`

retrieve the status of the "add request" having the id [id]


### display last "add request" status 

`bchain add last [nb] <--userName [userName]> <--errorOnly>`

retrieve the [nb] last "add request" status
Arguments:
- nb: number of returned status
- --userName return only the status belonging to the user UserName
- --errorOnly: return only the status on error


# API

AntBlockchain is usable using Go api API github.com/freignat91/blockchain/api

### Usage

```
        import "github.com/freignat91/blockchain/api"
        ...
        bcApi := api.New("localhost:30103")
        bcApi.setUser("aUser", "~/.config/antblockchain/private.key")
        list, err := bcApi.NodeInfo()
        ...
```

### func (api *BchainAPI) UserSignup(name string, keyPath string) error

Create a new user and write the privateKey to authenticate the user in the file keyPath

Argument
- name: the user name to create
- keyPath: the created user private key file


### func (api *BchainAPI) UserRemove(name string) error

Remove a user

Argument
- name: the user name to remove, format userName:token


### func (api *BchainAPI) SetUser(user string, keyPath string) error

Set the current user and authenticate it with his privateKey. This action is mandatory before any other command on the blockchain.

Arguments:
- user: the user name
- keyPath: the path of the file where the user's privateKey is. It can be the temporary mounted USB key for more security


### func (api *AgridAPI) NodePing(node string, debugTrace bool) (string, error)

Ping a node
Arguments:
- node: node name to ping
- debugTrace: if true, trace the message especially in the node logs.

### func (api *AgridAPI) NodeInfo() ([]string, error)

List the node of the cluster with information

### func (api *BchainAPI) GetTree(labels []string, blocks bool, entries bool, callback interface{}) error 

get the whole blockchain tree of part of the tree and execute a callback function for each block

- labels: the branch of the tree the function starts to read, if empty, the whole tree is read
- blocks: if true get all blocks, if false get only branch blocks
- entries: not used, on this version all the entries in the blocks are read at the same time
- callback function call at each block read: `function(id string, blockType string, block *gnode.TreeBlock) error` where blocktype can be either "branch" or "block" (see cli_display.go file as usage sample on it)

### func (api *BchainAPI) AddRequestStatus(id string) (*gnode.RequestStatus, error) 

get the requestStatus having the id "id" and return an instance of RequestStatus
Argument:
- id, the request id (returned by function addEntry or addBranch)

### func (api *BchainAPI) LastAddRequestStatus(nb int, userName string, errorOnly bool, callback interface{}) error

get the nb last add requests, for each execute a callback function
Arguments:
- nb: number of request to get
- userName: if != "", return only the status beloging to the user
- errorOnly: if true, return only the status on error
- callback function exectuted for each status, proto: function(status *RequestStatus) error

# tests

execute tests: 
- start a new blockchain: make start
- execute tests: make test
- see resulting blockchain tree: bchain display --entries


# versions

## release 0.0.1 (first version)

- a antblockchain docker service starting with a given number of nodes. The service can't scale in/out in v0.0.1
- each node etablishes GRPC connections with part of the other nodes accordlying to the grid parameters, establishing a ready to work node network communication based on ant behavior.
- each node creates a random RSA key paire at startup, keeps its private one in memory only and sends its public one to all nodes which keep it in memory only also.
- a remote api and an antblockchain CLI based on it, called "bchain" are available
- user can be created with RSA keys paire, the public one is send to the nodes, the private one is kept by the user. Users can be removed.
- the blockchain tree manage a Merkel tree to store blockchain entries and ensure its integrity. The tree is replicated on each node and can be extends adding new branches. Each branch is defined by a label name=value list, one at each tree depth/branch See ./docs/blochchain-tree.pptx
- a new blockchain entry can be added at the end of any branch, found using a list of label name=value), the new entry in sign by the user who execute the command and verified by all the nodes to be accepted only if majority of nodes answer it's ok. 
- Each node verify the user signature, the others nodes signatures, the entry signature, the Merkel tree integrity, the root hash, before validating any entry or branch adding request
- the blockchain tree can be display with debug information
- the service nodes list can be display with status and information

On the v0.0.1 the blockchain works completly with the folowing limitations which have to be resolved in the next versions:
- we can't have several blockchain transactions (add branch, add entry) at the same time on the cluster. One has to be finished to start another one:
- -> need to handle sequencial queue with immediat error when max number of waiting request is reached
- the blockchain service is resilient to nodes crash/restart, but can't scale in/out
- -> need to be able to remove definitelly a node or add a new one
- the blockchain tree is fully replicated on each node
- -> needs to add optional sharding based on branches labels
- the blochchain entry payload is a byte array with no extensible behavior or dedicated verifications 
- -> need dedicated interfaces to make possible to extend it
- the remote api is GRPC only
- -> needs REST API
- a blockhchain entry payload is stored in the Merkel tree and has a maximum size
- -> needs to keep only hashs and id on the blockchain tree and move the payloads on a dedicated database with integrity guaranty based on hash.
- several antblockchain can't be associated to created a shared meta blockhain
- -> needs to allow antblockchain services inter-communication to form a meta blockchain working as a regular blockchain sharded by blockchain service name
- several antblockchain can't be associated to created a regular blockchain
- -> needs to allow antblockchain services inter-communication to form a meta blockchain working as a regular blockchain completly replicated on all services nodes.


## version 0.0.2 target

- add transactionnal request queue to garanty the "add" requests are sequantially executed and garanty their order for the same user


## License

antblockchain is licensed under the Apache License, Version 2.0. See https://github.com/freignat91/blockchain/blob/master/LICENSE
for the full license text.

