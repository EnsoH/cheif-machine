package account

import (
	"context"
	"crypto/ecdsa"
	"cw/config"
	"cw/globals"
	"cw/logger"
	"cw/utils"
	"errors"
	"os"
	"strings"
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

	if module == "CexWithdrawer" {
		accounts, err := processAddresses(inputs)
		if err != nil {
			return nil, err
		}
		return accounts, nil
	}

	accounts, err := processPrivateKeys(inputs)
	if err != nil {
		return nil, err
	}

	if module == "Collector" {
		destMap, defaultDest, err := readDestinations()
		if err != nil {
			return nil, err
		}
		applyDestinations(accounts, destMap, defaultDest)
	}

	return accounts, nil
}

func readDestinations() (map[common.Address]common.Address, common.Address, error) {
	destMap := make(map[common.Address]common.Address)
	var defaultDest common.Address

	if len(config.UserCfg.CollectorConfig.DestinationAddresses) > 0 {
		for mainAddr, destAddr := range config.UserCfg.CollectorConfig.DestinationAddresses {
			mainAddr = strings.ToLower(strings.TrimSpace(mainAddr))
			destAddr = strings.ToLower(strings.TrimSpace(destAddr))

			if !common.IsHexAddress(mainAddr) || !common.IsHexAddress(destAddr) {
				logger.GlobalLogger.Errorf("Некорректный адрес в ассоциации: %s -> %s", mainAddr, destAddr)
				os.Exit(1)
			}
			destMap[common.HexToAddress(mainAddr)] = common.HexToAddress(destAddr)
		}
	} else {
		dest := strings.ToLower(strings.TrimSpace(config.UserCfg.CollectorConfig.DestinationAddress))
		if common.IsHexAddress(dest) {
			defaultDest = common.HexToAddress(dest)
		} else {
			logger.GlobalLogger.Errorf("Неверный формат destination_address: %s", dest)
			os.Exit(1)
		}
	}

	return destMap, defaultDest, nil
}

func applyDestinations(accounts []*Account, destMap map[common.Address]common.Address, defaultDest common.Address) {
	for _, acc := range accounts {
		if !common.IsHexAddress(acc.DestinationAddr.Hex()) {
			logger.GlobalLogger.Errorf("Некорректный destination-адрес для аккаунта: %s", acc.Address.Hex())
			os.Exit(1)
		}

		if destAddr, ok := destMap[acc.Address]; ok {
			acc.DestinationAddr = destAddr
		} else {
			acc.DestinationAddr = defaultDest
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

func processAddresses(inputs []string) ([]*Account, error) {
	var (
		accounts = make([]*Account, 0, len(inputs))
	)

	for _, input := range inputs {
		input := input
		if !common.IsHexAddress(input) {
			return nil, errors.New("invalid address format: " + input)
		}
		addr := common.HexToAddress(input)
		accounts = append(accounts, NewAccount(WithAddress(addr)))
	}

	return accounts, nil
}
