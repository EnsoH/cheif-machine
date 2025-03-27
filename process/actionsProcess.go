package process

import (
	"cw/account"
	"cw/models"
	"cw/modules"
	"fmt"
)

type ActionCore struct {
	FunctionsMap   map[string]func(acc []*account.Account, mod *modules.Modules, name string) error
	withdrawAction map[string]models.WithdrawAction
}

func NewActionCore() (*ActionCore, error) {
	ac := &ActionCore{
		withdrawAction: make(map[string]models.WithdrawAction),
	}

	ac.FunctionsMap = map[string]func(acc []*account.Account, mod *modules.Modules, name string) error{
		"CexWithdrawer":   ac.withdrawProcess,
		"Bridger":         ac.bridgeProcess,
		"Cex_Bridger":     ac.cexBridgeProcess,
		"WalletGenerator": ac.walletGeneratorAction,
		"Сollector":       ac.CollectorAction,
	}
	return ac, nil
}

func (a *ActionCore) ActionsProcess(accounts []*account.Account, mods *modules.Modules, selectModule string) error {
	if mods == nil || selectModule == "" {
		return fmt.Errorf("ошибка запуска процессинга, один из параметров нулевой")
	}

	function, ok := a.FunctionsMap[selectModule]
	if !ok {
		return fmt.Errorf("ошибка при выборе генератора действий")
	}

	return function(accounts, mods, selectModule)
}
