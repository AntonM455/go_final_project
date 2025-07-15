package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB // global variable for connecting to the database

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(128) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string, logger *log.Logger) error {
	// Check if the file exists
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		install = true // create a file if it does not exist
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// if the file is new, we create a table and index
	if install {
		logger.Println("Create a database and a scheduler table...")
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
