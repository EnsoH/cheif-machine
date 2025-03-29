package modules

import (
	"cw/config"
	"cw/ethClient"
	"cw/logger"
	"cw/modules/bridgeAdapters"
	"cw/modules/exchangeAdapters"
	"cw/modules/walletGeneratorAdapters"
	"fmt"
	"os"
	"strings"

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
		if err := validateApiKeys("bybit"); err != nil {
			logger.GlobalLogger.Error(err)
			os.Exit(1)
		}

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
		if err := validateApiKeys("binance"); err != nil {
			logger.GlobalLogger.Error(err)
			os.Exit(1)
		}

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
		if err := validateApiKeys("mexc"); err != nil {
			logger.GlobalLogger.Error(err)
			os.Exit(1)
		}

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
		if err := validateApiKeys("kucoin"); err != nil {
			logger.GlobalLogger.Error(err)
			os.Exit(1)
		}

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
		if err := validateApiKeys("okx"); err != nil {
			logger.GlobalLogger.Error(err)
			os.Exit(1)
		}

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

// ############## API KEY CHECKER #################
func validateApiKeys(module string) error {
	switch module {
	case "bybit":
		if strings.TrimSpace(config.Cfg.CEXConfigs.BybitCfg.API_key) == "" || strings.TrimSpace(config.Cfg.CEXConfigs.BybitCfg.API_secret) == "" {
			return fmt.Errorf("ключи для bybit в конфиге пустые")
		}
	case "binance":
		if strings.TrimSpace(config.Cfg.CEXConfigs.BinanceCfg.API_key) == "" || strings.TrimSpace(config.Cfg.CEXConfigs.BinanceCfg.API_secret) == "" {
			return fmt.Errorf("ключи для binance в конфиге пустые")
		}
	case "okx":
		if strings.TrimSpace(config.Cfg.CEXConfigs.OkxCfg.API_key) == "" || strings.TrimSpace(config.Cfg.CEXConfigs.OkxCfg.API_secret) == "" || strings.TrimSpace(config.Cfg.CEXConfigs.OkxCfg.Password) == "" {
			return fmt.Errorf("ключи для okx в конфиге пустые")
		}
	case "mexc":
		if strings.TrimSpace(config.Cfg.CEXConfigs.MexcCfg.API_key) == "" || strings.TrimSpace(config.Cfg.CEXConfigs.MexcCfg.API_secret) == "" {
			// logger.GlobalLogger.Errorf("ключ в конфиге пустой")
			return fmt.Errorf("ключи для mexc в конфиге пустые")
		}
	case "kucoin":
		if strings.TrimSpace(config.Cfg.CEXConfigs.KucoinCfg.API_key) == "" || strings.TrimSpace(config.Cfg.CEXConfigs.KucoinCfg.API_secret) == "" || strings.TrimSpace(config.Cfg.CEXConfigs.KucoinCfg.Password) == "" {
			return fmt.Errorf("ключи для kucoin в конфиге пустые")
		}
	}

	return nil
}
