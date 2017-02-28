package main

import (
	"fmt"
	"log"
	"net/http"
	//"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func main() {
	router := NewRouter()
	//time.Sleep(time.Second * 10)
	fmt.Println("starting to connect to DB...")
	db = ConnectToRDS()
	fmt.Println("intermediate check")
	SetupDB(db)
	fmt.Println("connected")
	go Run()
	log.Fatal(http.ListenAndServe(":8080", router))
}
