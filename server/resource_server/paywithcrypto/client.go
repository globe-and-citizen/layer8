package paywithcrypto

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"globe-and-citizen/layer8/server/resource_server/paywithcrypto/rpcclient"
	"log"
	"math/big"
	"os"
	"strconv"
)

type Client struct {
	smartContractRpcClient *rpcclient.PayWithCryptoSmartContractClient
	ethClient              *ethclient.Client

	privateKey     *ecdsa.PrivateKey
	accountAddress common.Address
	gasLimit       uint64
}

func NewClient() *Client {
	ethClient, err := ethclient.Dial(os.Getenv("CHAIN_RPC_URL"))
	if err != nil {
		log.Fatalf("failed to set up the pay with crypto client: %e", err)
	}

	smartContractAddressStr := os.Getenv("SMART_CONTRACT_ADDRESS")
	if smartContractAddressStr == "" {
		log.Fatalf("error: SMART_CONTRACT_ADDRESS value is empty in .env")
	}
	smartContractAddress := common.HexToAddress(smartContractAddressStr)

	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("failed to parse the ecdsa PRIVATE_KEY from .env")
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("failed to cast public key to ECDSA format")
	}

	accountAddress := crypto.PubkeyToAddress(*publicKey)

	smartContractRpcClient, err := rpcclient.NewPayWithCryptoSmartContractClient(smartContractAddress, ethClient)
	if err != nil {
		log.Fatal(err)
	}

	gasLimit, err := strconv.ParseUint(os.Getenv("GAS_LIMIT"), 10, 64)
	if err != nil {
		log.Fatalf("failed to parse GAS_LIMIT variable from .env: %e", err)
	}

	return &Client{
		smartContractRpcClient: smartContractRpcClient,
		ethClient:              ethClient,
		accountAddress:         accountAddress,
		privateKey:             privateKey,
		gasLimit:               gasLimit,
	}
}

func (c *Client) AddTrafficUsage(clientId string, trafficInBytes int) error {
	// 16 bytes + 32 bytes = 48 bytes
	// 46 kb = 46000 bytes => for around 1000 clients and more this wouldn't work
	auth, err := c.getTransactOpts()
	if err != nil {
		return err
	}

	tx, err := c.smartContractRpcClient.AddTrafficUsage(auth, clientId, big.NewInt(int64(trafficInBytes)))
	if err != nil {
		return err
	}

	fmt.Println(tx.Hash().Hex())
	//  TODO: how to handle errors?

	return nil
}

func (c *Client) AddClient(clientId string, rate int) error {
	auth, err := c.getTransactOpts()
	if err != nil {
		return err
	}

	tx, err := c.smartContractRpcClient.AddClient(auth, clientId, big.NewInt(int64(rate)))
	if err != nil {
		return err
	}

	fmt.Println(tx.Hash().Hex())

	return nil
}

func (c *Client) SetClientRate(clientId string, newRate int) error {
	auth, err := c.getTransactOpts()
	if err != nil {
		return err
	}

	tx, err := c.smartContractRpcClient.SetClientRate(auth, clientId, big.NewInt(int64(newRate)))
	if err != nil {
		return err
	}

	fmt.Println(tx.Hash().Hex())

	return nil
}

func (c *Client) getTransactOpts() (*bind.TransactOpts, error) {
	nonce, err := c.ethClient.PendingNonceAt(context.Background(), c.accountAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := c.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	chainId, err := c.ethClient.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainId)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.Value = big.NewInt(0)
	auth.GasLimit = c.gasLimit

	return auth, nil
}
