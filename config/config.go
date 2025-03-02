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
var Cfg models.Config
var WithdrawCfg models.WithdrawConfig

func InitConfigs() error {
	// продумать механизм, который будет прокидывать массив названий нужных конфигов
	// массив должен заполняться самостоятельно при инициализации или что-то в этом роде, ПРОДУМАТЬ. Это позволит избавиться от ветвления и использовать цикл
	if err := getConfig(globals.Configuration, &Cfg); err != nil {
		return err
	}

	if err := getConfig(globals.Withdraw, &WithdrawCfg); err != nil {
		return err
	}

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
