package sql_helper

import (
	"fmt"
	"testing"
)

const path_to_db = "./social-network.db"
const path_to_nothing = "./nothing"

func Test_try(t *testing.T) {
	fmt.Println("RIGHT PATH", path_to_db)
	fmt.Println(Is_database_exist(path_to_db))
	fmt.Println("\nWRONG PATH", path_to_nothing)
	fmt.Println(Is_database_exist(path_to_nothing))
	fmt.Println(Get_all_table_names(path_to_db))
	fmt.Println(Get_all_table_names(path_to_nothing))
	fmt.Println("\n update a value  ")
	fmt.Println("how many rows ?")
	fmt.Println(Get_ids(path_to_db, "users"))
	user_id := "3"
	fmt.Println("\nwhat about the number", user_id)
	fmt.Println(Get_one_row(path_to_db, "users", "id", user_id))
	fmt.Println(Get_ids(path_to_db, "users"))
	fmt.Println("update output", Update_value(path_to_db, "users", user_id, "session", "same"))
	fmt.Println(Get_one_row(path_to_db, "users", "id", user_id))
	user_id = "1"
	fmt.Println("i\nwhat about the number", user_id)
	fmt.Println(Get_one_row(path_to_db, "users", "id", user_id))
	fmt.Println(Get_ids(path_to_db, "users"))
	fmt.Println("update output", Update_value(path_to_db, "users", user_id, "session", "same"))
	fmt.Println(Get_one_row(path_to_db, "users", "id", user_id))

}
