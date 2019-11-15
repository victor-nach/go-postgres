package main

import (
	"net/http"
	"encoding/json"
	"log"
	// "encoding/json"
	// "fmt"

	"github.com/gorilla/mux"
)

// define the book type and its properties
type Book struct {
	ID		int		`json:"id,		omitempty"`
	Name 	string	`json:"name,	omitempty"`
	Type 	string	`json:"type,	omitempty"`
	Author  string	`json:"author,	omitempty"`
}

type AllBooks []Book

// our own dummy database 
var books = AllBooks{
	{
		ID: 1,
		Name: "Into the badlands",
		Type: "adventure",
		Author: "Victor Iheanacho",
	},
	{
		ID: 1,
		Name: "50 shades of grey",
		Type: "romance",
		Author: "Victor Iheanacho",
	},
	{
		ID: 1,
		Name: "Charlie and the chocolate factory",
		Type: "adventure",
		Author: "Emmanuel Iheanacho",
	},
}

type APIResponse struct {
	Status		int				`json:"status,omitempty"`
	Message		string			`json:"message,omitempty"`
	Data		interface{}		`json:"data,omitempty"`
}

func main() {

	// retuns a *mux.Router instance
	// we can then attach routes(*mux.Route) to it
	// the returned router implements the http.Handler interface so it can as well be served as a handler
	router := mux.NewRouter().StrictSlash(true)

	// perfoms in a similar way to http.handleFunc
	// it is one of the methods that returns a route (*Route) that the mux.Router attaches to it's map
	// it is also a method on the mux.Route interface
	router.HandleFunc("/", welcome)
	router.HandleFunc("/books", GetAllBooks)

	PORT := "3000"
	log.Println("serving on port:", PORT)

	// func ListenAndServe(addr string, handler Handler) error {}
	// ListenAndServe listens on the TCP network address addr and then calls
	// Serve with handler to handle requests on incoming connections.
	// it takes in the port and the handler, if the handler is nil, the DefaultServeMux is used
	// it always returns a non-nil error
	log.Fatal(http.ListenAndServe(":" + PORT, router))
}

// sample handler function, takes in a response writer for the response and a http.Request
// you can write headers and also the content of the response to the writer
func welcome(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Status: 200,
		Message: "Welcome !",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response);
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Status: 200,
		Message: "All books",
		Data: books,
	}
	json.NewEncoder(w).Encode(response)
}
