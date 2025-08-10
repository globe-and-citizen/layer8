package service

import (
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/constants"
	"globe-and-citizen/layer8/server/entities"
	"globe-and-citizen/layer8/server/internals/repository"
	"globe-and-citizen/layer8/server/models"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	utilities "github.com/globe-and-citizen/layer8-utils"
	"golang.org/x/oauth2"

	rs_utils "globe-and-citizen/layer8/server/resource_server/utils"
)

type ServiceInterface interface {
	GetUserByToken(token string) (*models.User, error)
	LoginUser(username, password string) (map[string]interface{}, error)
	GenerateAuthorizationURL(config *oauth2.Config, userID int64) (*entities.AuthURL, error)
	GenerateAuthJwtCode(config *oauth2.Config, userID int64) (string, error)
	ExchangeCodeForToken(config *oauth2.Config, code string) (*oauth2.Token, error)
	AccessResourcesWithToken(token string) (map[string]interface{}, error)
	GetClient(id string) (*models.Client, error)
	VerifyToken(token string) (isvalid bool, err error)
	CheckClient(backendURL string) (*models.Client, error)
	SaveX509Certificate(clientID string, certificate string) error
	DecodeAuthorizationCode(secret string, code string) (*utilities.AuthCodeClaims, error)
	AuthenticateClient(uuid string, secret string) error
	GenerateAccessToken(authClaims *utilities.AuthCodeClaims, clientID string, clientSecret string) (string, error)
	ValidateAccessToken(clientSecret string, accessToken string) (*entities.ClientClaims, error)
	GetZkUserMetadata(scopesStr string, userID int64) (*entities.ZkMetadataResponse, error)
	AddTestClient() (*models.Client, error)
}

type Service struct {
	Repo repository.Repository
}

func NewService(repo repository.Repository) ServiceInterface {
	return &Service{
		Repo: repo,
	}
}

// GetUserByToken returns the user associated with the given token
func (u *Service) GetUserByToken(token string) (*models.User, error) {
	// verify token
	userID, err := utilities.VerifyUserToken(config.SECRET_KEY, token)
	if err != nil {
		return nil, err
	}
	// get user
	user, err := u.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *Service) LoginUser(username, password string) (map[string]interface{}, error) {
	// TODO: use SCRAM authentication here
	user, err := u.Repo.GetUser(username)
	if err != nil {
		return nil, err
	}

	token, err := utilities.GenerateUserToken(config.SECRET_KEY, int64(user.ID))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
		"user":  user,
	}, nil
}

// GenerateAuthorizationURL generates an authorization URL for the user to visit
// and authorize the application to access their account.
func (u *Service) GenerateAuthorizationURL(config *oauth2.Config, userID int64) (*entities.AuthURL, error) {
	// first, check that both client and user exist
	client, err := u.GetClient(config.ClientID)
	if err != nil {
		return nil, fmt.Errorf("could not get client: %v", err)
	}
	user, err := u.Repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %v", err)
	}

	state, stateErr := utilities.GenerateRandomString(24)
	if stateErr != nil {
		return nil, fmt.Errorf("could not generate random state: %v", stateErr)
	}

	// generate the auth code
	scopes := ""
	for _, scope := range config.Scopes {
		scopes += scope + ","
	}
	code, err := utilities.GenerateAuthCode(client.Secret, &utilities.AuthCodeClaims{
		ClientID:    config.ClientID,
		UserID:      int64(user.ID),
		RedirectURI: config.RedirectURL,
		Scopes:      scopes,
		ExpiresAt:   time.Now().Add(time.Minute * 5).Unix(),
	})
	if err != nil {
		return nil, fmt.Errorf("could not generate auth code: %v", err)
	}

	return &entities.AuthURL{
		URL: config.AuthCodeURL(
			state,
			oauth2.SetAuthURLParam("code", code),
		),
		Code:  code,
		State: state,
	}, nil
}

