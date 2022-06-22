package config // github.com/splunk/splunk-discord-bot

import (
	"encoding/json"
	"io/ioutil" //it will be used to help us read our config.json file.
	"os"
)

type WebhookConfig struct {
	ID      string `json:"id"`
	Channel string `json:"channel"`
}

type Config struct {
	Token              string           `json:"token"`
	HecEndpoint        string           `json:"hec_endpoint"`
	HecToken           string           `json:"hec_token"`
	InsecureSkipVerify bool             `json:"hec_insecure_skip_verify"`
	HecIndex           string           `json:"hec_index"`
	WebHooks           []*WebhookConfig `json:"webhooks"`
	WebhookListenAddr  string           `json:"listen_addr"`
}

func ReadConfig() (*Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "./config.json"
	}
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var cfg *Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil

}
