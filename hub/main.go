package main


import (
	"encoding/json"
	//"encoding/json"
	"log"
	"net/http"
	//"math/rand"
	//"strconv"
	"github.com/gorilla/mux"
)



// create message slice
var messages []string


func recieveMessage(w http.ResponseWriter, r *http.Request) {
 	w.Header().Set("Content-Type", "application/json")
	var message string
	_ = json.NewDecoder(r.Body).Decode(&message)
	messages = append(messages, message)
}

func getMessages(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}


func main() {
	// Init Router
	router := mux.NewRouter()


	// Route Handlers / Endpoints
	router.HandleFunc("/", recieveMessage).Methods("POST")
	router.HandleFunc("/", getMessages).Methods("GET")


	log.Fatal(http.ListenAndServe(":8080", router))

}