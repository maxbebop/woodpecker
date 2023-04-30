package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Slack slack `yaml:"slack"`
}

type slack struct {
	OAuthToken string `yaml:"oauth_token"`
	AppToken   string `yaml:"app_token"`
	AppUserId  string `yaml:"app_user_id"`
}

func New(filename string) *Config {

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

	if config.Slack.AppUserId == "" {
		log.Fatal("error: app_user_id is empty")
		return nil
	}

	return &config
}
