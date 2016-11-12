package handler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	RPong struct {
		Pong int `json:"pong"`
	}
)

var master_client *http.Client

func initInsecure() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// client := &http.Client{Transport: tr}
	master_client = &http.Client{Transport: tr}
}

//
//
//
//
// VIEW
// VIEW
//
func RenderQuorum(c *gin.Context) {
	c.HTML(http.StatusOK, "quorum.tmpl", gin.H{
		"title": "Quorum Checker",
	})
}

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

func RenderPing(c *gin.Context) {
	// get all bank tujuan
	ips := usaldo.GetAllUser()
	c.HTML(http.StatusOK, "ping.tmpl", gin.H{
		"title": "Try Ping",
		"ips":   ips,
	})
}

//
//
//
//
// HANDLER
// HANDLER
//

func HandleTransfer(c *gin.Context) {
	// get param from form input
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
	// get param from form input
	var p MyParam
	p.Id = c.PostForm("user_id")
	p.Nama = c.PostForm("nama")
	p.Ip = c.PostForm("ip_domisili")
	log.Println("[CHECK] Register id:", p.Id, "nama:", p.Nama, "ip:", p.Ip)

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
	// Get param from form input
	var p MyParam
	p.Id = c.PostForm("user_id")

	sld := RequestSaldo(p, "nindyatama")

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
	// Get param from form input
	var p MyParam
	p.Id = c.PostForm("user_id")

	sld := RequestTotal(p, "nindyatama")

	// get all user
	users := usaldo.GetRegisteredUser()
	if sld.Nilai == -1 {
		// ERROR -> Quorum !
		c.HTML(http.StatusOK, "totalsaldo.tmpl", gin.H{
			"users":   users,
			"title":   "Cek Total Saldo",
			"error":   1,
			"message": "Failed Quorum Check :(",
		})
	} else {
		c.HTML(http.StatusOK, "totalsaldo.tmpl", gin.H{
			"users":       users,
			"title":       "Cek Total Saldo",
			"success":     1,
			"nilai_saldo": sld.Nilai,
		})
	}

}

func RequestTotal(p MyParam, ns string) Saldo {
	var sld Saldo

	resp, err := sendRequest(ns, "getTotalSaldo", p)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleTotalSaldo Request ->", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleTotalSaldo Read resp Body ->", err)
	}

	err = json.Unmarshal(body, &sld)
	if err != nil {
		log.Println("[CHECK] Error body ->", string(body))
		log.Println("[ERROR] Handler Client HandleTotalSaldo json Unmarshal ->", err)
	}

	return sld
}

func HandlePing(c *gin.Context) {
	ns := c.PostForm("selectbasic")

	ips := usaldo.GetAllUser()
	var msg string

	urlstr := fmt.Sprintf("https://%v.sisdis.ui.ac.id/ewallet/ping", ns)
	req, err := http.NewRequest("POST", urlstr, bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")

	// Request ping
	// client := &http.Client{}
	client := master_client
	resp, err := client.Do(req)
	if err != nil {
		// fail to ping
		msg = fmt.Sprintf("Error Request -> %v", err)

	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// fail read body
			msg = fmt.Sprintf("Error Read body -> %v", err)
		} else {
			msg = string(body)
		}
	}

	c.HTML(http.StatusOK, "ping.tmpl", gin.H{
		"title":   "Try Ping",
		"ips":     ips,
		"message": msg,
	})
}

func RequestSaldo(p MyParam, ns string) Saldo {
	var sld Saldo

	resp, err := sendRequest(ns, "getSaldo", p)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleSaldo Request ->", err)
		sld.Nilai = -1
		return sld
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleSaldo Read resp Body ->", err)
		log.Println("[CHECK] Isi body", string(body))
	}

	log.Println("[CHECK] RequestSaldo body ->", string(body))
	err = json.Unmarshal(body, &sld)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleSaldo json Unmarshal ->", err)
	}

	return sld
}

func RequestRegister(p MyParam, ns string) Response {
	var rs Response

	resp, err := sendRequest(ns, "register", p)
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
	var st StatusTransfer

	resp, err := sendRequest(ns, "transfer", p)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleTransfer Read resp Body ->", err)
	}

	err = json.Unmarshal(body, &st)
	if err != nil {
		log.Println("[ERROR] Handler Client HandleTransfer json Unmarshal ->", err)
		log.Println("[CHECK] Recieve body", string(body))
	}

	return st
}

func sendRequest(ns, action string, p MyParam) (*http.Response, error) {
	urlstr := fmt.Sprintf("https://%v.sisdis.ui.ac.id/ewallet/%v", ns, action)
	log.Println("[CHECK] Request", action, " to ->", ns)

	pb, err := json.Marshal(p)
	if err != nil {
		log.Println("[ERROR] Handler Client Request", action, " json Marshal ->", err)
	}

	req, err := http.NewRequest("POST", urlstr, bytes.NewBuffer(pb))
	req.Header.Set("Content-Type", "application/json")

	client := master_client
	// client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}

//
//
//
//
// QUORUM
// QUORUM
//
func QuorumCheck(c *gin.Context) {
	tpong, sukses := cekQuorum()

	c.HTML(http.StatusOK, "quorum.tmpl", gin.H{
		"title":  "Quorum Checker",
		"pong":   tpong,
		"sukses": sukses,
	})
}

func cekQuorum() (int, string) {
	ips := [8]string{"71", "76", "85", "95", "96", "97", "99", "104"}
	tpong := int(0)
	sukses := ""
	for _, ip := range ips {
		var pong RPong
		urlstr := fmt.Sprintf("http://192.168.75.%v/ewallet/ping", ip)
		log.Println("[CHECK] Ping to ->", urlstr)

		req, err := http.NewRequest("POST", urlstr, bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("[ERROR] Handler Client QuorumCheck Request ->", err)
			continue
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("[ERROR] Handler Client QuorumCheck Read resp Body ->", err)
		}

		err = json.Unmarshal(body, &pong)
		if err != nil {
			log.Println("[ERROR] Handler Client QuorumCheck json Unmarshal ->", err)
		}

		if pong.Pong == 1 {
			tpong = tpong + 1
			sukses = sukses + "192.168.75." + ip + "; \n"
		}
	}

	return tpong, sukses
}

func QuorumCheckFront(c *gin.Context) {
	nss := usaldo.NsKelompok
	tpong := int(0)
	for _, ns := range nss {
		var pong RPong
		urlstr := fmt.Sprintf("https://%v.sisdis.ui.ac.id/ewallet/ping", ns)

		req, err := http.NewRequest("POST", urlstr, bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("[ERROR] Handler Client HandleTotalSaldo Request ->", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("[ERROR] Handler Client HandleTotalSaldo Read resp Body ->", err)
		}

		err = json.Unmarshal(body, &pong)
		if err != nil {
			log.Println("[ERROR] Handler Client HandleSaldo json Unmarshal ->", err)
		}

		if pong.Pong == 1 {
			tpong = tpong + 1
		}
	}

	c.JSON(200, gin.H{
		"pong": tpong,
	})
}
