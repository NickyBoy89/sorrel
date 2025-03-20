package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
)

var (
	priv = ""
	pub  = ""
)

func InitializeVAPIDKeys() {
	var err error
	priv, pub, err = webpush.GenerateVAPIDKeys()
	if err != nil {
		panic(err)
	}
}

func handleVAPIDPublicKeyRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	w.Write([]byte(pub))
}

func handlePushSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	fmt.Println(io.ReadAll(r.Body))
}
