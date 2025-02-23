package config

import (
	"encoding/json"
	"os"
)

func LoadConfig(path string, parceReseult interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	// var config Config
	if err := json.Unmarshal(data, &parceReseult); err != nil {
		return nil
	}

	return nil
}
