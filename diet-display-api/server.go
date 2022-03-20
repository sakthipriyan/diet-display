package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var db *sql.DB

// main function to boot up everything
func main() {
	var err error
	db, err = OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer CloseDatabase(db)

	router := mux.NewRouter()
	router.HandleFunc("/api/v0/diet/", PostRecords).Methods("POST")
	router.HandleFunc("/api/v0/diet/", GetRecords).Methods("GET")
	router.HandleFunc("/api/v0/diet/{id}", ToDo).Methods("GET")
	router.HandleFunc("/api/v0/diet/{id}", ToDo).Methods("PUT")
	router.HandleFunc("/api/v0/diet/{id}", ToDo).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func GetRecords(w http.ResponseWriter, r *http.Request) {
	records, err := ReadRecords(3)
	if err != nil {
		internalServerError(w, err)
		return
	}
	response := DietResponse{Data: records, Header: defaultHeader}
	okRequest(w, response)
}

func PostRecords(w http.ResponseWriter, r *http.Request) {
	var request DietRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		badRequest(w, APIResponse{Status: 400, Message: "Failed to process" + err.Error()})
		return
	}
	err = CreateRecords(db, request.Data)
	if err != nil {
		internalServerError(w, err)
		return
	}
	okRequest(w, APIResponse{Status: 200, Message: "To Do: Not yet implemented"})
}

// CreateUser creates required user
func ToDo(w http.ResponseWriter, r *http.Request) {
	j := APIResponse{Status: 501, Message: "To Do: Not yet implemented"}
	okRequest(w, j)
}

func okRequest(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
}

func badRequest(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(r)
}

func internalServerError(w http.ResponseWriter, err error) {
	fmt.Println(err.Error())
	r := APIResponse{Status: 500, Message: "Internal Server Error"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(r)
}
