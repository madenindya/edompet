package handler

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"ewallet/src/usaldo"
	"github.com/gin-gonic/gin"
)

type (
	TransferField struct {
		Id    string `json:"user_id"`
		Nilai int64  `json:"nilai"`
	}

	RegisterField struct {
		Id   string `json:"user_id"`
		Nama string `json:"nama"`
		Ip   string `json:"ip_domisili"`
	}
)

func RenderTransfer(c *gin.Context) {
	// get all user
	users := usaldo.GetRegisteredUser()
	// get all bank tujuan
	ips := usaldo.GetAllUser()

	c.HTML(http.StatusOK, "transfer.tmpl", gin.H{
		"title": "Transfer Saldo",
		"ips":   ips,
		"users": users,
	})
}

func HandleTransfer(c *gin.Context) {
	id := c.PostForm("user_id")
	ns := c.PostForm("selectbasic")
	nilai_str := c.PostForm("nilai_saldo")
	nilai, err := strconv.ParseInt(nilai_str, 10, 64)
	if err != nil {
		log.Println("[ERROR] ParseInt Transfer", err)
	}

	log.Println("[CHECK] id", id, " ns", ns, " nilai", nilai)

	// urlstr := fmt.Sprintf("https://%v.sisdis.ui.ac.id/ewallet/transfer", ns)
	urlstr := "http://localhost:8080/transfer"

	// transfer request from client
	// Check saldo before transfer
	var resp Response
	var st StatusTransfer

	can := usaldo.CheckTransfer(id, nilai)
	if can == 0 {
		resp.Error = 1
		resp.Message = "Saldo not enough"
	} else if can == -1 {
		resp.Error = 1
		resp.Message = "User not exist"
	} else {
		res, err := http.PostForm(urlstr, url.Values{"user_id": {id}, "nilai": {nilai_str}})
		if err != nil {
			log.Println("[ERROR] Handler HandleTransfer PostForm", err)
		}
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)
		_ = json.Unmarshal(body, &st)

		// reduce saldo if transfer success
		if st.Status == 0 {
			usaldo.ReduceSaldo(id, nilai)
			resp.Success = 1
			resp.Message = "Succcess Transfer"
		} else {
			resp.Error = 1
			resp.Message = "Failed to Transfer"
		}
	}

	// BACK TO ORIGINAL PAGE
	users := usaldo.GetRegisteredUser()
	ips := usaldo.GetAllUser()
	if resp.Error == 0 && st.Status == 0 {
		c.HTML(http.StatusOK, "transfer.tmpl", gin.H{
			"title":   "Transfer Saldo",
			"ips":     ips,
			"users":   users,
			"success": resp.Success,
			"message": resp.Message,
		})
	} else {
		c.HTML(http.StatusOK, "transfer.tmpl", gin.H{
			"title":   "Transfer Saldo",
			"ips":     ips,
			"users":   users,
			"error":   resp.Error,
			"message": resp.Message,
		})
	}
}

func RenderRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tmpl", gin.H{
		"title": "Register",
	})
}

func HandleRegister(c *gin.Context) {
	id := c.PostForm("user_id")
	nama := c.PostForm("nama")
	ip := c.PostForm("ip_domisili")

	// urlstr := "https://nindyatama.sisdis.ui.ac.id/ewallet/register"
	urlstr := "http://localhost:8080/register"

	res, _ := http.PostForm(urlstr, url.Values{"user_id": {id}, "nama": {nama}, "ip_domisili": {ip}})
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var rs Response
	_ = json.Unmarshal(body, &rs)

	// BACK TO ORIGINAL PAGE
	if rs.Success == 1 {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title":   "Register",
			"success": rs.Success,
			"message": rs.Message,
		})
	} else {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title":   "Register",
			"error":   rs.Error,
			"message": rs.Message,
		})
	}
}

func RenderSaldo(c *gin.Context) {
	// get all user
	users := usaldo.GetRegisteredUser()
	c.HTML(http.StatusOK, "saldo.tmpl", gin.H{
		"users": users,
		"title": "Cek Saldo",
	})
}

func HandleSaldo(c *gin.Context) {
	id := c.PostForm("user_id")
	urlstr := fmt.Sprintf("http://localhost:8080/getSaldo/%v", id)

	resp, _ := http.Get(urlstr)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var sld Saldo
	_ = json.Unmarshal(body, &sld)

	// get all user
	users := usaldo.GetRegisteredUser()
	if sld.Nilai != -1 {
		c.HTML(http.StatusOK, "saldo.tmpl", gin.H{
			"users":       users,
			"title":       "Cek Saldo",
			"success":     1,
			"nilai_saldo": sld.Nilai,
		})
	} else {
		c.HTML(http.StatusOK, "saldo.tmpl", gin.H{
			"users":   users,
			"title":   "Cek Saldo",
			"error":   1,
			"message": "User not found",
		})
	}
}

func RenderTotalSaldo(c *gin.Context) {
	// get all user
	users := usaldo.GetRegisteredUser()
	c.HTML(http.StatusOK, "totalsaldo.tmpl", gin.H{
		"users": users,
		"title": "Cek Total Saldo",
	})
}

func HandleTotalSaldo(c *gin.Context) {
	id := c.PostForm("user_id")
	urlstr := fmt.Sprintf("http://localhost:8080/getTotalSaldo/%v", id)

	resp, _ := http.Get(urlstr)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var sld Saldo
	_ = json.Unmarshal(body, &sld)

	// get all user
	users := usaldo.GetRegisteredUser()
	c.HTML(http.StatusOK, "saldo.tmpl", gin.H{
		"users":       users,
		"title":       "Cek Total Saldo",
		"success":     1,
		"nilai_saldo": sld.Nilai,
	})

}
