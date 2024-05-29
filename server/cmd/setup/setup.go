package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/repository"
	"io"
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

	err := CopyFile(".env.dev", ".env")
	if err != nil {
		logrus.Fatal("failed to copy .dev.env file :", err)
	}

	logrus.Debug("the configuration was copied successfully.")

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("failed to read configuration: ", err)
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
	cmd := exec.Command("docker", "compose", "version")
	_, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Fatal("Docker Compose is not installed. Please install Docker Compose before running this script.")
	}

	// Check if Docker is running
	cmd = exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		logrus.Fatal("Docker is not running. Please start Docker before running this script.")
	}
}

func SetupPG() {
	err := RunDockerCompose("docker-compose-pg.yml")
	if err != nil {
		logrus.Fatal("failed to run postgresql database container", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
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
		"file://../migrations",
		dsn,
	)

	if err != nil {
		logrus.Fatal("failed to create postgresql migration instance: ", err)
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Fatal("failed to latest schema to postgresql instance: ", err)
	}

	config.InitDB()
	resourceRepository := repository.NewRepository(config.DB)

	if os.Getenv("CREATE_TEST_USER") == "true" {
		logrus.Debug("creating test user...")

		resourceRepository.RegisterUser(
			dto.RegisterUserDTO{
				Email:       os.Getenv("TEST_USER_EMAIL"),
				Password:    os.Getenv("TEST_USER_PASSWORD"),
				Username:    os.Getenv("TEST_USER_USERNAME"),
				FirstName:   os.Getenv("TEST_USER_FIRST_NAME"),
				LastName:    os.Getenv("TEST_USER_LAST_NAME"),
				DisplayName: os.Getenv("TEST_USER_DISPLAY_NAME"),
				Country:     os.Getenv("TEST_USER_COUNTRY"),
			},
		)

		logrus.Debug("test user created successfully.")
	}

	if os.Getenv("CREATE_TEST_CLIENT") == "true" {
		logrus.Debug("creating test client...")

		resourceRepository.RegisterClient(
			dto.RegisterClientDTO{
				Password:    os.Getenv("TEST_CLIENT_PASSWORD"),
				Username:    os.Getenv("TEST_CLIENT_USERNAME"),
				RedirectURI: os.Getenv("TEST_CLIENT_REDIRECT_URI"),
				BackendURI:  os.Getenv("TEST_CLIENT_BACKEND_URI"),
			},
		)

		logrus.Debug("test client created successfully.")
	}
}

func SetupInfluxDB() {
	err := RunDockerCompose("docker-compose-influxdb.yml")
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
	err := RunDockerCompose("docker-compose-telegraf.yml")
	if err != nil {
		logrus.Fatal("failed to run telegraf as sidecar container to collect metrics from opentelemetry - ", err)
	}
}

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func RunDockerCompose(dockerComposeFile string) error {
	cmd := exec.Command("docker", "compose", "-f", GetFullInfraPath(dockerComposeFile), "--env-file", ".env", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func GetFullInfraPath(fileName string) string {
	return fmt.Sprintf("../infra/local/" + fileName)
}
