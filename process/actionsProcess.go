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

func ActionsProcess(addresses []string, exchange modules.Exchanges, cex string) error {
	actions, err := WithdrawFactory(addresses)
	if err != nil {
		return err
	}

	if err := validateActions(actions, exchange, cex); err != nil {
		return err
	}

	return withdrawProcess(actions, exchange, cex)
}

func withdrawProcess(actions []models.WithdrawAction, exchange modules.Exchanges, cex string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, _ := errgroup.WithContext(ctx)

	for _, act := range actions {
		act := act
		g.Go(func() error {
			amount, err := calculateAmount(act.Currency, act.Amount, exchange, cex)
			if err != nil {
				return err
			}

			logger.GlobalLogger.Infof("Sleep before withdraw %v", act.TimeRange)
			time.Sleep(time.Second * time.Duration(act.TimeRange))
			return exchange.Withdraw(cex, act.Currency, act.Address, act.Chain, amount)
		})
	}

	return g.Wait()
}

func validateActions(actions []models.WithdrawAction, exchange modules.Exchanges, cex string) error {
	sums := getTokenSums(actions)

	for token, amount := range sums {
		balance, err := exchange.GetBalances(cex, token)
		if err != nil {
			return err
		}

		tokenPrice, err := exchange.GetPrices(cex, token)
		if err != nil {
			return err
		}

		// Add 1% to the balance to account for withdrawal fees
		if (balance * 1.01) <= amount/tokenPrice {
			return fmt.Errorf("there is not enough balance in the token: %s. Total amount: %v, CEX account balance: %v", token, sums[token], balance)
		}
	}

	return nil
}

func calculateAmount(token string, amount float64, exchange modules.Exchanges, cex string) (float64, error) {
	tickerPrice, err := exchange.GetPrices(cex, token)
	if err != nil {
		return 0.0, err
	}

	return (amount / tickerPrice), nil
}
