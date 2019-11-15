package main

import (
	"net/http"
	"log"
	// "encoding/json"
	// "fmt"

	"github.com/gorilla/mux"
)

func main() {

	// retuns a *mux.Router instance
	// we can then attach routes(*mux.Route) to it
	// the returned router implements the http.Handler interface so it can as well be served as a handler
	router := mux.NewRouter()

	// perfoms in a similar way to http.handleFunc
	// it is one of the methods that returns a route (*Route) that the mux.Router attaches to it's map
	// it is also a method on the mux.Route interface
	router.HandleFunc("/", welcome)

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

	// the write method here takes in a byte array
	// writes the array to the stream and waits for the stream to be finished reading
	w.Write([]byte("welcome"))

	// another way to write to w using Fprint
	// Fprint writes to w after formatting the operand using it's default formatter
	// fmt.Fprint(w, "welcome")

	// we can also use json.Encode to write json to the response
	// first we create a new encoder using the response stream
	// the Encoder type interface has an Encode method that writes a json interface to a stream
	// the data interface can be a struct, it would use the json tags in the struct to marshall the data
	// json.NewEncoder(w).Encode(json interface{})
}
