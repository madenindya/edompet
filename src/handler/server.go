package handler

import (
	"fmt"
	"log"
	"strconv"

	"ewallet/src/usaldo"

	"github.com/gin-gonic/gin"
)

func GetTotalSaldo(c *gin.Context) {
	id := c.Param("user_id")

	sld, err := usaldo.GetTotalSaldo(id)
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

	sld, err := usaldo.GetSaldo(id)
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
