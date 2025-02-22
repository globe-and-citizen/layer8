package handlers

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	utilities "github.com/globe-and-citizen/layer8-utils"

	interfaces "globe-and-citizen/layer8/server/resource_server/interfaces"
	"globe-and-citizen/layer8/server/resource_server/utils"
)

var (
	meter = otel.GetMeterProvider().Meter("layer8")

	TotalByteTransferredMetrics, _ = meter.Int64Counter(
		"total_byte_transferred",
		metric.WithDescription("The total number of bytes transferred"),
	)

	TotalRequestMetrics, _ = meter.Int64Counter(
		"total_request",
		metric.WithDescription("The total number of requests"),
	)

	TotalSuccessMetrics, _ = meter.Int64Counter(
		"total_success",
		metric.WithDescription("The total number of successful requests"),
	)

	TotalTunnelInitiated, _ = meter.Int64Counter(
		"total_tunnel_initiated",
		metric.WithDescription("The total number of tunnel initiated"),
	)
)

// Tunnel forwards the request to the service provider's backend
func InitTunnel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\n*************")
	fmt.Println(r.Method) // > GET  | > POST
	fmt.Println(r.URL)    // (http://localhost:5000/api/v1 ) > /api/v1

	backend := r.URL.Query().Get("backend")
	if backend == "" {
		res := utils.BuildErrorResponse("Failed to get User. Malformed query string.", "", utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}

	backendWithoutProtocol := utils.RemoveProtocolFromURL(backend)

	srv := r.Context().Value("service").(interfaces.IService)
	client, err := srv.GetClientDataByBackendURL(backendWithoutProtocol)
	if err != nil {
		fmt.Println("Error getting client data:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	mpJWT, err := utilities.GenerateStandardToken(os.Getenv("MP_123_SECRET_KEY"))
	if err != nil {
		fmt.Println("Error generating mpJWT:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// reqData := utilities.ReadResponseBody(r.Body)
	// b64PubJWK := string(reqData)
	// fmt.Println("b64PubJWK: ", b64PubJWK)
	// fmt.Println("x-ecdh-init: ", r.Header.Get("x-ecdh-init"))

	// create the request
	req, err := http.NewRequest(r.Method, backend, r.Body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header = r.Header
	req.Header.Add("x-tunnel", "true")
	req.Header.Add("mp-jwt", mpJWT)

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

	fmt.Println("\nReceived response from backend:", backend, " of code: ", res.StatusCode)

	upJWT, err := utils.GenerateUPTokenJWT(os.Getenv("UP_999_SECRET_KEY"), client.ID)
	if err != nil {
		fmt.Println("Error generating upJWT:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Successfully generated upJWT: %s***************\n", upJWT[0:10])

	server_pubKeyECDH, err := utilities.B64ToJWK(string(resBodyTempBytes))
	// server_pubKeyECDH, err := utilities.B64ToJWK(res.Header.Get("server_pubKeyECDH"))
	if err != nil {
		fmt.Println("Error acquiring server_pubKeyECDH from Headers:", err)
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
		fmt.Println("Error marshalling data:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(datatoSend)

	TotalTunnelInitiated.Add(r.Context(), 1,
		metric.WithAttributes(
			attribute.String("client_id", client.ID),
		),
	)
}

func Tunnel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\n*************")
	fmt.Println(r.Method) // > GET  | > POST
	fmt.Println(r.URL)    // (http://localhost:5000/api/v1 ) > /api/v1
	fmt.Println("Protocol: ", r.Header.Get("X-Forwarded-Proto"))
	fmt.Println("Host: ", r.Header.Get("X-Forwarded-Host"))

	// This operation is blocking, let's confirm if the tunnel is established
	if wsIdentifier := r.Header.Get("upgrade"); strings.EqualFold(wsIdentifier, "websocket") {
		// we need to strip the host and origin headers
		r.Header.Del("Host")
		r.Header.Del("Origin")

		wsTunnel(w, r)
		return
	}

	httpTunnel(w, r)
}

type WsRoundtripEnvelope struct {
	WebSocket WsPayload
}

type WsPayload struct {
	Payload  string `json:"payload"`
	MetaData any    `json:"metadata"`
}

func wsTunnel(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Printf("error accepting websocket: %v", err)
		return
	}
	defer c.CloseNow()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// we expect a msg from the client, we already set the timeout so we don't expect blocking
	_, msg, err := c.Read(ctx)
	if err != nil {
		fmt.Printf("error reading from client: %v", err)
		return
	}

	logrus.Info("Read message from client")
	logrus.Info("Message from client: ", string(msg))

	var data WsRoundtripEnvelope
	if err = json.Unmarshal(msg, &data); err != nil {
		fmt.Printf("error unmarshalling message-: %v", err)
		return
	}

	metadata, ok := data.WebSocket.MetaData.(map[string]any)
	if !ok || metadata["backend_url"] == "" {
		fmt.Println("metadata is malformed --", data.WebSocket.MetaData)
		return
	}

	backendWithoutProtocol := utils.RemoveProtocolFromURL(metadata["backend_url"].(string))

	srv := r.Context().Value("service").(interfaces.IService)
	client, err := srv.GetClientDataByBackendURL(backendWithoutProtocol)
	if err != nil {
		fmt.Println("Error getting client data:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	mpJWT, err := utilities.GenerateStandardToken(os.Getenv("MP_123_SECRET_KEY"))
	if err != nil {
		fmt.Println("Error generating mpJWT:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serviceProviderSoc, _, err := websocket.Dial(ctx, metadata["backend_url"].(string), &websocket.DialOptions{HTTPHeader: r.Header})
	if err != nil {
		fmt.Printf("error dialing backend: %v\n", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "ws service provider closed unexpectedly")

	// initialize the handshake
	{
		initTunnelPayload, err := json.Marshal(WsRoundtripEnvelope{
			WebSocket: WsPayload{
				MetaData: map[string]interface{}{
					"x-tunnel":      true,
					"mp-jwt":        mpJWT,
					"x-ecdh-init":   metadata["x-ecdh-init"],
					"x-client-uuid": metadata["x-client-uuid"],
				},
			},
		})
		if err != nil {
			fmt.Printf("error marshalling init payload: %v\n", err)
			return
		}

		if err = serviceProviderSoc.Write(ctx, websocket.MessageText, initTunnelPayload); err != nil {
			fmt.Printf("error writing to backend: %v\n", err)
			return
		}

		// we expect an ack from the service provider if the tunnel is established
		_, msg, err := serviceProviderSoc.Read(ctx)
		if err != nil {
			fmt.Printf("error reading from service provider: %v\n", err)
			return
		}

		var backendResp WsRoundtripEnvelope
		if err = json.Unmarshal(msg, &backendResp); err != nil {
			fmt.Printf("error unmarshalling message: %v", err)
			return
		}

		metadata, ok := backendResp.WebSocket.MetaData.(map[string]any)
		if !ok {
			fmt.Println("metadata is malformed")
			return
		}

		var dataToSend []byte
		if val, ok := metadata["server_pubKeyECDH"]; !ok || val == "" {
			fmt.Println("tunnel not established; the server_pubKeyECDH is missing")
			dataToSend = msg
		} else {
			upJWT, err := utils.GenerateUPTokenJWT(os.Getenv("UP_999_SECRET_KEY"), client.ID)
			if err != nil {
				fmt.Println("Error generating upJWT:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			server_pubKeyECDH, err := utilities.B64ToJWK(metadata["server_pubKeyECDH"].(string))
			if err != nil {
				fmt.Println("Error acquiring server_pubKeyECDH from Headers:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Make a json response of server_pubKeyECDH and up_JWT and send it back to client
			data := map[string]interface{}{
				"server_pubKeyECDH": server_pubKeyECDH,
				"up-JWT":            upJWT,
			}

			fmt.Println("Data returning to the user from the Service Provider: ", data)

			if dataToSend, err = json.Marshal(&WsRoundtripEnvelope{
				WebSocket: WsPayload{
					MetaData: data,
				}}); err != nil {
				fmt.Println("Error marshalling data:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if err = c.Write(ctx, websocket.MessageText, dataToSend); err != nil {
			fmt.Printf("error writing to client: %v\n", err)
			return
		}

		fmt.Println("Tunnel established successfully")
	}

	var (
		fromServiceProvider = make(chan []byte, 5)
		fromClient          = make(chan []byte, 5)
	)

	go read_from_client(ctx, c, fromClient)
	go read_from_service_provider(serviceProviderSoc, fromServiceProvider)

	for {
		select {
		case data := <-fromClient:
			var data_ WsRoundtripEnvelope
			if err = json.Unmarshal(data, &data_); err != nil {
				fmt.Printf("error unmarshalling message: %v", err)
				return
			}

			if err = serviceProviderSoc.Write(ctx, websocket.MessageText, data); err != nil {
				if isNormal(err) {
					return
				}

				fmt.Printf("error writing to backend: %v\n", err)
				return
			}

		case data := <-fromServiceProvider:
			var data_ WsRoundtripEnvelope
			if err = json.Unmarshal(data, &data_); err != nil {
				fmt.Printf("error unmarshalling message: %v", err)
				return
			}

			if err = c.Write(ctx, websocket.MessageText, data); err != nil {
				if isNormal(err) {
					return
				}

				fmt.Printf("error writing to client: %v\n", err)
				return
			}
		}
	}
}

func read_from_client(ctx context.Context, upstream *websocket.Conn, data chan []byte) {
	for {
		_, msg, err := upstream.Read(ctx)
		if err != nil {
			fmt.Printf("error reading from client: %v", err)
			return
		}

		data <- msg
	}
}

func read_from_service_provider(service_provider *websocket.Conn, data chan []byte) {
	for {
		_, msg, err := service_provider.Read(context.Background())
		if err != nil {
			fmt.Printf("error reading from service provider: %v", err)
			return
		}

		data <- msg
	}
}

func isNormal(err error) bool {
	return websocket.CloseStatus(err) == websocket.StatusNormalClosure
}

func httpTunnel(w http.ResponseWriter, r *http.Request) {
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

	upJWTClaims, err := utils.ValidateUPTokenJWT(upJWT, os.Getenv("UP_999_SECRET_KEY"))
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

	// copy response back
	for k, v := range res.Header {
		w.Header()[k] = v
		fmt.Println("header pairs from SP: ", k, v)
	}

	w.WriteHeader(res.StatusCode)
	n, err := io.Copy(w, res.Body)
	if err != nil {
		fmt.Println("Error copying response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Copied", n, "bytes from response body to client")
	fmt.Println("w.Headers 2: ", w.Header())

	TotalRequestMetrics.Add(r.Context(), 1,
		metric.WithAttributes(
			attribute.String("client_id", upJWTClaims.Audience[0]),
		),
	)

	TotalSuccessMetrics.Add(r.Context(), 1,
		metric.WithAttributes(
			attribute.String("client_id", upJWTClaims.Audience[0]),
		),
	)

	TotalByteTransferredMetrics.Add(r.Context(), int64(binary.Size(bodyBytes)+binary.Size(w.Header())),
		metric.WithAttributes(
			attribute.String("client_id", upJWTClaims.Audience[0]),
		),
	)
}

func TestError(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("this is a test error")
	fmt.Println("Test error endpoint:", err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
