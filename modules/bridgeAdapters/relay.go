package bridgeAdapters

import (
	"cw/account"
	"cw/ethClient"
	"cw/globals"
	"cw/httpClient"
	"cw/logger"
	"cw/models"
	"cw/utils"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Relay struct {
	EthClients map[string]*ethClient.EthClient
	HttpClient *httpClient.HttpClient
	Endpoint   string
}

func NewRelay(endpoint string, clients map[string]*ethClient.EthClient, hc *httpClient.HttpClient) (*Relay, error) {
	if endpoint == "" || clients == nil {
		return nil, fmt.Errorf("missing parameter for init Relay")
	}

	return &Relay{
		EthClients: clients,
		HttpClient: hc,
		Endpoint:   endpoint,
	}, nil
}

func (r *Relay) Bridge(fromChain, destChain, fromToken, toToken string, amount *big.Int, acc *account.Account) error {
	txData, err := r.getTxData(fromChain, destChain, fromToken, toToken, amount, acc, r.HttpClient)
	if err != nil {
		return err
	}

	client, ok := r.EthClients[fromChain]
	if !ok {
		return fmt.Errorf("eth client for chain %q not found", fromChain)
	}

	if _, err := client.ApproveTx(txData.TokenIn, txData.To, acc, txData.Value, false); err != nil {
		return err
	}

	return client.SendTransaction(
		acc.PrivateKey,
		acc.Address,
		txData.To,
		client.GetNonce(acc.Address),
		txData.Value,
		txData.Data,
	)
}

func (r *Relay) getTxData(fromChain, destChain, tokenIn, tokenOut string, amountIn *big.Int, acc *account.Account, client *httpClient.HttpClient) (*models.RelayTransactionData, error) {
	quoteData, err := r.getQuoteData(fromChain, destChain, tokenIn, tokenOut, amountIn, acc, client)
	if err != nil {
		return nil, err
	}

	return r.prepareData(quoteData)

}

func (r *Relay) getQuoteData(fromChain, destChain, tokenIn, tokenOut string, amountIn *big.Int, acc *account.Account, client *httpClient.HttpClient) (*models.RelayResponse, error) {
	relayParams, err := r.getRelayQuoteParams(fromChain, destChain, tokenIn, tokenOut)
	if err != nil {
		return nil, err
	}

	request := models.RelayRequest{
		User:                 acc.Address.Hex(),
		OriginChainId:        relayParams.OriginChainId,
		DestinationChainId:   relayParams.DestinationChainId,
		OriginCurrency:       relayParams.TokenIn,
		DestinationCurrency:  relayParams.TokenOut,
		Recipient:            acc.Address.Hex(),
		TradeType:            "EXACT_INPUT",
		Amount:               amountIn.String(),
		Referrer:             "relay.link/swap",
		UseExternalLiquidity: false,
		UseDepositAddress:    false,
	}
	// log.Printf("req: %+v", request)
	var result models.RelayResponse
	if err := client.SendJSONRequest(r.Endpoint, "POST", request, &result, nil); err != nil {
		if contextErr, critical := utils.IsCriticalError(err); critical {
			logger.GlobalLogger.Error(contextErr)
		}
		// log.Printf("err: %v", err)
		return nil, err
	}

	return &result, nil
}

func (r *Relay) getRelayQuoteParams(fromChain, destChain, tokenIn, tokenOut string) (*models.RelayQuoteModel, error) {
	originClient, ok := ethClient.GlobalETHClient[fromChain]
	if !ok {
		return nil, fmt.Errorf("unknown chain: %s", fromChain)
	}
	originChainId, err := originClient.GetChainID()
	if err != nil {
		return nil, err
	}

	destClient, ok := ethClient.GlobalETHClient[destChain]
	if !ok {
		return nil, fmt.Errorf("unknown chain: %s", destChain)
	}
	destinationChainId, err := destClient.GetChainID()
	if err != nil {
		return nil, err
	}

	tokenInContract, ok := globals.TokenContracts[originChainId][tokenIn]
	if !ok {
		return nil, fmt.Errorf("incorrect 'token from' parameter for bridge")
	}
	tokenOutContract, ok := globals.TokenContracts[destinationChainId][tokenOut]
	if !ok {
		return nil, fmt.Errorf("incorrect 'token out' parameter for bridge")
	}

	return &models.RelayQuoteModel{
		TokenIn:            tokenInContract,
		TokenOut:           tokenOutContract,
		OriginChainId:      originChainId,
		DestinationChainId: destinationChainId,
	}, nil
}

func (r *Relay) prepareData(quoteData *models.RelayResponse) (*models.RelayTransactionData, error) {
	if err := r.validateQuoteData(quoteData); err != nil {
		return nil, err
	}

	stepData := quoteData.Steps[0].Items[0].Data

	value, ok := new(big.Int).SetString(stepData.Value, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert value to big.Int: %s", stepData.Value)
	}
	decodedData, err := hex.DecodeString(strings.TrimPrefix(stepData.Data, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	relayTx := &models.RelayTransactionData{
		Value:   value,
		Data:    decodedData,
		To:      common.HexToAddress(stepData.To),
		TokenIn: common.HexToAddress(quoteData.Details.CurrencyIn.Currency.Address),
	}

	return relayTx, nil
}

func (r *Relay) validateQuoteData(quoteData *models.RelayResponse) error {
	if quoteData == nil {
		return fmt.Errorf("quoteData is nil")
	}

	if quoteData.Message != "" && quoteData.ErrCode != "" {
		return fmt.Errorf("API error: %s (code: %s)", quoteData.Message, quoteData.ErrCode)
	}

	if quoteData.Details.Impact.Percent == "" {
		return fmt.Errorf("impact percent is missing in response")
	}

	if p, err := strconv.ParseFloat(quoteData.Details.Impact.Percent, 64); err != nil || p > globals.MaxPercent {
		return fmt.Errorf("превышен максимально допустимый импакт, импакт - %s", quoteData.Details.Impact.Percent)
	}

	if len(quoteData.Steps) == 0 || len(quoteData.Steps[0].Items) == 0 {
		return fmt.Errorf("invalid response format: no steps/items found")
	}

	if quoteData.Steps[0].Items[0].Data.Value == "" {
		return fmt.Errorf("missing value in response data")
	}

	if quoteData.Steps[0].Items[0].Data.Data == "" {
		return fmt.Errorf("missing data field in response")
	}

	return nil
}
