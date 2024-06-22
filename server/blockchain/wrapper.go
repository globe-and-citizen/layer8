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
	GetContracts(ctx context.Context) ([]PayAsYouGoAgreement, error)
	GetContractByID(ctx context.Context, contractID string) (PayAsYouGoAgreement, error)
	BulkAddBillToContract(ctx context.Context, billings []PayAsYouGoBillingInput) error
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
	sign, err := w.generateSign(ctx)
	if err != nil {
		log.Printf("Failed to sign transaction %s", err)
		return nil, err
	}

	tx, err := w.c.NewContract(
		sign,
		rate,
		clientId,
	)

	if err != nil {
		log.Printf("failed to create contract %s", err)
		return nil, err
	}

	receipt, err := bind.WaitMined(context.Background(), w.rpcClient, tx)
	if err != nil {
		log.Printf("failed to wait for transaction to be mined %s", err)
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

func (w *PayAsYouGoWrapperImpl) BulkAddBillToContract(
	ctx context.Context,
	billings []PayAsYouGoBillingInput,
) error {
	sign, err := w.generateSign(ctx)
	if err != nil {
		log.Printf("Failed to sign transaction %s", err)
		return err
	}

	_, err = w.c.BulkAddBillToContract(sign, billings)
	if err != nil {
		log.Printf("failed to bulk add bill to contract %s", err)
		return err
	}

	return nil
}

func (w *PayAsYouGoWrapperImpl) GetContracts(
	ctx context.Context,
) ([]PayAsYouGoAgreement, error) {
	sign, err := w.generateSign(ctx)
	if err != nil {
		log.Printf("Failed to sign transaction %s", err)
		return nil, err
	}

	contracts, err := w.c.GetContracts(&bind.CallOpts{
		From: sign.From,
	})
	if err != nil {
		log.Printf("Failed to get contracts %s", err)
		return nil, err
	}

	return contracts, nil
}

func (w *PayAsYouGoWrapperImpl) GetContractByID(ctx context.Context, contractID string) (PayAsYouGoAgreement, error) {
	chainId, err := w.rpcClient.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain id %s", err)
		return PayAsYouGoAgreement{}, err
	}

	sign, err := w.signer.CreateSign(chainId)
	if err != nil {
		log.Printf("Failed to sign transaction %s", err)
		return PayAsYouGoAgreement{}, err
	}

	decodedContractId, err := hex.DecodeString(contractID)
	if err != nil {
		log.Printf("Failed to decode contract id %s", err)
		return PayAsYouGoAgreement{}, err
	}

	contract, err := w.c.GetContractById(&bind.CallOpts{
		From: sign.From,
	}, [32]byte(decodedContractId))
	if err != nil {
		log.Printf("Failed to get contract %s", err)
		return PayAsYouGoAgreement{}, err
	}

	return contract, nil
}

func (w *PayAsYouGoWrapperImpl) generateSign(ctx context.Context) (*bind.TransactOpts, error) {
	chainId, err := w.rpcClient.ChainID(ctx)
	if err != nil {
		log.Printf("Failed to get chain id %s", err)
		return nil, err
	}

	sign, err := w.signer.CreateSign(chainId)
	if err != nil {
		log.Printf("Failed to sign transaction %s", err)
		return nil, err
	}

	return sign, nil
}
