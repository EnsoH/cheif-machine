package modules

import (
	"cw/config"
	"cw/globals"
	"cw/modules/exchangeAdapters"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

var optionsMap = map[string]Option{
	"bybit":   withBybit(),
	"binance": withBinance(),
	"mexc":    withMexc(),
}

type Option func(*Exchanges)

func withBybit() Option {
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

func withBinance() Option {
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

func withMexc() Option {
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

func withTokensMap() Option {
	return func(e *Exchanges) {
		e.Tokens = globals.TokenNamesMap
	}
}

func withDecimals() Option {
	return func(e *Exchanges) {
		e.Decimals = globals.DecimalsMap
	}
}
