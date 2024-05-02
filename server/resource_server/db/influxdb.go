package db

import (
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var influxDBClient influxdb2.Client

func InitInfluxDBClient() {
	influxDBClient = influxdb2.NewClient(
		os.Getenv("INFLUXDB_URL"),
		os.Getenv("INFLUXDB_TOKEN"),
	)
}

func GetInfluxDBClient() influxdb2.Client {
	return influxDBClient
}
