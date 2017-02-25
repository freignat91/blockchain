package main

import ()

// build vars
var (
	Version string
	Build   string
	bCLI    = &bchainCLI{}
	config  = &CliConfig{}
)

func main() {
	config.init(Version, Build)
	cli()
}
