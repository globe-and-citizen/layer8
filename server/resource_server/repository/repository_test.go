package repository

import (
	"database/sql"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/utils"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const id uint = 1

const userId uint = 1
const username = "test_username"
const usernameWithSpecialChars = "user_with_special_chars!@#$%^&*()"
const userFirstName = "first_name"
const userLastName = "last_name"
const userSalt = "user_salt"
const userDisplayName = "display_name"
const userColor = "color"
const userBio = "bio"

const verificationCode = "123456"

const clientId = "1"
const clientUsername = "client_username"
const clientName = "test_client"
const clientSecret = "client_secret"
const redirectUri = "https://gcitizen.com/callback"
const backendUri = "https://gcitizen.com/backend"
const clientSalt = "client_salt"
const clientIterationCount = 4096
const zkKeyPairId uint = 2

var timestamp = time.Date(2024, time.May, 24, 14, 0, 0, 0, time.UTC)

var emailProof = []byte("AbcdfTs")
var provingKey = []byte("proving key")
var verifyingKey = []byte("verifying key")

var publicKey = make([]byte, 33)
var storedKey = string(make([]byte, 33))
var serverKey = string(make([]byte, 33))

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

func TestGetClientData_ClientDoesNotExist(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "clients" WHERE name = $1 ORDER BY "clients"."id" LIMIT $2`),
	).WithArgs(
		clientName, 1,
	).WillReturnRows(
		sqlmock.NewRows([]string{"id", "secret", "name", "redirect_uri"}),
	)

	_, err := repository.GetClientData(clientName)

	assert.NotNil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetClientData_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "clients" WHERE name = $1 ORDER BY "clients"."id" LIMIT $2`),
	).WithArgs(
		clientName, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "secret", "name", "redirect_uri", "username", "salt"},
		).AddRow(
			clientId, clientSecret, clientName, redirectUri, clientUsername, clientSalt,
		),
	)

	client, err := repository.GetClientData(clientName)

	assert.Nil(t, err)

	assert.Equal(t, clientId, client.ID)
	assert.Equal(t, clientSecret, client.Secret)
	assert.Equal(t, clientName, client.Name)
	assert.Equal(t, redirectUri, client.RedirectURI)
	assert.Equal(t, clientUsername, client.Username)
	assert.Equal(t, clientSalt, client.Salt)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestProfileUser_UserIsNotFoundInTheUsersTable(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`),
	).WithArgs(
		userId, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "username", "first_name", "last_name", "password", "salt"},
		),
	)

	_, _, err := repository.ProfileUser(userId)

	assert.NotNil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestProfileUser_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`),
	).WithArgs(
		userId, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "username", "first_name", "last_name", "salt"},
		).AddRow(
			userId, username, userFirstName, userLastName, userSalt,
		),
	)

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "user_metadata" WHERE id = $1`),
	).WithArgs(
		userId,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "is_email_verified", "is_phone_number_verified", "display_name", "color", "bio"},
		).AddRow(
			userId, true, true, userDisplayName, userColor, userBio,
		),
	)

	user, userMetadata, err := repository.ProfileUser(userId)

	assert.Nil(t, err)
	assert.Equal(t, userId, user.ID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, userDisplayName, userMetadata.DisplayName)
	assert.Equal(t, userColor, userMetadata.Color)
	assert.Equal(t, userBio, userMetadata.Bio)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUserMetadata_TableUpdateFails(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "user_metadata" SET "display_name"=$1,"color"=$2,"bio"=$3,"updated_at"=$4 WHERE id = $5`),
	).WithArgs(
		userDisplayName, userColor, userBio, sqlmock.AnyArg(), userId,
	).WillReturnError(
		fmt.Errorf("failed to update user's metadata"),
	)

	mock.ExpectRollback()

	err := repository.UpdateUserMetadata(
		userId,
		dto.UpdateUserMetadataDTO{
			DisplayName: userDisplayName,
			Color:       userColor,
			Bio:         userBio,
		},
	)

	assert.NotNil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUserMetadata_AllFieldsPopulated_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "user_metadata" SET "display_name"=$1,"color"=$2,"bio"=$3,"updated_at"=$4 WHERE id = $5`),
	).WithArgs(
		userDisplayName, userColor, userBio, sqlmock.AnyArg(), userId,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	mock.ExpectCommit()

	err := repository.UpdateUserMetadata(
		userId,
		dto.UpdateUserMetadataDTO{
			DisplayName: userDisplayName,
			Color:       userColor,
			Bio:         userBio,
		},
	)

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUserMetadata_MetadataUpdatedPartially_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "user_metadata" SET "display_name"=$1,"color"=$2,"updated_at"=$3 WHERE id = $4`),
	).WithArgs(
		userDisplayName, userColor, sqlmock.AnyArg(), userId,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	mock.ExpectCommit()

	err := repository.UpdateUserMetadata(
		userId,
		dto.UpdateUserMetadataDTO{
			DisplayName: userDisplayName,
			Color:       userColor,
		},
	)

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestProfileClient_ClientNotFound(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "clients" WHERE username = $1 ORDER BY "clients"."id" LIMIT $2`),
	).WithArgs(
		clientUsername, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "secret", "name", "redirect_uri", "username"},
		),
	)

	_, err := repository.ProfileClient(clientUsername)

	assert.NotNil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestProfileClient_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "clients" WHERE username = $1 ORDER BY "clients"."id" LIMIT $2`),
	).WithArgs(
		clientUsername, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "secret", "name", "redirect_uri", "username"},
		).AddRow(
			clientId, clientSecret, clientName, redirectUri, clientUsername,
		),
	)

	client, err := repository.ProfileClient(clientUsername)

	assert.Nil(t, err)

	assert.Equal(t, clientId, client.ID)
	assert.Equal(t, clientUsername, client.Username)
	assert.Equal(t, clientName, client.Name)
	assert.Equal(t, clientSecret, client.Secret)
	assert.Equal(t, redirectUri, client.RedirectURI)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestFindUser_UserDoesNotExist(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`),
	).WithArgs(
		userId, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "username", "password", "first_name",
				"last_name", "salt", "email_proof", "verification_code"},
		),
	)

	_, err := repository.FindUser(userId)

	assert.NotNil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestFindUser_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`),
	).WithArgs(
		userId, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "username", "password", "first_name",
				"last_name", "salt", "email_proof", "verification_code"},
		).AddRow(userId, username, "password", "first_name",
			"last_name", "salt", "", "",
		),
	)

	user, err := repository.FindUser(userId)

	if err != nil {
		t.Fatalf("Error while retrieving user: %v", err)
	}

	assert.Equal(t, userId, user.ID)
	assert.Equal(t, username, user.Username)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetEmailVerificationData_VerificationDataForTheGivenUserDoesNotExist(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "email_verification_data" WHERE user_id = $1 ORDER BY "email_verification_data"."id" LIMIT $2`,
		),
	).WithArgs(
		userId, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "user_id", "verification_code", "expires_at"},
		),
	)

	_, err := repository.GetEmailVerificationData(userId)

	assert.NotNil(t, err)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetEmailVerificationData_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "email_verification_data" WHERE user_id = $1 ORDER BY "email_verification_data"."id" LIMIT $2`,
		),
	).WithArgs(
		userId, 1,
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

func TestSaveEmailVerificationData_RowWithUserIdDoesNotExist(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "email_verification_data" WHERE "email_verification_data"."user_id" = $1 ORDER BY "email_verification_data"."id" LIMIT $2`),
	).WithArgs(
		userId, 1,
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
			`SELECT * FROM "email_verification_data" WHERE "email_verification_data"."user_id" = $1 ORDER BY "email_verification_data"."id" LIMIT $2`,
		),
	).WithArgs(
		userId, 1,
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

func TestSaveProofOfEmailVerification_UsersTableUpdateFails(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "users" SET "email_proof"=$1,"verification_code"=$2,"zk_key_pair_id"=$3 WHERE id = $4`),
	).WithArgs(
		emailProof, verificationCode, zkKeyPairId, userId,
	).WillReturnError(
		fmt.Errorf(""),
	)

	mock.ExpectRollback()

	err = repository.SaveProofOfEmailVerification(userId, verificationCode, emailProof, zkKeyPairId)

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
		regexp.QuoteMeta(`UPDATE "users" SET "email_proof"=$1,"verification_code"=$2,"zk_key_pair_id"=$3 WHERE id = $4`),
	).WithArgs(
		emailProof, verificationCode, zkKeyPairId, userId,
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

	err = repository.SaveProofOfEmailVerification(userId, verificationCode, emailProof, zkKeyPairId)

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
		regexp.QuoteMeta(`UPDATE "users" SET "email_proof"=$1,"verification_code"=$2,"zk_key_pair_id"=$3 WHERE id = $4`),
	).WithArgs(
		emailProof, verificationCode, zkKeyPairId, userId,
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
		regexp.QuoteMeta(`UPDATE "user_metadata" SET "is_email_verified"=$1,"updated_at"=$2 WHERE id = $3`),
	).WithArgs(
		true, sqlmock.AnyArg(), userId,
	).WillReturnError(
		fmt.Errorf(""),
	)

	mock.ExpectRollback()

	err = repository.SaveProofOfEmailVerification(userId, verificationCode, emailProof, zkKeyPairId)

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
		regexp.QuoteMeta(`UPDATE "users" SET "email_proof"=$1,"verification_code"=$2,"zk_key_pair_id"=$3 WHERE id = $4`),
	).WithArgs(
		emailProof, verificationCode, zkKeyPairId, userId,
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
		regexp.QuoteMeta(`UPDATE "user_metadata" SET "is_email_verified"=$1,"updated_at"=$2 WHERE id = $3`),
	).WithArgs(
		true, sqlmock.AnyArg(), userId,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)

	mock.ExpectCommit()

	err = repository.SaveProofOfEmailVerification(userId, verificationCode, emailProof, zkKeyPairId)

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestSaveZkSnarksKeyPair_FailedToSaveZkKeyPair(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "zk_snarks_key_pairs" ("proving_key","verifying_key") VALUES ($1,$2) RETURNING "id"`,
		),
	).WithArgs(
		provingKey, verifyingKey,
	).WillReturnError(
		fmt.Errorf(""),
	)

	mock.ExpectRollback()

	_, err = repository.SaveZkSnarksKeyPair(
		models.ZkSnarksKeyPair{
			ProvingKey:   provingKey,
			VerifyingKey: verifyingKey,
		},
	)

	assert.NotNil(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestSaveZkSnarksKeyPair_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "zk_snarks_key_pairs" ("proving_key","verifying_key") VALUES ($1,$2) RETURNING "id"`,
		),
	).WithArgs(
		provingKey, verifyingKey,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id"},
		).AddRow(zkKeyPairId),
	)

	mock.ExpectCommit()

	actualZkTableId, err := repository.SaveZkSnarksKeyPair(
		models.ZkSnarksKeyPair{
			ProvingKey:   provingKey,
			VerifyingKey: verifyingKey,
		},
	)

	assert.Nil(t, err)
	assert.Equal(t, zkKeyPairId, actualZkTableId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetZkSnarksKeys_FailedToGetNewestZkSnarksKeys(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "zk_snarks_key_pairs" ORDER BY "zk_snarks_key_pairs"."id" DESC LIMIT $1`),
	).WithArgs(1).WillReturnError(
		fmt.Errorf(""),
	)

	_, err = repository.GetLatestZkSnarksKeys()

	assert.NotNil(t, err)
}

func TestGetZkSnarksKeys_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "zk_snarks_key_pairs" ORDER BY "zk_snarks_key_pairs"."id" DESC LIMIT $1`),
	).WithArgs(1).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "proving_key", "verifying_key"},
		).AddRow(
			zkKeyPairId, provingKey, verifyingKey,
		),
	)

	zkKeyPair, err := repository.GetLatestZkSnarksKeys()

	assert.Nil(t, err)
	assert.True(t, utils.Equal(provingKey, zkKeyPair.ProvingKey))
	assert.True(t, utils.Equal(verifyingKey, zkKeyPair.VerifyingKey))
}

