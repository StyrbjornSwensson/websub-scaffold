package main


import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)


type SubRequestResponse struct{
	HubTopic	 	string 	`url:"hub.topic"`
	HubMode 		string 	`url:"hub.mode"`
	HubChallenge	string	`url:"hub.challenge"`
	HubLease		string	`url:"hub.lease_seconds"`
}


func SubRequest(w http.ResponseWriter, r *http.Request) {
 	w.Header().Set("Content-Type", "application/json")

 	buf := new(bytes.Buffer)
 	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		log.Fatal(err)
		return
	}

 	newStr := buf.String()
 	decodedString, err := url.QueryUnescape(newStr)

	if err != nil {
		log.Fatal(err)
		return
	}

	params, err := url.ParseQuery(decodedString)

	subRequestResponse := SubRequestResponse{
		HubTopic: params["hub.topic"][0],
		HubMode: params["hub.mode"][0],
		HubChallenge: createRandomString(),
		HubLease: "3600",
	}

	confirmationParams := url.Values{}
	confirmationParams.Add("hub.topic", subRequestResponse.HubTopic)
	confirmationParams.Add("hub.mode", subRequestResponse.HubMode)
	confirmationParams.Add("hub.challenge", subRequestResponse.HubChallenge)
	confirmationParams.Add("hub.lease_seconds", subRequestResponse.HubLease)

	u := &url.URL{
		RawQuery: confirmationParams.Encode(),
	}

	subConfirmationResponse, err := http.Get(params["hub.callback"][0]+u.String())

	fmt.Println(subConfirmationResponse)
	fmt.Println(params["hub.callback"][0] + u.String())
}



func createRandomString() string {

		rand.Seed(time.Now().UnixNano())
		digits := "0123456789"
		specials := "~=+%^*/()[]{}/!@#$?|"
		all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"abcdefghijklmnopqrstuvwxyz" +
			digits + specials
		length := 8
		buf := make([]byte, length)
		buf[0] = digits[rand.Intn(len(digits))]
		buf[1] = specials[rand.Intn(len(specials))]
		for i := 2; i < length; i++ {
			buf[i] = all[rand.Intn(len(all))]
		}
		rand.Shuffle(len(buf), func(i, j int) {
			buf[i], buf[j] = buf[j], buf[i]
		})
		str := string(buf)
		return str
}

func main() {
	// Init Router
	router := mux.NewRouter()


	// Route Handlers / Endpoints
	router.HandleFunc("/", SubRequest).Methods("POST")


	log.Fatal(http.ListenAndServe(":8080", router))

}