package sql_helper

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_insert(t *testing.T) {
	path_to_nothing := "./nothing"
	fmt.Println("\n__________________________")
	fmt.Println("\nINSERT 1 ROW")
	fmt.Println("ok")
	user_one, _ := Get_row_one_cond(path_to_db, "users", "id", "5")
	user_x := user_one
	delete(user_x, "id")
	user_x["email"] += randomstring(23)
	user_x["password"] = "ladsa"
	user_x["session"] += randomstring(20)
	user_x["first_name"] = ""
	user_x["private"] = "1"
	fmt.Println(Insert_one_row(path_to_db, "users", user_x))
	fmt.Println("wrong")
	fmt.Println(Insert_one_row(path_to_nothing, "users", user_x))
	fmt.Println(Insert_one_row(path_to_db, "usesadrs", user_x))
	delete(user_x, "email")
	fmt.Println(Insert_one_row(path_to_db, "users", user_x))
	user_x["email"] = "wobi@e.w2"
	fmt.Println(Insert_one_row(path_to_db, "users", user_x))
	fmt.Println("..................")
	fmt.Println("experiment insert NULL as value")
	user_a := map[string]string{}
	email_a := randomstring(23)
	user_a["email"] = email_a
	user_a["password"] = "modepasse"
	user_a["session"] = "NULL"
	user_a["first_name"] = ""
	user_a["private"] = "1"
	fmt.Println(Insert_one_row(path_to_db, "users", user_a))
	fmt.Println(Get_row_one_cond(path_to_db, "users", "email", email_a))
	fmt.Println("experiment insert NULL as value")
	user_b := map[string]string{}
	email_b := randomstring(23)
	user_b["email"] = email_b
	user_b["password"] = "mo2passe"
	user_b["first_name"] = ""
	user_b["private"] = "1"
	fmt.Println(Insert_one_row(path_to_db, "users", user_b))
	fmt.Println(Get_row_one_cond(path_to_db, "users", "email", email_b))
	// fmt.Println("experiment insert null value")

}
func randomstring(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = "nadwqg0-23-2108hskndseiuo39482 dfwq2387 klrqwjbvg&"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
