package service

import (
	"context"
	"globe-and-citizen/layer8/server/blockchain"
	"globe-and-citizen/layer8/server/resource_server/repository"
	"log"
	"time"
)

type ClientService interface {
	RefreshClientBill(ctx context.Context) error
}

type clientService struct {
	statRepo          repository.StatRepository
	blockchainWrapper blockchain.PayAsYouGoWrapper
	getNow            func() time.Time
}

func NewClientService(
	statRepo repository.StatRepository,
	blockchainWrapper blockchain.PayAsYouGoWrapper,
	nowFunc func() time.Time,
) ClientService {
	return &clientService{
		statRepo:          statRepo,
		blockchainWrapper: blockchainWrapper,
		getNow:            nowFunc,
	}
}

func (c *clientService) RefreshClientBill(ctx context.Context) error {
	log.Println("Refreshing client bill")
	clients, err := c.blockchainWrapper.GetContracts(ctx)
	if err != nil {
		return err
	}

	usages := make([]blockchain.PayAsYouGoBillingInput, 0, len(clients))
	for _, client := range clients {
		latestFetchTime := time.Unix(int64(client.LastUsageFetchTime), 0)
		currentTime := c.getNow()

		clientUsageInBytes, err := c.statRepo.GetTotalByDateRangeByClient(
			ctx,
			latestFetchTime,
			currentTime,
			client.ClientId,
		)

		if err != nil {
			return err
		}

		if clientUsageInBytes == 0 {
			continue
		}

		usages = append(usages, blockchain.PayAsYouGoBillingInput{
			ContractId: client.ContractId,
			Amount:     uint64(clientUsageInBytes / 1000000),
			Timestamp:  uint64(currentTime.Unix()),
		})
	}

	if len(usages) == 0 {
		log.Println("No client usage to bill")
		return nil
	}

	if err := c.blockchainWrapper.BulkAddBillToContract(ctx, usages); err != nil {
		return err
	}

	log.Println("Client bill refreshed")
	return nil
}
