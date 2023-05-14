package config_test

import (
	"testing"
	config "woodpecker/internal/configs"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := config.New("../../slack.config.yml")
	expected := createConfigMock()

	require.Equal(t, expected, cfg, "config")
}

func createConfigMock() *config.Config {
	cfg := config.Config{}
	cfg.Slack.AppToken = "xxx"
	cfg.Slack.OAuthToken = "xxx"

	return &cfg
}
