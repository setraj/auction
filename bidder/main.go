package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	register() // register bidder as soon as it is up
	http.HandleFunc("/", BidHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Println("** Bidder Started on Port " + port + " **")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

type Response struct {
	ID    string `json:"bidder_id"`
	Value int    `json:"bid_value"`
}

type Bidder struct {
	ID       string `json:"bidder_id"`
	Endpoint string `json:"bidder_endpoint"`
}

func register() {
	log.Println("register")
	bidder := Bidder{
		//ID:       os.Getenv("ID"),
		//Endpoint: "http://localhost:" + os.Getenv("PORT"),
		ID:       "alpha",
		Endpoint: "http://localhost:8081/",
	}
	data, err := json.Marshal(bidder)
	if err != nil {
		log.Fatal("Failed to marshal")
	}
	//req, err := http.NewRequest(http.MethodPost, os.Getenv("A_URL"), bytes.NewReader(data))
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/bidder", bytes.NewReader(data))
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error:", err.Error, os.Getenv("A_URL"))
		log.Fatal(err.Error())
	}
	resp.Body.Close()
	log.Println("Bidder " + os.Getenv("ID") + " registered!")
}
func BidHandler(w http.ResponseWriter, r *http.Request) {
	delay, err := strconv.Atoi("510")

	// delay, err := strconv.Atoi("10")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Duration(delay) * time.Millisecond)
	rand.Seed(time.Now().UnixNano())
	resp := Response{
		// ID: os.Getenv("ID"),
		ID:    "alpha",
		Value: rand.Int() % 100,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("Failed to marshal")
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
