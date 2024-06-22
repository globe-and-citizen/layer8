package blockchain

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PayAsYouGoWrapper interface {
	CreateContract(ctx context.Context, rate uint8, clientId string) (*string, error)
	GetContractByID(ctx context.Context, contractID string) (PayAsYouGoAgreement, error)
}

type PayAsYouGoWrapperImpl struct {
	c         *PayAsYouGo
	signer    TransactionSigner
	rpcClient *ethclient.Client
}

func NewPayAsYouGoWrapper(
	smartContractAddress string,
	rpcClient *ethclient.Client,
	signer TransactionSigner,
) PayAsYouGoWrapper {

	c, err := NewPayAsYouGo(
		common.HexToAddress(smartContractAddress),
		rpcClient,
	)
	if err != nil {
		log.Fatalf("Failed to create PayAsYouGo contract wrapper %s", err)
	}

	return &PayAsYouGoWrapperImpl{
		c:         c,
		signer:    signer,
		rpcClient: rpcClient,
	}
}

func (w *PayAsYouGoWrapperImpl) CreateContract(ctx context.Context, rate uint8, clientId string) (*string, error) {
	chainId, err := w.rpcClient.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain id %s", err)
		return nil, err
	}

	sign, err := w.signer.CreateSign(chainId)
	if err != nil {
		log.Fatalf("Failed to sign transaction %s", err)
		return nil, err
	}

	tx, err := w.c.NewContract(
		sign,
		rate,
		clientId,
	)

	if err != nil {
		log.Fatalf("failed to create contract %s", err)
		return nil, err
	}

	receipt, err := bind.WaitMined(context.Background(), w.rpcClient, tx)
	if err != nil {
		log.Fatalf("failed to wait for transaction to be mined %s", err)
		return nil, err
	}

	if receipt.Status != 1 {
		return nil, errors.New("transaction failed")
	}

	parsedLog, err := w.c.PayAsYouGoFilterer.ParseContractCreated(*receipt.Logs[0])
	if err != nil {
		return nil, err
	}

	contractId := fmt.Sprintf("%x", parsedLog.ContractId[:])

	return &contractId, nil
}

func (w *PayAsYouGoWrapperImpl) GetContractByID(ctx context.Context, contractID string) (PayAsYouGoAgreement, error) {
	chainId, err := w.rpcClient.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain id %s", err)
		return PayAsYouGoAgreement{}, err
	}

	sign, err := w.signer.CreateSign(chainId)
	if err != nil {
		log.Fatalf("Failed to sign transaction %s", err)
		return PayAsYouGoAgreement{}, err
	}

	decodedContractId, err := hex.DecodeString(contractID)
	if err != nil {
		log.Fatalf("Failed to decode contract id %s", err)
		return PayAsYouGoAgreement{}, err
	}

	contract, err := w.c.GetContractById(&bind.CallOpts{
		From: sign.From,
	}, [32]byte(decodedContractId))
	if err != nil {
		log.Fatalf("Failed to get contract %s", err)
		return PayAsYouGoAgreement{}, err
	}

	return contract, nil
}
