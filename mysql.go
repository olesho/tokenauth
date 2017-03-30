// tokenauth project tokenauth.go
package tokenauth

import (
	"database/sql"
	//	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConfig interface {
	GetDbName() string
	GetDbAddress() string
	GetDbUser() string
	GetDbPassword() string
	GetDbTable() string
}

type MysqlStorage struct {
	db   *sql.DB
	conf MysqlConfig
}

func (s *MysqlStorage) CreateUser(name string, passwordHash string, additional map[string]interface{}) (*User, error) {
	addKeys := ""
	qMarks := ""
	vals := make([]interface{}, len(additional)+2)
	vals[0] = name
	vals[1] = passwordHash
	if additional != nil {
		i := 2
		for k, v := range additional {
			addKeys += ", " + k
			qMarks += ", ?"
			vals[i] = v
			i++
		}
	}

	stmt, err := s.db.Prepare("INSERT INTO " + s.conf.GetDbTable() + " (name, passwordHash" + addKeys + ") VALUES (?, ?" + qMarks + ")")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(vals...)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &User{
		id:   id,
		name: name,
	}, nil
}
func (s *MysqlStorage) ReadUser(id int64) (*User, error) {
	rows, err := s.db.Query("SELECT * FROM "+s.conf.GetDbTable()+" WHERE id = ? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		res := &User{}
		err = rows.Scan(&res.id, &res.name, &res.passwordHash, &res.recoveryState, &res.created)
		return res, err
	}
	return nil, nil
}

func (s *MysqlStorage) ReadUserByName(name string) (*User, error) {
	rows, err := s.db.Query("SELECT * FROM "+s.conf.GetDbTable()+" WHERE name = ? LIMIT 1", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		res := &User{}
		err = rows.Scan(&res.id, &res.name, &res.passwordHash, &res.recoveryState, &res.created)
		return res, err
	}
	return nil, nil
}
func (s *MysqlStorage) UpdateUser(user *User) error {
	stmt, err := s.db.Prepare("UPDATE " + s.conf.GetDbTable() + " SET name=?, passwordHash=?, recoveryState=? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.name, user.passwordHash, user.recoveryState, user.id)
	return err
}
func (s *MysqlStorage) DeleteUser(user *User) error {
	stmt, err := s.db.Prepare("DELETE FROM " + s.conf.GetDbTable() + " WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.id)
	return err
}

func NewMysqlStorage(conf MysqlConfig) (*MysqlStorage, error) {
	db, err := sql.Open("mysql",
		conf.GetDbUser()+":"+conf.GetDbPassword()+"@tcp("+conf.GetDbAddress()+")/"+conf.GetDbName())
	storage := &MysqlStorage{db, conf}
	return storage, err
}
