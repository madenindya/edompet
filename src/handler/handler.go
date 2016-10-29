package handler

import (
	"fmt"
	"log"

	"ewallet/src/usaldo"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	USER   = "postgres"
	PASS   = "postgres"
	DBNAME = "ewallet"
)

type (
	RegisterField struct {
		Id   string `json:"user_id"`
		Nama string `json:"nama"`
		Ip   string `json:"ip_domisili"`
	}

	Saldo struct {
		Nilai int64 `json:"nilai_saldo"`
	}

	Response struct {
		Error   int    `json:"error"`
		Message string `json:"message"`
	}

	StatusTransfer struct {
		Status int `json:"status_transfer"`
	}
)

func Init() {
	// initial database
	data_src := fmt.Sprintf("user=%v password='%v' dbname=%v sslmode=disable", USER, PASS, DBNAME)
	db, err := sqlx.Connect("postgres", data_src)
	if err != nil {
		log.Fatalln(err)
	}

	usaldo.Init(db)
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"pong": 1,
	})
}
