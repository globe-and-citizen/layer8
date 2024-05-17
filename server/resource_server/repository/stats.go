package repository

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/models"
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

func (s *StatRepository) GetTotalRequestsInLastXDaysByClient(ctx context.Context, days int, clientID string) (models.Statistics, error) {
	result := make([]models.UsageStatisticPerDate, 0)

	queryAPI := s.influxDBClient.QueryAPI("layer8")

	query := fmt.Sprintf(`from(bucket: "layer8")
	|> range(start: -%dd)
	|> filter(fn: (r) => r["_measurement"] == "total_byte_transferred")
	|> filter(fn: (r) => r["_field"] == "counter")
	|> filter(fn: (r) => r["client_id"] == "%s")
	|> group(columns: ["client_id"])
	|> aggregateWindow(every: 1d, fn: sum, createEmpty: true)
	|> yield(name: "sum")`, days, clientID)

	rawDataFromInflux, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return models.Statistics{}, err
	}

	var totalRequest float64
	for rawDataFromInflux.Next() {
		rawDataPointer := rawDataFromInflux.Record()
		unparsedTotal := rawDataPointer.ValueByKey("_value")
		decimalValueTotal, err := strconv.ParseFloat(fmt.Sprint(unparsedTotal), 64)
		if err != nil {
			decimalValueTotal = 0
		}

		var totalForThisPeriod float64
		if decimalValueTotal > 0 {
			totalRequest += decimalValueTotal / 1000000000
			totalForThisPeriod = decimalValueTotal / 1000000000
		}

		at := rawDataPointer.ValueByKey("_time").(time.Time)
		result = append(result, models.UsageStatisticPerDate{
			Date:  at.Format("Mon, 02 Jan 2006"),
			Total: totalForThisPeriod,
		})
	}

	var averageRequest float64
	if totalRequest > 0 {
		averageRequest = totalRequest / float64(len(result))
	}

	return models.Statistics{
		Total:            totalRequest,
		Average:          averageRequest,
		StatisticDetails: result,
	}, nil
}

func (s *StatRepository) GetTotalByDateRangeByClient(ctx context.Context, start time.Time, end time.Time, clientID string) (float64, error) {
	queryAPI := s.influxDBClient.QueryAPI("layer8")

	query := fmt.Sprintf(`
	from(bucket: "layer8")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r["_measurement"] == "total_byte_transferred")
	|> filter(fn: (r) => r["client_id"] == "%s")
	|> filter(fn: (r) => r["_field"] == "counter")
	|> group(columns: ["client_id"])
	|> sum()`, start.Format(time.RFC3339), end.Format(time.RFC3339), clientID)

	rawDataFromInflux, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return 0, err
	}

	var decimalValueTotal float64
	for rawDataFromInflux.Next() {
		rawDataPointer := rawDataFromInflux.Record()
		unparsedTotal := rawDataPointer.ValueByKey("_value")
		decimalValueTotal, err = strconv.ParseFloat(fmt.Sprint(unparsedTotal), 64)
		if err != nil {
			decimalValueTotal = 0
		}
	}

	return decimalValueTotal, err
}
