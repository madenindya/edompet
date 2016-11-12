package handler

import (
	"fmt"
	"log"

	"ewallet/src/usaldo"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	USER   = "postgres"
	PASS   = "postgres"
	DBNAME = "ewallet"
)

func Init() {
	// init insecure
	initInsecure()

	// initial database
	data_src := fmt.Sprintf("user=%v password='%v' dbname=%v sslmode=disable", USER, PASS, DBNAME)
	db, err := sqlx.Connect("postgres", data_src)
	if err != nil {
		log.Fatalln(err)
	}

	usaldo.Init(db)
}
