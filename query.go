package sql_helper

import (
	"database/sql"
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
	if len(tables) == 0 {
		return map[string]bool{}, fmt.Errorf("no table found")
	}
	return tables, err
}

func Get_row_one_cond(path_to_database string, table_name, field_key, value_key string) (map[string]string, error) {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	value_key = single_quote_text(value_key)
	query := query_rows_one_cond(table_name, field_key, value_key)
	breadcrumb(verbose, "query:", query)
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	mb_fields := mb.Bind_x_o_e(mb_rows, mb_rows.Value.Columns)
	if mb_fields.Is_error() {
		return *new(map[string]string), mb_fields.Error
	}
	return only_one_row(mb_rows.Value, mb_fields.Value)
}

func only_one_row(rows *sql.Rows, fields []string) (map[string]string, error) {
	var empty map[string]string
	breadcrumb(verbose, "fields", fields)
	maybe_values := make([]sql.NullString, len(fields))
	pointers_v := make([]any, len(fields))
	for i := range maybe_values {
		pointers_v[i] = &maybe_values[i]
	}
	rows_as_slice := []map[string]string{}
	if rows.Next() {
		err := rows.Scan(pointers_v...)
		if err != nil {
			return empty, err
		}
		ascerted_values := make([]string, len(maybe_values))
		for i, m_v := range maybe_values {
			if m_v.Valid {
				ascerted_values[i] = m_v.String
			}
		}
		breadcrumb(verbose, "values", ascerted_values)
		row_as_map, err := zip_map(fields, ascerted_values)
		if err != nil {
			return empty, err
		}
		rows_as_slice = append(rows_as_slice, row_as_map)
	}
	if len(rows_as_slice) != 1 {
		return empty, fmt.Errorf("less or more than one row found")
	}
	return rows_as_slice[0], nil
}

func rows_sorted(rows *sql.Rows, fields []string) ([]map[string]string, error) {
	var empty []map[string]string
	breadcrumb(verbose, "fields", fields)
	maybe_values := make([]sql.NullString, len(fields))
	pointers_v := make([]any, len(fields))
	for i := range maybe_values {
		pointers_v[i] = &maybe_values[i]
	}
	rows_as_slice := []map[string]string{}
	if rows.Next() {
		err := rows.Scan(pointers_v...)
		if err != nil {
			return empty, err
		}
		ascerted_values := make([]string, len(maybe_values))
		for i, m_v := range maybe_values {
			if m_v.Valid {
				ascerted_values[i] = m_v.String
			}
		}
		breadcrumb(verbose, "values", ascerted_values)
		row_as_map, err := zip_map(fields, ascerted_values)
		if err != nil {
			return empty, err
		}
		rows_as_slice = append(rows_as_slice, row_as_map)
	}
	return rows_as_slice, nil
}

func rows_by_id(rows *sql.Rows, fields []string) (map[string]map[string]string, error) {
	var empty map[string]map[string]string
	breadcrumb(verbose, "fields", fields)
	maybe_values := make([]sql.NullString, len(fields))
	pointers_v := make([]any, len(fields))
	for i := range maybe_values {
		pointers_v[i] = &maybe_values[i]
	}
	rows_as_map := map[string]map[string]string{}
	if rows.Next() {
		err := rows.Scan(pointers_v...)
		if err != nil {
			return empty, err
		}
		ascerted_values := make([]string, len(maybe_values))
		for i, m_v := range maybe_values {
			if m_v.Valid {
				ascerted_values[i] = m_v.String
			}
		}
		breadcrumb(verbose, "values", ascerted_values)
		row_as_map, err := zip_map(fields, ascerted_values)
		if err != nil {
			return empty, err
		}
		rows_as_map[row_as_map["id"]] = row_as_map
	}
	if len(rows_as_map) == 0 {
		return empty, fmt.Errorf("no table found of that name or table empty")
	}
	return rows_as_map, nil
}

func only_ids(rows *sql.Rows) (map[string]bool, error) {
	var empty map[string]bool
	ids_from_rows := map[string]bool{}
	if rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return empty, err
		}
		ids_from_rows[id] = true
	}
	return ids_from_rows, nil
}

func only_one_id(rows *sql.Rows) (string, error) {
	var empty string
	ids_from_rows := []string{}
	if rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return empty, err
		}
		ids_from_rows = append(ids_from_rows, id)
	}
	if len(ids_from_rows) != 1 {
		return empty, fmt.Errorf("more or less than 1 id found")
	}
	return ids_from_rows[0], nil
}

