package sql_helper

import (
	"database/sql"
	"errors"
	"fmt"
	mb "github.com/g-ameline/maybe"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

/* some terminology :
if we take the user entity we have
a table users
with column (fields) like username or email
with also rows (records) like id : 2 , username : jojo , email : john@gmail.com, etc.
id is the key and "jojo" is the value matching a column and a row

note:
unless we need to do operation on queried values, it will be handle as string
*/

var database *sql.DB

const database_driver = "sqlite3"

const comma = " , "

const verbose = true

func Is_database_exist(path_to_db string) bool {
	info, err := os.Stat(path_to_db)
	fmt.Println("info stat", info)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func Get_tables_name(path_to_database string) (map[string]bool, error) {
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

func Get_row_two_one_cond(path_to_database string, table_name, field_key, value_key string) (map[string]string, error) {
	database, err := sql.Open(database_driver, path_to_database)
	defer database.Close() // good practice
	var query string
	value_key = single_quote_text(value_key)
	query = query_rows_one_cond(table_name, field_key, value_key)
	breadcrumb(verbose, "query:", query)
	rows, err := database.Query(query)
	defer rows.Close()
	breadcrumb(verbose, "rows from  query :", rows)
	fields, err3 := rows.Columns()
	if err3 != nil { // catch here
		return nil, err
	}
	maybe_values := make([]sql.NullString, len(fields))
	pointers_v := make([]any, len(fields))
	for i := range maybe_values {
		pointers_v[i] = &maybe_values[i]
	}
	if rows.Next() {
		err := rows.Scan(pointers_v...)
		breadcrumb(verbose, "fields", fields)
		ascerted_values := make([]string, len(maybe_values))
		for i, m_v := range maybe_values {
			if m_v.Valid {
				ascerted_values[i] = m_v.String
			}
		}
		breadcrumb(verbose, "values", ascerted_values)
		row_as_map, err := zip_map(fields, ascerted_values)
		return row_as_map, err
	}
	return *new(map[string]string), errors.New("error during row scanning")
}

func Get_rows(path_to_database string, table_name string) (map[string]map[string]string, error) {
	table_as_map := make(map[string]map[string]string)
	m_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(m_db, m_db.Value.Close) // good practice
	m_query := mb.Convey[*sql.DB, string](m_db, func() string { return query_rows(table_name) })
	breadcrumb(verbose, "query for getting all rows", m_query)
	m_rows := mb.Convey[string, *sql.Rows](m_query, func() (*sql.Rows, error) { return m_db.Value.Query(m_query.Value) })
	defer mb.Bind_x_x_e(m_rows, m_rows.Value.Close)
	m_fields := mb.Bind_x_o_e(m_rows, m_rows.Value.Columns)
	breadcrumb(verbose, "fields from table", m_fields.Value)
	if m_fields.Is_error() {
		return *new(map[string]map[string]string), m_fields.Error
	}
	rows := m_rows.Ascertain()
	fields := m_fields.Ascertain()
	var err error
	for rows.Next() {
		var row_as_map map[string]string
		maybe_values := make([]sql.NullString, len(fields))
		pointers_v := make([]any, len(fields))
		for i := range maybe_values {
			pointers_v[i] = &maybe_values[i]
		}
		err = rows.Scan(pointers_v...)
		ascerted_values := make([]string, len(fields))
		for i, maybe_value := range maybe_values {
			if maybe_value.Valid {
				ascerted_values[i] = maybe_value.String
			}
		}
		if_wrong(err, "error during scanning of a row"+" "+table_name)
		row_as_map, _ = zip_map(fields, ascerted_values)
		table_as_map[row_as_map["id"]] = row_as_map
	}
	return table_as_map, err
}

func Get_rows_sorted(path_to_database string, table_name, sorting_field string) ([]map[string]string, error) {
	var table_as_slice []map[string]string
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error opening/creating database")
	defer database.Close() // good practice
	var query string
	query = query_rows_sorted(table_name, sorting_field)
	breadcrumb(verbose, "query for getting all rows", query)
	rows, err := database.Query(query)
	if_wrong(err, "error when fetching all rows")
	defer rows.Close()
	var fields []string
	fields, err = rows.Columns()
	if_wrong(err, "issue when fetching columns names")
	breadcrumb(verbose, "fields from table", fields)
	for rows.Next() {
		var row_as_map map[string]sql.NullString
		values := make([]string, len(fields))
		pointers_v := make([]any, len(fields))
		for i := range values {
			pointers_v[i] = &values[i]
		}
		err = rows.Scan(pointers_v...)
		if_wrong(err, "error during scanning of a row"+" "+table_name+" "+sorting_field)
		breadcrumb(verbose, "fields", fields)
		breadcrumb(verbose, "values", values)
		//deal with null values
		treated_rows := map[string]string{}
		for k, v := range row_as_map {
			if v.Valid {
				treated_rows[k] = v.String
			} else {
				treated_rows[k] = ""
			}
		}
		treated_rows, _ = zip_map(fields, values)
		table_as_slice = append(table_as_slice, treated_rows)
	}
	return table_as_slice, err
}

func query_rows_sorted(table_name, sorting_field string) string {
	return fmt.Sprintln("SELECT * FROM", table_name, "ORDER BY", sorting_field)
}
func query_rows(table_name string) string {
	return fmt.Sprintln("SELECT id FROM", table_name)
}

func query_rows_one_cond(table_name, field, value string) string {
	return fmt.Sprintln("SELECT * FROM", table_name, "WHERE", field, "=", value)
}

func query_rows_two_cond(table_name, field, value, other_field, other_value string) string {
	return fmt.Sprintln("SELECT * FROM", table_name, "WHERE", field, "=", value, "AND", other_field, "=", other_value)
}

func Get_id_one_cond(path_to_database string, table_name, field_key, value_key string) (string, error) {
	var row_as_map map[string]string
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_rows_one_cond(table_name, field_key, single_quote_text(value_key))
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	rows := mb_rows.Ascertain()
	defer rows.Close()
	fields, err := rows.Columns()
	if_wrong(err, "error while reading row")
	values := make([]string, len(fields))
	pointers_v := make([]any, len(fields))
	for i := range values {
		pointers_v[i] = &values[i]
	}
	rows.Next()
	err = rows.Scan(pointers_v...)
	if_wrong(err, "error during scanning of single row to get Id"+" "+table_name+" "+field_key+" "+value_key)
	row_as_map, _ = zip_map(fields, values)
	return row_as_map["Id"], err
}

func Get_id_two_cond(path_to_database string, table_name, field_key, value_key, other_field, other_value string) (string, error) {
	s := single_quote_text
	m_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(m_db, m_db.Value.Close) // good practice
	querying := query_ids_two_cond(table_name, field_key, s(value_key), other_field, s(other_value))
	m_query := mb.Convey[*sql.DB, string](m_db, querying)
	breadcrumb(verbose, "cond ids query:", m_query)
	m_rows := mb.Convey[string, *sql.Rows](m_query, func() (*sql.Rows, error) { return m_db.Value.Query(m_query.Value) })
	// m_rows.Print("inside rows")
	defer mb.Bind_x_x_e(m_rows, m_rows.Value.Close)
	rows := m_rows.Ascertain()
	var err error
	ids := map[string]bool{}
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		ids[id] = true
		if_wrong(err, "error during scanning of a row"+" "+table_name)
	}
	fmt.Println("ids", ids)
	if len(ids) != 1 {
		return "", fmt.Errorf("there was less or more than one result")
	}
	var id string
	for k := range ids {
		id = k
	}
	return id, err
}
func Get_ids_two_cond(path_to_database string, table_name, field_key, value_key, other_field, other_value string) (map[string]bool, error) {
	m_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(m_db, m_db.Value.Close) // good practice
	querying := query_ids_two_cond(table_name, field_key, value_key, other_field, other_value)
	m_query := mb.Convey[*sql.DB, string](m_db, querying)
	breadcrumb(verbose, "cond ids query:", m_query)
	m_rows := mb.Convey[string, *sql.Rows](m_query, func() (*sql.Rows, error) { return m_db.Value.Query(m_query.Value) })
	defer mb.Bind_x_x_e(m_rows, m_rows.Value.Close)
	rows := m_rows.Ascertain()
	var err error
	ids := map[string]bool{}
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		ids[id] = true
		if_wrong(err, "error during scanning of a row"+" "+table_name)
	}
	return ids, err
}

func query_ids_two_cond(table_name, field, value, other_field, other_value string) string {
	return fmt.Sprintln("SELECT id FROM", table_name, "WHERE", field, "=", value, "AND", other_field, "=", other_value)
}
func Is_record_one_cond(path_to_database string, table, field_1, value_1 string) (bool, error) { // for likes or dislikes
	query := fmt.Sprintln("SELECT 1 FROM", table, "WHERE", field_1, "=", value_1)
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error accessing database")
	defer database.Close() // good practice
	rows, err := database.Query(query)
	if_wrong(err, "error while querying all row/record")
	defer rows.Close()
	return rows.Next(), err
}

func Is_record_two_cond(path_to_database string, table, field_1, value_1, field_2, value_2 string) (bool, error) {
	s := single_quote_text
	query := fmt.Sprintln("SELECT 1 FROM", table, "WHERE", field_1, "=", s(value_1), "AND", field_2, "=", s(value_2))
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error accessing database")
	defer database.Close() // good practice
	rows, err := database.Query(query)
	if_wrong(err, "error while querying all row/record")
	defer rows.Close()
	return rows.Next(), err
}

func Count_all_rows(path_to_database string, table_name string) (int, error) {
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error accessing database")
	defer database.Close() // good practice
	query := query_ids(table_name)
	breadcrumb(verbose, "counting statement:", query)
	rows, err := database.Query(query)
	if_wrong(err, "error while querying all row/record")
	defer rows.Close()
	var counter int
	for rows.Next() {
		counter++
	}
	return counter, err
}
func Get_ids(path_to_database string, table_name string) (map[string]bool, error) {
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error accessing database")
	defer database.Close() // good practice
	query := query_ids(table_name)
	breadcrumb(verbose, "counting statement:", query)
	rows, err := database.Query(query)
	if_wrong(err, "error while querying all row/record")
	defer rows.Close()
	ids := map[string]bool{}
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		ids[id] = true
		if_wrong(err, "error during scanning of a row"+" "+table_name)
	}
	return ids, err
}

