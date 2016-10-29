package usaldo

func Register(id, nama, ip string) error {
	var err error
	u_saldo, err = getUser(id)
	if err == nil {
		// error nil -> user exist (by id)
		return err
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

func getUser(id string) (Usaldo, error) {
	query := `
        select *
        from usaldo
        where user_id = $1`
	var usr Usaldo
	row := db_main.QueryRowx(query, id)
	err := row.StructScan(&usr)
	if err != nil {
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
	_, err = tx.Exec(query, id, nama, ip)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}
