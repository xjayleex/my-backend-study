package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("mysql", "jay:sql@tcp(localhost:3306)/testdb")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("INSERT users SET name=?, age=?, location=?")
	CheckErr(err)

	_, err = stmt.Exec("jay", "29", "Seoul")
	CheckErr(err)

}
