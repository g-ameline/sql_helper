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
const NULL = "NULL"

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
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	mb_fields := mb.Bind_x_o_e(mb_rows, mb_rows.Value.Columns)
	mb_rows_by_id := mb.Bind_x_o_e(mb_rows, func() (map[string]string, error) { return only_one_row(mb_rows.Value, mb_fields.Value) })
	return mb.Relinquish(mb_rows_by_id)
}

func only_one_row(rows *sql.Rows, fields []string) (map[string]string, error) {
	var empty map[string]string
	maybe_values := make([]sql.NullString, len(fields))
	pointers_v := make([]any, len(fields))
	for i := range maybe_values {
		pointers_v[i] = &maybe_values[i]
	}
	rows_as_slice := []map[string]string{}
	for rows.Next() {
		err := rows.Scan(pointers_v...)
		if err != nil {
			return empty, err
		}
		row_as_map, err := zip_nullables_map(fields, maybe_values)
		if err != nil {
			return empty, err
		}
		rows_as_slice = append(rows_as_slice, row_as_map)
		if err := rows.Err(); err != nil {
			return rows_as_slice[0], err
		}
	}
	if len(rows_as_slice) != 1 {
		return empty, fmt.Errorf("less or more than one row found")
	}
	return rows_as_slice[0], nil
}

func rows_sorted(rows *sql.Rows, fields []string) ([]map[string]string, error) {
	var empty []map[string]string
	maybe_values := make([]sql.NullString, len(fields))
	pointers_v := make([]any, len(fields))
	for i := range maybe_values {
		pointers_v[i] = &maybe_values[i]
	}
	rows_as_slice := []map[string]string{}
	for rows.Next() {
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
		row_as_map, err := zip_map(fields, ascerted_values)
		if err != nil {
			return empty, err
		}
		rows_as_slice = append(rows_as_slice, row_as_map)
		if err := rows.Err(); err != nil {
			return rows_as_slice, err
		}
	}
	return rows_as_slice, nil
}

func rows_by_id(rows *sql.Rows, fields []string) (map[string]map[string]string, error) {
	var empty map[string]map[string]string
	maybe_values := make([]sql.NullString, len(fields))
	pointers_v := make([]any, len(fields))
	for i := range maybe_values {
		pointers_v[i] = &maybe_values[i]
	}
	rows_as_map := map[string]map[string]string{}
	for rows.Next() {
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
		row_as_map, err := zip_map(fields, ascerted_values)
		if err != nil {
			return empty, err
		}
		rows_as_map[row_as_map["id"]] = row_as_map
		if err := rows.Err(); err != nil {
			return rows_as_map, err
		}
	}
	if len(rows_as_map) == 0 {
		return empty, fmt.Errorf("no table found of that name or table empty")
	}
	return rows_as_map, nil
}

func only_ids(rows *sql.Rows) (map[string]bool, error) {
	var empty map[string]bool
	ids_from_rows := map[string]bool{}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return empty, err
		}
		ids_from_rows[id] = true
		if err := rows.Err(); err != nil {
			return ids_from_rows, err
		}
	}
	return ids_from_rows, nil
}

func only_one_id(rows *sql.Rows) (string, error) {
	var empty string
	ids_from_rows := []string{}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return empty, err
		}
		if err := rows.Err(); err != nil {
			return id, err
		}
		ids_from_rows = append(ids_from_rows, id)
	}
	if len(ids_from_rows) != 1 {
		return empty, fmt.Errorf("more or less than 1 id found")
	}
	return ids_from_rows[0], nil
}

func is_at_least_one_value(rows *sql.Rows) (answer bool, err error) {
	for rows.Next() {
		var thingamajig any
		err = rows.Scan(&thingamajig)
		if err != nil {
			return answer, err
		}
		if err = rows.Err(); err != nil {
			return answer, err
		}
		answer = true
	}
	return answer, err
}
func is_only_one_value(rows *sql.Rows) (answer bool, err error) {
	var thingies []any
	for rows.Next() {
		var thingamajig any
		err = rows.Scan(&thingamajig)
		if err != nil {
			return answer, err
		}
		if err = rows.Err(); err != nil {
			return answer, err
		}
		thingies = append(thingies, thingamajig)
	}
	if len(thingies) == 1 {
		answer = true
	}
	return answer, err
}

func how_many_rows(rows *sql.Rows) (counter int, err error) {
	for rows.Next() {
		var thingamajig any
		err = rows.Scan(&thingamajig)
		if err != nil {
			return counter, err
		}
		if err = rows.Err(); err != nil {
			return counter, err
		}
		counter++
	}
	return counter, err
}

