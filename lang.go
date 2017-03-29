// tokenauth project tokenauth.go
package tokenauth

import (
	"encoding/json"
	"io/ioutil"
)

type Lang struct {
	ERROR_UNKNOWN          string
	ERROR_DECODE_JSON      string
	ERROR_LOGIN            string
	ERROR_CREATE_USER      string
	ERROR_FINISH_SIGNUP    string
	ERROR_UPDATING_ACCOUNT string
	ERROR_GET_ACCOUNT      string
	ERROR_SEND_RECOVERY    string
	ERROR_NEW_PASSWORD     string
	ERROR_NO_TOKEN         string
	ERROR_TOKEN_INVALID    string

	SUCCESS_LOGIN            string
	SUCCESS_CREATE_USER      string
	SUCCESS_FINISH_SIGNUP    string
	SUCCESS_UPDATING_ACCOUNT string
	SUCCESS_GET_ACCOUNT      string
	SUCCESS_SEND_RECOVERY    string
	SUCCESS_NEW_PASSWORD     string
}

func NewLang(fileName string) (*Lang, error) {
	res := &Lang{}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}
