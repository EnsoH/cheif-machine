package main

import (
	"cw/account"
	"cw/config"
	"cw/ethClient"
	"cw/globals"
	"cw/logger"
	"cw/modules"
	"cw/process"
	"cw/utils"
)

func main() {
	// utils.SetENV()
	globals.SetInit()
	selectModule := utils.UserChoice()
	if selectModule == "" || selectModule == "Exit" {
		logger.GlobalLogger.Warnf("user choice 'exit'.")
		return
	}

	if err := config.InitConfigs(selectModule); err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	if err := ethClient.EthClientFactory(config.Cfg.Rpc); err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	accs, err := account.AccsFactory(selectModule)
	if err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	modules, err := modules.ModulesInit(selectModule, config.SelectModules...)
	if err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	actionCore, err := process.NewActionCore()
	if err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	if err := actionCore.ActionsProcess(accs, modules, selectModule); err != nil {
		logger.GlobalLogger.Error(err)
		return
	}
}
