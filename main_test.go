package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestEncodeEmptyTimestamp(t *testing.T) {
	var bb bytes.Buffer

	if err := json.NewEncoder(&bb).Encode(Timestamp{}); err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(bb.String())
	want := `{"unix":null,"natural":null}`

	if got != want {
		t.Log("Wanted   :", want)
		t.Log("Got      :", got)
		t.Fatal("Mismatch")
	}
}

func TestEncodeTimestamp(t *testing.T) {
	var bb bytes.Buffer
	var unix int64 = 1450137600
	natural := `2015-12-15 00:00:00 +0000 UTC`
	timestamp := Timestamp{&unix, &natural}
	if err := json.NewEncoder(&bb).Encode(timestamp); err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(bb.String())
	want := `{"unix":1450137600,"natural":"2015-12-15 00:00:00 +0000 UTC"}`

	if got != want {
		t.Log("Wanted  :", want)
		t.Log("Got     :", got)
		t.Fatal("Mismatch")
	}
}

func TestParseTimestamp(t *testing.T) {
	var unix int64 = 1450137600

	got, err := parseTimestamp(unix)
	if err != nil {
		t.Fatal(err)
	}

	var wantUnix int64 = 1450137600
	wantNat := `2015-12-15 00:00:00 +0000 UTC`
	if (*got.Unix != wantUnix) || (*got.Natural != wantNat) {
		t.Log("Wanted   :", wantUnix, wantNat)
		t.Log("Got      :", *got.Unix, *got.Natural)

		t.Fatal("Mismatch")
	}
}

func TestParseInvalidTimestamp(t *testing.T) {
	var unix int64 = -1

	timestamp, got := parseTimestamp(unix)
	if !reflect.DeepEqual(timestamp, Timestamp{}) {
		t.Fatal("Returned Timestamp should be empty: Timestamp{}")
	}

	want := `Invalid Timestamp: less than 0`
	if got.Error() != want {
		t.Log("Wanted   :", want)
		t.Log("Got      :", got)

		t.Fatal("Mismatch")
	}
}

func TestParseTimestring(t *testing.T) {
	timestring := `December 15, 2015`

	got, err := parseTimestring(timestring)
	if err != nil {
		t.Fatal(err)
	}

	var wantUnix int64 = 1450137600
	wantNat := `2015-12-15 00:00:00 +0000 UTC`
	if (*got.Unix != wantUnix) || (*got.Natural != wantNat) {
		t.Log("Wanted   :", wantUnix, wantNat)
		t.Log("Got      :", *got.Unix, *got.Natural)

		t.Fatal("Mismatch")
	}
}

func TestParseInvalidTimestring(t *testing.T) {
	timestring := `Dec 15, 2015`

	timestamp, got := parseTimestring(timestring)

	if !reflect.DeepEqual(timestamp, Timestamp{}) {
		t.Fatal("Returned Timestamp should be empty: Timestamp{}")
	}

	if !strings.Contains(got.Error(), `cannot parse`) {
		t.Fatal("Invalid date should be be parsed")
	}
}

func TestIndexRoute(t *testing.T) {
	ts := httptest.NewServer(app())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := `Timestamp Microservice`
	got := string(b)

	if !strings.Contains(got, want) {
		t.Log("Wanted:", want)
		t.Fatal("Incorrect Title")
	}
}

func TestTimestampRouteHeader(t *testing.T) {
	ts := httptest.NewServer(app())
	defer ts.Close()

	res, err := http.Get(ts.URL + `/1450137600`)
	if err != nil {
		t.Fatal(err)
	}

	got := res.Header.Get(`Content-Type`)
	want := `application/json; charset=UTF-8`
	if got != want {
		t.Log("Invalid Content-Type Header")
		t.Log("Wanted   :", want)
		t.Log("Got      :", got)
		t.Fatal("Mismatch")
	}
}

func TestTimestampRoute(t *testing.T) {
	ts := httptest.NewServer(app())
	defer ts.Close()

	res, err := http.Get(ts.URL + `/1450137600`)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"unix":1450137600,"natural":"2015-12-15 00:00:00 +0000 UTC"}`
	got := strings.TrimSpace(string(b))

	if got != want {
		t.Log("Wanted   :", want)
		t.Log("Got      :", got)
		t.Fatal("Mismatch")
	}
}

func TestTimestringRoute(t *testing.T) {
	ts := httptest.NewServer(app())
	defer ts.Close()

	res, err := http.Get(ts.URL + `/December%2015,%202015`)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"unix":1450137600,"natural":"2015-12-15 00:00:00 +0000 UTC"}`
	got := strings.TrimSpace(string(b))

	if got != want {
		t.Log("Wanted   :", want)
		t.Log("Got      :", got)
		t.Fatal("Mismatch")
	}
}

func TestInvalidTimestringRoute(t *testing.T) {
	ts := httptest.NewServer(app())
	defer ts.Close()

	res, err := http.Get(ts.URL + `/Dec%2015,%202015`)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"unix":null,"natural":null}`
	got := strings.TrimSpace(string(b))

	if got != want {
		t.Log("Wanted   :", want)
		t.Log("Got      :", got)
		t.Fatal("Mismatch")
	}
}
