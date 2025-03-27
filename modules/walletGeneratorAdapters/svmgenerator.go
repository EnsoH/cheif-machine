package walletGeneratorAdapters

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/mr-tron/base58"
	"github.com/tyler-smith/go-bip39"
)

type SolAdapter struct{}

func NewSolAdapter() *SolAdapter { return &SolAdapter{} }

func (s *SolAdapter) GenerateWallet() (privateKey string, address string, mnemonic string, err error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err = bip39.NewMnemonic(entropy)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	seed := bip39.NewSeed(mnemonic, "")
	if len(seed) < 32 {
		return "", "", "", fmt.Errorf("seed too short")
	}

	seed32 := seed[:32]
	privKey := ed25519.NewKeyFromSeed(seed32)
	pubKey := privKey.Public().(ed25519.PublicKey)
	address = base58.Encode(pubKey)
	privateKey = hex.EncodeToString(privKey)

	return privateKey, address, mnemonic, nil
}
