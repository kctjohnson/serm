package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Service struct {
	Name       string   `json:"name"`
	Bin        string   `json:"bin"`
	Args       []string `json:"args"`
	Retry      bool     `json:"retry"`
	RetryCount int      `json:"retryCount"`
	OnStartup  bool     `json:"onStartup"`
}

type Config struct {
	LogPath  string    `json:"logPath"`
	Services []Service `json:"services"`
}

func NewConfig() Config {
	c := Config{}
	c.init()
	return c
}

func (c *Config) init() {
	// Load the config into viper, otherwise create it
	c.loadConfig(false)

	// Unmarshal the json file's contents into the config struct
	viper.Unmarshal(&c)

	// Check to see that all of the services have unique names
	foundServices := make(map[string]bool)
	for _, serv := range c.Services {
		if ok := foundServices[serv.Name]; ok {
			panic("Duplicate service name entry found! Please only provide unique service names!")
		}
		foundServices[serv.Name] = true
	}
}

func (c Config) loadConfig(retry bool) {
	configPath := "/.config/serm/"
	// Set up the config
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME" + configPath)
	err := viper.ReadInConfig()

	// If it failed to read in the config, create a new blank one
	if err != nil && retry == false {
		log.Printf("Creating blank config file at $HOME%s\n", configPath)
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		err = os.MkdirAll(homeDir+configPath, os.ModePerm)
		if err != nil {
			panic(err)
		}

		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatalf("ERROR: %#v\n", err)
		}

		// Attempt to load the new default config
		c.loadConfig(true)
	}
}
