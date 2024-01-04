package sql_helper

import (
	"database/sql"
	"fmt"
	mb "github.com/g-ameline/maybe"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

func Update_value(path_to_database, table, id, column, new_value string) error {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	table = single_quote_text(table)
	new_value = single_quote_text(new_value)
	mb_statement := mb.Convey[*sql.DB, string](mb_db, func() string { return statement_update_value(table, column, new_value, id) })
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
	mb_result := mb.Convey[*sql.DB, sql.Result](mb_db, func() (sql.Result, error) { return mb_db.Value.Exec(mb_statement.Value) })
	// mb_id_int := mb.Bind_i_o_e(mb_result, sql.Result.LastInsertId)
	// mb_id_string := mb.Convey[int64, string](mb_id_int, func() string { return strconv.FormatInt(mb_id_int.Value, 10) })
	return mb_result.Error
}

func statement_update_value_simple(table, column, old_value, new_value string) string {
	query := fmt.Sprintln("UPDATE", table, "SET", column, "=", new_value, "WHERE", column, "IS", old_value)
	return query
}

func Update_row(path_to_database, table, column_id, id string, new_values_by_fields map[string]string) error {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	table = single_quote_text(table)
	single_quote_text_values(new_values_by_fields)
	mb_statement := mb.Convey[*sql.DB, string](mb_db, func() string { return statement_update_row(table, column_id, id, new_values_by_fields) })
	mb_result := mb.Convey[*sql.DB, sql.Result](mb_db, func() (sql.Result, error) { return mb_db.Value.Exec(mb_statement.Value) })
	return mb_result.Error
}

func Upsert_row(path_to_database, table string, new_values_by_fields map[string]string) error {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	table = single_quote_text(table)
	single_quote_text_values(new_values_by_fields)
	mb_statement := mb.Convey[*sql.DB, string](mb_db, func() string { return statement_upsert_row(table, new_values_by_fields) })
	mb_result := mb.Convey[*sql.DB, sql.Result](mb_db, func() (sql.Result, error) { return mb_db.Value.Exec(mb_statement.Value) })
	return mb_result.Error
}

func statement_update_row(table_name string, column_id string, id string, values_by_fields map[string]string) string {
	var statement string
	statement += fmt.Sprintln("UPDATE", table_name, "SET")
	update_part := []string{}
	for field, value := range values_by_fields {
		update_part = append(update_part, field+" = "+value)
	}
	statement += strings.Join(update_part, ",\n")
	statement += fmt.Sprintln("\nWHERE", column_id, "IS", id)
	return statement
}

func statement_upsert_row(table_name string, values_by_fields map[string]string) string {
	var statement string
	statement += statement_insert_row(table_name, values_by_fields)
	statement += "ON CONFLICT DO UPDATE\n SET "
	update_part := []string{}
	for field, value := range values_by_fields {
		update_part = append(update_part, field+" = "+value)
	}
	statement += strings.Join(update_part, ",\n")
	statement += ";\n"
	return statement
}
