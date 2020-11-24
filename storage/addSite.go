package storage

import (
	"time"
)

func AddSite(url string, dbConn *Connection) error {
	sql := "insert into sites values (null, ?, 1, 0, ?, ?)"

	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = statement.Exec(url, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	return err
}
