package modules

import (
	"cw/config"
	"cw/ethClient"
	"cw/logger"
	"cw/modules/bridgeAdapters"
	"cw/modules/exchangeAdapters"
	"cw/modules/walletGeneratorAdapters"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

// ########################### EXCHANGE OPTIONS ########################################
var exhangeOptionsMap = map[string]ExchangeOption{
	"bybit":   withBybit(),
	"binance": withBinance(),
	"mexc":    withMexc(),
	"kucoin":  withKucoin(),
	"okx":     withOkx(),
}

type ExchangeOption func(*Exchanges)

func withBybit() ExchangeOption {
	return func(e *Exchanges) {
		bybit := ccxt.NewBybit(map[string]interface{}{
			"apiKey":          config.Cfg.CEXConfigs.BybitCfg.API_key,
			"secret":          config.Cfg.CEXConfigs.BybitCfg.API_secret,
			"enableRateLimit": true,
		})
		bybit.HttpProxy = config.Cfg.IpAddresses[0]
		e.CEXs["bybit"] = &exchangeAdapters.BybitAdapter{Client: &bybit}
	}
}

func withBinance() ExchangeOption {
	return func(e *Exchanges) {
		binance := ccxt.NewBinance(map[string]interface{}{
			"apiKey":          config.Cfg.CEXConfigs.BinanceCfg.API_key,
			"secret":          config.Cfg.CEXConfigs.BinanceCfg.API_secret,
			"enableRateLimit": true,
			"options": map[string]interface{}{
				"defaultType": "spot",
			},
		})
		binance.HttpProxy = config.Cfg.IpAddresses[0]
		e.CEXs["binance"] = &exchangeAdapters.BinanceAdapter{Client: &binance}
	}
}

func withMexc() ExchangeOption {
	return func(e *Exchanges) {
		mexc := ccxt.NewMexc(map[string]interface{}{
			"apiKey":          config.Cfg.CEXConfigs.MexcCfg.API_key,
			"secret":          config.Cfg.CEXConfigs.MexcCfg.API_secret,
			"enableRateLimit": true,
		})
		mexc.HttpProxy = config.Cfg.IpAddresses[0]
		e.CEXs["mexc"] = &exchangeAdapters.MexcAdapter{Client: &mexc}
	}
}

func withKucoin() ExchangeOption {
	return func(e *Exchanges) {
		kucoin := ccxt.NewKucoin(map[string]interface{}{
			"apiKey":          config.Cfg.CEXConfigs.KucoinCfg.API_key,
			"secret":          config.Cfg.CEXConfigs.KucoinCfg.API_secret,
			"enableRateLimit": true,
			"password":        config.Cfg.CEXConfigs.KucoinCfg.Password,
		})
		kucoin.HttpProxy = config.Cfg.IpAddresses[0]
		e.CEXs["kucoin"] = &exchangeAdapters.KucoinAdapter{Client: &kucoin}

	}
}

func withOkx() ExchangeOption {
	return func(e *Exchanges) {
		okx := ccxt.NewOkx(map[string]interface{}{
			"apiKey":          config.Cfg.CEXConfigs.OkxCfg.API_key,
			"secret":          config.Cfg.CEXConfigs.OkxCfg.API_secret,
			"enableRateLimit": true,
			"password":        config.Cfg.CEXConfigs.OkxCfg.Password,
		})
		okx.HttpProxy = config.Cfg.IpAddresses[0]
		e.CEXs["okx"] = &exchangeAdapters.OkxAdapter{Client: &okx}
	}
}

// func withTokensMap() ExchangeOption {
// 	return func(e *Exchanges) {
// 		e.Tokens = globals.TokenNamesMap
// 	}
// }

// func withDecimals() ExchangeOption {
// 	return func(e *Exchanges) {
// 		// e.Decimals = globals.DecimalsMap
// 	}
// }

// ########################### BRIDGE OPTIONS ##########################################
type BridgeOptions func(*Bridger)

var bridgeOptionsMap = map[string]BridgeOptions{
	"relay": withRelay(),
}

func withRelay() BridgeOptions {
	return func(b *Bridger) {
		relay, err := bridgeAdapters.NewRelay(config.Cfg.Endpoints["relay"], ethClient.GlobalETHClient, b.HttpClient)
		if err != nil {
			logger.GlobalLogger.Error(err)
			return
		}
		b.BridgesMap["relay"] = relay
	}
}

// ########################### WALLET GENERATORS OPTIONS ##########################################
type WalletGenOptions func(w *WalletGenerator)

var walletGeneratorOptionsMap = map[string]WalletGenOptions{
	"evm": withEvm(),
	"sol": withSol(),
	"btc": withBtc(),
}

func withEvm() WalletGenOptions {
	return func(w *WalletGenerator) {
		w.WalletGenMap["evm"] = walletGeneratorAdapters.NewEvmAdapter()
	}
}

func withBtc() WalletGenOptions {
	return func(w *WalletGenerator) {
		w.WalletGenMap["btc"] = walletGeneratorAdapters.NewBtcAdapter()
	}
}

func withSol() WalletGenOptions {
	return func(w *WalletGenerator) {
		w.WalletGenMap["sol"] = walletGeneratorAdapters.NewSolAdapter()
	}
}
