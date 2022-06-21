package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type People struct {
	ID        int     `json:"id"`
	UID       string  `json:"uid"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Username  string  `json:"username"`
	Address   Address `json:"address"`
}

type Address struct {
	City          string      `json:"city"`
	StreetName    string      `json:"street_name"`
	StreetAddress string      `json:"street_address"`
	ZipCode       string      `json:"zip_code"`
	State         string      `json:"state"`
	Country       string      `json:"country"`
	Coordinates   Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/random-data", getDataRandom)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func getDataRandom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	response, err := http.Get("https://random-data-api.com/api/users/random_user?size=10")

	if err != nil {
		fmt.Println("error to get data")
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("failed to read data")
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	var people []People

	errJSON := json.Unmarshal(responseData, &people)

	if errJSON != nil {
		fmt.Println("Error unmarshal", errJSON)
		return
	}

	jsonData, errMarsh := json.Marshal(&people)

	if errMarsh != nil {
		fmt.Println("Error marshal", errMarsh)
		return
	}

	_, err = w.Write(jsonData)

	if err != nil {
		fmt.Println("Error write data", err)
		return
	}
}
