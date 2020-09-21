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
	_, err := sql.Open("mysql", "jay:sql@tcp(testbed:3306)/test")
	CheckErr(err)
}
