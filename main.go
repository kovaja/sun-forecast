package main

import (
	"fmt"
	"kovaja/sun-forecast/db"

	"github.com/joho/godotenv"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func main() {
	checkError(godotenv.Load())
	checkError(db.InitializeDatabase())
	checkError(InitializeServer())
}
