package process

import (
	"cw/account"
	"cw/config"
	"cw/logger"
	"cw/modules"
	"time"
)

func (a *ActionCore) cexBridgeProcess(accounts []*account.Account, mods *modules.Modules, selectModule string) error {
	if err := a.withdrawProcess(accounts, mods, selectModule); err != nil {
		return err
	}

	logger.GlobalLogger.Infof("Спим перед Bridge %v минуты", config.Cfg.SleepAfterWithdraw)
	time.Sleep(time.Minute * time.Duration(config.Cfg.SleepAfterWithdraw))

	return a.bridgeProcess(accounts, mods, selectModule)
}