func Get_rows(path_to_database string, table_name string) (map[string]map[string]string, error) {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_rows(table_name)
	breadcrumb(verbose, "query for getting all rows", query)
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	mb_fields := mb.Bind_x_o_e(mb_rows, mb_rows.Value.Columns)
	breadcrumb(verbose, "fields from table", mb_fields.Value)
	if mb_fields.Is_error() {
		return *new(map[string]map[string]string), mb_fields.Error
	}
	return rows_by_id(mb_rows.Value, mb_fields.Value)
}

func Get_rows_sorted(path_to_database string, table_name, sorting_field string) ([]map[string]string, error) {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_rows_sorted(table_name, sorting_field)
	breadcrumb(verbose, "query for getting all rows", query)
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	mb_fields := mb.Bind_x_o_e(mb_rows, mb_rows.Value.Columns)
	if mb_fields.Is_error() {
		return []map[string]string{}, mb_fields.Error
	}
	breadcrumb(verbose, "fields from table", mb_fields)
	return rows_sorted(mb_rows.Value, mb_fields.Value)
}

func query_rows_sorted(table_name, sorting_field string) string {
	return fmt.Sprintln("SELECT * FROM", table_name, "ORDER BY", sorting_field)
}
func query_rows(table_name string) string {
	return fmt.Sprintln("SELECT * FROM", table_name)
}

func query_rows_one_cond(table_name, field, value string) string {
	return fmt.Sprintln("SELECT * FROM", table_name, "WHERE", field, "=", value)
}

func query_rows_two_cond(table_name, field, value, other_field, other_value string) string {
	return fmt.Sprintln("SELECT * FROM", table_name, "WHERE", field, "=", value, "AND", other_field, "=", other_value)
}

func Get_id_one_cond(path_to_database string, table_name, field_key, value_key string) (string, error) {
	s := single_quote_text
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_ids_one_cond(table_name, field_key, s(value_key))
	breadcrumb(true, query)
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	if mb_rows.Is_error() {
		return "", mb_rows.Error
	}
	return only_one_id(mb_rows.Value)
}

func Get_id_two_cond(path_to_database string, table_name, field_key, value_key, other_field, other_value string) (string, error) {
	s := single_quote_text
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_ids_two_cond(table_name, field_key, s(value_key), other_field, s(other_value))
	breadcrumb(verbose, "cond ids query:", query)
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	if mb_rows.Is_error() {
		return "", mb_rows.Error
	}
	return only_one_id(mb_rows.Value)
}
func Get_ids_two_cond(path_to_database string, table_name, field_key, value_key, other_field, other_value string) (map[string]bool, error) {
	s := single_quote_text
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_ids_two_cond(table_name, field_key, s(value_key), other_field, s(other_value))
	mb_query := mb.Convey[*sql.DB, string](mb_db, query)
	breadcrumb(verbose, "cond ids query:", mb_query)
	mb_rows := mb.Convey[string, *sql.Rows](mb_query, func() (*sql.Rows, error) { return mb_db.Value.Query(mb_query.Value) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	if mb_rows.Is_error() {
		return map[string]bool{}, mb_rows.Error
	}
	return only_ids(mb_rows.Value)
}

func query_ids_one_cond(table_name, field, value string) string {
	return fmt.Sprintln("SELECT id FROM", table_name, "WHERE", field, "=", value)
}
func query_ids_two_cond(table_name, field, value, other_field, other_value string) string {
	return fmt.Sprintln("SELECT id FROM", table_name, "WHERE", field, "=", value, "AND", other_field, "=", other_value)
}
func Is_record_one_cond(path_to_database string, table, field_1, value_1 string) (bool, error) { // for likes or dislikes
	s := single_quote_text
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := fmt.Sprintln("SELECT 1 FROM", table, "WHERE", field_1, "=", s(value_1))
	fmt.Println("query checking 1 record exist", query)
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	if mb_rows.Is_error() {
		return false, mb_rows.Error
	}
	rows := mb_rows.Value
	exist := rows.Next()
	var thingamajig any
	return exist, rows.Scan(&thingamajig)
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
	if len(value)<1 {
		return value
	}
	if value[0] == 39 && value[len(value)-1] == 39 {
		return value
	}
	_, err_a := strconv.Atoi(value) // if can be inferred to an int then it is an int
	if err_a != nil {
		return "'" + value + "'"
	}
	return value
}
