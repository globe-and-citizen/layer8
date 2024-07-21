package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"globe-and-citizen/layer8/server/resource_server/db"
	"globe-and-citizen/layer8/server/resource_server/dto"
	"globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/models"
	"globe-and-citizen/layer8/server/resource_server/repository"
	"globe-and-citizen/layer8/server/resource_server/utils"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/public/welcome.html")
}
func LoginUserPage(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/src/pages/user_portal/login.html")
}
func RegisterUserPage(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/src/pages/user_portal/register.html")
}
func ClientProfilePage(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/src/pages/client_portal/profile.html")
}
func UserHandler(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/src/pages/user_portal/profile.html")
}
func ClientHandler(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/src/pages/client_portal/register.html")
}
func LoginClientPage(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/src/pages/client_portal/login.html")
}

func InputYourEmailPage(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/src/pages/user_portal/email/input-your-email.html")
}

func InputVerificationCodePage(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/src/pages/user_portal/email/input-verification-code.html")
}

func ServeFileHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	utils.ParseHTML(w, filePath, map[string]interface{}{
		"ProxyURL": os.Getenv("PROXY_URL"),
	})
}

func LoginClientHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.LoginClientDTO](w, r.Body)
	if err != nil {
		return
	}

	tokenResp, err := newService.LoginClient(request)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(tokenResp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}
}

func ClientProfileHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:]
	clientClaims, err := utils.ValidateClientToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile, invalid token", err)
		return
	}

	profileResp, err := newService.ProfileClient(clientClaims.UserName)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile, user not found", err)
		return
	}

	if err := json.NewEncoder(w).Encode(profileResp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile, error encoding response", err)
		return
	}
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.RegisterUserDTO](w, r.Body)
	if err != nil {
		return
	}

	err = newService.RegisterUser(request)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to register user", err)
		return
	}

	res := utils.BuildResponse(true, "OK!", "User registered successfully")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to register user", err)
		return
	}
}

func RegisterClientHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.RegisterClientDTO](w, r.Body)
	if err != nil {
		return
	}

	err = newService.RegisterClient(request)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to register client", err)
		return
	}

	res := utils.BuildResponse(true, "OK!", "Client registered successfully")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to register client", err)
		return
	}
}

// LoginPrecheckHandler handles login precheck requests and get the salt of the user from the database using the username from the request URL
func LoginPrecheckHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.LoginPrecheckDTO](w, r.Body)
	if err != nil {
		return
	}

	loginPrecheckResp, err := newService.LoginPreCheckUser(request)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(loginPrecheckResp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.LoginUserDTO](w, r.Body)
	if err != nil {
		return
	}

	tokenResp, err := newService.LoginUser(request)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(tokenResp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile", err)
		return
	}

	profileResp, err := newService.ProfileUser(userID)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(profileResp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile", err)
		return
	}
}

func GetClientData(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	clientName := r.Header.Get("Name")

	clientModel, err := newService.GetClientData(clientName)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(clientModel); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}
}

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Request failed: invalid authorization token", err)
		return
	}

	request, err := utils.DecodeJsonFromRequest[dto.VerifyEmailDTO](w, r.Body)
	if err != nil {
		return
	}

	err = newService.VerifyEmail(userID, request.Email)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to verify email", err)
		return
	}

	response := utils.BuildResponse(true, "OK!", "Verification email sent")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Internal error happened", err)
	}
}

func CheckEmailVerificationCode(w http.ResponseWriter, r *http.Request) {
	service := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to verify user's token", err)
		return
	}

	request, err := utils.DecodeJsonFromRequest[dto.CheckEmailVerificationCodeDTO](w, r.Body)
	if err != nil {
		return
	}

	err = service.CheckEmailVerificationCode(userID, request.Code)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to verify code", err)
		return
	}

	user, e := service.FindUser(userID)
	if e != nil {
		utils.HandleError(w, http.StatusBadRequest, "User with provided id does not exist", e)
		return
	}

	zkProof, e := service.GenerateZkProofOfEmailVerification(user, request)
	if e != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to generate zk proof of email verification", e)
		return
	}

	e = service.SaveProofOfEmailVerification(userID, request.Code, zkProof)
	if e != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to save proof of the email verification procedure", e)
		return
	}

	response := utils.BuildResponse(
		true,
		"Your email was successfully verified!",
		"Email verified!",
	)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Internal error happened", err)
	}
}

func UpdateDisplayNameHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to update display name", err)
		return
	}

	request, err := utils.DecodeJsonFromRequest[dto.UpdateDisplayNameDTO](w, r.Body)
	if err != nil {
		return
	}

	err = newService.UpdateDisplayName(userID, request)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to update display name", err)
		return
	}

	resp := utils.BuildResponse(true, "OK!", "Display name updated successfully")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to update display name", err)
		return
	}
}

func GetUsageStats(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		utils.HandleError(w, http.StatusUnauthorized, "failed to show client usage statistics", errors.New("missing jwt token"))
		return
	}

	authToken = authToken[7:]
	clientClaims, err := utils.ValidateClientToken(authToken)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "failed to show client usage statistics", errors.New("jwt token invalid"))
		return
	}

	statRepo := repository.NewStatRepository(db.GetInfluxDBClient())

	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayOfNextMonth := time.Date(firstDayOfMonth.Year(), firstDayOfMonth.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfCurrentMonth := firstDayOfNextMonth.Add(-24 * time.Hour)
	totalDaysInMonth := lastDayOfCurrentMonth.Day()
	totalDaysBeforeNextMonth := totalDaysInMonth - now.Day()

	thirtyDaysStatistic, err := statRepo.GetTotalRequestsInLastXDaysByClient(r.Context(), 30, clientClaims.ClientID)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get last thrthy days usage statistic", err)
		return
	}

	monthToDateTotal, err := statRepo.GetTotalByDateRangeByClient(r.Context(), firstDayOfMonth, firstDayOfNextMonth, clientClaims.ClientID)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get month to date usage statistic", err)
		return
	}

	finalResponse := models.UsageStatisticResponse{
		MonthToDate: models.MonthToDateStatistic{
			Month: firstDayOfMonth.Month().String(),
		},
		LastThirtyDaysStatistic: thirtyDaysStatistic,
		MetricType:              "data_transferred",
		UnitOfMeasurement:       "GB",
	}

	if monthToDateTotal > 0 {
		finalResponse.MonthToDate.MonthToDateUsage = monthToDateTotal / 1000000000
		finalResponse.MonthToDate.ForecastedEndOfMonthUsage = (monthToDateTotal / 1000000000) + float64(totalDaysBeforeNextMonth)*thirtyDaysStatistic.Average
	}

	resp := utils.BuildResponse(true, "OK!", finalResponse)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get stats", err)
		return
	}
}

func CheckBackendURI(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.CheckBackendURIDTO](w, r.Body)
	if err != nil {
		return
	}

	response, err := newService.CheckBackendURI(request.BackendURI)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to check backend url", err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to check backend url", err)
		return
	}
}
