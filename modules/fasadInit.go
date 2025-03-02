package modules

import (
	"context"
	"cw/config"
	"cw/httpClient"
	"log"
	"sync"

	"golang.org/x/sync/errgroup"
)

type ModulesFasad interface {
	Withdraw(token, address, network string, amount float64) error
	GetBalances(token string) (float64, error)
	GetPrices(token string) (float64, error)
}

type ModuleFactory func() (ModulesFasad, error)

func ModulesInit() (map[string]ModulesFasad, error) {
	log.Printf(config.Cfg.IpAddresses[0])
	hc, err := httpClient.NewHttpClient(
		httpClient.WithHttp2(),
		httpClient.WithProxy(config.Cfg.IpAddresses[0]),
	)
	if err != nil {
		return nil, err
	}

	modules := map[string]ModuleFactory{
		"bybit": func() (ModulesFasad, error) {
			return NewBybitModule(
				config.Cfg.CEXConfigs.BybitCfg.BalanceEndpoint,
				config.Cfg.CEXConfigs.BybitCfg.TickersEndpoint,
				config.Cfg.CEXConfigs.BybitCfg.API_key,
				config.Cfg.CEXConfigs.BybitCfg.API_secret,
				config.Cfg.IpAddresses[0],
				hc,
			)
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	defer cancel()

	var (
		mu     sync.Mutex
		result = make(map[string]ModulesFasad, len(modules))
	)

	for name, factory := range modules {
		name, factory := name, factory

		g.Go(func() error {
			module, err := factory()
			if err != nil {
				return err
			}
			mu.Lock()
			result[name] = module
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}
