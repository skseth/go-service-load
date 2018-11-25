package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Chunk struct {
    Id        string   `json:"id,omitempty"`
    Value string   `json:"value,omitempty"`
}

var chunks []Chunk

func initChunks() {
	chunks = append(chunks, Chunk{Id: "1", Value: "abcdedffff"})
	chunks = append(chunks, Chunk{Id: "2", Value: "cdedffffdddd"})
	chunks = append(chunks, Chunk{Id: "3", Value: "ffffdddddddddd"})
}

func GetChunk(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range chunks {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Chunk{})
}

func main() {
	router := mux.NewRouter()
	initChunks()
	router.HandleFunc("/chunk/{id}", GetChunk).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

