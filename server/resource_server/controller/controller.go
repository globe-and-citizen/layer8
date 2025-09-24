package controller

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	utils.ParseHTML(w, http.StatusOK,
		"assets-v1/templates/src/pages/client_portal/profile.html",
		map[string]interface{}{
			"ProxyURL":             os.Getenv("PROXY_URL"),
			"SmartContractAddress": os.Getenv("SMART_CONTRACT_ADDRESS"),
		},
	)
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

func ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/src/pages/user_portal/password_reset/reset-password-page.html")
}

func VerifyPhoneNumberPage(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/src/pages/user_portal/cellphone/verify-phone-number.html")
}

func InputPhoneNumberVerificationCodePage(w http.ResponseWriter, r *http.Request) {
	ServeFileHandler(w, r, "assets-v1/templates/src/pages/user_portal/cellphone/input-verification-code.html")
}

func ServeFileHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	utils.ParseHTML(w, http.StatusOK, filePath, map[string]interface{}{
		"ProxyURL": os.Getenv("PROXY_URL"),
	})
}

func LoginClientHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.LoginClientDTO](w, r.Body)
	if err != nil {
		return
	}

	serverSignatureResp, err := newService.LoginClient(request)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to perform login", err)
		return
	}

	response := utils.BuildResponse(w, http.StatusOK, "Login successful", serverSignatureResp)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func RegisterClientPrecheckHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService, ok := r.Context().Value("service").(interfaces.IService)
	if !ok {
		utils.HandleError(w, http.StatusInternalServerError, "Service not found in context", nil)
		return
	}

	iterCount, err := strconv.Atoi(os.Getenv("SCRAM_ITERATION_COUNT"))
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Invalid iteration count configuration", err)
		return
	}

	request, err := utils.DecodeJsonFromRequest[dto.RegisterClientPrecheckDTO](w, r.Body)
	if err != nil {
		return
	}

	salt, err := newService.RegisterClientPrecheck(request, iterCount)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to register user", err)
		return
	}

	registerClientPrecheckResp := models.RegisterClientPrecheckResponseOutput{
		Salt:           salt,
		IterationCount: iterCount,
	}

	resp := utils.BuildResponse(w, http.StatusCreated, "Client is successfully registered", registerClientPrecheckResp)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal Server Error",
			err,
		)
	}
}

func LoginClientPrecheckHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.LoginPrecheckDTO](w, r.Body)
	if err != nil {
		return
	}

	loginPrecheckResp, err := newService.LoginPrecheckClient(request)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to perform precheck, service error", err)
		return
	}

	response := utils.BuildResponse(w, http.StatusOK, "Precheck successful", loginPrecheckResp)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func ClientProfileHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodGet) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:]
	clientClaims, err := utils.ValidateClientToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	profileResp, err := newService.ProfileClient(clientClaims.UserName)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile, user not found", err)
		return
	}

	if err := json.NewEncoder(w).Encode(profileResp); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func RegisterClientHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

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

	res := utils.BuildResponseWithNoBody(w, http.StatusCreated, "Client registered successfully")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func LoginPrecheckHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.LoginPrecheckDTO](w, r.Body)
	if err != nil {
		return
	}

	loginPrecheckResp, err := newService.LoginPrecheckUser(request)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to perform precheck, service error", err)
		return
	}

	response := utils.BuildResponse(w, http.StatusOK, "Precheck successful", loginPrecheckResp)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.LoginUserDTO](w, r.Body)
	if err != nil {
		return
	}

	serverSignatureResp, err := newService.LoginUser(request)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to perform login", err)
		return
	}

	response := utils.BuildResponse(w, http.StatusOK, "Login successful", serverSignatureResp)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodGet) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	profileResp, err := newService.ProfileUser(userID)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get user profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(profileResp); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func GetClientData(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodGet) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)
	clientName := r.Header.Get("Name")

	clientModel, err := newService.GetClientData(clientName)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client profile", err)
		return
	}

	if err := json.NewEncoder(w).Encode(clientModel); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Authentication error: invalid token", err)
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

	response := utils.BuildResponseWithNoBody(w, http.StatusOK, "Verification email sent")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Internal error happened", err)
	}
}

func CheckEmailVerificationCode(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	service := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Authentication error: invalid token", err)
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

	user, err := service.FindUser(userID)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "User with provided id does not exist", err)
		return
	}

	zkProof, zkKeyPairId, err := service.GenerateZkProof(user, request.Email, request.Code)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to generate zk proof of email verification", err)
		return
	}

	err = service.SaveProofOfEmailVerification(userID, request.Code, zkProof, zkKeyPairId)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to save proof of the email verification procedure", err)
		return
	}

	response := utils.BuildResponseWithNoBody(w, http.StatusOK, "Your email was successfully verified!")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Internal error happened", err)
	}
}

func UpdateUserMetadataHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	request, err := utils.DecodeJsonFromRequest[dto.UpdateUserMetadataDTO](w, r.Body)
	if err != nil {
		return
	}

	err = newService.UpdateUserMetadata(userID, request)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to update user's metadata", err)
		return
	}

	resp := utils.BuildResponseWithNoBody(w, http.StatusOK, "User's metadata updated successfully")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func GetUsageStats(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodGet) {
		return
	}

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

	resp := utils.BuildResponse(w, http.StatusOK, "Client usage statistics", finalResponse)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func CheckBackendURI(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

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
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func RegisterUserPrecheck(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService, ok := r.Context().Value("service").(interfaces.IService)
	if !ok {
		utils.HandleError(w, http.StatusInternalServerError, "Service not found in context", nil)
		return
	}

	iterCount, err := strconv.Atoi(os.Getenv("SCRAM_ITERATION_COUNT"))
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Invalid iteration count configuration", err)
		return
	}

	request, err := utils.DecodeJsonFromRequest[dto.RegisterUserPrecheckDTO](w, r.Body)
	if err != nil {
		return
	}

	salt, err := newService.RegisterUserPrecheck(request, iterCount)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to register user", err)
		return
	}

	registerUserPrecheckResp := models.RegisterUserPrecheckResponseOutput{
		Salt:           salt,
		IterationCount: iterCount,
	}

	resp := utils.BuildResponse(w, http.StatusCreated, "User is successfully registered", registerUserPrecheckResp)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal Server Error",
			err,
		)
	}
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

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

	res := utils.BuildResponseWithNoBody(w, http.StatusCreated, "User registered successfully")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func ResetPasswordPrecheckHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.ResetPasswordPrecheckDTO](w, r.Body)
	if err != nil {
		return
	}

	user, err := newService.GetUserForUsername(request.Username)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "User does not exist!", err)
		return
	}

	resetPasswordPrecheckResp := models.ResetPasswordPrecheckResponseOutput{
		Salt:           user.Salt,
		IterationCount: user.IterationCount,
	}

	response := utils.BuildResponse(w, http.StatusOK, "User does exist!", resetPasswordPrecheckResp)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal Server Error",
			err,
		)
	}
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.ResetPasswordDTO](w, r.Body)
	if err != nil {
		return
	}

	user, err := newService.GetUserForUsername(request.Username)
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "User does not exist!", err)
		return
	}

	err = newService.ValidateSignature("Sign-in with Layer8", request.Signature, user.PublicKey)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Signature is invalid!", err)
		return
	}

	err = newService.UpdateUserPassword(user.Username, request.StoredKey, request.ServerKey)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Internal error: failed to update user", err)
		return
	}

	response := utils.BuildResponseWithNoBody(
		w, http.StatusCreated, "Your password was updated successfully!",
	)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to encode the response", err)
	}
}

func ClientUnpaidAmountHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)

	request, err := utils.DecodeJsonFromRequest[dto.ClientUnpaidAmountDTO](w, r.Body)
	if err != nil {
		return
	}

	unpaidAmount, err := newService.GetClientUnpaidAmount(request.ClientId)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Failed to get client's unpaid amount", err)
		return
	}

	clientUnpaidAmountResponseOutput := models.ClientUnpaidAmountResponseOutput{
		UnpaidAmount: unpaidAmount,
	}

	response := utils.BuildResponse(
		w,
		http.StatusOK,
		"successfully retrieved client's unpaid amount",
		clientUnpaidAmountResponseOutput,
	)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to encode the response", err)
	}
}

func GenerateTelegramSessionIDHandler(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)

	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	sessionID, err := newService.GenerateTelegramSessionID()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "failed to generate session id", err)
		return
	}

	err = newService.SaveTelegramSessionID(userID, sessionID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "failed to save the Telegram session id", err)
		return
	}

	sessionIdDTO := dto.TelegramSessionIdDTO{
		SessionID: base64.RawURLEncoding.EncodeToString(sessionID),
	}

	response := utils.BuildResponse(w, http.StatusOK, "session id generated", sessionIdDTO)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.HandleError(
			w,
			http.StatusInternalServerError,
			"Internal error: could not encode response into json",
			err,
		)
	}
}

