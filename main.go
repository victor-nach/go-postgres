package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

)

// define the book type and its properties
type Book struct {
	gorm.Model
	// ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Author string `json:"author,omitempty"`
}

type AllBooks []Book

// our own dummy database
var books = AllBooks{
	{
		// ID:     1,
		Name:   "Into the badlands",
		Type:   "adventure",
		Author: "Victor Iheanacho",
	},
	{
		// ID:     2,
		Name:   "50 shades of grey",
		Type:   "romance",
		Author: "Victor Iheanacho",
	},
	{
		// ID:     3,
		Name:   "Charlie and the chocolate factory",
		Type:   "adventure",
		Author: "Emmanuel Iheanacho",
	},
}

type APIResponse struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Error	string		`json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var db *gorm.DB

func main() {

	// Connect to the database
	connectionString := "postgres://xlvhtudc:vmE-X-O8YCDLByk_g1jwV_KacKb_dj2E@raja.db.elephantsql.com:5432/xlvhtudc"
	var err error
	db, err = gorm.Open("postgres", connectionString)

	if err != nil {
		log.Println("Could not connect to the db")
		log.Fatal(err)
	}
	log.Println("Succesfully connected to the db")
	// close connection after the call to the database connection
	defer db.Close()

	db.AutoMigrate(&Book{})
	log.Println("Migrated book table")


	// retuns a *mux.Router instance
	// we can then attach routes(*mux.Route) to it
	// the returned router implements the http.Handler interface so it can as well be served as a handler
	router := mux.NewRouter().StrictSlash(true)

	// perfoms in a similar way to http.handleFunc
	// it is one of the methods that returns a route (*Route) that the mux.Router attaches to it's map
	// it is also a method on the mux.Route interface
	router.HandleFunc("/", welcome).Methods("GET")
	router.HandleFunc("/books", GetAllBooks).Methods("GET")
	router.HandleFunc("/books/{id}", GetSingleBook).Methods("GET")
	router.HandleFunc("/books/{id}", UpdateBook).Methods("PATCH")
	router.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE")
	router.HandleFunc("/books", AddBook).Methods("POST")

	PORT := "3000"
	log.Println("serving on port:", PORT)

	// func ListenAndServe(addr string, handler Handler) error {}
	// ListenAndServe listens on the TCP network address addr and then calls
	// Serve with handler to handle requests on incoming connections.
	// it takes in the port and the handler, if the handler is nil, the DefaultServeMux is used
	// it always returns a non-nil error
	log.Fatal(http.ListenAndServe(":"+PORT, router))

}

// sample handler function, takes in a response writer for the response and a http.Request
// you can write headers and also the content of the response to the writer
func welcome(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Status:  200,
		Message: "Welcome !",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {

	// create empty slice to hold data
	var books []Book
	db.Find(&books) // SELECT * FROM books;

	response := APIResponse{
		Status:  http.StatusOK,
		Message: "All books",
		Data:    books,
	}
	json.NewEncoder(w).Encode(response)
}

func GetSingleBook(w http.ResponseWriter, r *http.Request) {

	// returns the request parameter (routes variables) for the request and nil if none
	// it takes in http.Requeust as parameter, and it returns a map
	// we can then pick the exact value we are looking for from the map
	id := mux.Vars(r)["id"]
	// find by ID 
	var book Book
	// this pattern only works where the id is the primary key and is also an integer 
	err := db.First(&book, id).Error // select all from users where id = id
	
	if err != nil {
		log.Println(err)
		
		// if book not found
		response := APIResponse{
			Status:  http.StatusNotFound,
			Error: "Not found",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// if a match is found return that match
	response := APIResponse{
		Status:  http.StatusOK,
		Message: "All books",
		Data:    book,
	}
	
	json.NewEncoder(w).Encode(response)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	var singleBook Book

	// the read all method reads from a reader until all the data has been read
	// it returns the data a a byte array
	// doesn't return error if the 
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	// unmarshall would only fail for an invalid request body
	// would not fail if there's an extra field or missing field
	// would not also fail if there's no recognizable field
	err = json.Unmarshal(reqBody, &singleBook)
	if err != nil {
		fmt.Println("error:", err)
	}

	// it automatically knows which table to put the data in
	err = db.Create(&singleBook).Error
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusCreated)

	response := APIResponse{
		Status:  http.StatusCreated,
		Message: "book successfuly created",
		Data:    singleBook,
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updatedBook Book
	bookID := mux.Vars(r)["id"]
	bookIDInt, _ := strconv.Atoi(bookID)

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &updatedBook)
	for i, singleBook := range books {
		if singleBook.ID == uint(bookIDInt) {
			singleBook.Name = updatedBook.Name
			singleBook.Type =  updatedBook.Type

			// this is how to update an element in a slice
			// you have to pass in an item belonging in the original slice
			books = append(books[:i], singleBook)

			// if a match is found return that match
			response := APIResponse{
				Status:  http.StatusOK,
				Message: "All books",
				Data:    books,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID := mux.Vars(r)["id"]
	bookIDInt, _ := strconv.Atoi(bookID)

	for i, singleBook := range books {
		if singleBook.ID == uint(bookIDInt) {
			books = append(books[:i], books[i+1:]...)
		}
		// if a match is found return that match
		response := APIResponse{
			Status:  http.StatusOK,
			Message: "events have been successfully deleted",
			Data:    books,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := APIResponse{
		Status:  http.StatusNotFound,
		Error: "Not found",
	}
	json.NewEncoder(w).Encode(response)
}