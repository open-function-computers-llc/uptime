package storage

import "github.com/open-function-computers-llc/uptime/models"

func DeleteWebhook(id int, c *Connection) error {
	sql := "DELETE FROM webhooks WHERE id = ?"
	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	if err != nil {
		return err
	}
	return statement.Close()
}

func StoreWebhook(name, url, method, hookType string, siteID int, c *Connection) error {
	newHookType := models.HooktypeEmergency
	if hookType == "standard" {
		newHookType = models.HooktypeStandard
	}

	sql := `INSERT INTO webhooks
			(website_id, hook_name, hook_url, hook_verb, hook_type) VALUES
			(?,          ?,         ?,        ?,         ?)`

	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = statement.Exec(siteID, name, url, method, newHookType)
	if err != nil {
		return err
	}

	return nil
}

func SiteWebhooks(siteID int, c *Connection) ([]*models.Webhook, error) {
	rows, err := c.DB.Query(`
        SELECT id, website_id, hook_name, hook_url, hook_verb, hook_type
        FROM webhooks
        WHERE website_id = ?
    `, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	output := []*models.Webhook{}

	for rows.Next() {
		w := &models.Webhook{}
		var hookType int // to scan into, then cast to models.Hooktype

		err := rows.Scan(&w.ID, &w.WebsiteID, &w.Name, &w.URL, &w.Verb, &hookType)
		if err != nil {
			return nil, err
		}

		w.Type = models.Hooktype(hookType)
		output = append(output, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return output, nil
}
