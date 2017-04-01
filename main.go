package main

import (
	"log"
	"net/http"
	"time"
)

// dd if=/dev/random bs=1 count=22| base64
const (
	HashKey  = "PelNm42cg29dpOtSUFhX0rPMbAILXQ=="
	BlockKey = "RuKqpsX7+cZESvai/cfzL9GDFaOROg=="
)

func main() {
	a, err := NewAPI([]byte(HashKey), []byte(BlockKey))
	if err != nil {
		log.Panic(err)
	}

	s := &http.Server{
		Addr:           ":8000",
		Handler:        a,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
