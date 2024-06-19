package main

import (
	"context"
	"globe-and-citizen/layer8/server/blockchain"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	rpcClient, err := blockchain.NewSepholiaClient()
	if err != nil {
		log.Fatalf("Failed to create Ethereum client %s", err)
	}

	signer, err := blockchain.NewTransactionSigner(os.Getenv("ETH_PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("Failed to create transaction signer %s", err)
	}

	wrapper := blockchain.NewPayAsYouGoWrapper(
		os.Getenv("PAYASYOUGO_CONTRACT_ADDRESS"),
		rpcClient,
		signer,
	)

	c, err := wrapper.GetContractByID(context.Background(), "352023bed1da998bfdbe5607e32683a5898659c7705e1516dd7e82bdb5fbfe73")
	if err != nil {
		log.Fatalf("Failed to create contract %s", err)
	}

	log.Println("contract : ", c.Rate)

}
