// tokenauth project tokenauth.go
package tokenauth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	//	"strconv"
)

type AuthConfig interface {
	GetLangFile() string
	GetSuccessRedirect() string
	GetFailRedirect() string
}

type Auth struct {
	PrivateAdapter  func(h http.Handler) http.Handler
	Login           func(res http.ResponseWriter, req *http.Request)
	GetToken        func(res http.ResponseWriter, req *http.Request)
	Logout          func(res http.ResponseWriter, req *http.Request)
	Signup          func(res http.ResponseWriter, req *http.Request)
	RecoverPassword func(res http.ResponseWriter, req *http.Request)
	ChangePassword  func(res http.ResponseWriter, req *http.Request)
}

//func (a *Auth) PrivateAdapter

func NewAuth(authInstance AuthApi, logger *log.Logger, config AuthConfig) *Auth {
	lang, err := NewLang(config.GetLangFile())
	if err != nil {
		panic(err)
	}

	return &Auth{
		PrivateAdapter: func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				// prevent browser from caching routes with restricted access
				res.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
				res.Header().Add("Expires", "0")

				// try to get token string from cookie
				tokenCookie, _ := req.Cookie("Token")
				tokenString := ""
				// try to get token string from cookie
				if tokenCookie != nil {
					tokenString = tokenCookie.Value
				}
				// try to get token string from HTTP header
				if tokenString == "" {
					tokenString = req.Header.Get("Authorization")
				}

				if tokenString == "" {
					res.WriteHeader(401)
					res.Write(resErrorMsg(lang.ERROR_NO_TOKEN))
					return
				}

				if user, err := authInstance.Authorized(tokenString); err == nil && user != nil {
					context.Set(req, "user", user)
					// serve next
					h.ServeHTTP(res, req)
					return
				}

				// otherwise send error response
				res.WriteHeader(401)
				res.Write(resErrorMsg(lang.ERROR_TOKEN_INVALID))
				return
			})
		},
		Login: func(res http.ResponseWriter, req *http.Request) {
			// decode request body (json)
			var name, password string
			name = req.FormValue("name")
			password = req.FormValue("password")

			// try get parsed token and status for credentials
			tokenString, err := authInstance.Login(name, password)
			if err != nil {
				logger.Println(err)
				http.Redirect(res, req, config.GetFailRedirect(), 301)
				return
			}

			// success response
			tokenCookie := http.Cookie{Name: "Token", Value: tokenString, HttpOnly: true, Path: "/"}
			identityCookie := http.Cookie{Name: "Identity", Value: name, Path: "/"}

			http.SetCookie(res, &tokenCookie)
			http.SetCookie(res, &identityCookie)
			http.Redirect(res, req, config.GetSuccessRedirect(), 301)

			return
		},
		GetToken: func(res http.ResponseWriter, req *http.Request) {
			var fields map[string]string
			err := json.NewDecoder(req.Body).Decode(&fields)
			if err != nil {
				res.Write(resError(err))
				return
			}
			name := fields["name"]
			password := fields["password"]

			// try get parsed token and status for credentials
			tokenString, err := authInstance.Login(name, password)
			if err != nil {
				logger.Println(err)
				res.Write(resErrorMsg(lang.ERROR_LOGIN))
				return
			}

			// success response
			payload := make(map[string]interface{})
			payload["name"] = name
			payload["token"] = tokenString
			res.Write(resSuccess(payload))
			return
		},
		Logout: func(res http.ResponseWriter, req *http.Request) {
			// prevent browser from caching routes with restricted access
			res.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
			res.Header().Add("Expires", "0")

			tokenCookie := http.Cookie{Name: "Token", Value: "", HttpOnly: true, Path: "/"}
			identityCookie := http.Cookie{Name: "Identity", Value: "", Path: "/"}

			http.SetCookie(res, &tokenCookie)
			http.SetCookie(res, &identityCookie)
			http.Redirect(res, req, config.GetFailRedirect(), 301)
		},
		Signup: func(res http.ResponseWriter, req *http.Request) {
			// decode request body (json)
			var fields map[string]string
			err := json.NewDecoder(req.Body).Decode(&fields)
			if err != nil {
				res.Write(resError(err))
				return
			}

			// signup and get token string
			err = authInstance.Signup(fields["name"], fields["password"])
			if err != nil {
				logger.Println(err)
				res.Write(resErrorMsg(lang.ERROR_CREATE_USER))
				return
			}

			res.Write(resSuccess(nil))
			return
		},
		RecoverPassword: func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Login"))
		},
		ChangePassword: func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Login"))
		},
	}
}

type SuccessResponse struct {
	Data interface{}
}

type ErrorResponse struct {
	Error string
}

func resError(err error) []byte {
	enc, _ := json.Marshal(&ErrorResponse{err.Error()})
	return enc
}

func resErrorMsg(msg string) []byte {
	enc, _ := json.Marshal(&ErrorResponse{msg})
	return enc
}

func resSuccess(data interface{}) []byte {
	enc, _ := json.Marshal(&SuccessResponse{data})
	return enc
}
