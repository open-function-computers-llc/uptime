package storage

// BootstrapSites - make sure that the DB is ready to go
func BootstrapSites(dbConn *Connection) error {
	err := createTables(dbConn)
	if err != nil {
		return err
	}
	return nil
}

func createTables(dbConn *Connection) error {
	// create sites table
	sql := `CREATE TABLE IF NOT EXISTS sites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url VARCHAR NOT NULL,
		is_up INTEGER DEFAULT 1,
		is_deleted INTEGER DEFAULT 0,
		created_at TEXT,
		last_checked TEXT,
		CONSTRAINT url_unique UNIQUE (url)
	);`
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		return err
	}
	// create
	statement.Exec()

	// create outages table
	sql = `CREATE TABLE IF NOT EXISTS outages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		website_id INTEGER NOT NULL,
		outage_start TEXT,
		outage_end TEXT
	);`
	statement, err = dbConn.DB.Prepare(sql)
	if err != nil {
		return err
	}
	// create
	statement.Exec()
	statement.Close()

	return nil
}
