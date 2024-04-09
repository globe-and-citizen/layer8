package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"globe-and-citizen/layer8/server/config"
	"globe-and-citizen/layer8/server/handlers"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	Ctl "globe-and-citizen/layer8/server/resource_server/controller"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/interfaces"

	rsRepo "globe-and-citizen/layer8/server/resource_server/repository"

	svc "globe-and-citizen/layer8/server/resource_server/service" // there are two services

	oauthSvc "globe-and-citizen/layer8/server/internals/service" // there are two services

	oauthRepo "globe-and-citizen/layer8/server/internals/repository"
)

// go:embed dist
var StaticFiles embed.FS

var workingDirectory string

func getPwd() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	workingDirectory = dir
}

func main() {

	// Use flags to set the port
	port := flag.Int("port", 8080, "Port to run the server on")
	jwtKey := flag.String("jwtKey", "secret", "Key to sign JWT tokens")
	MpKey := flag.String("MpKey", "secret", "Key to sign mpJWT tokens")
	UpKey := flag.String("UpKey", "secret", "Key to sign upJWT tokens")
	ProxyURL := flag.String("ProxyURL", "http://localhost:5001", "URL to populate go HTML templates")

	flag.Parse()

	// Use flags for using in-memory repository, otherwise app will use database
	if *port != 8080 && *jwtKey != "" && *MpKey != "" && *UpKey != "" && *ProxyURL != "" {
		os.Setenv("SERVER_PORT", strconv.Itoa(*port))
		os.Setenv("JWT_SECRET_KEY", *jwtKey)
		os.Setenv("MP_123_SECRET_KEY", *MpKey)
		os.Setenv("UP_999_SECRET_KEY", *UpKey)
		os.Setenv("PROXY_URL", *ProxyURL)
		repository := rsRepo.NewMemoryRepository()
		repository.RegisterUser(dto.RegisterUserDTO{
			Email:       "user@test.com",
			Username:    "tester",
			FirstName:   "Test",
			LastName:    "User",
			Password:    "12341234",
			Country:     "Antarctica",
			DisplayName: "test_user_mem",
		})
		service := svc.NewService(repository)
		fmt.Println("Running app with in-memory repository")
		Server(*port, service, repository) // Run server
	}

	// If the user has set a database user or password, init the database
	if os.Getenv("DB_USER") != "" || os.Getenv("DB_PASSWORD") != "" {
		config.InitDB()
	}

	proxyServerPort := os.Getenv("SERVER_PORT") // Port override

	proxyServerPortInt, err := strconv.Atoi(proxyServerPort)
	if err != nil {
		log.Fatal(err)
	}

	// Register repository
	rsRepository := rsRepo.NewRepository(config.DB)
	// Register service(usecase)
	service := svc.NewService(rsRepository)

	Server(proxyServerPortInt, service, rsRepository) // Run server (which never returns)

}

func Server(port int, service interfaces.IService, memoryRepository interfaces.IRepository) {

	// CHOOSE TO USE POSTRGRES OR IN_MEMORY IMPLEMENTATION BY COMMENTING / UNCOMMENTING

	// ** USE LOCAL POSTGRES DB **
	oauthRepository := oauthRepo.NewOauthRepository(config.DB)
	oauthService := &oauthSvc.Service{Repo: oauthRepository}

	// ** USE THE IN MEMORY IMPLEMENTATION **
	// oauthService := &oauthSvc.Service{Repo: memoryRepository}

	_, err := oauthService.AddTestClient()
	if err != nil {
		log.Fatal(err)
	}

	getPwd()

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "Oauthservice", oauthService))
			r = r.WithContext(context.WithValue(r.Context(), "service", service))

			staticFS, _ := fs.Sub(StaticFiles, "dist")
			httpFS := http.FileServer(http.FS(staticFS))

			if r.Header.Get("up-JWT") != "" {
				handlers.Tunnel(w, r)
				return
			}

			switch path := r.URL.Path; {

			// Authorization Server endpoints
			case path == "/login":
				handlers.Login(w, r)
			case path == "/authorize":
				handlers.Authorize(w, r)
			case path == "/error":
				handlers.Error(w, r)
			case path == "/api/oauth":
				handlers.OAuthToken(w, r)
			case path == "/api/user":
				handlers.UserInfo(w, r)
			case strings.HasPrefix(path, "/assets-v1"):
				http.StripPrefix("/assets-v1", http.FileServer(http.Dir("./assets-v1"))).ServeHTTP(w, r)

				// Resource Server endpoints
			case path == "/":
				Ctl.IndexHandler(w, r)
			case path == "/user":
				Ctl.UserHandler(w, r)
			case path == "/user-login-page":
				Ctl.LoginUserPage(w, r)
			case path == "/user-register-page":
				Ctl.RegisterUserPage(w, r)
			case path == "/client-register-page":
				Ctl.ClientHandler(w, r)
			case path == "/client-login-page":
				Ctl.LoginClientPage(w, r)
			case path == "/client-profile":
				Ctl.ClientProfilePage(w, r)
			case path == "/api/v1/register-user":
				Ctl.RegisterUserHandler(w, r)
			case path == "/api/v1/register-client":
				Ctl.RegisterClientHandler(w, r)
			case path == "/api/v1/getClient":
				Ctl.GetClientData(w, r)
			case path == "/api/v1/login-precheck":
				Ctl.LoginPrecheckHandler(w, r)
			case path == "/api/v1/login-user":
				Ctl.LoginUserHandler(w, r)
			case path == "/api/v1/login-client":
				Ctl.LoginClientHandler(w, r) // Login Client
			case path == "/api/v1/profile":
				Ctl.ProfileHandler(w, r)
			case path == "/api/v1/client-profile":
				Ctl.ClientProfileHandler(w, r)
			case path == "/api/v1/verify-email":
				Ctl.VerifyEmailHandler(w, r)
			case path == "/api/v1/change-display-name":
				Ctl.UpdateDisplayNameHandler(w, r)
			case path == "/favicon.ico":
				faviconPath := workingDirectory + "/dist/favicon.ico"
				http.ServeFile(w, r, faviconPath)
			case strings.HasPrefix(path, "/assets/"):
				httpFS.ServeHTTP(w, r)

			// Proxy Server endpoints
			case path == "/init-tunnel":
				handlers.InitTunnel(w, r)
			case path == "/error":
				handlers.TestError(w, r)
				// TODO: For later, to be discussed more
				// case path == "/tunnel":
				// 	handlers.Tunnel(w, r)
			}
		}),
	}
	log.Printf("Starting server on port %d...", port)
	log.Fatal(server.ListenAndServe())
}