func TestGetUserForUsername_UserNotFound(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`,
		),
	).WithArgs(
		username, 1,
	).WillReturnError(
		fmt.Errorf("user not found"),
	)

	_, err = repository.GetUserForUsername(username)

	assert.NotNil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetUserForUsername_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`,
		),
	).WithArgs(
		username, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"id", "username", "salt", "stored_key", "server_key"},
		).AddRow(
			userId, username, userSalt, storedKey, serverKey,
		),
	)

	user, err := repository.GetUserForUsername(username)

	assert.Nil(t, err)
	assert.Equal(t, userId, user.ID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, userSalt, user.Salt)
	assert.Equal(t, storedKey, user.StoredKey)
	assert.Equal(t, serverKey, user.ServerKey)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestRegisterPrecheckUser_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	req := dto.RegisterUserPrecheckDTO{
		Username: "test_user",
	}
	salt := "random_salt"
	iterCount := 4096

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "users" ("username","salt","email_proof","verification_code","zk_key_pair_id","public_key","iteration_count","server_key","stored_key") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`,
		),
	).WithArgs(
		req.Username, salt, sqlmock.AnyArg(), "", 0, sqlmock.AnyArg(), iterCount, sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(1),
	)

	mock.ExpectCommit()

	err := repository.RegisterPrecheckUser(req, salt, iterCount)

	assert.Nil(t, err, "Error should be nil")
	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}

func TestRegisterPrecheckUser_RepositoryError(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	req := dto.RegisterUserPrecheckDTO{
		Username: "test_user",
	}
	salt := "random_salt"
	iterCount := 4096

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "users" ("username","salt","email_proof","verification_code","zk_key_pair_id","public_key","iteration_count","server_key","stored_key") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`,
		),
	).WithArgs(
		req.Username, salt, sqlmock.AnyArg(), "", 0, sqlmock.AnyArg(), iterCount, sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnError(fmt.Errorf("failed to create user"))

	mock.ExpectRollback()

	err := repository.RegisterPrecheckUser(req, salt, iterCount)

	assert.NotNil(t, err, "Expected error due to database error")
	assert.Equal(t, "failed to create a new user: failed to create user", err.Error())

	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}

