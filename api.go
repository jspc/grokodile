package main

import (
	"fmt"
	"log"
	"net/http"

	cookie "github.com/gorilla/securecookie"
	"github.com/satori/go.uuid"
)

const (
	CookieName = "grokodile"
)

type Cookie struct {
	Session string
}

type API struct {
	SecureCookie *cookie.SecureCookie
}

func NewAPI(hashKey, blockKey []byte) (a API, err error) {
	a.SecureCookie = cookie.New(hashKey, blockKey)

	return
}

func setHeaders(w http.ResponseWriter, origin string) (wDup http.ResponseWriter) {
	wDup = w
	wDup.Header().Set("Access-Control-Allow-Headers", "requested-with, Content-Type, origin, authorization, accept, client-security-token, cache-control, x-api-key")
	wDup.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	wDup.Header().Set("Access-Control-Allow-Origin", origin)
	wDup.Header().Set("Access-Control-Allow-Credentials", "true")
	wDup.Header().Set("Access-Control-Max-Age", "10")
	wDup.Header().Set("Cache-Control", "no-cache")

	return
}

func logRequest(requestID string, r *http.Request) {
	log.Printf("%s -> %s :: %s %s",
		requestID,
		r.RemoteAddr,
		r.Method,
		r.URL.Path)
}

func (a API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.NewV4().String()

	logRequest(requestID, r)
	w = setHeaders(w, r.Header.Get("Origin"))

	uuid, err, isNew := a.CookieUUID(r)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		http.Error(w, "Could not assign a session", 500)

		return
	}

	// do something with request - I don't want preflight OPTIONS requests
	// pushing things into a database and muddling things
	if r.Method == "GET" {
		if isNew {
			log.Printf("%s -> minted new sessionID: %s", requestID, uuid)
		} else {
			log.Printf("%s -> sessionID: %s", requestID, uuid)
		}

		ShipRequest(requestID, uuid, r)
	}

	encoded, err := a.encode(uuid)
	if err == nil {
		http.SetCookie(w, a.FormCookie(encoded))
	}
}

func (a API) CookieUUID(r *http.Request) (string, error, bool) {
	cookie, err := r.Cookie(CookieName)

	if err != nil {
		return uuid.NewV4().String(), nil, true
	}

	s, err := a.decode(cookie)
	return s, err, false
}

func (a API) FormCookie(encoded string) *http.Cookie {
	return &http.Cookie{
		Name:  CookieName,
		Value: encoded,
		Path:  "/",
	}
}

func (a API) encode(uuid string) (string, error) {
	return a.SecureCookie.Encode(CookieName, Cookie{uuid})
}

func (a API) decode(cookie *http.Cookie) (uuid string, err error) {
	c := &Cookie{}
	err = a.SecureCookie.Decode(CookieName, cookie.Value, c)
	if err != nil {
		return
	}
	uuid = c.Session

	return
}
