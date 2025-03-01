package process

import "cw/models"

func getTokenSums(withdrawArr []models.WithdrawAction) map[string]float64 {
	sums := map[string]float64{}
	for _, act := range withdrawArr {
		sums[act.Currency] += act.Amount
	}
	return sums
}