func TestRegisterPrecheckUser_BeginTransactionFailure(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	req := dto.RegisterUserPrecheckDTO{
		Username: "test_user",
	}
	salt := "random_salt"
	iterCount := 4096

	mock.ExpectBegin().WillReturnError(fmt.Errorf("failed to begin transaction"))

	err := repository.RegisterPrecheckUser(req, salt, iterCount)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Contains(t, err.Error(), "failed to begin transaction", "Error message should contain 'failed to begin transaction'")

	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}

func TestRegisterUser_FailToGetUserRecord(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`,
		),
	).WithArgs(
		username, 1,
	).WillReturnError(
		fmt.Errorf("could not get user record"),
	)

	mock.ExpectRollback()

	userDto := dto.RegisterUserDTO{
		Username:  username,
		PublicKey: publicKey,
		StoredKey: storedKey,
		ServerKey: serverKey,
	}
	err = repository.RegisterUser(userDto)

	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}

func TestRegisterUser_FailToUpdateUserRecord(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`,
		),
	).WithArgs(
		username, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{
				"id", "username", "salt", "email_proof", "verification_code", "zk_key_pair_id", "public_key", "iteration_count", "server_key", "stored_key",
			},
		).AddRow(
			userId, username, userSalt, emailProof, verificationCode, 0, publicKey, 4096, serverKey, storedKey,
		),
	)

	mock.ExpectExec(
		regexp.QuoteMeta(
			`UPDATE "users" SET "public_key"=$1,"server_key"=$2,"stored_key"=$3 WHERE "id" = $4`,
		),
	).WithArgs(
		publicKey, serverKey, storedKey, userId,
	).WillReturnError(
		fmt.Errorf("could not update user record"),
	)

	mock.ExpectRollback()

	userDto := dto.RegisterUserDTO{
		Username:  username,
		PublicKey: publicKey,
		StoredKey: storedKey,
		ServerKey: serverKey,
	}
	err = repository.RegisterUser(userDto)

	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}

