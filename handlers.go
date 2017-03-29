// tokenauth project tokenauth.go
package tokenauth

import (
	"net/http"
)

type Auth struct {
	Login           func(res http.ResponseWriter, req *http.Request)
	Logout          func(res http.ResponseWriter, req *http.Request)
	Signup          func(res http.ResponseWriter, req *http.Request)
	RecoverPassword func(res http.ResponseWriter, req *http.Request)
	ChangePassword  func(res http.ResponseWriter, req *http.Request)
}

func (a *Auth) PrivateAdapter(h http.Handler) http.Handler {
	return nil
}

func NewAuth(authInstanse AuthApi) *Auth {
	return &Auth{
		Login: func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Login"))
		},
		Logout: func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Login"))
		},
		Signup: func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Login"))
		},
		RecoverPassword: func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Login"))
		},
		ChangePassword: func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Login"))
		},
	}
}
