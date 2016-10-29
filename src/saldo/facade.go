package saldo

import (
	"ewallet/src/user"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db_main *sqlx.DB
var saldo Saldo

func Init(db *sqlx.DB) {
	db_main = db
}

// -1 -> user not exist
func GetSaldo(id string) (int64, error) {
	if saldo.Id == id {
		return saldo.Nilai, nil
	}

	sld, err := getSaldo(id)
	if err == nil {
		saldo = sld
		return sld.Nilai, err
	}

	existu := user.IsExist(id)
	if existu {
		err = newVal(id, 0)
		return 0, err
	}

	return -1, err
}

// -1 -> user not exist
func GetTotalSaldo(id string) (int64, error) {

	if saldo.Id == id {
		return saldo.Nilai, nil
	}

	sld, err := getSaldo(id)
	if err == nil {
		saldo = sld
		return sld.Nilai, err
	}

	existu := user.IsExist(id)
	if existu {
		err = newVal(id, 0)
		return 0, err
	}

	return -1, err
}

// 0 -> can not transfer
// 1 -> can transfer
func CheckTransfer(id string, val int64) int {
	sld := saldo
	var err error

	if sld.Id != id {
		sld, err = getSaldo(id)
		if err != nil {
			return 0
		}
		saldo = sld
	}

	c := sld.moreThan(val)
	if !c {
		return 0
	}
	return 1
}

// 0 -> success transfer
// 1 -> fail transfer
func RecieveTransfer(id string, val int64) int {
	var sld Saldo
	var err error
	exist := true

	if saldo.Id != id {
		sld, err = getSaldo(id)
		if err != nil {
			exist = false
		}
	}

	if exist {
		err = sld.addVal(val)
		if err != nil {
			// gagal transfer
			return -1
		}
	} else {
		// cek udah register?
		existu := user.IsExist(id)
		if existu {
			err = newVal(id, val)
			if err != nil {
				return -1
			}
		} else {
			return -1
		}
	}

	// sukses transfer
	return 0
}
