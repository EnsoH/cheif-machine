package ethClient

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	// "cw/globals"
	"cw/logger"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
)

var GlobalETHClient map[string]*EthClient

type EthClient struct {
	Client *ethclient.Client
}

// EthClientFactory создает клиентов для всех переданных RPC-узлов.
func EthClientFactory(rpcs map[string]string) error {
	if len(rpcs) == 0 {
		return errors.New("RPC URLs map is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var (
		result = make(map[string]*EthClient)
		mu     sync.Mutex
		g, _   = errgroup.WithContext(ctx)
	)

	for name, rpc := range rpcs {
		name, rpc := name, rpc
		g.Go(func() error {
			client, err := ethclient.DialContext(ctx, rpc)
			if err != nil {
				return fmt.Errorf("error connecting to RPC %s: %v", name, err)
			}
			mu.Lock()
			result[name] = &EthClient{Client: client}
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	GlobalETHClient = result

	return nil
}

// CloseAllClients закрывает все клиенты.
func CloseAllClients(clients map[string]*EthClient) {
	for _, client := range clients {
		if client.Client != nil {
			client.Client.Close()
		}
	}
}

// BalanceCheck проверяет баланс токена или нативной монеты.
func (c *EthClient) BalanceCheck(owner, tokenAddr common.Address) (*big.Int, error) {
	if IsNativeToken(tokenAddr) {
		return c.Client.BalanceAt(context.Background(), owner, nil)
	}

	// data, err := globals.Erc20ABI.Pack("balanceOf", owner)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to pack data: %v", err)
	// }

	// result, err := c.CallCA(tokenAddr, data)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to call contract: %v", err)
	// }

	// var balance *big.Int
	// if err := globals.Erc20ABI.UnpackIntoInterface(&balance, "balanceOf", result); err != nil {
	// 	return nil, fmt.Errorf("failed to unpack result: %v", err)
	// }

	// return balance, nil
	return nil, nil
}

// CallCA выполняет вызов контракта.
func (c *EthClient) CallCA(toCA common.Address, data []byte) ([]byte, error) {
	callMsg := ethereum.CallMsg{To: &toCA, Data: data}
	return c.Client.CallContract(context.Background(), callMsg, nil)
}

// GetGasValues оценивает лимит газа и комиссии.
func (c *EthClient) GetGasValues(msg ethereum.CallMsg) (uint64, *big.Int, *big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return 0, nil, nil, fmt.Errorf("gas wait timeout exceeded")

		case <-ticker.C:
			header, err := c.Client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return 0, nil, nil, fmt.Errorf("error fetching block header: %w", err)
			}

			tipCap, err := c.Client.SuggestGasTipCap(context.Background())
			if err != nil {
				return 0, nil, nil, fmt.Errorf("error suggesting gas tip cap: %w", err)
			}

			feeCap := new(big.Int).Add(header.BaseFee, tipCap)
			gasLimit, err := c.Client.EstimateGas(context.Background(), msg)
			if err != nil {
				return 0, nil, nil, fmt.Errorf("gas estimation error: %w", err)
			}

			if feeCap.Cmp(big.NewInt(100_000_000_000)) > 0 {
				logger.GlobalLogger.Warnf("High gwei detected: %v", feeCap)
				continue
			}
			return gasLimit, tipCap, feeCap, nil
		}
	}
}

// GetNonce получает nonce для транзакции.
func (c *EthClient) GetNonce(address common.Address) uint64 {
	nonce, err := c.Client.PendingNonceAt(context.Background(), address)
	if err != nil {
		logger.GlobalLogger.Warnf("Failed to get nonce for address %s: %v", address.Hex(), err)
		return 0
	}
	return nonce
}

// GetChainID получает ChainID сети.
func (c *EthClient) GetChainID() (*big.Int, error) {
	return c.Client.NetworkID(context.Background())
}

// SendTransaction отправляет транзакцию.
func (c *EthClient) SendTransaction(privateKey *ecdsa.PrivateKey, ownerAddr, to common.Address, nonce uint64, value *big.Int, txData []byte) error {
	chainID, err := c.GetChainID()
	if err != nil {
		return fmt.Errorf("failed to get ChainID: %v", err)
	}

	gasLimit, tipCap, feeCap, err := c.GetGasValues(ethereum.CallMsg{
		From:  ownerAddr,
		To:    &to,
		Value: value,
		Data:  txData,
	})
	if err != nil {
		return fmt.Errorf("failed to estimate gas: %v", err)
	}

	dynamicTx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       gasLimit,
		To:        &to,
		Value:     value,
		Data:      txData,
	}

	signedTx, err := types.SignTx(types.NewTx(dynamicTx), types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	if err = c.Client.SendTransaction(context.Background(), signedTx); err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	logger.GlobalLogger.Infof("Transaction sent: https://blockscout.monad.com/tx/%s", signedTx.Hash().Hex())
	return c.waitForTransactionSuccess(signedTx.Hash(), 1*time.Minute)
}

// waitForTransactionSuccess ожидает подтверждения транзакции.
func (c *EthClient) waitForTransactionSuccess(txHash common.Hash, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("transaction wait timeout")
		case <-ticker.C:
			receipt, err := c.Client.TransactionReceipt(context.Background(), txHash)
			if err != nil {
				if isUnknownBlockError(err) {
					continue
				}
				return fmt.Errorf("error getting transaction receipt: %v", err)
			}

			if receipt.Status == 1 {
				return nil
			}
			return fmt.Errorf("transaction failed")
		}
	}
}

// isUnknownBlockError проверяет известные ошибки.
func isUnknownBlockError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "Unknown block") ||
		strings.Contains(errMsg, "not found") ||
		strings.Contains(errMsg, "free tier limits")
}
