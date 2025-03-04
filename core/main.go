package main

import (
	"cw/config"
	"cw/globals"
	"cw/logger"
	"cw/modules"
	"cw/process"
	"cw/utils"
	"os"

	"github.com/joho/godotenv"
)

func main() {
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

	exchange, err := modules.ModulesInit(config.WithdrawCfg.CEX) // for example we init 2 cex
	if err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	if err := process.ActionsProcess(addresses, *exchange, config.WithdrawCfg.CEX); err != nil {
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