func (u *Service) GenerateAuthJwtCode(config *oauth2.Config, userID int64) (string, error) {
	// first, check that both client and user exist
	client, err := u.GetClient(config.ClientID)
	if err != nil {
		return "", fmt.Errorf("could not get client: %v", err)
	}
	user, err := u.Repo.GetUserByID(userID)
	if err != nil {
		return "", fmt.Errorf("could not get user: %v", err)
	}

	// generate the auth code
	scopes := ""
	for _, scope := range config.Scopes {
		scopes += scope + ","
	}
	code, err := utilities.GenerateAuthCode(client.Secret, &utilities.AuthCodeClaims{
		ClientID:    config.ClientID,
		UserID:      int64(user.ID),
		RedirectURI: config.RedirectURL,
		Scopes:      scopes,
		ExpiresAt:   time.Now().Add(time.Minute * 5).Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("could not generate auth code: %v", err)
	}

	return code, nil
}

// ExchangeCodeForToken generates an access token from an authorization code.
func (u *Service) ExchangeCodeForToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	// ensure that the secret is specified
	if config.ClientSecret == "" {
		return nil, fmt.Errorf("client secret is not specified")
	}
	// verify the code
	claims, err := utilities.DecodeAuthCode(config.ClientSecret, code)
	if err != nil {
		return nil, err
	}
	// generating random token
	token, err := utilities.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	// save token and claims for 5 minutes
	b, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}
	err = u.Repo.SetTTL("token:"+token, b, time.Minute*10)
	if err != nil {
		return nil, err
	}
	// generate the access token
	return &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Minute * 10),
	}, nil
}

// AccessResourcesWithToken returns the resources that the client has access to
// with the given token.
func (u *Service) AccessResourcesWithToken(token string) (map[string]interface{}, error) {
	// get the claims
	res, err := u.Repo.GetTTL("token:" + token)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("could not get token")
	}
	var claims utilities.AuthCodeClaims
	err = json.Unmarshal(res, &claims)
	if err != nil {
		return nil, err
	}

	// fmt.Println("claims headerMap:", claims.HeaderMap)
	// get the resources
	scopes := strings.Split(claims.Scopes, ",")
	resources := make(map[string]interface{})
	for _, scope := range scopes {
		switch scope {
		case constants.READ_USER_SCOPE:

			isEmailVerified, err := u.Repo.GetUserMetadata(claims.UserID, constants.USER_EMAIL_VERIFIED_METADATA_KEY)
			if err != nil {
				return nil, err
			}
			resources["is_email_verified"] = isEmailVerified

		case constants.READ_USER_DISPLAY_NAME_SCOPE:
			displayNameMetaData, err := u.Repo.GetUserMetadata(claims.UserID, constants.USER_DISPLAY_NAME_METADATA_KEY)
			if err != nil {
				return nil, err
			}
			resources["display_name"] = displayNameMetaData

		case constants.READ_USER_COUNTRY_SCOPE:
			countryMetaData, err := u.Repo.GetUserMetadata(claims.UserID, constants.USER_COUNTRY_METADATA_KEY)
			if err != nil {
				return nil, err
			}
			resources["country_name"] = countryMetaData
		}
	}
	fmt.Println("resources check:", resources)
	return resources, nil
}

func (u *Service) VerifyToken(token string) (isvalid bool, err error) {
	// verify token
	jwtClaims, err := utilities.VerifyStandardToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return false, err
	}
	// check if the token is expired
	if jwtClaims.ExpiresAt < time.Now().Unix() {
		return false, fmt.Errorf("token is expired")
	}
	return true, nil
}