func TestRegisterUser_FailToInsertANewUserMetadataRecord(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`,
		),
	).WithArgs(
		username, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{
				"id", "username", "salt", "email_proof", "verification_code", "zk_key_pair_id", "public_key", "iteration_count", "server_key", "stored_key",
			},
		).AddRow(
			userId, username, userSalt, emailProof, verificationCode, 0, publicKey, 4096, serverKey, storedKey,
		),
	)

	mock.ExpectExec(
		regexp.QuoteMeta(
			`UPDATE "users" SET "public_key"=$1,"server_key"=$2,"stored_key"=$3 WHERE "id" = $4`,
		),
	).WithArgs(
		publicKey, serverKey, storedKey, userId,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "user_metadata" ("display_name","color","bio","is_email_verified","is_phone_number_verified","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "created_at","updated_at","id"`,
		),
	).WithArgs(
		"", "", "", false, false, userId,
	).WillReturnError(
		fmt.Errorf(""),
	)

	mock.ExpectRollback()

	userDto := dto.RegisterUserDTO{
		Username:  username,
		PublicKey: publicKey,
		StoredKey: storedKey,
		ServerKey: serverKey,
	}
	err = repository.RegisterUser(userDto)

	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}

func TestRegisterUser_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`,
		),
	).WithArgs(
		username, 1,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{
				"id", "username", "salt", "email_proof", "verification_code", "zk_key_pair_id", "public_key", "iteration_count", "server_key", "stored_key",
			},
		).AddRow(
			userId, username, userSalt, emailProof, verificationCode, 0, publicKey, 4096, serverKey, storedKey,
		),
	)

	mock.ExpectExec(
		regexp.QuoteMeta(
			`UPDATE "users" SET "public_key"=$1,"server_key"=$2,"stored_key"=$3 WHERE "id" = $4`,
		),
	).WithArgs(
		publicKey, serverKey, storedKey, userId,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "user_metadata" ("display_name","color","bio","is_email_verified","is_phone_number_verified","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "created_at","updated_at","id"`,
		),
	).WithArgs(
		"", "", "", false, false, userId,
	).WillReturnRows(
		sqlmock.NewRows(
			[]string{"created_at", "updated_at", "id"},
		).AddRow(
			time.Time{}, time.Time{}, userId,
		),
	)

	mock.ExpectCommit()

	userDto := dto.RegisterUserDTO{
		Username:  username,
		PublicKey: publicKey,
		StoredKey: storedKey,
		ServerKey: serverKey,
	}
	err = repository.RegisterUser(userDto)

	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}

