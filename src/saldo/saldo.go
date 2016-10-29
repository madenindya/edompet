package saldo

type (
	Saldo struct {
		Id    string `db:"user_id" json:"user_id"`
		Nilai int64  `db:"nilai_saldo" json:"nilai_saldo"`
	}
)

func getSaldo(id string) (Saldo, error) {
	query := `
        select *
        from esaldo
        where user_id = $1`
	var s Saldo
	row := db_main.QueryRowx(query, id)
	err := row.StructScan(&s)
	if err != nil {
		return Saldo{}, err
	}
	return s, nil
}

func (s *Saldo) moreThan(val int64) bool {
	return s.Nilai >= val
}

func (s *Saldo) addVal(val int64) error {
	newVal := s.Nilai + val
	err := s.updateVal(newVal)
	return err
}

func (s *Saldo) subVal(val int64) error {
	newVal := s.Nilai - val
	err := s.updateVal(newVal)
	return err
}

func (s *Saldo) updateVal(val int64) error {
	tx, err := db_main.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        update esaldo
        set nilai_saldo = $2
        where user_id = $1`
	_, err = tx.Exec(query, s.Id, val)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func newVal(id string, val int64) error {
	tx, err := db_main.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        insert into esaldo
        values ($1, $2)`
	_, err = tx.Exec(query, id, val)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
