package main

import (
	"fmt"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

//AgentConfig configuration parameters
type CliConfig struct {
	serverAddress string
	serverPort    string
	colorTheme    string
	userName      string
	keyPath       string
}

//update conf instance with default value and environment variables
func (cfg *CliConfig) init(version string) {
	cfg.setDefault()
	cfg.loadConfig()
	//cfg.displayConfig(version)
}

//Set default value of configuration
func (cfg *CliConfig) setDefault() {
	cfg.serverAddress = "127.0.0.1"
	cfg.serverPort = "30103"
	cfg.colorTheme = "dark"
	cfg.keyPath = "blockchain.yaml"
	homedir, err := homedir.Dir()
	if err != nil {
		return
	}
	cfg.keyPath = path.Join(homedir, ".config", "antblockchain", "private.key")
}

//display amp-pilot configuration
func (cfg *CliConfig) displayConfig(version string) {
	fmt.Printf("antblockchain version: %v\n", version)
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("Configuration:")
	fmt.Printf("antblockchain address: %s\n", cfg.serverAddress)
}

// InitConfig reads secret variables in conffile
func (cfg *CliConfig) loadConfig() {

	// Add matching environment variables - will take precedence over config files.
	viper.AutomaticEnv()

	// Add default config file search paths in order of decreasing precedence.
	viper.SetConfigName("blockchain")
	homedir, err := homedir.Dir()
	if err != nil {
		return
	}
	viper.AddConfigPath(path.Join(homedir, ".config/antblockchain"))
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: unable to load /.config/antblockchain/blockchain.yaml\n")
		return
	}
	//fmt.Printf("Viper: %+v\n", viper.AllSettings())

	// Save viper into config. workround  unmarshal bug
	cfg.serverAddress = viper.AllSettings()["serveraddress"].(string)
	cfg.serverPort = viper.AllSettings()["serverport"].(string)
	cfg.userName = viper.AllSettings()["username"].(string)
	cfg.keyPath = viper.AllSettings()["keypath"].(string)
	//cfg.colorTheme = viper.AllSettings()["colortheme"].(string)
	/*
		if err := viper.Unmarshal(cfg); err != nil {
			fmt.Println("Unmarshal antblockchain conffile error: %v\n", err)
			return
		}
	*/
	//fmt.Printf("Conf: %+v\n", cfg)
	//fmt.Printf("Amplifier conffile /.config/amp/amplifier.yaml loaded\n")
}
