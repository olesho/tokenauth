// tokenauth project tokenauth.go
package tokenauth

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config interface {
	GetSecretKey() string
	GetIssuer() string
}

/*
type Config struct {
	Host string
	Port string

	DbName     string
	DbAddress  string
	DbUser     string
	DbPassword string

	SmtpServer   string
	SmtpPort     string
	SmtpUser     string
	SmtpPassword string
	DefaultEmail string

	LinkedInClientId     string
	LinkedInClientSecret string

	SecretKey string
	Issuer    string
}
*/

type ConfigVars struct {
	FileKeyMap map[string]interface{}

	SecretKey       string
	Issuer          string
	LangFile        string
	LogFile         string
	SuccessRedirect string
	FailRedirect    string

	DbName     string
	DbAddress  string
	DbUser     string
	DbPassword string
	DbTable    string
}

func (c *ConfigVars) GetSecretKey() string {
	return c.SecretKey
}

func (c *ConfigVars) GetIssuer() string {
	return c.Issuer
}

func (c *ConfigVars) GetLangFile() string {
	return c.LangFile
}

func (c *ConfigVars) GetLogFile() string {
	return c.LogFile
}

func (c *ConfigVars) GetSuccessRedirect() string {
	return c.SuccessRedirect
}

func (c *ConfigVars) GetFailRedirect() string {
	return c.FailRedirect
}

func (c *ConfigVars) GetDbName() string {
	return c.DbName
}
func (c *ConfigVars) GetDbAddress() string {
	return c.DbAddress
}
func (c *ConfigVars) GetDbUser() string {
	return c.DbUser
}
func (c *ConfigVars) GetDbPassword() string {
	return c.DbPassword
}
func (c *ConfigVars) GetDbTable() string {
	return c.DbTable
}

func NewEnvConfig() *ConfigVars {
	return &ConfigVars{
		make(map[string]interface{}),
		os.Getenv("SECRET_KEY"),
		os.Getenv("ISSUER"),
		os.Getenv("LANG_FILE"),
		os.Getenv("LOG_FILE"),
		os.Getenv("SUCCESS_REDIRECT"),
		os.Getenv("FAIL_REDIRECT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_TABLE"),
	}
}

func NewFileConfig(fileName string) (*ConfigVars, error) {
	res := &ConfigVars{}

	/*
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &res)
	*/

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &res.FileKeyMap)

	for k, v := range res.FileKeyMap {
		switch k {
		case "SecretKey":
			res.SecretKey = v.(string)
			break
		case "Issuer":
			res.Issuer = v.(string)
			break

		case "DbName":
			res.DbName = v.(string)
			break

		case "DbAddress":
			res.DbAddress = v.(string)
			break

		case "DbUser":
			res.DbUser = v.(string)
			break

		case "DbPassword":
			res.DbPassword = v.(string)
			break

		case "DbTable":
			res.DbTable = v.(string)
			break

		case "LangFile":
			res.LangFile = v.(string)
			break
		case "SuccessRedirect":
			res.SuccessRedirect = v.(string)
			break
		case "FailRedirect":
			res.FailRedirect = v.(string)
			break
		}
	}

	return res, err
}
