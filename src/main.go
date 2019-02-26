// credit: https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo
// https://github.com/gorilla/mux

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// main function to boot up everything
func main() {
	router := mux.NewRouter()

	router.StrictSlash(false)
	router.Use(loggingMiddleware)

	router.HandleFunc("/json", getJsonResults).Methods("GET")
	router.HandleFunc("/json/{a:[0-9]+}/{b:[a-zA-Z]+}", getJsonResult).Methods("GET")
	router.HandleFunc("/echo/{*}", echo).Methods("GET") //.Queries("count")

	subRoute := router.PathPrefix("/api").Subrouter()
	subRoute.HandleFunc("/json", getJsonResults).Methods("GET")
	subRoute.HandleFunc("/json/{a:[0-9]+}/{b:[a-zA-Z]+}", getJsonResult).Methods("GET")

	router.PathPrefix("/").Handler(catchAllHandler())

	port := getPort()

	metaReader(router)
	fmt.Println("Listening On ", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func catchAllHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Anything:", r.RequestURI)
		fmt.Fprintf(w, "Other\nQuery %v\nPath %v", r.URL.Query(), mux.Vars(r))
	})
}

// middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Host)
		log.Println(r.RequestURI)
		log.Println(r.URL.Query())
		log.Println(mux.Vars(r))

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func metaReader(r *mux.Router) {
	// swagger like
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

type AB struct {
	A int    `json:"a,omitempty"`
	B string `json:"b,omitempty"`
}

// Display all from the people var
func getJsonResults(w http.ResponseWriter, r *http.Request) {

	results := []AB{
		AB{A: 1, B: "B1"},
		AB{A: 2, B: "B2"},
		AB{A: 3, B: "B3"},
	}
	json.NewEncoder(w).Encode(results)
}

// Display all from the people var
func echo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Count: %v\nPath %v\n", r.URL.Query()["count"], vars)
}

// Display all from the people var
func getJsonResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, _ := strconv.Atoi(vars["a"])
	result := AB{A: a, B: vars["b"]}
	json.NewEncoder(w).Encode(result)
}

func getPort() string {
	port := os.Getenv("GO_SERVER_PORT")
	if len(port) == 0 {
		port = "8085"
	}
	port = ":" + port
	return port
}
