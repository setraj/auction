package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	bidderMap = make(map[string]string)
)

func main() {

	http.HandleFunc("/", AuctionHandeler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("** Auctioner Started on Port " + port + " **")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

type Bidder struct {
	ID       string `json:"bidder_id"`
	Endpoint string `json:"bidder_endpoint"`
}
type Auction struct {
	ID string `json:"auction_id"`
}

func AuctionHandeler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	//get all registered bidders
	case "GET":
		var bidderEndpoints []string
		for _, v := range bidderMap {
			bidderEndpoints = append(bidderEndpoints, v)
		}
		jsonResp, err := json.Marshal(bidderEndpoints)
		if err != nil {
			log.Fatal("Marshling failed")
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	case "POST":
		log.Println(r.URL.Path)
		switch r.URL.Path {
		//register bidder
		case "/bidder":
			var bidder = Bidder{}
			if err := json.NewDecoder(r.Body).Decode(&bidder); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "\"%s\"", "invalid body")
			}
			bidderMap[bidder.ID] = bidder.Endpoint
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "\"%s\"", "registered")

		//run an auction
		case "/auction":

			var auction = Auction{}
			if err := json.NewDecoder(r.Body).Decode(&auction); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "\"%s\"", "invalid body")
				return
			}

			//this channel will enable the auctioner to respond without
			//having to wait for all the calls to bidders to get finished.
			receiveRes := make(chan BidResponse, 1)

			var wg sync.WaitGroup
			wg.Add(1)
			go accumulateResult(receiveRes, w, &wg)

			startAuction(auction.ID, receiveRes)
			wg.Wait()
		}
	}
}

// here we wait for the final result by listening on "receiveRes" channel and then respond.
func accumulateResult(receiveRes chan BidResponse, w http.ResponseWriter, wg *sync.WaitGroup) {
	defer wg.Done()
	res := <-receiveRes
	if len(res.ID) == 0 {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "\"%s\"", "No bids")
	} else {
		jsonResp, err := json.Marshal(res)
		if err != nil {
			log.Fatal("Marshling failed")
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}

type BidResponse struct {
	ID    string `json:"bidder_id"`
	Value int    `json:"bid_value"`
}

/*
	callBidder will make call to the bidder and send the response
	to the "receive" channel.
*/
func callBidder(receive chan BidResponse, endpoint string, wg *sync.WaitGroup) {
	defer wg.Done()
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)

	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	http.DefaultClient.Timeout = time.Duration(1) * time.Minute
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	bidResp := BidResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&bidResp); err != nil {
		log.Fatal(err.Error())
	}
	receive <- bidResp
	resp.Body.Close()
}

/*

 */
func startAuction(id string, receiveRes chan BidResponse) {
	receive := make(chan BidResponse, 10)
	defer func() {
		close(receive)
	}()
	var wg sync.WaitGroup

	for _, v := range bidderMap {
		wg.Add(1)
		go callBidder(receive, v, &wg)
	}

	var result []BidResponse

	/*
		this for-select loop will wait for 200ms for all the bidders to get completed.
		and select the highest bid.

	*/
breakLoop:
	for {
		select {
		case res := <-receive: // keep receiving the response from bidders
			result = append(result, res)
		case <-time.After(time.Duration(200) * time.Millisecond): //timeout after 200ms and select the highest bid.
			log.Println("Auction over!")
			if len(result) == 0 {
				resp := BidResponse{
					ID:    "",
					Value: -1,
				}
				receiveRes <- resp
			} else {
				var max = -1
				var id = ""
				for _, v := range result {
					if v.Value > max {
						max = v.Value
						id = v.ID
					}
				}

				resp := BidResponse{
					ID:    id,
					Value: max,
				}
				receiveRes <- resp
			}
			break breakLoop
		}
	}
	wg.Wait()
}
