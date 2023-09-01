package sql_helper

import (
	"fmt"
	"testing"
)

func Test_insert(t *testing.T) {
	path_to_nothing := "./nothing"
	fmt.Println("\n__________________________")
	fmt.Println("\nINSERT 1 ROW")
	fmt.Println("ok")
	user_one, _ := Get_row_one_cond(path_to_db, "users", "id", "2")
	user_x := user_one
	delete(user_x, "id")
	user_x["email"] += "balek"
	user_x["password"] = "ladsa"
	user_x["first_name"] = "BOOOOBB"
	user_x["private"] = "1"
	fmt.Println(Insert_one_row(path_to_db, "users", user_x))
	fmt.Println(Insert_row(path_to_db, "users", user_x))
	fmt.Println("wrong")
	fmt.Println(Insert_one_row(path_to_nothing, "users", user_x))
	fmt.Println(Insert_one_row(path_to_db, "usesadrs", user_x))
	delete(user_x, "email")
	fmt.Println(Insert_one_row(path_to_db, "users", user_x))
	user_x["email"] = "wobi@e.w2"
	fmt.Println(Insert_one_row(path_to_db, "users", user_x))

}
