package process

import (
	"cw/account"
	"cw/config"
	"cw/ethClient"
	"cw/logger"
	"cw/models"
	"cw/modules"
	"cw/utils"
	"fmt"
	"sync"
	"time"
)

func (a *ActionCore) bridgeProcess(accounts []*account.Account, mods *modules.Modules, moduleName string) error {
	if len(accounts) == 0 || mods == nil || moduleName == "" {
		return fmt.Errorf("переданы пустые параметры в bridgeProcess")
	}

	actions, err := a.bridgeFactory(accounts, moduleName)
	if err != nil {
		return err
	}

	if len(actions) > 0 {
		loggingBridgeAction(actions)
	} else {
		logger.GlobalLogger.Infof("Нет действий для логирования.")
	}

	return a.executeBridgeActions(accounts, actions, mods)
}

func (a *ActionCore) executeBridgeActions(accounts []*account.Account, actions []*models.BridgeAction, mods *modules.Modules) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.Cfg.Threads)
	accMap := buildAccountMap(accounts)

	for _, act := range actions {
		wg.Add(1)

		go func(action *models.BridgeAction) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if action == nil {
				return
			}

			acc, ok := accMap[act.Address]
			if !ok {
				logger.GlobalLogger.Errorf("[%s] Ошибка при получении аккаунта для действия", action.Address)
				return
			}

			if !validateAction(acc, action) {
				logger.GlobalLogger.Infof("[%s] Недостаточно средств для действия", action.Address)
				return
			}

			logger.GlobalLogger.Infof("[%s] Спим перед бриджем %v", action.Address, action.Time)
			time.Sleep(time.Second * time.Duration(action.Time))

			if err := a.performBridgeAction(action, mods, acc); err != nil {
				logger.GlobalLogger.Errorf("[%s] Ошибка при бридже: %v", action.Address, err)
			}
		}(act)
	}

	wg.Wait()
	return nil
}

func (a *ActionCore) performBridgeAction(action *models.BridgeAction, mods *modules.Modules, acc *account.Account) error {
	bridgeModule, ok := mods.Bridges.BridgesMap[config.UserCfg.BridgeConfig.Bridge]
	if !ok {
		return fmt.Errorf("модуль бриджа %s не найден", config.UserCfg.BridgeConfig.Bridge)
	}

	return bridgeModule.Bridge(action.FromChain, action.ToChain, action.FromCurrency, action.ToCurrency, action.Amount, acc)
}

func (a *ActionCore) bridgeFactory(accounts []*account.Account, moduleName string) ([]*models.BridgeAction, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("нет списка адресов")
	}

	actions := make([]*models.BridgeAction, 0, len(accounts))
	for _, acc := range accounts {
		withdrawAction, err := a.getWithdrawActionForAccount(acc, moduleName)
		if err != nil {
			logger.GlobalLogger.Errorf("Ошибка при получении withdrawAction для аккаунта %s: %v", acc.Address.Hex(), err)
			continue
		}

		action, err := initBridgeAction(acc, withdrawAction, moduleName)
		if err != nil {
			logger.GlobalLogger.Errorf("Ошибка инициализации действия для аккаунта %s: %v", acc.Address.Hex(), err)
			continue
		}

		if action != nil {
			actions = append(actions, action)
		}
	}

	if len(actions) == 0 {
		return nil, fmt.Errorf("нет действующих действий для бриджа")
	}

	return actions, nil
}

func (a *ActionCore) getWithdrawActionForAccount(acc *account.Account, moduleName string) (*models.WithdrawAction, error) {
	if moduleName == "Cex_Bridger" {
		withdrawAction, ok := a.withdrawAction[acc.Address.Hex()]
		if !ok {
			return nil, fmt.Errorf("нет действия для аккаунта %s", acc.Address.Hex())
		}
		return &withdrawAction, nil
	}
	return nil, nil
}

