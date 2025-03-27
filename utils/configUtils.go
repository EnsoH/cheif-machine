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
		globals.AppConfiguration: "config/data/configuration.json",
		globals.UserConfig:       "config/data/user_config.json",
		globals.Addresses:        "config/data/withdraw_addresses.txt",
		// globals.Withdraw:         "config/data/withdraw_configuration.json",
		// globals.Proxy:            "account/proxy.txt",
		// globals.PrivateKeys:      "config/data/private_keys.txt",
		// globals.Bridge:       "config/data/bridge_configuration.json",
		// globals.WallGen:      "config/data/wallet_generator.json",
		// globals.Destinations: "config/data/destination_address.json",
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
