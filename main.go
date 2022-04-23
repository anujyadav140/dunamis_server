package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model

	Name string
	Email string `gorm:"typevarchar(100);unique_index"`
	Books []Book
}

type Book struct {
	gorm.Model

	Title string
	Author string
	ISBN int `gorm:"unique_index"`
	PersonID int 
}

// var (
// 	person = &Person{Name: "Anuj", Email: "anuj@x.com",}
// 	books = []Book{
// 		{Title: "Moby Dick, or; the whale", Author: "Herman Melville", ISBN: 123456, PersonID: 1},
// 		{Title: "Billy Budd", Author: "Herman Melville", ISBN: 134516, PersonID: 1},
// 	}
// )

var db *gorm.DB
var err error

func main() {

	//loading environment variables

	// host := os.Getenv("HOST")
	// port := os.Getenv("PORT")
	// user := os.Getenv("USER")
	// name := os.Getenv("NAME")
	// password := os.Getenv("PASSWORD")

	//database connection

	dsn := "host=localhost user=postgres password=postgres dbname=dunamis port=5432 sslmode=disable TimeZone=Asia/Calcutta"
	
	//opening connection to database

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to the databbase")
	}

	//close connection to database

	// defer postgres.Close()

	//make migration to the database if they have already been made

	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	router := mux.NewRouter()

	router.HandleFunc("/people", getPeople).Methods("GET")
	router.HandleFunc("/books", getBooks).Methods("GET")

	http.ListenAndServe(":4000", router)

}

func getPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person
	
	db.Find(&people)

	json.NewEncoder(w).Encode(&people)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book

	db.Find(&books)

	json.NewEncoder(w).Encode(&books)
}