package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Timestamp represents the timestamp returned to the user
type Timestamp struct {
	Unix    int64  `json:"unix"`
	Natural string `json:"natural"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/{date:[0-9]+}", unixTimeStamp)
	r.HandleFunc("/{date}", naturalTimeStamp)

	fmt.Println("Server listening on port: " + port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	path := path.Join("index.html")

	temp, _ := template.ParseFiles(path)
	temp.Execute(w, nil)
}

// unixTimeStamp will handle a request to to process a unix style integer timestamp
func unixTimeStamp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := parseTimestamp(vars["date"])
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(timestamp); err != nil {
		panic(err)
	}
}

// naturalTimeStamp will process a human readable timestamp in the format
// January 2, 2017
func naturalTimeStamp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	t, err := parseTimestring(vars["date"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid date format provided, must be of format January, 15, 2017")
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(t); err != nil {
			panic(err)
		}
	}
}

// parseTimestamp will parse a unix timestamp into a Timestamp struct to return
// as json
func parseTimestamp(timestamp string) Timestamp {
	tInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	t := time.Unix(tInt, 0)
	return Timestamp{Unix: t.Unix(), Natural: t.UTC().String()}
}

// parseTimestring will parse a human readable date format matching the specified
// layout into a Timestamp struct
func parseTimestring(timestamp string) (Timestamp, error) {
	layout := "January _2, 2006"
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		return Timestamp{}, err
	}
	return Timestamp{Unix: t.Unix(), Natural: t.UTC().String()}, nil
}
