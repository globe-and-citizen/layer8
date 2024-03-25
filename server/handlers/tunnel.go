package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"globe-and-citizen/layer8/server/resource_server/utils"

	utilities "github.com/globe-and-citizen/layer8-utils"
)

// Tunnel forwards the request to the service provider's backend
func InitTunnel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\n*************")
	fmt.Println(r.Method) // > GET  | > POST
	fmt.Println(r.URL)    // (http://localhost:5000/api/v1 ) > /api/v1
	params := r.URL.Query()
	protocol := r.Header.Get("X-Forwarded-Proto")
	var backend string
	if _, ok := params["backend"]; !ok {
		res := utils.BuildErrorResponse("Failed to get User. Malformed query string.", "", utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	} else {
		backend = params["backend"][0]
		backend = protocol + backend
		fmt.Println("User agent is attempting to initialize this backend SP: ", backend)
	}

	mpJWT, err := utilities.GenerateStandardToken(os.Getenv("MP_123_SECRET_KEY"))
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	reqData := utilities.ReadResponseBody(r.Body)
	b64PubJWK := string(reqData)
	fmt.Println("b64PubJWK: ", b64PubJWK)
	fmt.Println("x-ecdh-init: ", r.Header.Get("x-ecdh-init"))

	// create the request
	req, err := http.NewRequest(r.Method, backend, r.Body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// add headers
	for k, v := range r.Header {
		req.Header[k] = v
	}

	req.Header["x-tunnel"] = []string{"true"}
	req.Header["mp-jwt"] = []string{mpJWT}

	// send the request
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make a buffer to hold response body
	var resBodyTemp bytes.Buffer

	// Copy the response body to buffer
	resBodyTemp.ReadFrom(res.Body)

	// Convert resBodyTemp to []byte

	resBodyTempBytes := resBodyTemp.Bytes()

	// Make a copy of the response body to send back to client
	res.Body = io.NopCloser(bytes.NewBuffer(resBodyTemp.Bytes()))

	fmt.Println("\nReceived response from 8000:", backend, " of code: ", res.StatusCode)

	upJWT, err := utilities.GenerateStandardToken(os.Getenv("UP_999_SECRET_KEY"))
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server_pubKeyECDH, err := utilities.B64ToJWK(string(resBodyTempBytes))
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make a json response of server_pubKeyECDH and up_JWT and send it back to client
	data := map[string]interface{}{
		"server_pubKeyECDH": server_pubKeyECDH,
		"up-JWT":            upJWT,
	}

	fmt.Println("Data returning to the user from the Service Provider: ", data)

	datatoSend, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(datatoSend)

}

func Tunnel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\n*************")
	fmt.Println(r.Method) // > GET  | > POST
	fmt.Println(r.URL)    // (http://localhost:5000/api/v1 ) > /api/v1
	fmt.Println("Protocol: ", r.Header.Get("X-Forwarded-Proto"))
	fmt.Println("Host: ", r.Header.Get("X-Forwarded-Host"))

	protocol := r.Header.Get("X-Forwarded-Proto")
	host := r.Header.Get("X-Forwarded-Host")
	backendURL := fmt.Sprintf("%s://%s", protocol, host+r.URL.Path) // RAVI

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("length of byte slice returning S.P. => User: ", len(bodyBytes))

	// Create the request
	buff := bytes.NewBuffer(bodyBytes)
	req, err := http.NewRequest(r.Method, backendURL, buff)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// add headers
	for k, v := range r.Header {
		req.Header[k] = v
		fmt.Println("header pairs from client (Interceptor): ", k, v)
	}
	req.Header["x-tunnel"] = []string{"true"}

	// Get up-JWT from request header and verify it
	upJWT := r.Header.Get("up-jwt") // RAVI! LOOK HERE
	fmt.Println("up-jwt coming from client: ", upJWT)

	_, err = utilities.VerifyStandardToken(upJWT, os.Getenv("UP_999_SECRET_KEY"))
	if err != nil {
		fmt.Println("UP JWT verify error: ", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// send the request
	res, err := http.DefaultClient.Do(req) // Source of MapB Error

	if err != nil {
		fmt.Println("Error sending request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("\nReceived response from:", backendURL, " of code: ", res.StatusCode)

	// Get mp-JWT from response header and verify it
	mpJWT := res.Header.Get("mp-jwt")
	fmt.Println("mp-jwt coming from SP: ", mpJWT)

	_, err = utilities.VerifyStandardToken(mpJWT, os.Getenv("MP_123_SECRET_KEY"))
	if err != nil {
		fmt.Println("MP JWT verify error: ", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// copy response back
	for k, v := range res.Header {
		w.Header()[k] = v
		fmt.Println("header pairs from SP: ", k, v)
	}

	//w.WriteHeader(res.StatusCode)
	n, err := io.Copy(w, res.Body)
	if err != nil {
		fmt.Println("Error copying response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Copied", n, "bytes from response body to client")
	fmt.Println("w.Headers 2: ", w.Header())
}

func TestError(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("this is a test error")
	fmt.Println("Test error endpoint:", err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
