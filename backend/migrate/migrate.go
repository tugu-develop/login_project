package main

import (
	"fmt"
	"github.com/tugu-develop/login_project/db"
	"github.com/tugu-develop/login_project/model"
)


func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{})
}