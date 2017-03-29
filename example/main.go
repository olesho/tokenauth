// token-auth project main.go
package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olesho/tokenauth"
)

func main() {
	config := tokenauth.NewEnvConfig()
	userStorage, err := tokenauth.NewMysqlStorage(config)
	if err != nil {
		panic(err)
	}
	authInstance := tokenauth.NewDefaultAuthInstance(userStorage, config)
	auth := tokenauth.NewAuth(authInstance, "lang/english.json")

	r := mux.NewRouter()
	r.HandleFunc("/priv", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello"))
	})
	r.HandleFunc("/login", auth.Login).Methods("POST")
	r.HandleFunc("/logout", auth.Logout).Methods("GET")
	r.HandleFunc("/signup", auth.Signup).Methods("POST")
	r.HandleFunc("/recover_password", auth.RecoverPassword).Methods("GET")
	r.HandleFunc("/change_password", auth.ChangePassword).Methods("POST")
	http.ListenAndServe(":8080", r)
}
