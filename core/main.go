package main

import (
	"cw/config"
	"cw/globals"
	"cw/logger"
	"cw/modules"
	"cw/process"
	"cw/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	setENV()
	if err := config.InitConfigs(); err != nil {
		logger.GlobalLogger.Error(err)
		return
	}
	log.Printf("cfg: %v", config.Cfg)
	log.Printf("withd: %v", config.WithdrawCfg)
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

	actions, err := process.WithdrawFactory(addresses)
	if err != nil {
		logger.GlobalLogger.Error(err)

		return
	}

	modules, err := modules.ModulesInit()
	if err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	modules["bybit"].GetPrices("ETH")
	return
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
