package paywithcrypto

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/repository"
	"time"
)

type StatisticsUpdater struct {
	statRepository *repository.StatRepository
	repository     interfaces.IRepository
}

func NewStatisticsUpdater(
	statRepository *repository.StatRepository,
	repository interfaces.IRepository,
) *StatisticsUpdater {
	return &StatisticsUpdater{
		statRepository: statRepository,
		repository:     repository,
	}
}

func (s *StatisticsUpdater) Update(ctx context.Context, now time.Time) error {
	allClientStatistics, err := s.repository.GetAllClientStatistics()
	if err != nil {
		return err
	}

	for _, clientStat := range allClientStatistics {
		fmt.Printf("client %s; statistics from: %s; to: %s\n",
			clientStat.ClientId, clientStat.LastTrafficUpdateTimestamp.String(), now.String(),
		)

		consumedBytesFloat, err := s.statRepository.GetTotalByDateRangeByClient(
			ctx, clientStat.LastTrafficUpdateTimestamp, now, clientStat.ClientId,
		)
		if err != nil {
			return fmt.Errorf("failed to get traffic updates for client %s: %e", clientStat.ClientId, err)
		}

		if consumedBytesFloat == 0 {
			continue
		}

		consumedBytes := int(consumedBytesFloat)

		err = s.repository.AddClientTrafficUsage(clientStat.ClientId, consumedBytes, now)
		if err != nil {
			return err
		}
	}

	return nil
}
