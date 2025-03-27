package modules

import (
	"cw/utils"
	"fmt"
	"time"
)

type WalletGenerator struct {
	WalletGenMap map[string]WallGenModule
}

func NewWalletGen(walletTypes ...string) (*WalletGenerator, error) {
	wg := &WalletGenerator{
		WalletGenMap: make(map[string]WallGenModule),
	}

	for _, typeName := range walletTypes {
		if opt, ok := walletGeneratorOptionsMap[typeName]; ok {
			opt(wg)
		}
	}

	return wg, nil
}

func (w *WalletGenerator) GenerateWallets(walletType string, count int, csvFormat []string) error {
	generator, err := w.getAdapter(walletType)
	if err != nil {
		return err
	}

	fileName := w.generateFileName(walletType)
	headers := csvFormat

	for i := 0; i < count; i++ {
		pk, addr, mnemonic, err := generator.GenerateWallet()
		if err != nil {
			return err
		}

		row := buildWalletRow(headers, pk, addr, mnemonic)
		if err := utils.CsvWriter(headers, row, fileName); err != nil {
			return err
		}
	}
	return nil
}

func (w *WalletGenerator) getAdapter(walletType string) (WallGenModule, error) {
	generator, ok := w.WalletGenMap[walletType]
	if !ok {
		return nil, fmt.Errorf("нет генератора для этого типа кошелька: %v", walletType)
	}

	return generator, nil
}

func (w *WalletGenerator) generateFileName(walletType string) (fileName string) {
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	return fmt.Sprintf("%s-%s.csv", walletType, currentTime)
}

func buildWalletRow(headers []string, pk, addr, mnemonic string) []string {
	row := make([]string, 0, len(headers))
	for _, h := range headers {
		switch h {
		case "private_key":
			row = append(row, pk)
		case "address":
			row = append(row, addr)
		case "mnemonic":
			row = append(row, mnemonic)
		default:
			row = append(row, "")
		}
	}
	return row
}
