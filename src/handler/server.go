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
	tpong, _ := cekQuorum()
	var return_saldo Saldo

	if tpong == 8 {
		var p MyParam
		c.BindJSON(&p)
		log.Println("[CHECK] user_id", p.Id)

		sld := int64(0)

		is_here, ip := usaldo.IsHere(p.Id)

		if is_here == 1 {
			// user berdomisili di kantor cabang
			// request to all
			ns := usaldo.NsKelompok
			for _, n := range ns {
				tmp := RequestSaldo(p, n)
				if tmp.Nilai != -1 {
					sld = sld + tmp.Nilai
				}
			}
		} else if is_here == 0 {
			// user tidak berdomisli di kantor cabang ini
			// request getTotalSaldo to kantor cabang domisili
			var ns string
			switch {
			case ip == "152.118.33.76":
				ns = "raditya"
			case ip == "152.118.33.85":
				ns = "wicaksono"
			case ip == "152.118.33.95":
				ns = "joseph"
			case ip == "152.118.33.96":
				ns = "saga"
			case ip == "152.118.33.97":
				ns = "halim"
			case ip == "152.118.33.99":
				ns = "wijaya"
			case ip == "152.118.33.104":
				ns = "gylberth"
			}
			return_saldo = RequestTotal(p, ns)
		} else {
			// user tidak ditemukan
			return_saldo.Nilai = -1
		}
		return_saldo.Nilai = sld
	} else {
		return_saldo.Nilai = -1
	}

	c.JSON(200, return_saldo)
}

func GetSaldo(c *gin.Context) {
	tpong, _ := cekQuorum()
	var return_saldo Saldo

	if tpong >= 5 {

		var p MyParam
		c.BindJSON(&p)
		id := p.Id
		log.Println("[CHECK] Handler Server GetSaldo id ->", id)

		sld, _ := usaldo.GetSaldo(id)

		return_saldo.Nilai = sld
	} else {
		return_saldo.Nilai = -1
	}

	c.JSON(200, return_saldo)
}

func Transfer(c *gin.Context) {
	tpong, _ := cekQuorum()
	var status StatusTransfer

	if tpong >= 5 {
		var p MyParam
		c.BindJSON(&p)
		id := p.Id
		nilai := p.Nilai

		s := usaldo.RecieveTransfer(id, nilai)
		status.Status = s
	} else {
		status.Status = -1
	}
	c.JSON(200, status)
}

func Register(c *gin.Context) {
	tpong, _ := cekQuorum()
	var rs Response

	if tpong >= 5 {
		var p MyParam
		c.BindJSON(&p)
		id := p.Id
		nama := p.Nama
		ip := p.Ip

		err := usaldo.Register(id, nama, ip)
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
	} else {
		rs.Success = 0
		rs.Error = 1
		rs.Message = "Failed Check Quorum"
	}

	c.JSON(200, rs)
}
