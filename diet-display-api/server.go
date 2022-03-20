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

// APIResponse used for generic responses
type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// CreateUser creates required user
func ToDo(w http.ResponseWriter, r *http.Request) {
	j := APIResponse{Status: 501, Message: "To Do: Not yet implemented"}
	writeJson(w, j)
}

func writeJson(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
}
