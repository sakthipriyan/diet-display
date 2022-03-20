package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// main function to boot up everything
func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/v0/diet/", ToDo).Methods("POST")
	router.HandleFunc("/api/v0/diet/", ToDo).Methods("GET")
	router.HandleFunc("/api/v0/diet/{id}", ToDo).Methods("GET")
	router.HandleFunc("/api/v0/diet/{id}", ToDo).Methods("PUT")
	router.HandleFunc("/api/v0/diet/{id}", ToDo).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

type Response struct {
	status string
}

// CreateUser creates required user
func ToDo(w http.ResponseWriter, r *http.Request) {
	toJson(w, Response{status: "To Do"})
}

func toJson(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
}
