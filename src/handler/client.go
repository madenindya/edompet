package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "net/url"
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

func RenderRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tmpl", gin.H{
		"title": "Register",
	})
}

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

func RenderSaldo(c *gin.Context) {
	// get all user
	users := usaldo.GetRegisteredUser()
	c.HTML(http.StatusOK, "saldo.tmpl", gin.H{
		"users": users,
		"title": "Cek Saldo",
	})
}

func RenderTotalSaldo(c *gin.Context) {
	// get all user
	users := usaldo.GetRegisteredUser()
	c.HTML(http.StatusOK, "totalsaldo.tmpl", gin.H{
		"users": users,
		"title": "Cek Total Saldo",
	})
}

func HandleTransfer(c *gin.Context) {
	// get field from form input
	var p MyParam
	id := c.PostForm("user_id")
	p.Id = id
	ns := c.PostForm("selectbasic")
	nilai_str := c.PostForm("nilai_saldo")
	nilai, err := strconv.ParseInt(nilai_str, 10, 64)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleTransfer ParseInt nilai", err)
	}
	p.Nilai = nilai

	// Check saldo before transfer
	var resp Response
	var st StatusTransfer
	can := usaldo.CheckTransfer(id, nilai)
	if can == 0 {
		resp.Error = 1
		resp.Message = "Saldo not enough"
	} else if can == -1 {
		resp.Error = 1
		resp.Message = "User not exist in this Bank"
	} else {
		// Request Transfer
		st = RequestTranfer(p, ns)
		if st.Status == 0 {
			// SUCCESS: reduce saldo if transfer success
			usaldo.ReduceSaldo(id, nilai)
			resp.Success = 1
			resp.Message = "Succcess Transfer"
		} else {
			// FAIL: try to register
			us := usaldo.GetUserSaldo(id)
			p.Nama = us.Nama
			p.Ip = us.Ip
			resp = RequestRegister(p, ns)

			if resp.Success == 1 {
				// SUCCESS REGISTER: Transfer ulang
				st = RequestTranfer(p, ns)
				if st.Status == 0 {
					// SUCCESS: reduce saldo if transfer success
					usaldo.ReduceSaldo(id, nilai)
					resp.Success = 1
					resp.Message = "Succcess Transfer"
				} else {
					// FAIL REGISTER: fail! :(
					resp.Error = 1
					resp.Message = "Failed to Transfer"
				}
			} else {
				// FAIL REGISTER: fail! :(
				resp.Error = 1
				resp.Message = "Failed to Transfer"
			}

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

func HandleRegister(c *gin.Context) {

	// get field from form input
	var p MyParam
	p.Id = c.PostForm("user_id")
	p.Nama = c.PostForm("nama")
	p.Ip = c.PostForm("ip_domisili")

	rs := RequestRegister(p, "nindyatama")

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

func HandleSaldo(c *gin.Context) {
	// Get param from form
	var p MyParam
	p.Id = c.PostForm("user_id")

	sld := RequestSaldo(p)

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

func HandleTotalSaldo(c *gin.Context) {
	id := c.PostForm("user_id")
	urlstr := fmt.Sprintf("http://nindyatama.sisdis.ui.ac.id/ewallet/getTotalSaldo/%v", id)
	// urlstr := fmt.Sprintf("http://localhost:8080/getTotalSaldo/%v", id)

	resp, _ := http.Get(urlstr)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var sld Saldo
	_ = json.Unmarshal(body, &sld)

	// get all user
	users := usaldo.GetRegisteredUser()
	c.HTML(http.StatusOK, "totalsaldo.tmpl", gin.H{
		"users":       users,
		"title":       "Cek Total Saldo",
		"success":     1,
		"nilai_saldo": sld.Nilai,
	})

}

func RequestSaldo(p MyParam) Saldo {
	// urlstr := fmt.Sprintf("http://nindyatama.sisdis.ui.ac.id/ewallet/getSaldo/%v", id)
	urlstr := "http://localhost:8080/getSaldo"

	var sld Saldo

	pb, err := json.Marshal(p)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleSaldo json Marshal ->", err)
	}

	req, err := http.NewRequest("POST", urlstr, bytes.NewBuffer(pb))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleSaldo Request ->", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleSaldo Read resp Body ->", err)
	}

	err = json.Unmarshal(body, &sld)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleSaldo json Unmarshal ->", err)
	}

	return sld
}

func RequestRegister(p MyParam, ns string) Response {
	// urlstr := fmt.Sprintf("http://%v.sisdis.ui.ac.id/ewallet/%v", ns)
	urlstr := "http://localhost:8080/register"

	var rs Response

	pb, err := json.Marshal(p)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleRegister json Marshal ->", err)
	}

	req, err := http.NewRequest("POST", urlstr, bytes.NewBuffer(pb))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleRegister Request ->", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleRegister Read resp Body ->", err)
	}

	err = json.Unmarshal(body, &rs)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleRegister json Unmarshal ->", err)
	}

	return rs
}

func RequestTranfer(p MyParam, ns string) StatusTransfer {
	// urlstr := fmt.Sprintf("http://%v.sisdis.ui.ac.id/ewallet/transfer", ns)
	urlstr := "http://localhost:8080/transfer"

	var st StatusTransfer

	pb, err := json.Marshal(p)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleRegister json Marshal ->", err)
	}

	req, err := http.NewRequest("POST", urlstr, bytes.NewBuffer(pb))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleTransfer Request ->", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleTransfer Read resp Body ->", err)
	}

	err = json.Unmarshal(body, &st)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleTransfer json Unmarshal ->", err)
	}

	return st
}

//
//
//
//
// OLD
// Handle Saldo
// resp, _ := http.Get(urlstr)
//
// TRANSFER
// res, err := http.PostForm(urlstr, url.Values{"user_id": {id}, "nilai": {nilai_str}})
// 		if err != nil {
// 			log.Println("[ERROR] Handler HandleTransfer PostForm", err)
// 		}