func TestUpdateUserPassword_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "users" SET "server_key"=$1,"stored_key"=$2 WHERE username=$3`),
	).WithArgs(serverKey, storedKey, username).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repository.UpdateUserPassword(username, storedKey, serverKey)

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUserPassword_UpdateQueryFailed(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "users" SET "server_key"=$1,"stored_key"=$2 WHERE username=$3`),
	).WithArgs(serverKey, storedKey, username).
		WillReturnError(fmt.Errorf("database error"))

	mock.ExpectRollback()

	err := repository.UpdateUserPassword(username, storedKey, serverKey)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "database error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUserPassword_EdgeCaseUsername(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "users" SET "server_key"=$1,"stored_key"=$2 WHERE username=$3`),
	).WithArgs(serverKey, storedKey, usernameWithSpecialChars).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repository.UpdateUserPassword(usernameWithSpecialChars, storedKey, serverKey)

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestRegisterPrecheckClient_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	req := dto.RegisterClientPrecheckDTO{
		Username: clientUsername,
	}
	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(
			`INSERT INTO "clients" ("id","secret","name","redirect_uri","backend_uri","username","salt","iteration_count","server_key","stored_key","x509_certificate_bytes") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		),
	).WithArgs(
		"", "", "", "", "", clientUsername, clientSalt, clientIterationCount, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	mock.ExpectCommit()

	err := repository.RegisterPrecheckClient(req, clientSalt, clientIterationCount)

	assert.Nil(t, err, "Error should be nil")
	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}

func TestRegisterPrecheckClient_RepositoryError(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	req := dto.RegisterClientPrecheckDTO{
		Username: clientUsername,
	}

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(
			`INSERT INTO "clients" ("id","secret","name","redirect_uri","backend_uri","username","salt","iteration_count","server_key","stored_key","x509_certificate_bytes") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		),
	).WithArgs(
		"", "", "", "", "", clientUsername, clientSalt, clientIterationCount, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnError(fmt.Errorf("failed to create client"))

	mock.ExpectRollback()

	err := repository.RegisterPrecheckClient(req, clientSalt, clientIterationCount)

	assert.NotNil(t, err, "Expected error due to database error")
	assert.Equal(t, "failed to create a new client: failed to create client", err.Error())

	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}

