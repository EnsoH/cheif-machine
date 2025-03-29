package exchangeAdapters

import (
	"cw/globals"
	"cw/models"
	"cw/utils"
	"fmt"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

type KucoinAdapter struct {
	Client *ccxt.Kucoin
}

func (k *KucoinAdapter) Withdraw(token, address, network string, amount float64) (ccxt.Transaction, error) {
	networkName, ok := k.getNetworkName(network, "kucoin")
	if !ok {
		return ccxt.Transaction{}, fmt.Errorf("софт не поддерживает вывод в сеть %s с биржи Bybit", network)
	}

	return k.Client.Withdraw(token, amount, address, ccxt.WithdrawOptions(ccxt.WithWithdrawParams(map[string]interface{}{
		"network": networkName,
	})))
}

func (k *KucoinAdapter) GetBalance(symbol string) (ccxt.Balances, error) {
	return k.Client.FetchBalance(map[string]interface{}{
		"type": "spot",
	})
}

func (k *KucoinAdapter) GetPrices(symbol string) (float64, error) {
	k.Client.LoadMarkets()
	ticker, err := k.Client.FetchTicker(fmt.Sprintf("%s/USDT", symbol))
	if err != nil {
		return 0, err
	}
	return *ticker.Last, nil
}

func (k *KucoinAdapter) GetChains(token, withdrawChain string) (*models.ChainList, error) {
	networkName, ok := k.getNetworkName(withdrawChain, "kucoin")
	if !ok {
		return nil, fmt.Errorf("софт не поддерживает вывод в сеть %s с биржи Bybit", withdrawChain)
	}

	result := <-k.Client.LoadMarkets()
	if err, ok := result.(error); ok && err != nil {
		return nil, fmt.Errorf("failed to load markets: %w", err)
	}

	curRaw, exists := k.Client.Currencies[token]
	if !exists {
		return nil, fmt.Errorf("token %s not found", token)
	}

	var curParse *models.KucoinCurrencyList
	if err := utils.ResponseConvert(curRaw, &curParse); err != nil {
		return nil, err
	}
	// log.Printf("currs: %+v", curParse)

	var chainParams models.ChainList
	for _, param := range curParse.Networks {
		if param.Name == networkName {
			if param.Withdraw {
				chainParams.Chain = param.Id
				chainParams.WithdrawFee = param.Fee
				if withdMin, err := utils.ConvertToFloat(param.Limits.Withdraw.Min); err != nil {
					chainParams.WithdrawMin = withdMin
				}
			}
		}
	}

	return &chainParams, nil
}

func (k *KucoinAdapter) getNetworkName(network, cex string) (string, bool) {
	name, exists := globals.ChainNameToSymbolCEX[cex][network]
	return name, exists
}
