package main

import (
	"fmt"
	"kovaja/sun-forecast/core/db"
	"kovaja/sun-forecast/core/server"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func main() {
	db, err := db.InitializeDatabase()
	checkError(err)
	checkError(server.InitializeServer(db))
}
