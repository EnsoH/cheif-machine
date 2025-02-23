package modules

import (
	"context"
	"cw/models"
	"sync"

	"golang.org/x/sync/errgroup"
)

type ModulesFasad interface {
	Action() error
	GetBalances() error
	GetPrices() error
}

type ModuleFactory func(cfg *models.CexConfig) (ModulesFasad, error)

func ModulesInit(cfg *models.CexConfig) (map[string]ModulesFasad, error) {
	modules := map[string]ModuleFactory{
		"bybit": func(cfg *models.CexConfig) (ModulesFasad, error) {
			return NewBybitModule(
				cfg.BybitCfg.WithdrawEndpoint,
				cfg.BybitCfg.API_key,
				cfg.BybitCfg.API_secret,
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
			module, err := factory(cfg)
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
