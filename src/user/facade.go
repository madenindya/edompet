package user

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db_main *sqlx.DB

func Init(db *sqlx.DB) {
	db_main = db
}

func Register(id, nama, ip string) error {
	_, err := getUser(id)
	if err == nil {
		// error nil -> user exist (by id)
		return err
	}

	err = insertNew(id, nama, ip)
	return err
}

func IsExist(id string) bool {
	_, err := getUser(id)
	if err != nil {
		return false
	}

	return true
}
