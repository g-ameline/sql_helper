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

func Update_row(path_to_database, table, column, id string, new_values_by_fields map[string]string) error {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	table = single_quote_text(table)
	single_quote_text_values(new_values_by_fields)
	mb_statement := mb.Convey[*sql.DB, string](mb_db, func() string { return statement_update_row(table, column, id, new_values_by_fields) })
	mb_result := mb.Convey[*sql.DB, sql.Result](mb_db, func() (sql.Result, error) { return mb_db.Value.Exec(mb_statement.Value) })
	return mb_result.Error
}

func Upsert_row(path_to_database, table, column, id string, new_values_by_fields map[string]string) error {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	table = single_quote_text(table)
	single_quote_text_values(new_values_by_fields)
	mb_statement := mb.Convey[*sql.DB, string](mb_db, func() string { return statement_upsert_row(table, column, id, new_values_by_fields) })
	mb_result := mb.Convey[*sql.DB, sql.Result](mb_db, func() (sql.Result, error) { return mb_db.Value.Exec(mb_statement.Value) })
	return mb_result.Error
}

// func statement_upsert_value_simple(table, column, old_value, new_value string) string {
// 	query := fmt.Sprintln("UPDATE", table, "SET", column, "=", new_value, "WHERE", column, "IS", old_value)
// 	query := fmt.Sprintln("INSERT INTO", table, "SET", column, "=", new_value, "WHERE", column, "IS", old_value)
// 	return query
// }

// insert into T values (1, 'one'  ) on conflict(id) do update set val=excluded.val;
// UPDATE table SET
//
//	column1 = value1,
//	column2 = value2
//
// WHERE condition
func statement_update_row(table_name string, column string, id string, values_by_fields map[string]string) string {
	var statement string
	statement += fmt.Sprintln("UPDATE", table_name, "SET")
	update_part := []string{}
	for field, value := range values_by_fields {
		update_part = append(update_part, field+" = "+value)
	}
	statement += strings.Join(update_part, ",\n")
	statement += fmt.Sprintln("\nWHERE", column, "IS", id)
	return statement
}

func statement_upsert_row(table_name string, column string, id string, values_by_fields map[string]string) string {
	var statement string
	statement += fmt.Sprintln("UPDATE OR IGNORE", table_name, "SET")
	update_part := []string{}
	for field, value := range values_by_fields {
		update_part = append(update_part, field+" = "+value)
	}
	statement += strings.Join(update_part, ",\n")
	statement += fmt.Sprintln("\nWHERE", column, "IS", id)

	statement += "INSERT OR IGNORE INTO " + table_name
	fields_parts, values_parts := []string{}, []string{}
	for field, value := range values_by_fields {
		fields_parts = append(fields_parts, field)
		values_parts = append(values_parts, value)
	}
	fields_part := strings.Join(fields_parts, comma)
	values_part := strings.Join(values_parts, comma)
	fields_part = fmt.Sprintln("(", fields_part, ")")
	values_part = fmt.Sprintln("(", values_part, ")")

	statement += fmt.Sprintln(fields_part, "VALUES", values_part)

	return statement
}

// 	fields_part += ") "
// 	values_part += ") "
// 	// then values
// 	statement += fields_part
// 	statement += "VALUES "
// 	statement += values_part
// 	statement += fmt.Sprintln("ON CONFLICT(",conflict_col,") DO UPDATE SET )
// 	return statement
// }
