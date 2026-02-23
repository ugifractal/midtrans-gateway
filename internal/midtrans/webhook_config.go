package midtrans

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type URLItem struct {
	Code string `yaml:"code"`
	URL  string `yaml:"url"`
}

type Config struct {
	URLs []URLItem `yaml:"urls"`
}

var WebhookConfig Config

func LoadConfig() {
	location := os.Getenv("CONFIG_FILE")
	data, err := os.ReadFile(location)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	err = yaml.Unmarshal(data, &WebhookConfig)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML: %v", err)
	}
}
