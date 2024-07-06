package blockchain

import (
	"context"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewBlockchainClient() (*ethclient.Client, error) {
	client, err := ethclient.DialContext(context.Background(), os.Getenv("BLOCKCHAIN_RPC_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	return client, nil

}
