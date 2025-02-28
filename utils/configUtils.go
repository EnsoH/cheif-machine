package utils

import (
	"cw/globals"
	"fmt"
	"os"
	"strings"
)

func GetPath(path string) (string, error) {
	env := os.Getenv("ENV")
	basePath := map[string]string{
		globals.Configuration: "config/data/configuration.json",
		globals.Withdraw:      "config/data/withdraw_configuration.json",
		globals.Proxy:         "account/proxy.txt",
		globals.Addresses:     "config/data/withdraw_addresses.txt",
	}

	// Проверяем, существует ли путь
	filePath, exists := basePath[path]
	if !exists {
		return "", fmt.Errorf("unknown config key: %s", path)
	}

	// Если development, добавляем .dev
	if env == "development" {
		filePath = addDevSuffix(filePath)
	}

	return filePath, nil
}

func addDevSuffix(path string) string {
	if strings.HasSuffix(path, ".json") {
		return strings.TrimSuffix(path, ".json") + ".dev.json"
	} else if strings.HasSuffix(path, ".txt") {
		return strings.TrimSuffix(path, ".txt") + ".dev.txt"
	}
	return path
}
