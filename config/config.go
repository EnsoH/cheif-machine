package config

import (
	"cw/globals"
	"cw/models"
	"cw/utils"
	"encoding/json"
	"os"
)

// Singltons
// Configurations for CEX's and wihdraw params
var Cfg *models.AppConfig
var UserCfg *models.UserConfig

// var WithdrawCfg *models.WithdrawConfig
// var BridgeCfg *models.BridgeConfig
// var WallGenConfig *models.WalletGeneratorCfg

// Global array with modules names. Ex: 'relay', 'bybit', 'kucoin' etc...
var SelectModules = []string{}

// var (
// 	configMap = map[string][]string{
// 		"CexWithdrawer": {globals.Configuration, globals.Withdraw},
// 		"Bridger":       {globals.PrivateKeys},
// 		"Cex_Bridger":   {globals.Configuration, globals.Withdraw},
// 	}
// )

func InitConfigs(module string) error {
	// for _, module := range configMap {
	// }
	// продумать механизм, который будет прокидывать массив названий нужных конфигов
	// массив должен заполняться самостоятельно при инициализации или что-то в этом роде, ПРОДУМАТЬ. Это позволит избавиться от ветвления и использовать цикл
	if err := getConfig(globals.AppConfiguration, &Cfg); err != nil {
		return err
	}

	if err := getConfig(globals.UserConfig, &UserCfg); err != nil {
		return err
	}

	// if err := getConfig(globals.Bridge, &BridgeCfg); err != nil {
	// 	return err
	// }

	// if err := getConfig(globals.WallGen, &WallGenConfig); err != nil {
	// 	return err
	// }

	// SelectModules = append(SelectModules, WallGenConfig.WalletType)
	SelectModules = append(SelectModules, UserCfg.WithdrawConfig.CEX)
	SelectModules = append(SelectModules, UserCfg.BridgeConfig.Bridge)

	return nil
}

func getConfig(file string, configModel interface{}) error {
	path, err := utils.GetPath(file)
	if err != nil {
		return err
	}
	return loadConfig(path, &configModel)
}

func loadConfig(path string, parceReseult interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, parceReseult); err != nil {
		return err
	}

	return nil
}
