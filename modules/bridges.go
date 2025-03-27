package modules

import (
	"cw/account"
	"cw/httpClient"
	"fmt"
	"math/big"
)

type Bridger struct {
	BridgesMap map[string]BridgerModule
	HttpClient *httpClient.HttpClient
}

func NewBridgeModule(hc *httpClient.HttpClient, bridgeName ...string) (*Bridger, error) {
	if hc == nil {
		return nil, fmt.Errorf("http client is nil")
	}

	br := &Bridger{
		HttpClient: hc,
		BridgesMap: make(map[string]BridgerModule),
	}

	for _, bridge := range bridgeName {
		if _, ok := bridgeOptionsMap[bridge]; !ok {
			continue
		}

		bridgeOptionsMap[bridge](br)
	}

	return br, nil
}

func (b *Bridger) Bridge(bridgeName, fromChain, destChain, fromToken, toToken string, amount *big.Int, acc *account.Account) error {
	module, ok := b.BridgesMap[bridgeName]
	if !ok {
		return fmt.Errorf("failed to get bridge %q. Check config", bridgeName)
	}

	return module.Bridge(fromChain, destChain, fromToken, toToken, amount, acc)
}
