package service

import (
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/constants"
	"globe-and-citizen/layer8/server/entities"
	"globe-and-citizen/layer8/server/internals/repository"
	"globe-and-citizen/layer8/server/models"
	"globe-and-citizen/layer8/server/utils"
	"strings"
	"time"

	utilities "github.com/globe-and-citizen/layer8-utils"
	"golang.org/x/oauth2"
)

type Service struct {
	Repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
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

	userSalt, err := u.Repo.LoginUserPrecheck(username)
	if err != nil {
		return nil, err
	}

	HashedAndSaltedPass := utils.SaltAndHashPassword(password, userSalt)

	user, err := u.Repo.GetUser(username)
	if err != nil {
		return nil, err
	}

	if user.Password != HashedAndSaltedPass {
		return nil, fmt.Errorf("invalid password")
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
	// get the resources
	scopes := strings.Split(claims.Scopes, ",")
	resources := make(map[string]interface{})
	for _, scope := range scopes {
		switch scope {
		case constants.READ_USER_SCOPE:
			// 	user, err := u.Repo.GetUserByID(claims.UserID)
			// 	if err != nil {
			// 		return nil, err
			// 	}
			// 	resources["profile"] = user

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

func (u *Service) GetClient(id string) (*models.Client, error) {
	client, err := u.Repo.GetClient(fmt.Sprintf("client:%s", id))
	if err != nil {
		return nil, err
	}

	return client, nil
}

// this is only be used for testing purposes
func (u *Service) AddTestClient() (*models.Client, error) {
	client := &models.Client{
		ID:          "notanid",
		Secret:      "absolutelynotasecret!",
		Name:        "Ex-C",
		RedirectURI: "http://localhost:5173/oauth2/callback",
	}

	err := u.Repo.SetClient(client)
	if err != nil {
		return nil, err
	}
	return client, nil
}
