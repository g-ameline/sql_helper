package sql_helper

import (
	"fmt"
	"testing"
)

const path_to_db = "./social-network.db"

func Test_query(t *testing.T) {
	path_to_nothing := "./nothing"
	fmt.Println("RIGHT PATH", path_to_db)
	fmt.Println(Is_database_exist(path_to_db))
	fmt.Println("\nWRONG PATH", path_to_nothing)
	fmt.Println(Is_database_exist(path_to_nothing))
	fmt.Println(Get_tables_name(path_to_db))
	fmt.Println(Get_tables_name(path_to_nothing))
	fmt.Println("\n update a value  ")
	fmt.Println("how many rows ?")
	fmt.Println(Get_ids(path_to_db, "users"))
	user_id := "3"
	fmt.Println("\nwhat about the number", user_id)
	fmt.Println(Get_row_one_cond(path_to_db, "users", "id", user_id))
	fmt.Println(Get_ids(path_to_db, "users"))
	fmt.Println("update output", Update_value(path_to_db, "users", user_id, "session", "same"))
	fmt.Println(Get_row_one_cond(path_to_db, "users", "id", user_id))
	user_id = "1"
	fmt.Println("i\nwhat about the number", user_id)
	fmt.Println(Get_row_one_cond(path_to_db, "users", "id", user_id))
	fmt.Println(Get_ids(path_to_db, "users"))
	fmt.Println("TESTING HERE")
	fmt.Println(Get_id_one_cond(path_to_db, "users", "session", "same"))
	res, err := Get_id_one_cond(path_to_db, "users", "session", "prout")
	fmt.Println("res", res, "err", err)
	fmt.Println("update output", Update_value(path_to_db, "users", user_id, "session", "same"))
	fmt.Println(Get_row_one_cond(path_to_db, "users", "id", user_id))
	fmt.Println("\ntry the det 1 id 2 cond stuff")
	user_id = "3"
	fmt.Println(Get_id_two_cond(path_to_db, "users", "id", user_id, "session", "same"))
}
