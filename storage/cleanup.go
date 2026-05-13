package storage

import (
	"fmt"
)

// GetCurrentOutageCount returns the number of currently active outages.
// An outage is considered "active" if outage_end is the MySQL "zero" date.
func GetCurrentOutageCount(c *Connection) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM outages WHERE outage_end = "0000-00-00 00:00:00"`
	err := c.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count current outages: %w", err)
	}
	return count, nil
}

// CleanupOutages purges short blips and old history, then resets IDs.
// It does not require any arguments from the caller, just the connection.
func CleanupOutages(c *Connection) error {
	// 1. Delete short blips (<= 35 seconds)
	_, err := c.DB.Exec(`
        DELETE FROM outages
        WHERE TIMESTAMPDIFF(SECOND, outage_start, outage_end) <= 35
    `)
	if err != nil {
		return fmt.Errorf("failed to delete short blips: %w", err)
	}

	// 2. Delete old history (> 90 days)
	_, err = c.DB.Exec(`
        DELETE FROM outages
        WHERE outage_end < DATE_SUB(NOW(), INTERVAL 90 DAY)
    `)
	if err != nil {
		return fmt.Errorf("failed to delete old history: %w", err)
	}

	// 3. Reset Auto-Increment IDs
	// These must be executed sequentially as separate exec calls
	// Step 3a: Set variable
	_, err = c.DB.Exec("SET @count = 0")
	if err != nil {
		return fmt.Errorf("failed to set count variable: %w", err)
	}

	// Step 3b: Update IDs
	_, err = c.DB.Exec("UPDATE outages SET id = (@count:= @count + 1)")
	if err != nil {
		return fmt.Errorf("failed to update IDs: %w", err)
	}

	// Step 3c: Reset Auto Increment
	_, err = c.DB.Exec("ALTER TABLE outages AUTO_INCREMENT = 1")
	if err != nil {
		return fmt.Errorf("failed to reset auto-increment: %w", err)
	}

	return nil
}
