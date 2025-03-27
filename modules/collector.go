package modules

import (
	"cw/account"
	"cw/ethClient"
	"cw/logger"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Collectors struct{}

func NewCollector() *Collectors { return &Collectors{} }

func (c *Collectors) Collect(acc *account.Account, chains []string) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, chain := range chains {
		client, err := c.getClient(chain)
		if err != nil {
			return err
		}

		balance, err := c.validateBalance(acc.Address, client)
		if err != nil {
			// return err
			continue
		}

		if err := client.SendTransaction(acc.PrivateKey, acc.Address, acc.DestinationAddr, client.GetNonce(acc.Address), balance, nil); err != nil {
			// return err
			continue
		}
		randVal := r.Intn(60) + 1
		logger.GlobalLogger.Infof("[%s] Ждем %d сек перед следующей отправкой", acc.Address.Hex(), randVal)

		time.Sleep(time.Second * time.Duration(randVal))

	}
	return nil
}

func (c *Collectors) getClient(chain string) (*ethClient.EthClient, error) {
	if client, ok := ethClient.GlobalETHClient[chain]; ok {
		return client, nil
	}
	return nil, fmt.Errorf("сети %s нет в софте", chain)
}

func (c *Collectors) validateBalance(addr common.Address, client *ethClient.EthClient) (*big.Int, error) {
	balance, err := client.BalanceCheck(addr, common.Address{})
	if err != nil {
		return nil, err
	}

	if balance.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("недостаточно баланса, текущий баланс: %v", balance)
	}

	onePercent := new(big.Int).Div(balance, big.NewInt(100))
	finalBalance := new(big.Int).Sub(balance, onePercent)

	return finalBalance, nil
}
