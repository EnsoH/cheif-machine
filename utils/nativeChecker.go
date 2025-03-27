package utils

import "github.com/ethereum/go-ethereum/common"

func IsNativeToken(tokenAddr string) bool {
	nativeTokens := map[string]bool{
		common.Address{}.Hex(): true,
		"ETH":                  true,
		"BNB":                  true,
		"MATIC":                true,
		"AVAX":                 true,
		// globals.WETH: true,
		// globals.NULL: true,
	}

	return nativeTokens[tokenAddr]
}
