package webhook

import (
	"database/sql"
)

func StoreWebhook(db *sql.DB, name, url, method, hookType, siteID string) error {
	newHookType := emergency
	if hookType == "standard" {
		newHookType = standard
	}

	sql := "insert into webhooks values (null, ?, ?, ?, ?, ?)"

	statement, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = statement.Exec(siteID, name, url, method, newHookType)
	if err != nil {
		return err
	}

	return nil
}
