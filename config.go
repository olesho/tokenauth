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
	SecretKey string
	Issuer    string

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
		os.Getenv("SECRET_KEY"),
		os.Getenv("ISSUER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_TABLE"),
	}
}

func NewFileConfig(fileName string) (*ConfigVars, error) {
	res := &ConfigVars{}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}
