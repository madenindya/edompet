package main

import (
	hd "ewallet/src/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	hd.Init()

	// ganti jadi post semua
	r.POST("/ping", hd.Ping)
	r.GET("/pingall", hd.QuorumCheck)
	// Register User
	r.GET("/client/register", hd.RenderRegister)
	r.POST("/client/register", hd.HandleRegister)
	r.POST("/register", hd.Register)
	// Transfer Saldo
	r.GET("/client/transfer", hd.RenderTransfer)
	r.POST("/client/transfer", hd.HandleTransfer)
	r.POST("/transfer", hd.Transfer)
	// Get Saldo
	r.GET("/client/getSaldo", hd.RenderSaldo)
	r.POST("/client/getSaldo", hd.HandleSaldo)
	r.POST("/getSaldo", hd.GetSaldo)
	// Get Total Saldo
	r.GET("/client/getTotalSaldo", hd.RenderTotalSaldo)
	r.POST("/client/getTotalSaldo", hd.HandleTotalSaldo)
	r.POST("/getTotalSaldo", hd.GetTotalSaldo)
	// Render HTML
	r.LoadHTMLGlob("view/*")

	r.Run(":8080") // listen and server on 0.0.0.0:8080
}
