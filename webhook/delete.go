package webhook

import "database/sql"

func Delete(db *sql.DB, id int) error {
	sql := "DELETE FROM webhooks WHERE id = ?"
	statement, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	if err != nil {
		return err
	}
	return statement.Close()
}
