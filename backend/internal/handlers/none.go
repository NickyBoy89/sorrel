package handlers

import "net/http"

type NoAuthHandler struct{}

func (*NoAuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

}

func (*NoAuthHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
}
