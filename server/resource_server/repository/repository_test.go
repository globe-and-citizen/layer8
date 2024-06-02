package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

const id uint = 1
const userId uint = 1
const verificationCode = "12345"
const emailProof = "AbcdfTs"

var timestamp = time.Date(2024, time.May, 24, 14, 0, 0, 0, time.UTC)

var mockDB *sql.DB
var mock sqlmock.Sqlmock
var err error
var db *gorm.DB
var repository interfaces.IRepository

func SetUp(t *testing.T) {
	mockDB, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}

	db, err = gorm.Open(
		postgres.New(
			postgres.Config{
				Conn: mockDB,
			},
		),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	repository = NewRepository(db)
}

func TestRegisterUser(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the user repository with the mock database connection
	repo := NewRepository(db)

	// Define a test user DTO
	testUser := dto.RegisterUserDTO{
		Email:       "test@gcitizen.com",
		Username:    "test_user",
		FirstName:   "Test",
		LastName:    "User",
		Password:    "TestPass123",
		Country:     "Unknown",
		DisplayName: "user",
	}

	// Call the RegisterUser function
	repo.RegisterUser(testUser)
}

func TestRegisterClient(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the client repository with the mock database connection
	repo := NewRepository(db)

	// Define a test client DTO
	testClient := dto.RegisterClientDTO{
		Name:        "testclient",
		RedirectURI: "https://gcitizen.com/callback",
	}

	// Call the RegisterClient function
	repo.RegisterClient(testClient)
}

func TestGetClientData(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the client repository with the mock database connection
	repo := NewRepository(db)

	// Define a test client name
	testClientName := "testclient"

	// Expect a query to be executed and return a row
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "clients" WHERE name = $1 ORDER BY "clients"."id" LIMIT 1`)).WithArgs(testClientName).WillReturnRows(sqlmock.NewRows([]string{"id", "secret", "name", "redirect_uri"}).AddRow("notanid", "testsecret", "testclient", "https://gcitizen.com/callback"))

	// Call the GetClientData function
	repo.GetClientData(testClientName)

	// Check if the expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestLoginPreCheckUser(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the user repository with the mock database connection
	repo := NewRepository(db)

	// Define a test login precheck DTO
	testLoginPrecheck := dto.LoginPrecheckDTO{
		Username: "test_user",
	}

	// Expect a query to be executed and return a row
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT 1`)).WithArgs(testLoginPrecheck.Username).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "first_name", "last_name", "password", "salt"}).AddRow(1, "user@test.com", "test_user", "Test", "User", "testpass", "testsalt"))

	// Call the LoginPreCheckUser function
	repo.LoginPreCheckUser(testLoginPrecheck)

	// Check if the expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestProfileUser(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the user repository with the mock database connection
	repo := NewRepository(db)

	// Expect a query to be executed and return a row
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT 1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "first_name", "last_name", "password", "salt"}).AddRow(1, "user@test.com", "test_user", "Test", "User", "testpass", "testsalt"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_metadata" WHERE user_id = $1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "key", "value"}).AddRow(1, 1, "email_verified", "true").AddRow(2, 1, "display_name", "user").AddRow(3, 1, "country", "Unknown"))

	// Call the ProfileUser function
	repo.ProfileUser(1)

	// Check if the expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestUpdateDisplayName(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the user repository with the mock database connection
	repo := NewRepository(db)

	// Define a test update display name DTO
	testUpdateDisplayName := dto.UpdateDisplayNameDTO{
		DisplayName: "new_user",
	}

	// Call the UpdateDisplayName function
	repo.UpdateDisplayName(1, testUpdateDisplayName)
}

// Javokhir started the testing
func TestLoginClient(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the client repository with the mock database connection
	repo := NewRepository(db)

	// Define a test client username
	testUsername := "testuser"

	// Expect a query to be executed and return a row
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "clients" WHERE username = $1 ORDER BY "clients"."id" LIMIT 1`)).WithArgs(testUsername).WillReturnRows(sqlmock.NewRows([]string{"id", "secret", "username", "password"}).AddRow("notanid", "testsecret", "testuser", "testpassword"))

	// Call the LoginClient method
	client, err := repo.LoginClient(dto.LoginClientDTO{Username: testUsername, Password: "testpassword"})
	if err != nil {
		t.Fatalf("Error while logging in client: %v", err)
	}

	// Check if the returned client is correct
	if client.Username != testUsername {
		t.Errorf("Expected username: %s, got: %s", testUsername, client.Username)
	}

	// Check if the expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestProfileClient(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the client repository with the mock database connection
	repo := NewRepository(db)

	// Define a test user ID
	testUserID := "testuser"

	// Expect a query to be executed and return a row
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "clients" WHERE username = $1 ORDER BY "clients"."id" LIMIT 1`)).WithArgs(testUserID).WillReturnRows(sqlmock.NewRows([]string{"id", "secret", "name", "redirect_uri", "username"}).AddRow("notanid", "testsecret", "testuser", "https://gcitizen.com/callback", "testuser"))

	// Call the ProfileClient method
	client, err := repo.ProfileClient(testUserID)

	// Check for errors
	if err != nil {
		t.Fatalf("Error while retrieving client profile: %v", err)
	}

	// Check if the returned client is correct
	if client.Username != testUserID {
		t.Errorf("Expected user ID: %s, got: %s", testUserID, client.Name)
	}

	// Check if the expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

// Javokhir finished the testing

