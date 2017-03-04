package main

import ()

// build vars
var (
	Version string
	bCLI    = &bchainCLI{}
	config  = &CliConfig{}
)

func main() {
	config.init(Version)
	cli()
}
