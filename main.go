package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
)

var client *redis.Client

func init() {
	NewRedisClient()
}

func main() {
	port := os.Getenv("HTTP_PORT")

	r := mux.NewRouter()
	r.Path("/api/health").Methods("GET").HandlerFunc(HealthCheckHandler)
	r.Path("/api/hello").Methods("GET").HandlerFunc(HelloHandler)
	r.Path("/redis/get/{key}").Methods("GET").HandlerFunc(RedisGetHandler)
	r.Path("/redis/set/{key}/{value}").Methods("GET").HandlerFunc(RedisSetHandler)
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

// NewRedisClient new redis client
func NewRedisClient() {
	inst := redis.NewClient(&redis.Options{
		Addr:     "dell-vm7.7onetella.net:22221",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := inst.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	if err == nil {
		client = inst
	}
}

// RedisGetHandler redis get calls
func RedisGetHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	key := params["key"]
	fmt.Println("key=", key)

	val, err := client.Get(key).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(val))
}

// RedisSetHandler redis set call
func RedisSetHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	key := params["key"]
	value := params["value"]

	err := client.Set(key, value, 0).Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
