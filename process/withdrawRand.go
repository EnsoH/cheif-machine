package process

import (
	"context"
	"cw/models"
	"fmt"
	"math"
	"math/rand"
	"sync"

	"golang.org/x/sync/errgroup"
)

func WithdrawFactory(withdrawConfig *models.WithdrawConfig, addresses []string) ([]models.WithdrawAction, error) {
	if len(addresses) == 0 {
		return nil, fmt.Errorf("Нет списка адресов.")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	var (
		result = make([]models.WithdrawAction, len(addresses))
		mu     sync.Mutex
	)

	for i, address := range addresses {
		address := address
		g.Go(func() error {
			action := withdrawActionInit(withdrawConfig, address)

			mu.Lock()
			result[i] = *action
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func withdrawActionInit(withdrawConfig *models.WithdrawConfig, address string) *models.WithdrawAction {
	chain := getRandomChain(withdrawConfig.Chain)
	currency := getRandomChain(withdrawConfig.Currency)
	amount := getRandomAmount(withdrawConfig.AmountRange)

	return &models.WithdrawAction{
		Address:  address,
		CEX:      withdrawConfig.CEX,
		Chain:    chain,
		Currency: currency,
		Amount:   amount,
	}
}

func getRandomChain(chains []string) string {
	return chains[rand.Intn(len(chains))]
}

func getRandomAmount(amountArr []float64) float64 {
	if len(amountArr) == 0 {
		return amountArr[0]
	}

	min, max := amountArr[0], amountArr[1]
	if min > max {
		min, max = max, min
	}

	if min == max {
		return min
	}

	randoValue := min + rand.Float64()*(max-min)
	return math.Round(randoValue*100) / 100
}
