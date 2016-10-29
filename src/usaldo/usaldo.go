package usaldo

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type (
	Usaldo struct {
		Id    string `db:"user_id" json:"user_id"`
		Nama  string `db:"nama" json:"nama"`
		Ip    string `db:"ip_domisili" json:"ip_domisili"`
		Nilai int64  `db:"nilai_saldo" json:"nilai_saldo"`
	}
)

var db_main *sqlx.DB
var u_saldo Usaldo

func Init(db *sqlx.DB) {
	db_main = db
}
