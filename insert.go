package sql_helper

import (
	"database/sql"
	mb "github.com/g-ameline/maybe"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

func Insert_one_row(path_to_database string, table_name string, values_by_fields map[string]string) (string, error) {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	single_quote_text_values(values_by_fields)
	mb_statement := mb.Convey[*sql.DB, string](mb_db, func() string { return statement_insert_row(table_name, values_by_fields) })
	breadcrumb(verbose, "statement:", mb_statement)
	mb_result := mb.Convey[*sql.DB, sql.Result](mb_db, func() (sql.Result, error) { return mb_db.Value.Exec(mb_statement.Value) })
	mb_id_int := mb.Bind_i_o_e(mb_result, sql.Result.LastInsertId)
	mb_id_string := mb.Convey[int64, string](mb_id_int, func() string { return strconv.FormatInt(mb_id_int.Value, 10) })
	return mb.Relinquish(mb_id_string)
}

func statement_insert_row(table_name string, values_by_fields map[string]string) string {
	var statement string
	var fields_part, values_part string
	statement += "INSERT INTO " + table_name
	// column first
	fields_part += "( "
	values_part += "( "
	for field, value := range values_by_fields {
		fields_part += field + comma
		values_part += value + comma
	}
	if len(fields_part) > len(comma) && len(values_part) > len(comma) {
		fields_part = fields_part[:len(fields_part)-len(comma)] // remove last comma
		values_part = values_part[:len(values_part)-len(comma)] // remove last comma
	}
	fields_part += ") "
	values_part += ") "
	// then values
	statement += fields_part
	statement += "VALUES "
	statement += values_part
	return statement
}
func Create_database(path_to_db string) error {
	file, err := os.Create(path_to_db)
	defer file.Close()
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func Create_table(path_to_database string, table_name string, types_by_fields map[string]string, constraints map[string]string) error {
	database, err := sql.Open(database_driver, path_to_database)
	if err != nil {
		return err
	}
	defer database.Close() // in case of
	var statement string
	statement = statement_create_table(table_name, types_by_fields, constraints)
	query, err := database.Prepare(statement)
	if_wrong(err, "error creating table "+table_name)
	_, err = query.Exec()
	return err
}

func statement_create_table(table_name string, types_by_fields, constraints map[string]string) string {
	var statement string
	statement += "CREATE TABLE IF NOT EXISTS " + table_name
	statement += " ("
	for field, data_type := range types_by_fields {
		statement += field + " " + data_type + comma
	}
	// add constraints to statement
	for constraint_name, link := range constraints {
		statement += constraint_name + " " + link + comma
	}
	statement = statement[:len(statement)-len(comma)] // remove last comma
	statement += ") "

	breadcrumb(verbose, "create table statement", statement)
	return statement
}

func Add_column_to_table(path_to_database, table_name, field, type_of_field, constraint string) error {
	database, err := sql.Open(database_driver, path_to_database)
	if err != nil {
		return err
	}
	defer database.Close() // in case of
	var statement string
	statement = statement_add_column(table_name, field, type_of_field, constraint)
	query, err := database.Prepare(statement)
	if_wrong(err, "error when adding column to  "+table_name)
	_, err = query.Exec()
	return err
}
func statement_add_column(table_name, field, type_of_field, constraint string) string {
	var statement string
	statement += "ALTER TABLE " + table_name + " ADD COLUMN " + field + " " + type_of_field + " " + constraint
	breadcrumb(verbose, "add column", statement)
	return statement
}