func Get_rows(path_to_database string, table_name string) (map[string]map[string]string, error) {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_rows(table_name)
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	mb_fields := mb.Bind_x_o_e(mb_rows, mb_rows.Value.Columns)
	mb_rows_by_id := mb.Bind_x_o_e(mb_rows, func() (map[string]map[string]string, error) { return rows_by_id(mb_rows.Value, mb_fields.Value) })
	return mb.Relinquish(mb_rows_by_id)
}

func Get_rows_sorted(path_to_database string, table_name, sorting_field string) ([]map[string]string, error) {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_rows_sorted(table_name, sorting_field)
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	mb_fields := mb.Bind_x_o_e(mb_rows, mb_rows.Value.Columns)
	mb_sorted_rows := mb.Bind_x_o_e(mb_rows, func() ([]map[string]string, error) { return rows_sorted(mb_rows.Value, mb_fields.Value) })
	return mb.Relinquish(mb_sorted_rows)
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
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	return mb.Relinquish(mb.Bind_i_o_e(mb_rows, only_one_id))
}

func Get_ids_one_cond(path_to_database string, table_name, field_key, value_key string) (map[string]bool, error) {
	s := single_quote_text
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_ids_one_cond(table_name, field_key, s(value_key))
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	return mb.Relinquish(mb.Bind_i_o_e(mb_rows, only_ids))
}

func Get_id_two_cond(path_to_database string, table_name, field_key, value_key, other_field, other_value string) (string, error) {
	s := single_quote_text
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_ids_two_cond(table_name, field_key, s(value_key), other_field, s(other_value))
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	return mb.Relinquish(mb.Bind_i_o_e(mb_rows, only_one_id))
}
func Get_ids_two_cond(path_to_database string, table_name, field_key, value_key, other_field, other_value string) (map[string]bool, error) {
	s := single_quote_text
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_ids_two_cond(table_name, field_key, s(value_key), other_field, s(other_value))
	mb_query := mb.Convey[*sql.DB, string](mb_db, query)
	mb_rows := mb.Convey[string, *sql.Rows](mb_query, func() (*sql.Rows, error) { return mb_db.Value.Query(mb_query.Value) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	return mb.Relinquish(mb.Bind_i_o_e(mb_rows, only_ids))
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
	return mb.Relinquish(mb.Bind_i_o_e(mb_rows, is_at_least_one_value))
}

func Is_record_two_cond(path_to_database string, table, field_1, value_1, field_2, value_2 string) (bool, error) {
	s := single_quote_text
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := fmt.Sprintln("SELECT 1 FROM", table, "WHERE", field_1, "=", s(value_1), "AND", field_2, "=", s(value_2))
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	return mb.Relinquish(mb.Bind_i_o_e(mb_rows, is_at_least_one_value))
}

func Count_all_rows(path_to_database string, table_name string) (int, error) {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_ids(table_name)
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	return mb.Relinquish(mb.Bind_i_o_e(mb_rows, how_many_rows))
}
func Get_ids(path_to_database string, table_name string) (map[string]bool, error) {
	mb_db := mb.Mayhaps(sql.Open(database_driver, path_to_database))
	defer mb.Bind_x_x_e(mb_db, mb_db.Value.Close) // good practice
	query := query_ids(table_name)
	mb_rows := mb.Convey[*sql.DB, *sql.Rows](mb_db, func() (*sql.Rows, error) { return mb_db.Value.Query(query) })
	if mb_rows.Is_error() {
		return map[string]bool{}, mb_rows.Error
	}
	defer mb.Bind_x_x_e(mb_rows, mb_rows.Value.Close)
	return mb.Relinquish(mb.Bind_i_o_e(mb_rows, only_ids))
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

func zip_nullables_map(keys_slice []string, nullables_values_slice []sql.NullString) (map[string]string, error) {
	if len(keys_slice) != len(nullables_values_slice) {
		return nil, fmt.Errorf("different length of slices when zipping it")
	}
	if len(keys_slice) == 0 {
		return nil, fmt.Errorf("zero length slice of slices when zipping it")
	}
	keys_values := make(map[string]string)
	for i := 0; i < len(keys_slice); i++ {
		key := keys_slice[i]
		nullable_value := nullables_values_slice[i]
		if nullable_value.Valid {
			keys_values[key] = nullables_values_slice[i].String
		}
	}
	return keys_values, nil
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
	if value == NULL {
		return value
	}
	if value == "" {
		return NULL
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
