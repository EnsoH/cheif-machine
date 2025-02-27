package ethClient

import (
	"github.com/ethereum/go-ethereum/common"
)

func IsNativeToken(tokenAddr common.Address) bool {
	nativeTokens := map[common.Address]bool{
		// globals.WETH: true,
		// globals.NULL: true,
	}

	return nativeTokens[tokenAddr]
}
