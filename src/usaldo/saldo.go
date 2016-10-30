package usaldo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// -1 -> user not exist
func GetSaldo(id string) (int64, error) {
	var err error
	if u_saldo.Id != id {
		u_saldo, err = getUser(id)
	}
	if err != nil {
		return -1, err
	}
	return u_saldo.getSaldo(id), nil
}

type Saldo struct {
	Nilai int64 `json:"nilai_saldo"`
}

// -1 -> user not exist
func GetTotalSaldo(id string) (int64, error) {
	var err error
	total := int64(0)

	if u_saldo.Id != id {
		u_saldo, err = getUser(id)
	}
	if err == nil {
		// add saldo if user exist
		total = total + u_saldo.Nilai
	}
	log.Println("[CHK] Saldo before ip", total)

	// get from all cabang
	ips, _ := getAllIp()
	log.Println("[CHK] ips", ips)
	for _, ip := range ips {
		// url: ip/ewallet/getSaldo/user_id
		url := fmt.Sprintf("http://%v:8080/getSaldo/%v", ip, id)
		resp, _ := http.Get(url)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var sld Saldo
		_ = json.Unmarshal(body, &sld)
		log.Println("[CHK]", ip, " saldo", sld.Nilai)
		if sld.Nilai > 0 {
			total = total + sld.Nilai
		}
	}
	return total, nil
}

// 0 -> can not transfer
// 1 -> can transfer
// -1 -> user not exist
func CheckTransfer(id string, val int64) int {
	var err error
	if u_saldo.Id != id {
		u_saldo, err = getUser(id)
	}
	if err != nil {
		return -1
	}
	return u_saldo.moreThan(val)
}

func (s *Usaldo) getSaldo(id string) int64 {
	return s.Nilai
}

// 0 -> success transfer
// -1 -> fail transfer
func RecieveTransfer(id string, val int64) int {
	var err error

	if u_saldo.Id != id {
		u_saldo, err = getUser(id)
		if err != nil {
			log.Println("[ERR] ReviceTransfer getUser", err)
			// user not exist
			return -1
		}
	}

	err = u_saldo.addVal(val)
	if err != nil {
		// fail to transfer
		log.Println("[ERR] ReviceTransfer addVal", err)
		return -1
	}

	// sukses transfer
	return 0
}

func ReduceSaldo(id string, val int64) int {
	var err error
	if u_saldo.Id != id {
		u_saldo, err = getUser(id)
	}
	if err != nil {
		// user not exist
		return 0
	}

	err = u_saldo.subVal(val)
	if err != nil {
		// fail to reduce saldo
		return 0
	}

	return 1
}

// 1 -> more / equals val
// 0 -> kurang
func (s *Usaldo) moreThan(val int64) int {
	if s.Nilai >= val {
		return 1
	}
	return 0
}

func (s *Usaldo) addVal(val int64) error {
	newVal := s.Nilai + val
	err := s.updateVal(newVal)
	return err
}

func (s *Usaldo) subVal(val int64) error {
	newVal := s.Nilai - val
	err := s.updateVal(newVal)
	return err
}

func (s *Usaldo) updateVal(val int64) error {
	tx, err := db_main.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        update usaldo
        set nilai_saldo = $2
        where user_id = $1`
	_, err = tx.Exec(query, s.Id, val)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err == nil {
		// update for app variable
		u_saldo.Nilai = val
	}
	return err
}
