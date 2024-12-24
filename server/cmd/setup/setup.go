package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/zk"
	"globe-and-citizen/layer8/server/resource_server/repository"
	"globe-and-citizen/layer8/server/resource_server/service"
	"globe-and-citizen/layer8/server/resource_server/utils"
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
	ValidateDockerUpAndRunning()

	if err := CopyFile(".env.dev", ".env"); err != nil {
		logrus.Fatal("failed to copy .dev.env file :", err)
	}

	if err := AppendFileContent(".env.secret", ".env"); err != nil {
		logrus.Fatal("failed to append .secret.env file :", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("failed to read configuration: ", err)
	}

	if err := RunDockerCompose(); err != nil {
		logrus.Fatal("failed to run all necessary docker containers", err)
	}

	SetupPG()
	SetupInfluxDB()
}

func ValidateDockerUpAndRunning() {
	if _, err := exec.LookPath("docker"); err != nil {
		logrus.Fatal("Docker is not installed. Please install Docker before running this script.")
	}

	cmd := exec.Command("docker", "compose", "version")
	_, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Fatal("Docker Compose is not installed. Please install Docker Compose before running this script.")
	}

	cmd = exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		logrus.Fatal("Docker is not running. Please start Docker before running this script.")
	}
}

func SetupPG() {
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

		break
	}

	migration, err := migrate.New(
		"file://../migrations",
		dsn,
	)

	if err != nil {
		logrus.Fatal("failed to create postgresql migration instance: ", err)
	}

	_, err = db.Exec(`
		UPDATE schema_migrations
			 SET version = CASE 
                WHEN dirty = true THEN version - 1
                ELSE version
              END,
			  dirty = false;
	`)
	if err != nil {
		logrus.Fatal("failed to update schema_migrations table: ", err)
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Fatal("failed to latest schema to postgresql instance: ", err)
	}

	config.InitDB()

	resourceService := service.NewService(
		repository.NewRepository(config.DB),
		&verification.EmailVerifier{},
		&zk.ProofProcessor{},
	)

	if os.Getenv("CREATE_TEST_USER") == "true" {
		resourceService.RegisterUser(
			dto.RegisterUserDTO{
				Password:    os.Getenv("TEST_USER_PASSWORD"),
				Username:    os.Getenv("TEST_USER_USERNAME"),
				FirstName:   os.Getenv("TEST_USER_FIRST_NAME"),
				LastName:    os.Getenv("TEST_USER_LAST_NAME"),
				DisplayName: os.Getenv("TEST_USER_DISPLAY_NAME"),
				Country:     os.Getenv("TEST_USER_COUNTRY"),
				PublicKey:   make([]byte, 33),
			},
		)
	}

	if os.Getenv("CREATE_TEST_CLIENT") == "true" {
		resourceService.RegisterClient(
			dto.RegisterClientDTO{
				Name:        os.Getenv("TEST_CLIENT_NAME"),
				Password:    os.Getenv("TEST_CLIENT_PASSWORD"),
				Username:    os.Getenv("TEST_CLIENT_USERNAME"),
				RedirectURI: os.Getenv("TEST_CLIENT_REDIRECT_URI"),
				BackendURI:  utils.RemoveProtocolFromURL(os.Getenv("TEST_CLIENT_BACKEND_URI")),
			},
		)
	}
}

func SetupInfluxDB() {
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

func AppendFileContent(sourceFile, destFile string) error {
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		return nil
	}

	source, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", sourceFile, err)
	}
	defer source.Close()

	dest, err := os.OpenFile(destFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", destFile, err)
	}

	defer dest.Close()

	writer := bufio.NewWriter(dest)

	if _, err := writer.WriteString("\n\n"); err != nil {
		return fmt.Errorf("error writing new lines: %v", err)
	}

	reader := bufio.NewReader(source)
	if _, err := io.Copy(writer, reader); err != nil {
		return fmt.Errorf("error copying content: %v", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %v", err)
	}

	return nil
}

func RunDockerCompose() error {
	cmd := exec.Command("docker", "compose", "-f", GetFullInfraPath("docker-compose.yml"), "--env-file", ".env", "up", "-d")
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