func TestFindUser(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	var testUserId uint = 1
	testUserEmail := "example@test.com"

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT 1`),
	).WithArgs(
		testUserId,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "email", "username", "password", "first_name",
				"last_name", "salt", "email_proof", "verification_code"},
		).AddRow(testUserId, testUserEmail, "username", "password", "first_name",
			"last_name", "salt", "", "",
		),
	)

	user, err := repository.FindUser(testUserId)

	if err != nil {
		t.Fatalf("Error while retrieving user: %v", err)
	}

	assert.Equal(t, testUserId, user.ID)
	assert.Equal(t, testUserEmail, user.Email)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestSaveProofOfEmailVerification_UsersTableUpdateFails(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "users" SET "email_proof"=$1,"verification_code"=$2 WHERE id = $3`),
	).WithArgs(
		emailProof, verificationCode, userId,
	).WillReturnError(
		fmt.Errorf(""),
	)

	mock.ExpectRollback()

	err = repository.SaveProofOfEmailVerification(userId, verificationCode, emailProof)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestSaveProofOfEmailVerification_DeletingEmailVerificationDataFails(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "users" SET "email_proof"=$1,"verification_code"=$2 WHERE id = $3`),
	).WithArgs(
		emailProof, verificationCode, userId,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)

	mock.ExpectExec(
		regexp.QuoteMeta(`DELETE FROM "email_verification_data" WHERE user_id = $1`),
	).WithArgs(
		userId,
	).WillReturnError(
		fmt.Errorf(""),
	)

	mock.ExpectRollback()

	err = repository.SaveProofOfEmailVerification(userId, verificationCode, emailProof)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestSaveProofOfEmailVerification_SettingUserStatusAsVerifiedFails(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "users" SET "email_proof"=$1,"verification_code"=$2 WHERE id = $3`),
	).WithArgs(
		emailProof, verificationCode, userId,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)

	mock.ExpectExec(
		regexp.QuoteMeta(`DELETE FROM "email_verification_data" WHERE user_id = $1`),
	).WithArgs(
		userId,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "user_metadata" SET "value"=$1,"updated_at"=$2 WHERE user_id = $3 AND key = $4`),
	).WithArgs(
		"true",
		sqlmock.AnyArg(),
		userId,
		"email_verified",
	).WillReturnError(
		fmt.Errorf(""),
	)

	mock.ExpectRollback()

	err = repository.SaveProofOfEmailVerification(userId, verificationCode, emailProof)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestSaveProofOfEmailVerification_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "users" SET "email_proof"=$1,"verification_code"=$2 WHERE id = $3`),
	).WithArgs(
		emailProof, verificationCode, userId,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)

	mock.ExpectExec(
		regexp.QuoteMeta(`DELETE FROM "email_verification_data" WHERE user_id = $1`),
	).WithArgs(
		userId,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "user_metadata" SET "value"=$1,"updated_at"=$2 WHERE user_id = $3 AND key = $4`),
	).WithArgs(
		"true",
		sqlmock.AnyArg(),
		userId,
		"email_verified",
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)

	mock.ExpectCommit()

	err = repository.SaveProofOfEmailVerification(userId, verificationCode, emailProof)

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestSaveEmailVerificationData_RowWithUserIdDoesNotExist(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "email_verification_data" WHERE "email_verification_data"."user_id" = $1 ORDER BY "email_verification_data"."id" LIMIT 1`),
	).WithArgs(
		userId,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "user_id", "verification_code", "expires_at"},
		),
	)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "email_verification_data" ("user_id","verification_code","expires_at") VALUES ($1,$2,$3) RETURNING "id"`,
		),
	).WithArgs(
		userId, verificationCode, timestamp,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "user_id", "verification_code", "expires_at"},
		).AddRow(
			1, userId, verificationCode, timestamp,
		),
	)

	mock.ExpectCommit()

	err = repository.SaveEmailVerificationData(
		models.EmailVerificationData{
			UserId:           userId,
			VerificationCode: verificationCode,
			ExpiresAt:        timestamp,
		},
	)

	if err != nil {
		t.Fatalf("Error while saving email verification data: %v", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestSaveEmailVerificationData_RowWithUserIdExists(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "email_verification_data" WHERE "email_verification_data"."user_id" = $1 ORDER BY "email_verification_data"."id" LIMIT 1`,
		),
	).WithArgs(
		userId,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "user_id", "verification_code", "expires_at"},
		).AddRow(
			id, userId, verificationCode, timestamp,
		),
	)

	mock.ExpectExec(
		regexp.QuoteMeta(
			`UPDATE "email_verification_data" SET "expires_at"=$1,"user_id"=$2,"verification_code"=$3 WHERE "email_verification_data"."user_id" = $4 AND "id" = $5`,
		),
	).WithArgs(
		timestamp, userId, verificationCode, userId, id,
	).WillReturnResult(sqlmock.NewResult(2, 1))

	mock.ExpectCommit()

	err = repository.SaveEmailVerificationData(
		models.EmailVerificationData{
			UserId:           userId,
			VerificationCode: verificationCode,
			ExpiresAt:        timestamp,
		},
	)

	if err != nil {
		t.Fatalf("Error while saving email verification data: %v", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetEmailVerificationData(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "email_verification_data" WHERE user_id = $1 ORDER BY "email_verification_data"."id" LIMIT 1`,
		),
	).WithArgs(
		userId,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "user_id", "verification_code", "expires_at"},
		).AddRow(
			id, userId, verificationCode, timestamp,
		),
	)

	data, err := repository.GetEmailVerificationData(userId)

	if err != nil {
		t.Fatalf("Error while getting email verification data: %v", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, userId, data.UserId)
	assert.Equal(t, verificationCode, data.VerificationCode)
	assert.Equal(t, timestamp, data.ExpiresAt)
}
