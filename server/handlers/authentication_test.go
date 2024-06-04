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
)

/* TESTS THAT I WILL NEED TO WRITE
* 1) Query param next == ""
* 2) Query param next == "<?>"
* 3) request.Cookie("token") != nil
* 4) etc...
 */

func Test_Login(t *testing.T) {
	// Step x: Create a mock service provider if necessary
	mockedServiceProvider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		/*
		* Custom logic that you would
		* like the mock service provider to
		* accomplish.
		 */

		w.Write([]byte("b64PubJWK"))
	}))

	defer mockedServiceProvider.Close()

	// Step x, invoke an in_mem repository.Repository
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

	memoryService := service.NewService(memoryRepository)

	form := url.Values{}
	form.Add("username", "stravid")
	form.Add("password", "12341234")
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
	//assert.Equal(t, actual, expected)
}
