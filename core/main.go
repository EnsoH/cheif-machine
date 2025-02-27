package main

import (
	"cw/logger"
	"cw/models"
	"cw/modules"
	"cw/process"
	"cw/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		// logger.GlobalLogger.Error(err)
		return
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	logger.GlobalLogger.Infof("Current ENV: %v", env)

	var cexcfg *models.Config
	if err := utils.GetCexConfig("cex_config", &cexcfg); err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	// if err := process.InitSingltons("", map[string]string{}); err != nil {
	// 	logger.GlobalLogger.Error(err)
	// 	return
	// }

	var withdrawCfg *models.WithdrawConfig
	if err := utils.GetCexConfig("withdraw_config", &withdrawCfg); err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	addresses, err := utils.FileReader("/Users/ssq/Desktop/Softs/cex-machine/config/data/withdraw_addresses.dev.txt")
	if err != nil {
		logger.GlobalLogger.Error(err)
		return
	}
	log.Printf("addr len: %d", len(addresses))
	actions, err := process.WithdrawFactory(withdrawCfg, addresses)
	if err != nil {
		logger.GlobalLogger.Error(err)

		return
	}
	for _, act := range actions {
		log.Printf("addr: %s", act.Address)
		log.Printf("cex: %s", act.CEX)
		log.Printf("chain: %s", act.Chain)
		log.Printf("currency: %s", act.Currency)
		log.Printf("value: %s", act.Amount)
		log.Printf("_________")
	}

	// return
	modules, err := modules.ModulesInit(&cexcfg.CEXConfigs)
	if err != nil {
		logger.GlobalLogger.Error(err)
		return
	}

	modules["bybit"].GetPrices("ETH")
	return
}
