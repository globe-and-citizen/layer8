package handlers

import (
	"context"
	"globe-and-citizen/layer8/server/internals/service"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/repository"
	"net/http"
	"net/http/httptest"
	"net/url"

	"testing"

	"github.com/stretchr/testify/assert"
)

/* TESTS THAT I WILL NEED TO WRITE
* 1) Query param next == ""
* 2) Query param next == "<?>"
* 3) request.Cookie("token") != nil
*
 */

func Test_Login(t *testing.T) {
	// Step 1: Create a mock service provider
	mockedServiceProvider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		/*
		* Custom logic that you would
		* like the mock service provider to
		* accomplish.
		 */

		w.Write([]byte("b64PubJWK"))
	}))

	defer mockedServiceProvider.Close()

	// Step 1, invoke an in_mem repository.Repository
	memoryRepository := repository.NewMemoryRepository()
	memoryRepository.RegisterUser(dto.RegisterUserDTO{
		Email:       "auth@test.com",
		Username:    "tester_chester",
		FirstName:   "Tester",
		LastName:    "Chester",
		Password:    "12341234",
		Country:     "North Pole",
		DisplayName: "testa_chesta",
	})

	// In the actual program, memoryService would now be passed to the server
	memoryService := service.NewService(memoryRepository)

	form := url.Values{}
	form.Add("username", "stravid")
	form.Add("password", "12341234")
	//formBytes := []byte(form.Encode())
	//formAsIOReader := bytes.NewReader(formBytes)
	t.Log("why", form.Encode())
	// Create a mockedRequest
	reqToAuthentication := httptest.NewRequest("POST", "http://?"+form.Encode(), nil)

	// Add the in_memory service to a request using context
	reqToAuthentication = reqToAuthentication.WithContext(context.WithValue(reqToAuthentication.Context(), "Oauthservice", memoryService))

	// the response should be this html "assets-v1/templates/src/pages/oauth_portal/login.html"
	responseRecorder := httptest.NewRecorder()

	// NEXT UP: Dress up the request and change it to exercise all branches of this function and assert against them.
	// Unit undertest

	Login(responseRecorder, reqToAuthentication)

	// Run assertions on the recorded response.
	t.Log(responseRecorder.Code)
	t.Log(responseRecorder.Body)
	assert.Equal(t, responseRecorder.Code, 500)
}

// type StubbedRepository struct{}

// // service.GetUserByToken(token.Value) uses GetUserByID
// // GetUserByID gets a user by ID.
// func (ss StubbedRepository) GetUserByID(id int64) (*models.User, error) {
// 	return &models.User{
// 		ID:        1,
// 		Email:     "stubbedEmail@gmail.com",
// 		Username:  "tester",
// 		Password:  "12341234",
// 		FirstName: "Stubby",
// 		LastName:  "McGee",
// 		Salt:      "f23f113949201b90aaf6d634e3d5f5788fb3b708cc736b665c3b726f73414aae",
// 	}, nil
// }

// // service.LoginUser(username, password) uses:
// // u.Repo.LoginUserPrecheck(username)
// // u.Repo.GetUser(username)
// func (ss StubbedRepository) LoginUserPrecheck(username string) (string, error) {
// 	return "salty12341234", nil
// }

// // Get user from db by username
// func (ss StubbedRepository) GetUser(username string) (*models.User, error) {
// 	return &models.User{
// 		ID:        1,
// 		Email:     "stubbedEmail@gmail.com",
// 		Username:  "stubbyMcGee",
// 		Password:  "12341234",
// 		FirstName: "Stubby",
// 		LastName:  "McGee",
// 		Salt:      "salty12341234",
// 	}, nil
// }

// /* unnecessary */

// // GetUserMetadata gets a user metadata by key.
// func (ss StubbedRepository) GetUserMetadata(userID int64, key string) (*models.UserMetadata, error) {
// 	return nil, fmt.Errorf("Stubbed Repository. GetUserMetadata will always return this error")
// }

// // Set a client for testing purposes
// func (ss StubbedRepository) SetClient(client *models.Client) error {
// 	return fmt.Errorf("Stubbed Repository. SetClient will always return this error")
// }

// // Get a client by ID.
// func (ss StubbedRepository) GetClient(id string) (*models.Client, error) {
// 	return nil, fmt.Errorf("Stubbed Repository. GetClient will always return this error")
// }

// // SetTTL sets the value for the given key with a short TTL.
// func (ss StubbedRepository) SetTTL(key string, value []byte, ttl time.Duration) error {
// 	return fmt.Errorf("Stubbed Repository. SetTTL will always return this error")
// }

// // GetTTL gets the value for the given key which has a short TTL.
// func (ss StubbedRepository) GetTTL(key string) ([]byte, error) {
// 	return []byte{}, fmt.Errorf("Stubbed Repository. GetTTL will always return this error")
// }
