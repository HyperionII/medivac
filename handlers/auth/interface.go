package auth

import "net/http"

type Handler interface {
	Login(http.ResponseWriter, *http.Request) error
}
