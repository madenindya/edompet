package usaldo

import (
	"errors"
	"log"
)

func Register(id, nama, ip string) error {
	var err error
	u_saldo, err = getUser(id)
	if err == nil {
		return errors.New("User existed")
	}

	err = insertNew(id, nama, ip)
	return err
}

func IsExist(id string) bool {
	var err error
	u_saldo, err = getUser(id)
	if err != nil {
		return false
	}

	return true
}

func GetRegisteredUser() []Usaldo {
	regusers := make([]Usaldo, 0)
	var tmpusaldo Usaldo
	query := `
	select user_id, nama
	from usaldo`
	rows, _ := db_main.Queryx(query)
	for rows.Next() {
		err := rows.StructScan(&tmpusaldo)
		if err != nil {
			log.Println(err)
		}
		regusers = append(regusers, tmpusaldo)
	}
	return regusers
}

func GetAllUser() map[string]string {
	return ns_users
}

func getUser(id string) (Usaldo, error) {
	query := `
        select *
        from usaldo
        where user_id = $1`
	var usr Usaldo
	row := db_main.QueryRowx(query, id)
	err := row.StructScan(&usr)
	if err != nil {
		// log.Println("[ERROR] Usaldo getUser", id, ":", err)
		return Usaldo{}, err
	}
	return usr, nil
}

func insertNew(id, nama, ip string) error {
	tx, err := db_main.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        insert into usaldo
        values ($1, $2, $3, 0)`
	if id == "1306381622" {
		// Jika pemilik bank
		query = `
        insert into usaldo
        values ($1, $2, $3, 1000000)`
	}
	_, err = tx.Exec(query, id, nama, ip)
	if err != nil {
		log.Println("[ERROR] Usaldo insertNew", id, ":", err)
		return err
	}
	err = tx.Commit()
	return err
}
