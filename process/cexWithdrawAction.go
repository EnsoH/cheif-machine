package process

import (
	"cw/account"
	"cw/config"
	"cw/ethClient"
	"cw/logger"
	"cw/models"
	"cw/modules"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func (a *ActionCore) withdrawProcess(accs []*account.Account, mod *modules.Modules, selectModule string) error {
	actions, err := withdrawFactory(accs, *mod.Exchange)
	if err != nil {
		return err
	}

	loggingActions(actions)

	if selectModule == "Cex_Bridger" {
		a.cacheWithdrawActions(actions)
	}

	return a.executeWithdrawActions(actions, mod)
}

func (a *ActionCore) executeWithdrawActions(actions []models.WithdrawAction, mod *modules.Modules) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.Cfg.Threads)
	errCh := make(chan error, len(actions))

	for _, act := range actions {
		wg.Add(1)
		go func(act models.WithdrawAction) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := a.processSingleWithdrawAction(act, mod); err != nil {
				errCh <- err
			}
		}(act)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		logger.GlobalLogger.Warnf("Ошибка при выводе: %v", err)
	}

	return nil
}

func (a *ActionCore) processSingleWithdrawAction(act models.WithdrawAction, mod *modules.Modules) error {
	logger.GlobalLogger.Infof("[%s] Sleep before withdraw %v", act.Address, act.TimeRange)
	time.Sleep(time.Second * time.Duration(act.TimeRange))

	if err := mod.Exchange.Withdraw(act.CEX, act.Currency, act.Address, act.Chain, act.Amount); err != nil {
		return handleWithdrawError(err)
	}

	if err := ethClient.GlobalETHClient[act.Chain].WaitTokenDeposit(act.Chain, act.Currency, common.HexToAddress(act.Address)); err != nil {
		logger.GlobalLogger.Errorf("[%s] Ошибка ожидания депозита", act.Address)
		return nil
	}

	logger.GlobalLogger.Infof("[%s] Токен %s поступил на счет.", act.Address, act.Currency)
	return nil
}

func withdrawFactory(accounts []*account.Account, exch modules.Exchanges) ([]models.WithdrawAction, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("нет списка адресов")
	}

	tokenData := make(map[string]models.TokenInfo)

	for _, token := range config.UserCfg.WithdrawConfig.Currency {
		if info, ok := gatherTokenData(exch, token); ok {
			tokenData[token] = info
		}
	}

	if len(tokenData) == 0 {
		return nil, fmt.Errorf("нет доступных токенов для вывода")
	}

	actions := make([]models.WithdrawAction, 0, len(accounts))
	rand.Seed(time.Now().UnixNano())

	for _, acc := range accounts {
		withdrawUSD := getRandomAmount(config.UserCfg.WithdrawConfig.AmountRange)
		chosenToken, amountInToken, found := selectToken(tokenData, withdrawUSD)
		if !found {
			logger.GlobalLogger.Errorf("Нет токена с достаточным балансом для вывода %v$ для аккаунта %s", withdrawUSD, acc.Address.String())
			continue
		}

		info := tokenData[chosenToken]
		info.AvailableCapacity -= withdrawUSD
		info.Balance -= amountInToken
		tokenData[chosenToken] = info

		if len(info.Chains) == 0 {
			logger.GlobalLogger.Errorf("Нет доступных сетей для токена %s у аккаунта %s", chosenToken, acc.Address.String())
			continue
		}
		chosenChain := info.Chains[rand.Intn(len(info.Chains))]
		timeVal := getRandomAmount(config.UserCfg.WithdrawConfig.TimeRange)
		destinationChain := config.UserCfg.WithdrawConfig.DestinationChain
		destinationCurrency := config.UserCfg.WithdrawConfig.DestinationCurrency

		action := createWithdrawAction(config.UserCfg.WithdrawConfig.CEX, acc.Address.String(), chosenToken, chosenChain, destinationChain, destinationCurrency, amountInToken, timeVal)
		actions = append(actions, action)
	}

	if len(actions) == 0 {
		return nil, fmt.Errorf("ошибка при обработке выводов")
	}

	return actions, nil
}

func gatherTokenData(exch modules.Exchanges, token string) (models.TokenInfo, bool) {
	var info models.TokenInfo

	balance, err := exch.GetBalances(config.UserCfg.WithdrawConfig.CEX, token)
	if err != nil {
		logger.GlobalLogger.Errorf("Ошибка получения баланса для %s: %v", token, err)
		return info, false
	}

	price, err := exch.GetPrices(config.UserCfg.WithdrawConfig.CEX, token)
	if err != nil {
		logger.GlobalLogger.Errorf("Ошибка получения цены для %s: %v", token, err)
		return info, false
	}

	var maxAvailableUSD float64
	var usableChains []string

	for _, chain := range config.UserCfg.WithdrawConfig.Chain {
		withdrawParams, err := exch.GetChains(config.UserCfg.WithdrawConfig.CEX, token, chain)
		if err != nil {
			logger.GlobalLogger.Warnf("Ошибка получения параметров вывода для %s в сети %s: %v", token, chain, err)
			continue
		}
		if withdrawParams.Chain == "" {

			continue
		}
		if balance > withdrawParams.WithdrawFee {
			availableUSD := (balance - withdrawParams.WithdrawFee) * price
			if availableUSD > 0 {
				usableChains = append(usableChains, chain)
				if availableUSD > maxAvailableUSD {
					maxAvailableUSD = availableUSD
				}
			}
		}
	}

	if len(usableChains) == 0 {
		return info, false
	}

	info = models.TokenInfo{
		AvailableCapacity: maxAvailableUSD,
		Price:             price,
		Chains:            usableChains,
		Balance:           balance,
	}

	return info, true
}

func selectToken(tokenData map[string]models.TokenInfo, withdrawUSD float64) (string, float64, bool) {
	for token, info := range tokenData {
		amountInToken := withdrawUSD / info.Price
		if info.AvailableCapacity >= withdrawUSD && info.Balance >= amountInToken {
			return token, amountInToken, true
		}
	}
	return "", 0, false
}

func createWithdrawAction(cex, address, token, chain, destinationChain, destinationToken string, amount, time float64) models.WithdrawAction {
	return models.WithdrawAction{
		Address:             address,
		CEX:                 cex,
		Chain:               chain,
		Currency:            token,
		Amount:              amount,
		TimeRange:           time,
		DestinationChain:    destinationChain,
		DestinationCurrency: destinationToken,
	}
}

func (a *ActionCore) cacheWithdrawActions(actions []models.WithdrawAction) {
	for _, act := range actions {
		a.withdrawAction[act.Address] = act
	}
}
