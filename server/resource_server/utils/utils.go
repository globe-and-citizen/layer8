package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/models"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/xdg-go/pbkdf2"
)

const SaltSize = 32
const SecretSize = 32

var WorkingDirectory string

// Response is used for static shape json return
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type EmptyObj struct{}

func GetPwd() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	WorkingDirectory = dir
}

func GenerateRandomSalt(saltSize int) string {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(salt[:])
}

func SaltAndHashPassword(password string, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 4096, 32, sha1.New)
	return hex.EncodeToString(dk[:])
}

func CheckPassword(password string, salt string, hash string) bool {
	return SaltAndHashPassword(password, salt) == hash
}

func BuildResponse(status bool, message string, data interface{}) Response {
	res := Response{
		Status:  status,
		Message: message,
		Error:   nil,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	res := Response{
		Status:  false,
		Message: message,
		Error:   splittedError,
		Data:    data,
	}

	return res
}

func HandleError(w http.ResponseWriter, status int, message string, err error) {
	w.WriteHeader(status)
	res := BuildErrorResponse(message, err.Error(), EmptyObj{})
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func GenerateUUID() string {
	newUUID := uuid.New()

	return newUUID.String()
}
func GenerateSecret(secretSize int) string {
	var randomBytes = make([]byte, secretSize)

	_, err := rand.Read(randomBytes[:])

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(randomBytes[:])
}

func CompleteLogin(req dto.LoginUserDTO, user models.User) (models.LoginUserResponseOutput, error) {
	HashedAndSaltedPass := SaltAndHashPassword(req.Password, user.Salt)

	if user.Password != HashedAndSaltedPass {
		return models.LoginUserResponseOutput{}, fmt.Errorf("invalid password")
	}

	tokenString, err := GenerateToken(user)
	if err != nil {
		return models.LoginUserResponseOutput{}, err
	}

	resp := models.LoginUserResponseOutput{
		Token: tokenString,
	}
	return resp, nil
}

func CompleteClientLogin(req dto.LoginClientDTO, client models.Client) (models.LoginUserResponseOutput, error) {

	JWT_SECRET_STR := os.Getenv("JWT_SECRET")
	JWT_SECRET_BYTE := []byte(JWT_SECRET_STR)

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &models.ClientClaims{
		UserName: client.Username,
		UserID:   client.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "GlobeAndCitizen",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWT_SECRET_BYTE)
	if err != nil {
		return models.LoginUserResponseOutput{}, err
	}

	resp := models.LoginUserResponseOutput{
		Token: tokenString,
	}
	fmt.Println(resp)

	return resp, nil
}

func ValidateToken(tokenString string) (uint, error) {
	claims := &models.Claims{}
	JWT_SECRET_STR := os.Getenv("JWT_SECRET")
	JWT_SECRET_BYTE := []byte(JWT_SECRET_STR)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET_BYTE, nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	return claims.UserID, nil
}

func ValidateClientToken(tokenString string) (string, error) {
	claims := &models.ClientClaims{}
	JWT_SECRET_STR := os.Getenv("JWT_SECRET")
	JWT_SECRET_BYTE := []byte(JWT_SECRET_STR)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET_BYTE, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	return claims.UserName, nil
}

func GenerateToken(user models.User) (string, error) {
	JWT_SECRET_STR := os.Getenv("JWT_SECRET")
	JWT_SECRET_BYTE := []byte(JWT_SECRET_STR)

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &models.Claims{
		UserName: user.Username,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "GlobeAndCitizen",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWT_SECRET_BYTE)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
