package main

import (
	"github.com/iamruchy/finalexam/database"
	"github.com/iamruchy/finalexam/handler"
)

func main() {
	database.CreateTable()
	r := handler.NewRouter()
	r.Run(":2019")
}