func initBridgeAction(acc *account.Account, withdrawAction *models.WithdrawAction, moduleName string) (*models.BridgeAction, error) {
	action := &models.BridgeAction{
		Address: acc.Address,
		Bridge:  config.UserCfg.BridgeConfig.Bridge,
	}

	switch moduleName {
	case "Cex_Bridger":
		if withdrawAction == nil {
			return nil, fmt.Errorf("не передан withdrawAction для модуля %s", moduleName)
		}
		return initCexBridgeAction(action, withdrawAction)
	case "Bridger":
		return initDefaultBridgeAction(action, acc)
	default:
		return nil, fmt.Errorf("неподдерживаемый модуль: %s", moduleName)
	}
}

func initCexBridgeAction(action *models.BridgeAction, withdrawAction *models.WithdrawAction) (*models.BridgeAction, error) {
	initBridgeActionBase(action, withdrawAction.Chain, withdrawAction.DestinationChain, withdrawAction.Currency, withdrawAction.DestinationCurrency, withdrawAction.TimeRange)

	amount, err := getAmount(withdrawAction.Chain, withdrawAction.Currency, withdrawAction.Amount)
	if err != nil {
		return nil, err
	}

	action.Amount = amount
	return action, nil
}

func initDefaultBridgeAction(action *models.BridgeAction, acc *account.Account) (*models.BridgeAction, error) {
	initBridgeActionBase(action, config.UserCfg.BridgeConfig.FromChain, config.UserCfg.BridgeConfig.ToChain, config.UserCfg.BridgeConfig.FromCurrency, config.UserCfg.BridgeConfig.ToCurrency, getRandomAmount(config.UserCfg.BridgeConfig.TimeRange))
	action.Address = acc.Address

	tokenContract, ethClientInstance, decimals, tokenPrice, err := getTokenInfo(config.UserCfg.BridgeConfig.FromChain, config.UserCfg.BridgeConfig.FromCurrency)
	if err != nil {
		return nil, err
	}

	balanceWei, err := ethClientInstance.BalanceCheck(acc.Address, tokenContract)
	if err != nil {
		return nil, err
	}

	balanceTokens := utils.ConvertFromWei(balanceWei, int(decimals))
	balanceUSD := balanceTokens * tokenPrice

	minAmount, maxAmount := config.UserCfg.BridgeConfig.AmountRange[0], config.UserCfg.BridgeConfig.AmountRange[1]

	if balanceUSD < minAmount {
		logger.GlobalLogger.Errorf("[%s] Недостаточно баланса. Требуемый: %v$, текущий: %v$", acc.Address.Hex(), minAmount, balanceUSD)
		return nil, nil
	}

	amountUSD := getRandomAmount(config.UserCfg.BridgeConfig.AmountRange)
	if amountUSD > balanceUSD {
		amountUSD = balanceUSD
	} else if amountUSD > maxAmount {
		amountUSD = maxAmount
	}

	action.Amount, err = getAmount(action.FromChain, action.FromCurrency, amountUSD)
	if err != nil {
		return nil, err
	}
	return action, nil
}

func initBridgeActionBase(action *models.BridgeAction, fromChain, toChain, fromCurrency, toCurrency string, time float64) {
	action.FromChain = fromChain
	action.ToChain = toChain
	action.FromCurrency = fromCurrency
	action.ToCurrency = toCurrency
	action.Time = time
}

func validateAction(account *account.Account, resultAction *models.BridgeAction) bool {
	tokenContract, err := getTokenContract(resultAction.FromChain, resultAction.FromCurrency)
	if err != nil {
		logger.GlobalLogger.Error("[%s] Ошибка получения контракта токена: %v", account.Address, err)
		return false
	}

	weiAmount, err := ethClient.GlobalETHClient[resultAction.FromChain].BalanceCheck(account.Address, tokenContract)
	if err != nil {
		logger.GlobalLogger.Error("[%s] Ошибка проверки баланса:", account.Address, err)
		return false
	}

	if weiAmount.Cmp(resultAction.Amount) < 0 {
		logger.GlobalLogger.Errorf("[%s] Недостаточно средств: баланс - %v, требуется - %v", account.Address.Hex(), weiAmount, resultAction.Amount)
		return false
	}

	return true
}
