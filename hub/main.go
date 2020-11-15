package main


import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type FormValue struct {
	key 	string	`json:"hub.callback"`
	value []string	`json:"value"`
}

type Message struct {
	body string `json:"body"`
}

// create message slice
var messages []Message
var formvals []FormValue


func receiveMessage(w http.ResponseWriter, r *http.Request) {
 	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	var message Message
	_ = json.NewDecoder(r.Body).Decode(message)
	messages = append(messages, message)

	if err != nil {
		log.Fatal(err)
	}
	for key, value := range r.PostForm {
		var formval = FormValue{key, value}
		formvals = append(formvals, formval)
	}

}

func getMessages(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	_= json.NewEncoder(w).Encode(messages)
}
func getFormvals(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	_= json.NewEncoder(w).Encode(formvals)
}




func main() {
	// Init Router
	router := mux.NewRouter()


	// Route Handlers / Endpoints
	router.HandleFunc("/", receiveMessage).Methods("POST")
	router.HandleFunc("/", getMessages).Methods("GET")
	router.HandleFunc("/formvals/", getFormvals).Methods("GET")


	log.Fatal(http.ListenAndServe(":8080", router))

}