package handler

import (
	"fmt"
	"log"

	"ewallet/src/usaldo"
	"github.com/gin-gonic/gin"
)

type (
	Saldo struct {
		Nilai int64 `json:"nilai_saldo"`
	}

	Response struct {
		Success int    `json:"success"`
		Error   int    `json:"error"`
		Message string `json:"message"`
	}

	StatusTransfer struct {
		Status int `json:"status_transfer"`
	}
)

type MyParam struct {
	Id    string `json:"user_id"`
	Nama  string `json:"nama"`
	Ip    string `json:"ip_domisili"`
	Nilai int64  `json:"nilai"`
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"pong": 1,
	})
}

func GetTotalSaldo(c *gin.Context) {
	var p MyParam
	c.BindJSON(&p)

	sld := int64(0)

	// request to all
	ns := usaldo.NsKelompok
	for _, n := range ns {
		log.Println("[CHECK] Request GetSaldo to", n)
		tmp := RequestSaldo(p, n)
		if tmp.Nilai != -1 {
			sld = sld + tmp.Nilai
		}
	}

	return_saldo := Saldo{
		Nilai: sld,
	}
	c.JSON(200, return_saldo)
}

func GetSaldo(c *gin.Context) {
	var p MyParam
	c.BindJSON(&p)
	id := p.Id
	log.Println("[CHECK] Handler Server GetSaldo id ->", id)

	sld, err := usaldo.GetSaldo(id)
	if err != nil {
		log.Println("[ERROR] GetSaldo", id, ":", err)
	}

	return_saldo := Saldo{
		Nilai: sld,
	}
	c.JSON(200, return_saldo)
}

func Transfer(c *gin.Context) {
	var p MyParam
	c.BindJSON(&p)
	id := p.Id
	nilai := p.Nilai

	s := usaldo.RecieveTransfer(id, nilai)
	status := StatusTransfer{
		Status: s,
	}
	c.JSON(200, status)
}

func Register(c *gin.Context) {
	var p MyParam
	c.BindJSON(&p)
	id := p.Id
	nama := p.Nama
	ip := p.Ip

	err := usaldo.Register(id, nama, ip)
	var rs Response
	if err != nil {
		m := fmt.Sprintf("Failed to Register User %s", id)
		rs.Success = 0
		rs.Error = 1
		rs.Message = m
	} else {
		m := fmt.Sprintf("Success Register User %s", id)
		rs.Success = 1
		rs.Error = 0
		rs.Message = m
	}

	c.JSON(200, rs)
}

//
//
// OLD
// GET SALDO
// id := c.PostForm("user_id")
// log.Println("[CHECK] user id", id)
// TRANSFER
// id := c.PostForm("user_id")
// nilai_str := c.PostForm("nilai")
// log.Println("[CHECK] user id", id, " nilai", nilai_str)
// nilai, err := strconv.ParseInt(nilai_str, 10, 64)
// if err != nil {
// 	log.Println("[ERROR] Handler Transfer ParseInt", err)
// }
// REGISTER
// id := c.PostForm("user_id")
// nama := c.PostForm("nama")
// ip := c.PostForm("ip_domisili")
