package utils

import (
	"cw/config"
	"os"
)

func GetCexConfig(file string, cfg interface{}) error {
	if err := config.LoadConfig(GetPath(file), &cfg); err != nil {
		return err
	}
	return nil
}

func GetPath(path string) string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	basePath := map[string]string{
		"cex_config":      "config/data/cex_configuration.json",
		"withdraw_config": "config/data/withdraw_configuration.json",
		"proxy":           "account/proxy.txt",
	}

	if env == "development" {
		devPaths := make(map[string]string)
		for key, value := range basePath {
			devPaths[key] = addDevSuffix(value)
		}
		return devPaths[path]
	}
	return basePath[path]
}

func addDevSuffix(path string) string {
	extIndex := len(path) - len(".json")
	if extIndex > 0 && path[extIndex:] == ".json" {
		return path[:extIndex] + ".dev.json"
	}
	return path
}
