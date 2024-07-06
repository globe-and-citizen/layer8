package blockchain

import (
	"context"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PayAsYouGoWrapper interface {
	StoreClient(ctx context.Context, rate uint64, clientId string) error
	GetClients(ctx context.Context) ([]PayAsYouGoClient, error)
	GetClientByID(ctx context.Context, clientID string) (PayAsYouGoClient, error)
	BulkAddBillToClient(ctx context.Context, billings []PayAsYouGoBillingInput) error
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

func (w *PayAsYouGoWrapperImpl) StoreClient(ctx context.Context, rate uint64, clientId string) error {
	sign, err := w.generateSign(ctx)
	if err != nil {
		log.Printf("Failed to sign transaction %s", err)
		return err
	}

	tx, err := w.c.NewClient(
		sign,
		rate,
		clientId,
	)

	if err != nil {
		log.Printf("failed to create contract %s", err)
		return err
	}

	receipt, err := bind.WaitMined(context.Background(), w.rpcClient, tx)
	if err != nil {
		log.Printf("failed to wait for transaction to be mined %s", err)
		return err
	}

	if receipt.Status != 1 {
		return errors.New("transaction failed")
	}

	return nil
}

func (w *PayAsYouGoWrapperImpl) BulkAddBillToClient(
	ctx context.Context,
	billings []PayAsYouGoBillingInput,
) error {
	sign, err := w.generateSign(ctx)
	if err != nil {
		log.Printf("Failed to sign transaction %s", err)
		return err
	}

	_, err = w.c.BulkAddBillToClient(sign, billings)
	if err != nil {
		log.Printf("failed to bulk add bill to contract %s", err)
		return err
	}

	return nil
}

func (w *PayAsYouGoWrapperImpl) GetClients(
	ctx context.Context,
) ([]PayAsYouGoClient, error) {
	sign, err := w.generateSign(ctx)
	if err != nil {
		log.Printf("Failed to sign transaction %s", err)
		return nil, err
	}

	contracts, err := w.c.GetClients(&bind.CallOpts{
		From: sign.From,
	})
	if err != nil {
		log.Printf("Failed to get contracts %s", err)
		return nil, err
	}

	return contracts, nil
}

func (w *PayAsYouGoWrapperImpl) GetClientByID(ctx context.Context, clientID string) (PayAsYouGoClient, error) {
	chainId, err := w.rpcClient.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain id %s", err)
		return PayAsYouGoClient{}, err
	}

	sign, err := w.signer.CreateSign(chainId)
	if err != nil {
		log.Printf("Failed to sign transaction %s", err)
		return PayAsYouGoClient{}, err
	}

	contract, err := w.c.GetClientById(&bind.CallOpts{
		From: sign.From,
	}, clientID)
	if err != nil {
		log.Printf("Failed to get contract %s", err)
		return PayAsYouGoClient{}, err
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
