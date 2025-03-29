package exchangeAdapters

import (
	"cw/globals"
	"cw/models"
	"cw/utils"
	"fmt"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

type OkxAdapter struct {
	Client *ccxt.Okx
}

func (o *OkxAdapter) Withdraw(token, address, network string, amount float64) (ccxt.Transaction, error) {
	chainsParams, err := o.GetChains(token, network)
	if err != nil {
		return ccxt.Transaction{}, err
	}

	return o.Client.Withdraw(token, amount, address, ccxt.WithWithdrawParams(map[string]interface{}{
		"chainName": chainsParams.Network,
		"fee":       chainsParams.WithdrawFee,
		"pwd":       "-",
		"amt":       amount,
		"network":   chainsParams.Chain,
	}))
}

func (o *OkxAdapter) GetBalance(symbol string) (ccxt.Balances, error) {
	return o.Client.FetchBalance(map[string]interface{}{
		"type": "funding",
	})
}

func (o *OkxAdapter) GetPrices(symbol string) (float64, error) {
	o.Client.LoadMarkets()
	ticker, err := o.Client.FetchTicker(fmt.Sprintf("%s/USDT", symbol))
	if err != nil {
		return 0, err
	}

	return *ticker.Last, nil
}

func (o *OkxAdapter) GetChains(token, withdrawChain string) (*models.ChainList, error) {
	okxNetworkName, ok := o.getNetworkName(withdrawChain, "okx")
	if !ok {
		return nil, fmt.Errorf("софт не поддерживает вывод в сеть %s с биржи Okx", withdrawChain)
	}

	result := <-o.Client.LoadMarkets()
	if err, ok := result.(error); ok && err != nil {
		return nil, fmt.Errorf("failed to load market: %w", err)
	}

	curRaw, exist := o.Client.Currencies[token]
	if !exist {
		return nil, fmt.Errorf("token %s not found", token)
	}

	var curParse *models.OkxCurrenceList
	if err := utils.ResponseConvert(curRaw, &curParse); err != nil {
		return nil, err
	}
	// log.Printf("currs: %+v", curParse)

	var chainParams models.ChainList
	for networkName, networkInfo := range curParse.Networks {
		if networkName == okxNetworkName && networkInfo.Withdraw {
			if networkInfo.Withdraw {
				chainParams.Chain = networkName
				chainParams.WithdrawFee = networkInfo.Fee
				chainParams.Network = networkInfo.Id
				chainParams.WithdrawMin = networkInfo.Limits.Withdraw.Min
			}
		}
	}

	return &chainParams, nil
}

func (o *OkxAdapter) getNetworkName(network, cex string) (string, bool) {
	name, exists := globals.ChainNameToSymbolCEX[cex][network]
	return name, exists
}
