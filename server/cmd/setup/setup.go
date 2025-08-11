package main

import (
	"bufio"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/xdg-go/pbkdf2"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/emails/verification"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/zk"
	"globe-and-citizen/layer8/server/resource_server/repository"
	"globe-and-citizen/layer8/server/resource_server/service"
	"globe-and-citizen/layer8/server/resource_server/utils"
)

func main() {
	if !(len(os.Args) > 1 && os.Args[1] == "docker") {
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
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	logrus.Println("SetupPG pg dsn: ", dsn)
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

	logrus.Println("migration path: ", fmt.Sprintf("file://%s", os.Getenv("DB_SETUP_MIGRATION_PATH")))
	migration, err := migrate.New(
		fmt.Sprintf("file://%s", os.Getenv("DB_SETUP_MIGRATION_PATH")),
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

	iterCount, err := strconv.Atoi(os.Getenv("SCRAM_ITERATION_COUNT"))
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("CREATE_TEST_USER") == "true" {
		salt, err := resourceService.RegisterUserPrecheck(
			dto.RegisterUserPrecheckDTO{Username: os.Getenv("TEST_USER_USERNAME")},
			iterCount,
		)
		if err != nil {
			// log.Fatal(err)
			fmt.Println(err)
			return
		}

		storedKey, serverKey := computeHmacKeys(os.Getenv("TEST_USER_PASSWORD"), salt, iterCount)

		err = resourceService.RegisterUser(
			dto.RegisterUserDTO{
				Username:    os.Getenv("TEST_USER_USERNAME"),
				FirstName:   os.Getenv("TEST_USER_FIRST_NAME"),
				LastName:    os.Getenv("TEST_USER_LAST_NAME"),
				DisplayName: os.Getenv("TEST_USER_DISPLAY_NAME"),
				Country:     os.Getenv("TEST_USER_COUNTRY"),
				PublicKey:   make([]byte, 33),
				StoredKey:   storedKey,
				ServerKey:   serverKey,
			},
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	if os.Getenv("CREATE_TEST_CLIENT") == "true" {
		salt, err := resourceService.RegisterClientPrecheck(
			dto.RegisterClientPrecheckDTO{
				Username: os.Getenv("TEST_CLIENT_USERNAME"),
			},
			iterCount,
		)
		if err != nil {
			log.Fatal(err)
		}

		storedKey, serverKey := computeHmacKeys(os.Getenv("TEST_CLIENT_PASSWORD"), salt, iterCount)

		resourceService.RegisterClient(
			dto.RegisterClientDTO{
				Name:        os.Getenv("TEST_CLIENT_NAME"),
				Username:    os.Getenv("TEST_CLIENT_USERNAME"),
				RedirectURI: os.Getenv("TEST_CLIENT_REDIRECT_URI"),
				BackendURI:  utils.RemoveProtocolFromURL(os.Getenv("TEST_CLIENT_BACKEND_URI")),
				StoredKey:   storedKey,
				ServerKey:   serverKey,
			},
		)
	}
}

func computeHmacKeys(password string, salt string, iterCount int) (storedKey string, serverKey string) {
	saltBytes, _ := hex.DecodeString(salt)

	hashedPassword := pbkdf2.Key(
		[]byte(password),
		saltBytes,
		iterCount, 20, sha1.New,
	)

	clientKey := computeHmac256(hashedPassword, "Client Key")
	serverKeyBytes := computeHmac256(hashedPassword, "Server Key")
	storedKeyBytes := sha256.Sum256(clientKey)

	return hex.EncodeToString(storedKeyBytes[:]), hex.EncodeToString(serverKeyBytes)
}

func computeHmac256(input []byte, key string) []byte {
	hmacInstance := hmac.New(sha256.New, []byte(key))
	hmacInstance.Write(input)
	return hmacInstance.Sum(make([]byte, 0))
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