func query_ids(table_name string) string {
	return fmt.Sprintln("SELECT id FROM", table_name)
}

func zip_map(keys_slice []string, values_slice []string) (map[string]string, error) {
	if len(keys_slice) != len(values_slice) {
		return nil, fmt.Errorf("different length of slices when zipping it")
	}
	if len(keys_slice) == 0 {
		return nil, fmt.Errorf("zero length slice of slices when zipping it")
	}
	keys_values := make(map[string]string)
	for i := 0; i < len(keys_slice); i++ {
		keys_values[keys_slice[i]] = values_slice[i]
	}
	return keys_values, nil
}

func breadcrumb(v bool, helpers ...any) {
	if v {
		for _, h := range helpers {
			fmt.Print(h, " ")
		}
		fmt.Print("\n")
	}
}

func if_wrong(err error, message string) {
	if err != nil {
		println(message, err.Error())
	}
}
func is_wrong(err error, message string) bool {
	if err != nil {
		println(message, err.Error())
		return true
	}
	return false
}

func single_quote_text_values(values_by_fields map[string]string) {
	for field, value := range values_by_fields {
		if value != `''` {
			values_by_fields[field] = single_quote_text(value)
		}
	}
}

func single_quote_text(value string) string {
	_, err := strconv.Atoi(value) // if can be inferred to an int then it is an int
	if err != nil {
		return "'" + value + "'"
	}
	return value
}
