package account

import (
	"context"
	"crypto/ecdsa"
	"cw/globals"
	"cw/logger"
	"cw/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
)

type Account struct {
	Address         common.Address
	PrivateKey      *ecdsa.PrivateKey
	DestinationAddr common.Address
	Mu              sync.Mutex
}

type AccountOption func(*Account)

func WithAddress(addr common.Address) AccountOption {
	return func(a *Account) {
		a.Address = addr
	}
}

func WithPrivateKey(pk *ecdsa.PrivateKey) AccountOption {
	return func(a *Account) {
		a.PrivateKey = pk
	}
}

func WithDestination(dest common.Address) AccountOption {
	return func(a *Account) {
		a.DestinationAddr = dest
	}
}

func NewAccount(opts ...AccountOption) *Account {
	acc := &Account{}
	for _, opt := range opts {
		opt(acc)
	}
	return acc
}

func AccsFactory(module string) ([]*Account, error) {
	addrPath, err := utils.GetPath(globals.Addresses)
	if err != nil {
		logger.GlobalLogger.Error(err)
		return nil, err
	}

	inputs, err := utils.FileReader(addrPath)
	if err != nil || len(inputs) == 0 {
		logger.GlobalLogger.Error("Ошибка чтения данных или файл пуст")
		return nil, errors.New("input list is empty")
	}

	if module == "WalletGenerator" {
		return []*Account{}, nil
	}

	accounts, err := processPrivateKeys(inputs)
	if err != nil {
		return nil, err
	}

	if module == "Сollector" {
		destMap, err := readDestinations()
		if err != nil {
			return nil, err
		}
		applyDestinations(accounts, destMap)
	}

	return accounts, nil
}

func readDestinations() (map[common.Address]common.Address, error) {
	destPath, err := utils.GetPath(globals.Destinations)
	if err != nil {
		logger.GlobalLogger.Error(err)
		return nil, err
	}

	data, err := ioutil.ReadFile(destPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла destination адресов: %w", err)
	}

	var rawDestMap map[string]string
	if err := json.Unmarshal(data, &rawDestMap); err != nil {
		return nil, fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	destMap := make(map[common.Address]common.Address)
	for mainAddr, destAddr := range rawDestMap {
		if !common.IsHexAddress(mainAddr) || !common.IsHexAddress(destAddr) {
			logger.GlobalLogger.Warnf("Некорректный адрес в ассоциации: %s -> %s", mainAddr, destAddr)
			continue
		}
		destMap[common.HexToAddress(mainAddr)] = common.HexToAddress(destAddr)
	}

	return destMap, nil
}

func applyDestinations(accounts []*Account, destMap map[common.Address]common.Address) {
	for _, acc := range accounts {
		if destAddr, ok := destMap[acc.Address]; ok {
			acc.DestinationAddr = destAddr
		}
	}
}

func processPrivateKeys(inputs []string) ([]*Account, error) {
	var (
		mu       sync.Mutex
		accounts = make([]*Account, 0, len(inputs))
	)
	g, _ := errgroup.WithContext(context.Background())
	for _, input := range inputs {
		input := input
		g.Go(func() error {
			priv, err := utils.ParsePrivateKey(input)
			if err != nil {
				return err
			}
			addr, err := utils.DeriveAddress(priv)
			if err != nil {
				return err
			}
			mu.Lock()
			accounts = append(accounts, NewAccount(WithAddress(addr), WithPrivateKey(priv)))
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return accounts, nil
}
