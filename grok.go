package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type LogData struct {
	// A subset of http.Request fields
	Method     string
	URL        *url.URL
	Proto      string
	Host       string
	RemoteAddr string
	RequestURI string
	Header     http.Header

	// Split params
	Params url.Values

	// Useful/ necessary metadata
	RequestID string
	SessionID string
}

func ShipRequest(requestID, uuid string, r *http.Request) {
	ld := LogData{
		Header:     r.Header,
		Host:       r.Host,
		Method:     r.Method,
		Params:     r.URL.Query(),
		Proto:      r.Proto,
		RemoteAddr: r.RemoteAddr,
		RequestID:  requestID,
		RequestURI: r.RequestURI,
		SessionID:  uuid,
		URL:        r.URL,
	}
	d, err := json.Marshal(ld)

	if err != nil {
		log.Print(err)
	}
}
