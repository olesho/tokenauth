// tokenauth project tokenauth.go
package tokenauth

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConfig interface {
	GetDbName() string
	GetDbAddress() string
	GetDbUser() string
	GetDbPassword() string
}

type MysqlStorage struct {
	db *sql.DB
}

func (s *MysqlStorage) CreateUser(email string, passwordHash string) (*User, error) {
	stmt, err := s.db.Prepare("INSERT INTO user (email, passwordHash) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(email, passwordHash)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &User{
		id:    id,
		email: email,
	}, nil
}
func (s *MysqlStorage) ReadUser(id int64) (*User, error) {
	rows, err := s.db.Query("SELECT * FROM `user` WHERE \"id\" = ? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		res := &User{}
		err = rows.Scan(&res.id, &res.email, &res.passwordHash, &res.recoveryState, &res.created)
		return res, err
	}
	return nil, nil
}

func (s *MysqlStorage) ReadUserByName(name string) (*User, error) {
	rows, err := s.db.Query("SELECT * FROM `user` WHERE \"name\" = ? LIMIT 1", name)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		res := &User{}
		err = rows.Scan(&res.id, &res.email, &res.passwordHash, &res.recoveryState, &res.created)
		return res, err
	}
	return nil, nil
}
func (s *MysqlStorage) UpdateUser(user *User) error {
	stmt, err := s.db.Prepare("UPDATE `user` SET email=?, passwordHash=?, recoveryState=? WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.email, user.passwordHash, user.recoveryState, user.id)
	return err
}
func (s *MysqlStorage) DeleteUser(user *User) error {
	return nil
}

func NewMysqlStorage(conf MysqlConfig) (*MysqlStorage, error) {
	db, err := sql.Open("mysql",
		conf.GetDbUser()+":"+conf.GetDbPassword()+"@tcp("+conf.GetDbAddress()+")/"+conf.GetDbName())
	storage := &MysqlStorage{db}
	return storage, err
}
