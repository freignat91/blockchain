# AGRID

AntBlockChain v0.0.0

# Purpose

This project is an experimental project aiming to build a full secured block chain service implementing the following ideas:

- lower the needed node interconnections using as Agrid project https://github.com/freignat91/agrid ants behavior to handle nodes communication, preversing the block chain security and integrity.

- be able to do not have the full block-chain stored on each node, but have blockchain branches creation/append on sub node network keeping their security and integrity. These branches stays related to block chain root, but stored on a random part of the network. a node should this wat handle several branches, but not all (all should stay possible)

- let open the object stored in the block with size limit and have capability to extend the object logical behavior. For instance object is a financial transaction and one of the logical behavior is to do not have the same user spending two times the same amount of money. So have a clear separation between the block-chain way of working and the object stored own logical constraints.

- technically developed in Go, using Docker services and able to dynamically interconnected several AntBlocChain services to extend the network with free services interconnection topologies.


At v0.0.0 stage, no sure that the propositions above are possible. This project will move to v1.0.0 if it prove to work as expected, stop if not.


# Install

It needsfirst Docker 1.13 installed, then

- clone the git project: https://github.com/freignat91/blockchain
- execute `make install` to build the antblockchain CLI executable
- execute `make build` to create a image freignat91/agrid:latest
- execute the command `make start` to initialize swarm, create an overlay network and start the antblockchain service on this network

antblockchain can't be used as a single container, it needs to be started as a service on a swarm machine (manager or worker).


antblockchain take more time to start than the "ready" docker status. Every CLI command executed before will be rejected with a message "Node not yet ready"


# Configuration using System Variables:


- GRPCPORT:               grpc server port used by nodes (default 30103)
- NB_LINE_CONNECT:        number of "line" type connection in grid: default 0 means computed automatically
- NB_CROSS_CONNECT:       number of "cross" type connection in grid: default 0 means computed automatically
- DATA_PATH:              path in container where file data is stored: default: /data (should be mapped on host using mount docker argument (--mount type=bind,source=/[hostpath],target=/data)


# Resilience

For resilience reason, it's mandatory to have a separated disk file system for each node (each node on its own VM), but for test reason it's possible to have several nodes on the same file system or have architecture with nodes spreaded on a less number of several VMs.


## Grid simulation

To simulate nodes connections using different parameters as, node number, line connections, cross connections,and see the grid topology, use the cli command:

`bchain grid simul [nodes] <--line xx > <--cross yy>`
- [nodes] the number of nodes
- <--line xx> optionally: xx the number of line connections 
- <--cross yy> optionally: yy the number of cross connections 

this command as not effect on the real cluster grid connections, see:
- ./docs/Agrid-grid-building.pptx
- ./docs/Agrid-Ant-net.pptx


# tests

execute: make test


# version 0.0.0 target

- Have antblockchain docker service starting with a given number of nodes. 
- Each node etablish GRPC connections with part of the other nodes accordlying to the grid parameters, establishing a ready to work node network communication based on ant behavior.
- Each node create a random RSA key paire, keep its private one in memory only and send its public one to all nodes which keep them in memory only also.
- have a antblockchain CLI called bchain with the following commands available:
    - node ls: to display node status


## License

antblockchain is licensed under the Apache License, Version 2.0. See https://github.com/freignat91/blockchain/blob/master/LICENSE
for the full license text.
