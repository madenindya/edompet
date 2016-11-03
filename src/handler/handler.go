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
