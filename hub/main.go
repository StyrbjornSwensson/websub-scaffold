package main


import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type Subscriber struct {
	SubCallback	string
	SubTopic	string
	SubSecret	string
}
var activeSubscribers []Subscriber

func SubRequest(w http.ResponseWriter, r *http.Request) {

 	buf := new(bytes.Buffer)
 	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		log.Fatal(err)
		return
	}

 	queryBuffer := buf.String()
 	decodedString, err := url.QueryUnescape(queryBuffer)

	if err != nil {
		log.Fatal(err)
		return
	}

	params, err := url.ParseQuery(decodedString)
	if err != nil {
		log.Fatal(err)
		return
	}

	confirmationParams := url.Values{}
	confirmationParams.Add("hub.topic", params["hub.topic"][0])
	confirmationParams.Add("hub.mode", params["hub.mode"][0])
	confirmationParams.Add("hub.challenge", createRandomString())
	confirmationParams.Add("hub.lease_seconds", "3600")

	responseParams := &url.URL{
		RawQuery: confirmationParams.Encode(),
	}

	subConfirmationResponse, err := http.Get(params["hub.callback"][0]+ responseParams.String())

	fmt.Println(subConfirmationResponse.StatusCode, "Client Subscribed")

	if subConfirmationResponse.StatusCode == 200 {
		activeSubscriber := Subscriber{
			params["hub.callback"][0],
			params["hub.topic"][0],
			params["hub.secret"][0],
		}
		activeSubscribers = append(activeSubscribers, activeSubscriber)
	}

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

func PublishData(w http.ResponseWriter, r *http.Request) {


	data := createRandomString()
	timeout:= time.Duration(5* time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	
	for _, subscriber := range activeSubscribers {
		requestBody, err := json.Marshal(map[string]string{
			"data": data,
		})
		if err != nil {
			log.Fatal(err)
			return
		}
		request, err := http.NewRequest("POST", subscriber.SubCallback, bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatal(err)
			return
		}

		encryptedSecret := hmac.New(sha256.New, []byte(subscriber.SubSecret))
		fmt.Println("subsecret:" + subscriber.SubSecret)

		encryptedSecret.Write(requestBody)

		signatureValue := "sha256=" + hex.EncodeToString(encryptedSecret.Sum(nil))

		fmt.Println(signatureValue)
		request.Header.Set("Content-type", "application/json")
		request.Header.Add("X-Hub-Signature", signatureValue)

		response, err := client.Do(request)

		fmt.Println(response.StatusCode, "Publish successful")
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}


func main() {
	// Init Router
	router := mux.NewRouter()


	// Route Handlers / Endpoints
	router.HandleFunc("/", SubRequest).Methods("POST")
	router.HandleFunc("/publish", PublishData).Methods("GET")


	log.Fatal(http.ListenAndServe(":8080", router))

}