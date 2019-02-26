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
	router.HandleFunc("/json", getJsonResults).Methods("GET")
	router.HandleFunc("/json/{a:[0-9]+}/{b:[a-zA-Z]+}", getJsonResult).Methods("GET")

	// swagger like
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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
	port := getPort()

	fmt.Println("Listening On ", port)
	log.Fatal(http.ListenAndServe(port, router))
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
