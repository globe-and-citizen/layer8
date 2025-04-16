package paywithcrypto

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"globe-and-citizen/layer8/server/resource_server/interfaces"
	"log"
	"math/big"
	"os"
	"strings"
)

type EventListener struct {
	repository interfaces.IRepository
}

func NewEventListener(repository interfaces.IRepository) *EventListener {
	return &EventListener{
		repository: repository,
	}
}

type TrafficPaidEvent struct {
	ClientId string
	From     common.Address
	Amount   *big.Int
}

const trafficPaidEventABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TrafficPaid\",\"type\":\"event\"}]"

func (l *EventListener) Start() {
	contractABI, err := abi.JSON(strings.NewReader(trafficPaidEventABI))
	if err != nil {
		log.Fatalf("failed to parse the smart contract abi: %e", err)
	}

	smartContractAddressStr := os.Getenv("SMART_CONTRACT_ADDRESS")
	if smartContractAddressStr == "" {
		log.Fatalf("error: SMART_CONTRACT_ADDRESS value is empty in .env")
	}
	smartContractAddress := common.HexToAddress(smartContractAddressStr)

	websocketNodeURL := os.Getenv("WEBSOCKET_NODE_URL")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{smartContractAddress},
	}

	for {
		err := l.readEvents(websocketNodeURL, contractABI, query)
		if err != nil {
			log.Println(err)
		}
	}
}

func (l *EventListener) readEvents(
	websocketNodeURL string,
	contractABI abi.ABI,
	query ethereum.FilterQuery,
) error {
	client, err := ethclient.Dial(websocketNodeURL)
	if err != nil {
		log.Fatalf("failed to create an ethereum client: %e", err)
	}

	logChannel := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logChannel)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("failed to read blockchain event log: %e", err)
		case eventLog := <-logChannel:
			var event TrafficPaidEvent

			err = contractABI.UnpackIntoInterface(&event, "TrafficPaid", eventLog.Data)
			if err != nil {
				log.Fatalf("failed to decode a traffic paid event: %e", err)
			}

			err = l.handleTrafficPaidEvent(&event)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (l *EventListener) handleTrafficPaidEvent(event *TrafficPaidEvent) error {
	return l.repository.PayClientTrafficUsage(event.ClientId, int(event.Amount.Int64()))
}
