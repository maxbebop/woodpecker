package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Slack Slack `yaml:"slack"`
}

type Slack struct {
	OAuthToken string `yaml:"oauth_token"`
	AppToken   string `yaml:"app_token"`
}

func New(filename string) *Config {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		return nil
	}

	if config.Slack.OAuthToken == "" {
		log.Fatal("error: oauth_token is empty")
		return nil
	}

	if config.Slack.AppToken == "" {
		log.Fatal("error: app_token is empty")
		return nil
	}

	return &config
}
