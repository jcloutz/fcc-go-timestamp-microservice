package main

import (
	"context"
	"encoding/json"
	"errors"
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

type key int

const (
	timestampType key = iota
)

// Timestamp represents the timestamp returned to the user
type Timestamp struct {
	Unix    *int64  `json:"unix"`
	Natural *string `json:"natural"`
}

func jsonMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		h(w, r)
	}
}

func dateParserMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get date from route parameter
		vars := mux.Vars(r)
		timeArg := vars["date"]
		if timeArg == "" {
			fmt.Println(`Missing required date argument`)
			errorHandler(w, r)
			return
		}

		// Declare t to hold return
		var t Timestamp
		var e error

		// check for integer value
		if tInt, err := strconv.ParseInt(timeArg, 10, 64); err == nil {
			// if integer parse timestamp
			t, e = parseTimestamp(tInt)
		} else {
			// if string parse string
			t, e = parseTimestring(timeArg)
		}

		// If either parse function returned an error, call the error handler and
		// return early
		if e != nil {
			errorHandler(w, r)
			return
		}

		// create new context with Timestamp
		ctx := context.WithValue(r.Context(), timestampType, &t)
		r = r.WithContext(ctx)

		h(w, r)
	}
}

func app() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/{date:[0-9]+}", jsonMiddleware(dateParserMiddleware(unixTimeStamp)))
	r.HandleFunc("/{date}", jsonMiddleware(dateParserMiddleware(naturalTimeStamp)))

	return r
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "localhost:8080"
	} else {
		port = ":" + port
	}
	fmt.Println("Server listening on port: " + port)
	if err := http.ListenAndServe(port, app()); err != nil {
		log.Println(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	path := path.Join("index.html")

	temp, _ := template.ParseFiles(path)
	temp.Execute(w, nil)
}

// unixTimeStamp will handle a request to to process a unix style integer timestamp
func unixTimeStamp(w http.ResponseWriter, r *http.Request) {
	t := r.Context().Value(timestampType).(*Timestamp)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(*t); err != nil {
		log.Println(err)
	}
}

// naturalTimeStamp will process a human readable timestamp in the format
// January 2, 2017
func naturalTimeStamp(w http.ResponseWriter, r *http.Request) {
	t := r.Context().Value(timestampType).(*Timestamp)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(*t); err != nil {
		log.Println(err)
	}
}

// errorHandler handles responses to bad requests to the api endpoints.
func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Status-Reason", "Invalid date time format, refer to documentation.")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(Timestamp{}); err != nil {
		log.Println(err)
	}
}

// parseTimestamp will parse a unix timestamp into a Timestamp struct to return
// as json
func parseTimestamp(timestamp int64) (Timestamp, error) {
	if timestamp < 0 {
		return Timestamp{}, errors.New("Invalid Timestamp: less than 0")
	}
	t := time.Unix(timestamp, 0)
	unix, natural := t.Unix(), t.UTC().String()
	return Timestamp{&unix, &natural}, nil
}

// parseTimestring will parse a human readable date format matching the specified
// layout into a Timestamp struct
func parseTimestring(timestamp string) (Timestamp, error) {
	layout := "January _2, 2006"
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		return Timestamp{nil, nil}, err
	}
	unix, natural := t.Unix(), t.UTC().String()
	return Timestamp{&unix, &natural}, nil
}
