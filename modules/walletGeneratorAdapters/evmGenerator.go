package walletGeneratorAdapters

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type EvmAdapter struct{}

func NewEvmAdapter() *EvmAdapter { return &EvmAdapter{} }

func (e *EvmAdapter) GenerateWallet() (privateKey string, address string, mnemonic string, err error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err = bip39.NewMnemonic(entropy)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}
	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate master key: %w", err)
	}

	purpose, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to derive purpose key: %w", err)
	}
	coinType, err := purpose.NewChildKey(bip32.FirstHardenedChild + 60)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to derive coin type key: %w", err)
	}
	account, err := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to derive account key: %w", err)
	}
	change, err := account.NewChildKey(0)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to derive change key: %w", err)
	}
	addressIndex, err := change.NewChildKey(0)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to derive address index key: %w", err)
	}

	privateKeyBytes := addressIndex.Key
	privateKeyECDSA, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to convert to ECDSA: %w", err)
	}

	publicKeyECDSA := privateKeyECDSA.Public().(*ecdsa.PublicKey)
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	privateKey = fmt.Sprintf("0x%x", privateKeyBytes)

	return privateKey, address, mnemonic, nil
}
