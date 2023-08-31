package sql_helper

import (
	"database/sql"
	"fmt"

	mb "github.com/g-ameline/maybe"
	_ "github.com/mattn/go-sqlite3"
)

func Delete_rows(path_to_database string, table_name, field, value string) error {
	// database, err := sql.Open(database_driver, path_to_database)
	// defer database.Close() // good practice
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	value = single_quote_text(value)
	single_quote_text(value)
	statement := statement_delete_rows(table_name, field, value)
	breadcrumb(verbose, "deletion statement:", statement)
	mb_result := mb.Convey[*sql.DB, sql.Result](mb_db, func() (sql.Result, error) { return mb_db.Value.Exec(statement) })
	return mb_result.Error
}

func statement_delete_rows(table_name, field, value string) string {
	return fmt.Sprintln("DELETE FROM ", table_name, "WHERE", field, "=", value)
}
