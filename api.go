// tokenauth project tokenauth.go
package tokenauth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthApi interface {
	Authorized(token string) bool
	Login(user, password string) (token string, status int, err error)
	Logout(uid int64) error
	Signup(user, password string) error
	RecoverPasswordToken(user string) (token string, err error)
	ChangePassword(token string, newPassword string) error
}

type DefaultAuthInstance struct {
	userStorage UserStorage
	config      Config
}

func NewDefaultAuthInstance(userStorage UserStorage, config Config) *DefaultAuthInstance {
	return &DefaultAuthInstance{userStorage, config}
}

func (a *DefaultAuthInstance) Authorized(token string) bool {
	return false
}

func (a *DefaultAuthInstance) Login(user, password string) (token string, status int, err error) {
	// creates token based on user ID
	var getToken = func(uid int64) (string, error) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, &Claims{
			a.config.GetIssuer(),
			uid,
			// token valid for 72 hours
			time.Now().Add(time.Hour * 72).Unix(),
			-1, // recovery state -1 means this token isn't used for password recovery
		})
		return token.SignedString([]byte(a.config.GetSecretKey()))
	}

	err = ValidateEmail(user)
	if err != nil {
		return "", -1, err
	}

	err = ValidatePassword(password)
	if err != nil {
		return "", -1, err
	}

	u, err := a.userStorage.ReadUserByName(user)
	if err != nil {
		return "", -1, err
	}
	if u == nil {
		return "", -1, errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password))
	if err != nil {
		return "", -1, err
	}
	token, err = getToken(u.id)
	if err != nil {
		return "", -1, err
	}
	return token, u.status, nil
}

func (a *DefaultAuthInstance) Logout(uid int64) error {
	return nil
}
func (a *DefaultAuthInstance) Signup(user, password string) error {
	return nil
}
func (a *DefaultAuthInstance) RecoverPasswordToken(user string) (token string, err error) {
	return "", nil
}
func (a *DefaultAuthInstance) ChangePassword(token string, newPassword string) error {
	return nil
}

type Claims struct {
	Iss           string
	Uid           int64
	Exp           int64
	RecoveryState int64
}

func (c *Claims) Valid() error {
	if time.Now().Unix() > c.Exp {
		return errors.New("Token expired")
	}
	return nil
}