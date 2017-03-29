// tokenauth project tokenauth.go
package tokenauth

import (
	"encoding/json"
	"net/http"
	//	"strconv"
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

func NewAuth(authInstance AuthApi, langFile string) *Auth {
	lang, err := NewLang(langFile)
	if err != nil {
		panic(err)
	}

	return &Auth{
		Login: func(res http.ResponseWriter, req *http.Request) {
			// decode request body (json)
			var fields map[string]string
			err := json.NewDecoder(req.Body).Decode(&fields)
			if err != nil {
				res.Write(resError(err))
				return
			}

			// try get parsed token and status for credentials
			tokenString, err := authInstance.Login(fields["email"], fields["password"])
			if err != nil {
				//log.Println(fields["email"], err.Error())
				res.Write(resErrorMsg(lang.ERROR_LOGIN))
				return
			}

			// success response
			tokenCookie := http.Cookie{Name: "Token", Value: tokenString, HttpOnly: true, Path: "/"}
			identityCookie := http.Cookie{Name: "Identity", Value: fields["name"], Path: "/"}
			//			statusCookie := http.Cookie{Name: "Status", Value: strconv.Itoa(status), Path: "/"}
			http.SetCookie(res, &tokenCookie)
			http.SetCookie(res, &identityCookie)
			//			http.SetCookie(res, &statusCookie)

			payload := make(map[string]interface{})
			payload["name"] = fields["name"]
			payload["token"] = tokenString
			res.Write(resSuccess(payload))
			return
		},
		Logout: func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Login"))
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
