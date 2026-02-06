package authentication

import "net/http"

type BaseAuthSystem interface {
	Login(w http.ResponseWriter, r *http.Request) bool
	Logout(w http.ResponseWriter, r *http.Request)
}

type CodeAuthSystem interface {
	SendCode(w http.ResponseWriter, r *http.Request)
	VerifyCode(w http.ResponseWriter, r *http.Request) bool
}
