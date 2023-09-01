package sql_helper

import (
	"fmt"
	"testing"
)

const path_to_db = "./social-network.db"

func Test_query(t *testing.T) {
	path_to_nothing := "./nothing"

	fmt.Println("\n__________________________")
	fmt.Println("\nTEST GET TABLES NAME")
	fmt.Println("ok")
	fmt.Println(Get_tables_name(path_to_db))
	fmt.Println("wrong")
	fmt.Println(Get_tables_name(path_to_nothing))

	fmt.Println("\n__________________________")
	fmt.Println("\nTEST GET ONE ROW ONE COND")
	fmt.Println("ok")
	fmt.Println(Get_row_one_cond(path_to_db, "users", "id", "2"))
	fmt.Println("wrong")
	fmt.Println(Get_row_one_cond(path_to_db, "users", "id", "ert"))
	fmt.Println(Get_row_one_cond(path_to_nothing, "users", "id", "ert"))
	fmt.Println(Get_row_one_cond(path_to_db, "users", "idasd", "2"))

	fmt.Println("\n__________________________")
	fmt.Println("\nTEST GET ROWS")
	fmt.Println("ok")
	fmt.Println(Get_rows(path_to_db, "users"))
	fmt.Println("wrong")
	fmt.Println(Get_rows(path_to_db, "usaaaaars"))
	fmt.Println(Get_rows(path_to_nothing, "users"))

	fmt.Println("\n__________________________")
	fmt.Println("\nIS RECORD EXIST")
	fmt.Println("ok")
	user_id := "3"
	fmt.Println(Is_record_one_cond(path_to_db, "users", "id", user_id))
	fmt.Println("wrong")
	fmt.Println(Is_record_one_cond(path_to_db, "users", "idsd", user_id))
	fmt.Println(Is_record_one_cond(path_to_nothing, "users", "id", user_id))
	fmt.Println(Is_record_one_cond(path_to_db, "users", "id", "0908099"))
	fmt.Println(Is_record_one_cond(path_to_db, "users", "id", "NDLSBsd34"))

	fmt.Println("\n__________________________")
	fmt.Println("\nGET ID 1 COND")
	fmt.Println("ok")
	fmt.Println(Get_id_one_cond(path_to_db, "users", "email", "wobi@e.w2"))
	fmt.Println("wrong")
	fmt.Println(Get_id_one_cond(path_to_db, "users", "email", "wobi@e.w2sad"))
	fmt.Println(Get_id_one_cond(path_to_db, "users", "emsdail", "wobi@e.w2sad"))
	fmt.Println(Get_id_one_cond(path_to_db, "userasds", "email", "wobi@e.w2sad"))
	fmt.Println(Get_id_one_cond(path_to_nothing, "users", "email", "wobi@e.w2"))

	fmt.Println("\n__________________________")
	fmt.Println("\nGET ID 2 COND")
	fmt.Println("ok")
	fmt.Println(Get_id_two_cond(path_to_db, "users", "email", "wobi@e.w2", "first_name", "uz"))
	fmt.Println("wrong")
	fmt.Println(Get_id_two_cond(path_to_db, "users", "email", "wobi@e.w2", "first_name", "uasdz"))
	fmt.Println(Get_id_two_cond(path_to_db, "users", "email", "wobi@e.w2", "firadst_name", "uz"))
	fmt.Println(Get_id_two_cond(path_to_db, "users", "email", "wosdbi@e.w2", "first_name", "uz"))
	fmt.Println(Get_id_two_cond(path_to_db, "users", "esadmail", "wobi@e.w2", "first_name", "uz"))
	fmt.Println(Get_id_two_cond(path_to_db, "usdasders", "email", "wobi@e.w2", "first_name", "uz"))
	fmt.Println(Get_id_two_cond(path_to_nothing, "users", "email", "wobi@e.w2", "first_name", "uz"))

	fmt.Println("\n__________________________")
	fmt.Println("\nINSERT 1 ROW")
	fmt.Println("ok")

	a_group := map[string]string{}
	a_group["title"] = "PATATE"
	a_group["creator_id"] = "1"
	fmt.Println(Insert_one_row(path_to_db, "groups", a_group))

	// user_one, _ := Get_row_one_cond(path_to_db, "users", "id", "2")
	// user_x := user_one
	// delete(user_x, "id")
	user_x := map[string]string{}
	user_x["email"] = "balek"
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
