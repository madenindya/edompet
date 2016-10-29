package user

type (
	User struct {
		Id   string `db:"user_id" json:"user_id"`
		Nama string `db:"nama" json:"nama"`
		Ip   string `db:"ip_domisili" json:"ip_domisili"`
	}
)

func getUser(id string) (User, error) {
	query := `
        select *
        from euser
        where user_id = $1`
	var usr User
	row := db_main.QueryRowx(query, id)
	err := row.StructScan(&usr)
	if err != nil {
		return User{}, err
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
        insert into euser
        values ($1, $2, $3)`
	_, err = tx.Exec(query, id, nama, ip)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
