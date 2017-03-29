// tokenauth project tokenauth.go
package tokenauth

// basic CRUD operations
type UserStorage interface {
	CreateUser(name string, passwordHash string, additional map[string]interface{}) (*User, error)
	ReadUser(id int64) (*User, error)
	ReadUserByName(name string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(user *User) error
}

type User struct {
	id           int64
	name         string
	passwordHash string
	status       int
	// recovery state is used in "password recovery token"
	// whenever user sets new password this state is changed randomly to invalidate old token
	recoveryState int64
	created       string
	additional    map[string]interface{}
}
