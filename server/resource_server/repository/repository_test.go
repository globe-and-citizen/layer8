package repository

import (
	"globe-and-citizen/layer8/server/resource_server/dto"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func TestVerifyEmail(t *testing.T) {
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

	// Call the VerifyEmail function
	repo.VerifyEmail(1)

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
