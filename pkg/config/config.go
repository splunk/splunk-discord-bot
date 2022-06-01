package config // github.com/splunk/splunk-discord-bot

import (
	"encoding/json"
	"io/ioutil" //it will be used to help us read our config.json file.
)

type Config struct {
	Token              string `json:"token"`
	HecEndpoint        string `json:"hec_endpoint"`
	HecToken           string `json:"hec_token"`
	InsecureSkipVerify bool   `json:"hec_insecure_skip_verify"`
	HecIndex           string `json:"hec_index"`
}

func ReadConfig() (*Config, error) {
	file, err := ioutil.ReadFile("./config.json")
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
