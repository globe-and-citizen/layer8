package blockchain

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewSepholiaClient() (*ethclient.Client, error) {
	client, err := ethclient.DialContext(context.Background(), "https://rpc.sepolia.org")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	return client, nil

}
