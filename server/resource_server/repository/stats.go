package repository

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/models"
	"log"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type StatRepository struct {
	influxDBClient influxdb2.Client
}

func NewStatRepository(
	influxDBClient influxdb2.Client,
) *StatRepository {
	return &StatRepository{
		influxDBClient: influxDBClient,
	}
}

func (s *StatRepository) GetTotalRequestsInLastXDays(ctx context.Context, days int) ([]models.UsageStatistic, error) {
	result := make([]models.UsageStatistic, 0)

	queryAPI := s.influxDBClient.QueryAPI("layer8")

	now := time.Now().UTC()
	xDaysAgo := now.AddDate(0, 0, -days)

	query := fmt.Sprintf(`from(bucket: "layer8")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r["_measurement"] == "total_request")
	|> filter(fn: (r) => r["_field"] == "counter")
	|> aggregateWindow(every: 1d, fn: sum, createEmpty: true)
	|> yield(name: "sum")`,
		xDaysAgo.Format(time.RFC3339),
		now.Format(time.RFC3339),
	)

	log.Println(query)

	rawDataFromInflux, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return result, err
	}

	for rawDataFromInflux.Next() {
		rawDataPointer := rawDataFromInflux.Record()
		unparsedTotal := rawDataPointer.ValueByKey("_value")
		intValueTotal, err := strconv.ParseInt(fmt.Sprintf("%v", unparsedTotal), 10, 64)
		if err != nil {
			intValueTotal = 0
		}

		at := rawDataPointer.ValueByKey("_time").(time.Time)

		result = append(result, models.UsageStatistic{
			Date:  at.Format("Mon, 02 Jan 2006"),
			Total: intValueTotal,
		})
	}

	return result, nil
}
