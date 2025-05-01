package main

import (
	"fmt"
	"login_project/infrastructure"
	"login_project/model"
)


func main() {
	dbConn := infrastructure.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer infrastructure.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{})
}