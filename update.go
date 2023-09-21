package sql_helper

import (
	"database/sql"
	"fmt"
	mb "github.com/g-ameline/maybe"
	_ "github.com/mattn/go-sqlite3"
)

func Update_value(path_to_database, table, id, column, new_value string) error {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	table = single_quote_text(table)
	new_value = single_quote_text(new_value)
	mb_statement := mb.Convey[*sql.DB, string](mb_db, func() string { return statement_update_value(table, column, new_value, id) })
	breadcrumb(verbose, "statement:", mb_statement)
	mb_result := mb.Convey[*sql.DB, sql.Result](mb_db, func() (sql.Result, error) { return mb_db.Value.Exec(mb_statement.Value) })
	// mb_id_int := mb.Bind_i_o_e(mb_result, sql.Result.LastInsertId)
	// mb_id_string := mb.Convey[int64, string](mb_id_int, func() string { return strconv.FormatInt(mb_id_int.Value, 10) })
	return mb_result.Error
}

func statement_update_value(table, column, value, id string) string {
	query := fmt.Sprintln("UPDATE", table, "SET", column, "=", value, "WHERE id IS", id)
	return query
}
func Update_value_simple(path_to_database, table, column, old_value, new_value string) error {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	table = single_quote_text(table)
	old_value = single_quote_text(old_value)
	new_value = single_quote_text(new_value)
	mb_statement := mb.Convey[*sql.DB, string](mb_db, func() string { return statement_update_value_simple(table, column, old_value, new_value) })
	breadcrumb(verbose, "statement:", mb_statement)
	mb_result := mb.Convey[*sql.DB, sql.Result](mb_db, func() (sql.Result, error) { return mb_db.Value.Exec(mb_statement.Value) })
	// mb_id_int := mb.Bind_i_o_e(mb_result, sql.Result.LastInsertId)
	// mb_id_string := mb.Convey[int64, string](mb_id_int, func() string { return strconv.FormatInt(mb_id_int.Value, 10) })
	return mb_result.Error
}

func statement_update_value_simple(table, column, old_value, new_value string) string {
	query := fmt.Sprintln("UPDATE", table, "SET", column, "=", new_value, "WHERE", column, "IS", old_value)
	return query
}
