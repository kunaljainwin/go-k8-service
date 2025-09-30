// @title        My Go Service API
// @version      1.0
// @description  This is my service API.

// @contact.name  API Support
// contact.url   http://www.example.com/support
// @contact.email kunaljainwin@gmail.com

// license.name  MIT
// license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/kunaljainwin/go-service/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// structure definitions
type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
	Details   any       `json:"details,omitempty"`
}
type ReadyResponse struct {
	Hostname string `json:"hostname"`
}

// @Summary      Provides Healthcheck
// @Description  Returns a JSON of HealthResponse
// @Tags         Stateless
// @Produce      json
// @Success 200 {object} HealthResponse
// @Router       /healthcheck [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {

	var response HealthResponse
	response.Service = "service"
	response.Status = "running" //"ok","unhealthy"
	response.Timestamp = time.Now()

	statusCode := 200 // 2XX and 3XX response codes are defined as success

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// @Summary      Echoes the Request
// @Description  Returns the Request
// @Tags         Stateless
// @Success      200 {object} map[string]interface{}
// @Router       /echo [post]
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	io.Copy(w, r.Body)
}

// @Summary 	 Shared app info
// @Description  Returns the info
// @Tags 		 Stateless
// @Success 	 200
// @Router 		 /readyz [get]
func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	var response ReadyResponse
	hostname, _ := os.Hostname()
	response.Hostname = hostname

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific methods
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// Allow headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func getEnvPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
func getEnvApp() string {
	env := os.Getenv("APP_ENV")
	switch env {
	case "":
		log.Println("Running in local dev mode")
	case "dev":
		log.Println("Running in dev environment")
	case "staging":
		log.Println("Running in staging environment")
	case "prod":
		log.Println("Running in production")
	default:
		env = ""
		log.Println("Unknown environment, defaulting to local")
	}
	return env
}

// main function
func main() {
	mux := http.NewServeMux()
	base := "/api/v1"
	// base = ""
	mux.HandleFunc(base+"/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	mux.HandleFunc(base+"/healthcheck", HealthHandler)
	//used as livenessProbe (container is alive).
	mux.HandleFunc(base+"/healthz", HealthHandler)
	//used as readinessProbe (service is ready to accept traffic).
	mux.HandleFunc(base+"/readyz", ReadyHandler)

	//used for echo the message.
	mux.HandleFunc(base+"/echo", EchoHandler)

	//read environment
	APP_ENV := getEnvApp()
	// APP_VERSION := "1.0.0"
	// COMMIT_HASH := "none"
	// BUILD_TIME := "none"
	//swagger
	if APP_ENV != "prod" {
		mux.Handle("/swagger/", httpSwagger.WrapHandler)
	}

	//# Database dependent APIs

	//TODO: used for counter fetch
	//TODO: used for incrementing counter
	//TODO: used for websockets
	//read environment
	port := getEnvPort()
	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, enableCORS(mux)))
}
