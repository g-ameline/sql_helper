package database

import (
	"database/sql"
	// _ "github.com/mattn/go-sqlite3"
	"os"
)

func Create_database_if_not_already(path_to_database string) error {
	if !(Is_databse_exist(path_to_database)) {
		Create_database(path_to_database)
	} // check if all ogod
	return Try_open_close_database(path_to_database)
}

func Is_databse_exist(path_to_db string) bool {
	info, err := os.Stat(path_to_db)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func Try_open_close_database(path_to_database string) error {
	database, err := sql.Open(database_driver, path_to_database)
	if err != nil {
		return err
	}
	return database.Close() // good practice
}
