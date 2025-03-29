package exchangeAdapters

import (
	"cw/globals"
	"cw/models"
	"cw/utils"
	"fmt"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

type MexcAdapter struct {
	Client *ccxt.Mexc
}

func (m *MexcAdapter) Withdraw(token, address, network string, amount float64) (ccxt.Transaction, error) {
	networkName, ok := m.getNetworkName(network, "mexc")
	if !ok {
		return ccxt.Transaction{}, fmt.Errorf("софт не поддерживает вывод в сеть %s с биржи Binance", network)
	}

	return m.Client.Withdraw(token, amount, address, ccxt.WithWithdrawParams(map[string]interface{}{
		"network": networkName,
	}))
}

func (m *MexcAdapter) GetBalance(symbol string) (ccxt.Balances, error) {
	return m.Client.FetchBalance()
}

func (m *MexcAdapter) GetPrices(symbol string) (float64, error) {
	m.Client.LoadMarkets()
	ticker, err := m.Client.FetchTicker(fmt.Sprintf("%s/USDT", symbol))
	if err != nil {
		return 0, err
	}
	return *ticker.Last, nil
}

func (m *MexcAdapter) GetChains(token, withdrawChain string) (*models.ChainList, error) {
	networkName, ok := m.getNetworkName(withdrawChain, "mexc")
	if !ok {
		return nil, fmt.Errorf("софт не поддерживает вывод в сеть %s с биржи Mexc", withdrawChain)
	}

	result := <-m.Client.LoadMarkets()
	if err, ok := result.(error); ok && err != nil {
		return nil, fmt.Errorf("failed to load markets: %w", err)
	}

	curRaw, exists := m.Client.Currencies[token]
	if !exists {
		return nil, fmt.Errorf("token %s not found", token)
	}

	var curParse *models.MexcCurrencyList
	if err := utils.ResponseConvert(curRaw, &curParse); err != nil {
		return nil, err
	}

	// log.Printf("currs: %+v", curParse)
	var chainParams models.ChainList
	for _, param := range curParse.Networks {
		if param.Info.NetWork == networkName {
			if param.Active {
				chainParams.Chain = param.Info.NetWork
				if withdrawFee, err := utils.ConvertToFloat(param.Info.WithdrawFee); err == nil {
					chainParams.WithdrawFee = withdrawFee
				}
				if withdrawMin, err := utils.ConvertToFloat(param.Limits.Withdraw.Min); err == nil {
					chainParams.WithdrawMin = withdrawMin
				}
			}
		}
	}

	return &chainParams, nil
}

func (m *MexcAdapter) getNetworkName(network, cex string) (string, bool) {
	name, exists := globals.ChainNameToSymbolCEX[cex][network]
	return name, exists
}
