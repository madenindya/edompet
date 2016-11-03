package handler

import (
	"fmt"
	"log"
	"strconv"

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

func GetTotalSaldo(c *gin.Context) {
	id := c.Param("user_id")

	sld, err := usaldo.GetTotalSaldo(id)
	if err != nil {
		log.Println("[ERROR] GetTotalSaldo", id, ":", err)
	}

	return_saldo := Saldo{
		Nilai: sld,
	}
	c.JSON(200, return_saldo)
}

func GetSaldo(c *gin.Context) {
	id := c.Param("user_id")

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
	id := c.PostForm("user_id")
	nilai_str := c.PostForm("nilai")
	log.Println("[CHECK] user id", id, " nilai", nilai_str)
	nilai, err := strconv.ParseInt(nilai_str, 10, 64)
	if err != nil {
		log.Println("[ERROR] Handler Transfer ParseInt", err)
	}

	s := usaldo.RecieveTransfer(id, nilai)
	status := StatusTransfer{
		Status: s,
	}
	c.JSON(200, status)
}

func Register(c *gin.Context) {
	id := c.PostForm("user_id")
	nama := c.PostForm("nama")
	ip := c.PostForm("ip_domisili")

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
