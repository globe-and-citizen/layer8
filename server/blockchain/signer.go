package blockchain

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
)

type TransactionSigner interface {
	CreateSign(chainId *big.Int) (*bind.TransactOpts, error)
}

type TransactionSignerImpl struct {
	privateKey *ecdsa.PrivateKey
	maxGasFee  uint64
}

func NewTransactionSigner(privateKeyStr string) (TransactionSigner, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	maxGasFee, err := strconv.Atoi(os.Getenv("BLOCKCHAIN_MAX_GAS_FEE"))
	if err != nil {
		log.Fatalf("Failed to load max gas fee: %v", err)
	}

	return &TransactionSignerImpl{
		privateKey: privateKey,
		maxGasFee:  uint64(maxGasFee),
	}, nil
}

func (ts *TransactionSignerImpl) CreateSign(chainId *big.Int) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(ts.privateKey, chainId)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
		return nil, err
	}

	auth.GasLimit = ts.maxGasFee

	return auth, nil

}