func (u *Service) GetClient(id string) (*models.Client, error) {
	client, err := u.Repo.GetClient(fmt.Sprintf("client:%s", id))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (u *Service) CheckClient(backendURL string) (*models.Client, error) {
	client, err := u.Repo.GetClientByURL(backendURL)
	if err != nil {
		return nil, fmt.Errorf("could not get client: %v", err)
	}
	if client == nil {
		return nil, fmt.Errorf("client not found")
	}
	return client, nil
}

func (u *Service) SaveX509Certificate(clientID string, certificate string) error {
	return u.Repo.SaveX509Certificate(clientID, certificate)
}

func (u *Service) DecodeAuthorizationCode(secret string, code string) (*utilities.AuthCodeClaims, error) {
	// Decode the auth code
	// verify the code
	claims, err := utilities.DecodeAuthCode(secret, code)
	if err != nil {
		return nil, fmt.Errorf("failed to decode auth code: %v", err)
	}
	return claims, nil
}

func (u *Service) AuthenticateClient(uuid string, secret string) error {
	client, err := u.Repo.GetClient(uuid)
	if err != nil {
		return fmt.Errorf("failed to authenticate client: %e", err)
	}

	if client.Secret != secret {
		return fmt.Errorf("failed to authenticate client: provided secret value is invalid")
	}

	return nil
}

func (u *Service) GenerateAccessToken(
	authClaims *utilities.AuthCodeClaims,
	clientID string,
	clientSecret string,
) (string, error) {
	claims := entities.ClientClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Globe and Citizen",
			IssuedAt:  time.Now().UTC().Unix(),
			Subject:   clientID,
			ExpiresAt: time.Now().Add(constants.AccessTokenValidityMinutes * time.Minute).UTC().Unix(),
		},
		Scopes: authClaims.Scopes,
		UserID: authClaims.UserID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (u *Service) ValidateAccessToken(clientSecret string, accessToken string) (*entities.ClientClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&entities.ClientClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(clientSecret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("jwt token is invalid")
	}

	claims, ok := token.Claims.(*entities.ClientClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse client claims")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, fmt.Errorf("access token is expired")
	}

	return claims, nil
}

func (u *Service) GetZkUserMetadata(scopesStr string, userID int64) (*entities.ZkMetadataResponse, error) {
	if scopesStr == "" {
		return &entities.ZkMetadataResponse{}, fmt.Errorf("no access scopes granted")
	}

	scopes := strings.Split(scopesStr, ",")

	var zkMetadata entities.ZkMetadataResponse

	for _, scope := range scopes {
		switch scope {
		case "read:user:country":
			countryMetadata, err := u.Repo.GetUserMetadata(userID, constants.USER_COUNTRY_METADATA_KEY)
			if err != nil {
				return &entities.ZkMetadataResponse{}, fmt.Errorf("failed to get country metadata: %e", err)
			}

			zkMetadata.Country = countryMetadata.Value
		// case "email_verified":
		// 	emailMetadata, err := u.Repo.GetUserMetadata(userID, constants.USER_EMAIL_VERIFIED_METADATA_KEY)
		// 	if err != nil {
		// 		return &entities.ZkMetadataResponse{}, fmt.Errorf("failed to get email metadata: %e", err)
		// 	}

		// 	zkMetadata.IsEmailVerified = emailMetadata.Value == "true"
		case "read:user:display_name":
			displayNameMetadata, err := u.Repo.GetUserMetadata(userID, constants.USER_DISPLAY_NAME_METADATA_KEY)
			if err != nil {
				return &entities.ZkMetadataResponse{}, fmt.Errorf("failed to get display name metadata: %e", err)
			}

			zkMetadata.DisplayName = displayNameMetadata.Value
		case "read:user:color":
			// TODO: implement
			zkMetadata.Color = "red"
		}
	}

	emailMetadata, err := u.Repo.GetUserMetadata(userID, constants.USER_EMAIL_VERIFIED_METADATA_KEY)
	if err != nil {
		return &entities.ZkMetadataResponse{}, fmt.Errorf("failed to get email metadata: %e", err)
	}

	zkMetadata.IsEmailVerified = emailMetadata.Value == "true"

	return &zkMetadata, nil
}

// this is only be used for testing purposes
func (u *Service) AddTestClient() (*models.Client, error) {
	rmSalt := rs_utils.GenerateRandomSalt(rs_utils.SaltSize)
	client := &models.Client{
		ID:          "notanid",
		Secret:      "absolutelynotasecret!",
		Name:        "Ex-C",
		RedirectURI: "http://localhost:5173/oauth2/callback",
		BackendURI:  os.Getenv("TEST_CLIENT_BACKEND_URL"),
		Username:    "layer8",
		Salt:        rmSalt,
		// BackendURI:  "localhost:8000",
	}

	err := u.Repo.SetClient(client)
	if err != nil {
		return nil, err
	}
	return client, nil
}
