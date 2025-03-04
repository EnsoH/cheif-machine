package process

import (
	"cw/logger"
	"cw/models"
)

func getTokenSums(withdrawArr []models.WithdrawAction) map[string]float64 {
	sums := map[string]float64{}
	for _, act := range withdrawArr {
		sums[act.Currency] += act.Amount
	}
	return sums
}

func loggingActions(actions []models.WithdrawAction) {
	for _, act := range actions {
		logger.GlobalLogger.Infof("[%s] | %s | | %s | | %.8f|", act.Address, act.Chain, act.Currency, act.Amount)
	}
	return
}
