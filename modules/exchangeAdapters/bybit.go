package exchangeAdapters

import (
	"cw/models"
	"cw/utils"
	"fmt"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

type BybitAdapter struct {
	Client *ccxt.Bybit
}

func (ba *BybitAdapter) Withdraw(token, address, network string, amount float64) (ccxt.Transaction, error) {
	return ba.Client.Withdraw(token, amount, address, ccxt.WithdrawOptions(ccxt.WithWithdrawParams(map[string]interface{}{
		"forceChain": 1,
		"network":    network,
	})))
}

func (ba *BybitAdapter) GetBalance(symbol string) (ccxt.Balances, error) {
	return ba.Client.FetchBalance(map[string]interface{}{
		"type": "funding",
	})
}

func (ba *BybitAdapter) GetPrices(symbol string) (float64, error) {
	if result := <-ba.Client.LoadMarkets(); result != nil {
		if err, ok := result.(error); ok && err != nil {
			return 0, fmt.Errorf("failed to load markets: %w", err)
		}
	}

	ticker, err := ba.Client.FetchTicker(fmt.Sprintf("%s/USDT", symbol))
	if err != nil {
		return 0, err
	}
	if ticker.Last == nil {
		return 0, fmt.Errorf("ticker last price is nil")
	}
	return *ticker.Last, nil
}

func (b *BybitAdapter) GetChains(token, withdrawChain string) (*models.ChainList, error) {
	result := <-b.Client.LoadMarkets()
	if err, ok := result.(error); ok && err != nil {
		return nil, fmt.Errorf("failed to load markets: %w", err)
	}

	curRaw, exists := b.Client.Currencies[token]
	if !exists {
		return nil, fmt.Errorf("token %s not found", token)
	}

	var curParse *models.BybitCurrencyList
	if err := utils.ResponseConvert(curRaw, &curParse); err != nil {
		return nil, err
	}

	var chainParams models.ChainList
	for _, param := range curParse.Info.Chains {
		// log.Printf("chain:%+v", chain.ChainType)
		if param.ChainType == withdrawChain {
			chainParams.Chain = param.Chain
			if withdrawMin, err := utils.ConvertToFloat(param.WithdrawMin); err == nil {
				chainParams.WithdrawFee = withdrawMin
			}
		}
	}

	return &chainParams, nil
}
