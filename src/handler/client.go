package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"ewallet/src/usaldo"
	"github.com/gin-gonic/gin"
)

func HandleTransfer(c *gin.Context) {
	id := c.PostForm("user_id")
	nilai_str := c.PostForm("nilai")
	nilai, err := strconv.ParseInt(nilai_str, 10, 64)
	if err != nil {
		log.Println("[ERR] ParseInt Transfer", err)
	}
	url := c.PostForm("url")

	// transfer request from client
	// Check saldo before transfer
	var resp Response
	can := usaldo.CheckTransfer(id, nilai)
	if can == 0 {
		resp.Error = 1
		resp.Message = "Saldo not enough"
	} else if can == -1 {
		resp.Error = 1
		resp.Message = "User not exist"
	} else {
		// transfer
		res, _ := http.Get(url)
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		var st StatusTransfer
		_ = json.Unmarshal(body, &st)

		// reduce saldo if transfer success
		if st.Status == 0 {
			usaldo.ReduceSaldo(id, nilai)
			resp.Error = 0
			resp.Message = "Succcess Transfer"
		}
	}

	c.JSON(200, resp)
}
