package process

import (
	"context"
	"cw/logger"
	"cw/models"
	"cw/modules"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func ActionsProcess(addresses []string, mod map[string]modules.ModulesFasad, cex string) error {
	actions, err := WithdrawFactory(addresses)
	if err != nil {
		return err
	}

	if err := validateActions(actions, mod, cex); err != nil {
		return err
	}

	return withdrawProcess(actions, mod)
}

func withdrawProcess(actions []models.WithdrawAction, mod map[string]modules.ModulesFasad) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	for _, act := range actions {
		act := act
		g.Go(func() error {
			amount, err := calculateAmount(act.Currency, act.Amount, mod[act.CEX])
			if err != nil {
				return err
			}

			logger.GlobalLogger.Infof("Sleep before withdraw %v", act.TimeRange)
			time.Sleep(time.Second * time.Duration(act.TimeRange))
			return mod[act.CEX].Withdraw(act.Currency, act.Address, act.Chain, amount)
		})
	}

	return g.Wait()
}

func validateActions(actions []models.WithdrawAction, mod map[string]modules.ModulesFasad, cex string) error {
	sums := getTokenSums(actions)

	for token, amount := range sums {
		balance, err := mod[cex].GetBalances(token)
		if err != nil {
			return err
		}

		tokenPrice, err := mod[cex].GetPrices(token)
		if err != nil {
			return err
		}

		// Add 1% to the balance to account for withdrawal fees
		if (balance * 1.01) <= amount/tokenPrice {
			return fmt.Errorf("There is not enough balance in the token: %s. Total amount: %v, CEX account balance: %v", token, sums[token], balance)
		}
	}

	return nil
}

func calculateAmount(token string, amount float64, mod modules.ModulesFasad) (float64, error) {
	tickerPrice, err := mod.GetPrices(token)
	if err != nil {
		return 0.0, err
	}

	return (amount / tickerPrice), nil
}
