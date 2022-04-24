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

	Name  string
	Email string `gorm:"type:VARCHAR(100);uniqueIndex"`
	Books []Book
}

type Book struct {
	gorm.Model

	Title    string
	Author   string
	ISBN     int `gorm:"uniqueIndex"`
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

	//make migration to the database if they have already been made

	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	router := mux.NewRouter()

	router.HandleFunc("/people", getPeople).Methods("GET")
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/people/{id}", getPerson).Methods("GET")
	router.HandleFunc("/people/create", createPerson).Methods("POST")
	router.HandleFunc("/people/delete/{id}", deletePerson).Methods("DELETE")
	router.HandleFunc("/books/delete/{id}", deleteBook).Methods("DELETE")
	router.HandleFunc("/people/books", postBooks).Methods("POST")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/people/change/{id}", changeName).Methods("PUT")

	log.Fatal(http.ListenAndServe(":4000", router))

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

func getPerson(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var person Person
	var books []Book

	db.First(&person, params["id"])
	db.Model(&person).Association("Books").Find(&books)

	person.Books = books

	json.NewEncoder(w).Encode(&person)

}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person Person

	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := db.Create(&person)
	err := createdPerson.Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&person)
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person

	db.First(&person, params["id"])
	db.Delete(&person)

	json.NewEncoder(w).Encode(&person)
}

func postBooks(w http.ResponseWriter, r *http.Request) {
	var books Book

	json.NewDecoder(r.Body).Decode(&books)

	createdBooks := db.Create(&books)
	err := createdBooks.Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&books)
	}

}

func getBook(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var book Book

	db.First(&book, params["id"])

	json.NewEncoder(w).Encode(&book)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book

	db.First(&book, params["id"])
	db.Delete(&book)

	json.NewEncoder(w).Encode(&book)
}

func changeName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person

	db.First(&person, params["id"])

	json.NewDecoder(r.Body).Decode(&person)

	db.Model(&person).Update("Name", &person.Name)

}
