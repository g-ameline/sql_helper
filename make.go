package database

import (
	"database/sql"
	"os"
)

// by default function will not overwrite existing element
// if needed some recreate element should be added

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
