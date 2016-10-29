package handler

import (
	"fmt"
	"log"
	"strconv"

	"ewallet/src/saldo"
	"ewallet/src/user"

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

	user.Init(db)
	saldo.Init(db)
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"pong": 1,
	})
}

func GetTotalSaldo(c *gin.Context) {
	id := c.Param("user_id")

	sld, err := saldo.GetTotalSaldo(id)
	if err != nil {
		log.Println("[ERR] GetTotalSaldo", id, ":", err)
	}

	return_saldo := Saldo{
		Nilai: sld,
	}
	c.JSON(200, return_saldo)
}

func GetSaldo(c *gin.Context) {
	id := c.Param("user_id")

	sld, err := saldo.GetTotalSaldo(id)
	if err != nil {
		log.Println("[ERR] GetSaldo", id, ":", err)
	}

	return_saldo := Saldo{
		Nilai: sld,
	}
	c.JSON(200, return_saldo)
}

func Transfer(c *gin.Context) {
	id := c.PostForm("user_id")
	nilai_str := c.PostForm("nilai")
	nilai, err := strconv.ParseInt(nilai_str, 10, 64)
	if err != nil {
		log.Println("[ERR] ParseInt Transfer", err)
	}

	s := saldo.RecieveTransfer(id, nilai)
	status := StatusTransfer{
		Status: s,
	}
	c.JSON(200, status)
}

func Register(c *gin.Context) {
	id := c.PostForm("user_id")
	nama := c.PostForm("nama")
	ip := c.PostForm("ip_domisili")

	err := user.Register(id, nama, ip)
	var resp Response
	if err != nil {
		resp.Error = 1
		message := fmt.Sprintf("Can not register user %v", id)
		resp.Message = message
	} else {
		resp.Error = 0
		message := fmt.Sprintf("Success to register user %v", id)
		resp.Message = message
	}

	c.JSON(200, resp)
}
