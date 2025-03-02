package main

import (
	"cw/config"
	"cw/globals"
	"cw/logger"
	"cw/modules"
	"cw/process"
	"cw/utils"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("START")
	setENV()

	if err := config.InitConfigs(); err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	addrPath, err := utils.GetPath(globals.Addresses)
	if err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	addresses, err := utils.FileReader(addrPath)
	if err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	modules, err := modules.ModulesInit()
	if err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	if err := process.ActionsProcess(addresses, modules, "bybit"); err != nil {
		logger.GlobalLogger.Error(err)
		return
	}
}

func setENV() string {
	if err := godotenv.Load(); err != nil {
		logger.GlobalLogger.Warnf("Not found ENV. Set default params(production)")
		os.Setenv("ENV", "production")
	}
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
		os.Setenv("ENV", "development")
	}
	return env
}
