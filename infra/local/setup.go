package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Debug("validating that Docker is up and running...")
	ValidateDockerUpAndRunning()
	logrus.Debug("docker validated")

	logrus.Debug("copy default configuration from root path to .env file...")

	cmd := exec.Command("cp", "../../.env.dev", ".env")
	err := cmd.Run()
	if err != nil {
		logrus.Fatal("failed to copy .dev.env file :", err)
	}

	logrus.Debug("the configuration was copied successfully.")

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("failed to read configuration")
	}

	logrus.Debug("setting up the PostgreSQL container...")
	SetupPG()
	logrus.Debug("the setup of PostgreSQL has been completed.")

	logrus.Debug("setting up the InfluxDB container...")
	SetupInfluxDB()
	logrus.Debug("the setup of InfluxDB has been completed.")

	logrus.Debug("setting up the Telegraf container...")
	SetupTelegraf()
	logrus.Debug("the setup of Telegraf has been completed.")
}

func ValidateDockerUpAndRunning() {
	// Check if Docker is installed
	if _, err := exec.LookPath("docker"); err != nil {
		logrus.Fatal("Docker is not installed. Please install Docker before running this script.")
	}

	// Check if Docker Compose is installed
	if _, err := exec.LookPath("docker-compose"); err != nil {
		logrus.Fatal("Docker Compose is not installed. Please install Docker Compose before running this script.")
	}

	// Check if Docker is running
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		logrus.Fatal("Docker is not running. Please start Docker before running this script.")
	}
}

func SetupPG() {
	cmd := exec.Command("docker-compose", "-f", "docker-compose-pg.yml", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		logrus.Fatal("failed to run postgresql database container", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.Fatal("failed to create postgresql connection instance")
	}

	defer db.Close()
	for {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := db.PingContext(ctx)
		if err != nil {
			logrus.Warn("failed to ping postgresql container, retrying...")
			time.Sleep(1 * time.Second)
			continue
		}

		logrus.Debug("postgresql verified to be up and running, migrating schemas...")
		break
	}

	migration, err := migrate.New(
		"file://../../migrations",
		dsn,
	)

	if err != nil {
		logrus.Fatal("failed to create postgresql migration instance")
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Fatal("failed to latest schema to postgresql instance: ", err)
	}
}

func SetupInfluxDB() {
	cmd := exec.Command("docker-compose", "-f", "docker-compose-influxdb.yml", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		logrus.Fatal("failed to run influxdb database container", err)
	}

	client := influxdb2.NewClient(os.Getenv("INFLUXDB_URL"), "")
	defer client.Close()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		pingResult, err := client.Ping(ctx)
		if err != nil {
			logrus.Warn("failed to ping influxdb2 container, retrying...")
			time.Sleep(1 * time.Second)
		}

		if pingResult {
			logrus.Debug("influxdb2 verified to be up and running, configuring credentials...")
			break
		}
	}

	if _, err := client.SetupWithToken(
		context.Background(),
		os.Getenv("INFLUXDB_USERNAME"),
		os.Getenv("INFLUXDB_PASSWORD"),
		os.Getenv("INFLUXDB_ORG"),
		os.Getenv("INFLUXDB_BUCKET"),
		0,
		os.Getenv("INFLUXDB_TOKEN"),
	); err != nil {
		logrus.Warn("failed to setup the layer8 token, ignore this if you have already set up the token before - ", err)
	}
}

func SetupTelegraf() {
	cmd := exec.Command("docker-compose", "-f", "docker-compose-telegraf.yml", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		logrus.Fatal("failed to run telegraf as sidecar container to collect metrics from opentelemetry - ", err)
	}
}
