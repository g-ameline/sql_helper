package sql_helper

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func Is_database_exist(path_to_db string) bool {
	info, err := os.Stat(path_to_db)
	fmt.Println("info stat", info)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func Get_all_table_names(path_to_database string) (map[string]bool, error) {
	tables := map[string]bool{}
	database, err := sql.Open(database_driver, path_to_database)
	if err != nil {
		return tables, err
	}
	defer database.Close() // good practice
	query := "SELECT name FROM sqlite_schema WHERE type='table'"
	fmt.Println(query)
	rows, err := database.Query(query)
	if_wrong(err, "error while querying all row/record")
	defer rows.Close()
	for rows.Next() {
		var table_name string
		err := rows.Scan(&table_name)
		tables[table_name] = true
		if err != nil {
			return tables, err
		}
	}
	return tables, err
}
