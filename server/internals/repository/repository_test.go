package repository_test

import (
	"globe-and-citizen/layer8/server/internals/repository"
	"globe-and-citizen/layer8/server/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLoginUserPrecheck(t *testing.T) {
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
	repo := repository.NewOauthRepository(db)

	// Make a mock loginPrecheck input
	username := "test_user"

	mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE username = (.+)").WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{"salt"}).AddRow("test_salt"))

	// Call the function and check the result
	salt, err := repo.LoginUserPrecheck(username)
	if err != nil {
		t.Fatal("Failed to call LoginUserPrecheck:", err)
	}

	assert.Equal(t, "test_salt", salt)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("Unmet expectations:", err)
	}
}

func TestGetUser(t *testing.T) {
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
	repo := repository.NewOauthRepository(db)

	// Make a mock getUser input
	username := "test_user"

	mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE username = (.+)").WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password", "salt", "first_name", "last_name"}).AddRow(1, "test_email", "test_user", "test_password", "test_salt", "test_first_name", "test_last_name"))

	// Call the function and check the result
	user, err := repo.GetUser(username)
	if err != nil {
		t.Fatal("Failed to call GetUser:", err)
	}

	assert.Equal(t, "test_user", user.Username)
	assert.Equal(t, "test_email", user.Email)
	assert.Equal(t, "test_password", user.Password)
	assert.Equal(t, "test_salt", user.Salt)
	assert.Equal(t, "test_first_name", user.FirstName)
	assert.Equal(t, "test_last_name", user.LastName)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("Unmet expectations:", err)
	}
}

func TestGetUserByID(t *testing.T) {
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
	repo := repository.NewOauthRepository(db)

	// Make a mock getUserByID input
	id := int64(1)

	mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE id = (.+)").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username", "password", "salt", "first_name", "last_name"}).AddRow(1, "test_email", "test_user", "test_password", "test_salt", "test_first_name", "test_last_name"))

	// Call the function and check the result
	user, err := repo.GetUserByID(id)
	if err != nil {
		t.Fatal("Failed to call GetUserByID:", err)
	}

	assert.Equal(t, "test_user", user.Username)
	assert.Equal(t, "test_email", user.Email)
	assert.Equal(t, "test_password", user.Password)
	assert.Equal(t, "test_salt", user.Salt)
	assert.Equal(t, "test_first_name", user.FirstName)
	assert.Equal(t, "test_last_name", user.LastName)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("Unmet expectations:", err)
	}
}

func TestGetUserMetadata(t *testing.T) {
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
	repo := repository.NewOauthRepository(db)

	// Make a mock getUserMetadata input
	userID := int64(1)
	key := "test_key"

	mock.ExpectQuery("SELECT (.+) FROM \"user_metadata\" WHERE user_id = (.+) AND key = (.+)").WithArgs(userID, key).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "key", "value"}).AddRow(1, 1, "test_key", "test_value"))

	// Call the function and check the result
	userMetadata, err := repo.GetUserMetadata(userID, key)
	if err != nil {
		t.Fatal("Failed to call GetUserMetadata:", err)
	}

	assert.Equal(t, "test_key", userMetadata.Key)
	assert.Equal(t, "test_value", userMetadata.Value)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("Unmet expectations:", err)
	}
}

func TestSetClient(t *testing.T) {
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
	repo := repository.NewOauthRepository(db)

	// Make a mock setClient input
	client := &models.Client{
		ID:          "test_id",
		Secret:      "test_secret",
		Name:        "test_name",
		RedirectURI: "test_redirect_uri",
	}

	mock.ExpectQuery("SELECT (.+) FROM \"clients\" WHERE id = (.+)").WithArgs(client.ID).WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO \"clients\"")).WithArgs(client.ID, client.Secret, client.Name, client.RedirectURI).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the function and check the result
	err = repo.SetClient(client)
	if err != nil {
		t.Fatal("Failed to call SetClient:", err)
	}

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("Unmet expectations:", err)
	}
}

func TestGetClient(t *testing.T) {
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
	repo := repository.NewOauthRepository(db)

	// Make a mock getClient input
	id := "test_id"

	mock.ExpectQuery("SELECT (.+) FROM \"clients\" WHERE id = (.+)").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "secret", "name", "redirect_uri"}).AddRow("test_id", "test_secret", "test_name", "test_redirect_uri"))

	// Call the function and check the result
	client, err := repo.GetClient(id)
	if err != nil {
		t.Fatal("Failed to call GetClient:", err)
	}

	assert.Equal(t, "test_id", client.ID)
	assert.Equal(t, "test_secret", client.Secret)
	assert.Equal(t, "test_name", client.Name)
	assert.Equal(t, "test_redirect_uri", client.RedirectURI)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("Unmet expectations:", err)
	}
}

func TestSetTTL(t *testing.T) {
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
	repo := repository.NewOauthRepository(db)

	// Call the function and check the result
	err = repo.SetTTL("test_key", []byte("test_value"), 1)
	if err != nil {
		t.Fatal("Failed to call SetTTL:", err)
	}
}

func TestGetTTL(t *testing.T) {
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
	repo := repository.NewOauthRepository(db)

	// Call the function and check the result
	_, err = repo.GetTTL("test_key")
	if err != nil {
		t.Fatal("Failed to call GetTTL:", err)
	}
}
