package process

import (
	"cw/account"
	"cw/config"
	"cw/logger"
	"cw/modules"
	"math/rand"
	"sync"
	"time"
)

func (ac *ActionCore) CollectorAction(accs []*account.Account, mod *modules.Modules, selectModule string) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.Cfg.Threads)

	for _, acc := range accs {
		wg.Add(1)

		go func(a *account.Account) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			randVal := r.Intn(260) + 1
			logger.GlobalLogger.Infof("[%s] Processing with delay %d sec", a.Address.Hex(), randVal)

			time.Sleep(time.Second * time.Duration(randVal))

			if err := mod.Collector.Collect(a, []string{"Base", "Optimism", "Arbitrum"}); err != nil {
				logger.GlobalLogger.Errorf("Error processing account %s: %v", a.Address.Hex(), err)
				// return
			}
		}(acc)
	}

	wg.Wait()
	return nil
}
