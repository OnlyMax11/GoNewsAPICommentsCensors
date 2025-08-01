package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ConJson Конфигурация приложения.
type ConJson struct {
	Period  int      `json:"request_period"`
	LinkArr []string `json:"rss"`
}

func NewConfig() []ConJson {
	configURL := filepath.Join("config.json")
	// loadConfiguration Чтение и раскодированние файла конфигурации.
	bytes, err := os.ReadFile(configURL)
	if err != nil {
		panic(err.Error())
	}
	var conf []ConJson
	_ = json.Unmarshal(bytes, &conf)
	if err != nil {
		panic(err.Error())
	}

	return conf
}