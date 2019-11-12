package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("HTTP_PORT")

	r := mux.NewRouter()
	r.Path("/api/health").Methods("GET").HandlerFunc(HealthCheckHandler)
	r.Path("/api/hello").Methods("GET").HandlerFunc(HelloHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	log.Fatal(srv.ListenAndServe())

}

// HealthCheckHandler handles health check
func HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	log.Println("/api/health")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// HelloHandler handles hello
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("HELLO"))
}
