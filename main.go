package main

import (
	hd "ewallet/src/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	hd.Init()

	r.GET("/ping", hd.Ping)
	r.GET("/getTotalSaldo/:user_id", hd.GetTotalSaldo)
	r.GET("/getSaldo/:user_id", hd.GetSaldo)
	r.POST("/register", hd.Register)
	r.POST("/transfer", hd.Transfer)

	r.Run(":8080") // listen and server on 0.0.0.0:8080
}
