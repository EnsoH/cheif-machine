package exchangeAdapters

import (
	"cw/globals"
	"cw/models"
	"cw/utils"
	"fmt"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

type BinanceAdapter struct {
	Client *ccxt.Binance
}

func (b *BinanceAdapter) Withdraw(token, address, network string, amount float64) (ccxt.Transaction, error) {
	networkName, ok := b.getNetworkName(network, "binance")
	if !ok {
		return ccxt.Transaction{}, fmt.Errorf("софт не поддерживает вывод в сеть %s с биржи Binance", network)
	}

	return b.Client.Withdraw(token, amount, address, ccxt.WithdrawOptions(ccxt.WithWithdrawParams(map[string]interface{}{
		"network": networkName,
	})))
}

func (b *BinanceAdapter) GetBalance(symbol string) (ccxt.Balances, error) {
	return b.Client.FetchBalance(map[string]interface{}{
		"type": "spot",
	})
}

func (b *BinanceAdapter) GetPrices(symbol string) (float64, error) {
	b.Client.LoadMarkets()
	ticker, err := b.Client.FetchTicker(fmt.Sprintf("%s/USDT", symbol))
	if err != nil {
		return 0, err
	}

	if ticker.Last == nil {
		return 0, fmt.Errorf("ticker last price is nil")
	}

	return *ticker.Last, nil
}

func (b *BinanceAdapter) GetChains(token, withdrawChain string) (*models.ChainList, error) {
	networkName, ok := b.getNetworkName(withdrawChain, "binance")
	if !ok {
		return nil, fmt.Errorf("софт не поддерживает вывод в сеть %s с биржи Binance", withdrawChain)
	}

	result := <-b.Client.LoadMarkets()
	if err, ok := result.(error); ok && err != nil {
		return nil, fmt.Errorf("failed to load markets: %w", err)
	}

	curRaw, exists := b.Client.Currencies[token]
	if !exists {
		return nil, fmt.Errorf("token %s not found", token)
	}

	var curParse *models.BinanceCurrencyList
	if err := utils.ResponseConvert(curRaw, &curParse); err != nil {
		return nil, err
	}
	// log.Printf("cur: %+v", curParse)

	var chainParams models.ChainList
	for _, param := range curParse.Networks {
		if param.Network == networkName {
			if !param.Info.Busy {
				chainParams.Chain = param.Id
				if fee, err := utils.ConvertToFloat(param.Fee); err == nil {
					chainParams.WithdrawFee = fee
				}
			}
		}
	}
	return &chainParams, nil
}

func (b *BinanceAdapter) getNetworkName(network, cex string) (string, bool) {
	name, exists := globals.ChainNameToSymbolCEX[cex][network]
	return name, exists
}
