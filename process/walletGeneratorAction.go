package process

import (
	"cw/account"
	"cw/config"
	"cw/modules"
)

func (a *ActionCore) walletGeneratorAction(accs []*account.Account, mod *modules.Modules, selectModule string) error {
	return mod.WalletGenerators.GenerateWallets(config.UserCfg.WalletGeneratorCfg.WalletType, int(config.UserCfg.WalletGeneratorCfg.WalletCount), config.UserCfg.WalletGeneratorCfg.CsvHeaders)
}