func VerifyPhoneNumberViaTelegramBotHandler(w http.ResponseWriter, r *http.Request) {
	newService := r.Context().Value("service").(interfaces.IService)

	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	user, err := newService.FindUser(userID)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "User not found for ID", err)
	}

	telegramAPIToken := os.Getenv("TELEGRAM_API_KEY")
	if telegramAPIToken == "" {
		log.Fatal("No Telegram API token provided")
	}

	baseURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s",
		url.PathEscape(telegramAPIToken),
	)

	var offset int64 = 0
	var telegramUserID int64 = 0

	for i := 0; i < 500; i++ {
		if i > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		updates, err := newService.RefreshTelegramMessages(baseURL, offset)
		if err != nil {
			log.Printf("failed to refresh messages from Telegram: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		foundStartMessage := false

		for _, u := range updates {
			offset = u.UpdateID + 1

			if u.Message == nil || u.Message.Chat.Type != "private" {
				continue
			}

			msg := u.Message

			if !strings.HasPrefix(msg.Text, "/start ") {
				continue
			}

			sessionID := msg.Text[7:]
			if sessionID == "" {
				continue
			}

			sessionIDBytes, err := base64.RawURLEncoding.DecodeString(sessionID)
			if err != nil {
				log.Printf("failed to decode session id %s, skipping\n", sessionID)
				continue
			}

			sessionIDHash := sha256.Sum256(sessionIDBytes)

			if bytes.Equal(user.TelegramSessionIDHash, sessionIDHash[:]) {
				log.Printf("Found matching session id\n")

				// Send a reply keyboard that *requests contact*.
				kb := dto.ReplyKeyboardMarkup{
					Keyboard: [][]dto.KeyboardButton{
						{
							{Text: "Share my phone", RequestContact: true},
						},
					},
					ResizeKeyboard:  true,
					OneTimeKeyboard: true,
				}
				text := "Hi! In order for us to verify your phone number, please tap the button below to allow Telegram sharing your phone number with us."

				err := newService.SendTelegramBotMessage(baseURL, dto.SendMessageRequestDTO{
					ChatID:      msg.Chat.ID,
					Text:        text,
					ParseMode:   "Markdown",
					ReplyMarkup: kb,
				})
				if err != nil {
					log.Printf("sendMessage error: %v", err)
				}

				foundStartMessage = true
				telegramUserID = msg.From.ID
				break
			}
		}

		if foundStartMessage {
			break
		}
	}

	for i := 0; i < 500; i++ {
		if i > 0 {
			time.Sleep(400 * time.Millisecond)
		}

		updates, err := newService.RefreshTelegramMessages(baseURL, offset)
		if err != nil {
			log.Printf("failed to refresh messages from Telegram: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, u := range updates {
			offset = u.UpdateID + 1

			if u.Message == nil || u.Message.Chat.Type != "private" {
				continue
			}

			msg := u.Message

			if msg.Contact == nil {
				continue
			}

			// Validate the contact belongs to the sender.
			if msg.From == nil || msg.Contact.UserID != msg.From.ID || msg.From.ID != telegramUserID {
				continue
			}

			phoneNumber := strings.TrimSpace(msg.Contact.PhoneNumber)

			verificationCode, err := newService.GeneratePhoneNumberVerificationCode(&user, phoneNumber)
			if err != nil {
				utils.HandleError(w, http.StatusInternalServerError, "failed to generate the verification code", err)
				return
			}

			// Remove keyboard after success.
			err = newService.SendTelegramBotMessage(baseURL, dto.SendMessageRequestDTO{
				ChatID: msg.Chat.ID,
				Text: fmt.Sprintf(
					"Thanks! Your verification code is: %s. You can go back to the Layer8 user portal now.",
					verificationCode,
				),
				ReplyMarkup: dto.ReplyKeyboardRemove{
					RemoveKeyboard: true,
				},
			})
			if err != nil {
				utils.HandleError(w, http.StatusInternalServerError, "failed to send message from the Telegram bot", err)
				log.Printf("failed to send message from the Telegram bot: %v", err)
				return
			}

			zkProof, zkPairID, err := newService.GenerateZkProof(user, phoneNumber, verificationCode)
			if err != nil {
				utils.HandleError(w, http.StatusInternalServerError, "failed to generate the zk proof of phone number verification", err)
				return
			}

			err = newService.SavePhoneNumberVerificationData(user.ID, verificationCode, zkProof, zkPairID)
			if err != nil {
				utils.HandleError(w, http.StatusInternalServerError, "failed to save proof of the phone number verification into the db", err)
				return
			}

			log.Println("Phone number successfully verified, exiting")

			apiResponse := utils.BuildResponseWithNoBody(w, http.StatusOK, "phone number is verified")

			if err := json.NewEncoder(w).Encode(apiResponse); err != nil {
				utils.HandleError(
					w,
					http.StatusInternalServerError,
					"Internal error: could not encode response into json",
					err,
				)
			}
			return
		}
	}
}

func CheckPhoneNumberVerificationCode(w http.ResponseWriter, r *http.Request) {
	if !validateHttpMethod(w, r.Method, http.MethodPost) {
		return
	}

	newService := r.Context().Value("service").(interfaces.IService)

	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	request, err := utils.DecodeJsonFromRequest[dto.CheckPhoneNumberVerificationCodeDTO](w, r.Body)
	if err != nil {
		return
	}

	verificationData, err := newService.GetPhoneNumberVerificationData(userID)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "user's verification data not found", err)
		return
	}

	err = newService.CheckPhoneNumberVerificationCode(request.VerificationCode, verificationData)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "failed to validate the provided verification code", err)
		return
	}

	err = newService.SaveProofOfPhoneNumberVerification(verificationData)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "failed to update phone number verification metadata in the db", err)
		return
	}

	response := utils.BuildResponseWithNoBody(w, http.StatusOK, "Your phone number is verified successfully! Congratulations!")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "failed to encode response into json", err)
	}
}

func validateHttpMethod(w http.ResponseWriter, actualMethod string, expectedMethod string) bool {
	if actualMethod != expectedMethod {
		errorMessage := fmt.Sprintf("Invalid http method. Expected %s", expectedMethod)
		utils.HandleError(
			w,
			http.StatusMethodNotAllowed,
			errorMessage,
			fmt.Errorf(errorMessage),
		)
		return false
	}

	return true
}