func TestRegisterPrecheckClient_BeginTransactionFailure(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	req := dto.RegisterClientPrecheckDTO{
		Username: clientUsername,
	}

	mock.ExpectBegin().WillReturnError(fmt.Errorf("failed to begin transaction"))

	err := repository.RegisterPrecheckClient(req, clientSalt, clientIterationCount)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Contains(t, err.Error(), "failed to begin transaction", "Error message should contain 'failed to begin transaction'")

	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}

func TestRegisterClient_FailToUpdateClientRecord(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(
			`UPDATE "clients" SET "backend_uri"=$1,"id"=$2,"name"=$3,"redirect_uri"=$4,"secret"=$5,"server_key"=$6,"stored_key"=$7 WHERE username = $8`,
		),
	).WithArgs(
		backendUri, clientId, clientName, redirectUri, clientSecret, serverKey, storedKey, clientUsername,
	).WillReturnError(
		fmt.Errorf("could not update client record"),
	)

	mock.ExpectRollback()

	clientDto := dto.RegisterClientDTO{
		Username:    clientUsername,
		Name:        clientName,
		RedirectURI: redirectUri,
		BackendURI:  backendUri,
		StoredKey:   storedKey,
		ServerKey:   serverKey,
	}
	err = repository.RegisterClient(clientDto, clientId, clientSecret)

	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}

func TestRegisterClient_Success(t *testing.T) {
	SetUp(t)
	defer mockDB.Close()

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta(
			`UPDATE "clients" SET "backend_uri"=$1,"id"=$2,"name"=$3,"redirect_uri"=$4,"secret"=$5,"server_key"=$6,"stored_key"=$7 WHERE username = $8`,
		),
	).WithArgs(
		backendUri, clientId, clientName, redirectUri, clientSecret, serverKey, storedKey, clientUsername,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	mock.ExpectCommit()

	clientDto := dto.RegisterClientDTO{
		Username:    clientUsername,
		Name:        clientName,
		RedirectURI: redirectUri,
		BackendURI:  backendUri,
		StoredKey:   storedKey,
		ServerKey:   serverKey,
	}
	err = repository.RegisterClient(clientDto, clientId, clientSecret)

	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "There were unfulfilled expectations!")
}
