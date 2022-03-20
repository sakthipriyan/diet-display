package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	router.HandleFunc("/api/v0/diet/{id}", ProcessGetRecord).Methods("GET")
	router.HandleFunc("/api/v0/diet/{id}", PutRecord).Methods("PUT")
	router.HandleFunc("/api/v0/diet/{id}", ProcessDeleteRecord).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func PutRecord(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		badRequest(w, "Failed to get id", err)
		return
	}
	var record *Record
	err = json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		badRequest(w, "Failed to process", err)
		return
	}
	if id != record.ID {
		s := fmt.Sprintf("Path variable %v and request body ID %v are different", id, record.ID)
		badRequest(w, "Invalid request", errors.New(s))
		return
	}
	record, err = UpdateRecord(db, *record)
	if err != nil {
		internalServerError(w, err)
		return
	}
	if record == nil {
		notFound(w, id)
		return
	}
	okRequest(w, record)
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
		badRequest(w, "Failed to process", err)
		return
	}
	err = CreateRecords(db, request.Data)
	if err != nil {
		internalServerError(w, err)
		return
	}
	okRequest(w, APIResponse{Status: 200, Message: "Created Diets"})
}

func ProcessGetRecord(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		badRequest(w, "Failed to get id", err)
		return
	}
	record, err := GetRecord(db, id)
	if err != nil {
		internalServerError(w, err)
		return
	}
	if record == nil {
		notFound(w, id)
		return
	}
	okRequest(w, record)

}

func ProcessDeleteRecord(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		badRequest(w, "Failed to get id", err)
		return
	}
	err = DeleteRecord(db, id)
	if err != nil {
		internalServerError(w, err)
		return
	}
	deleteRequest(w)
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

func deleteRequest(w http.ResponseWriter) {
	j := APIResponse{Status: 204, Message: "Deleted"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(j)
}

func badRequest(w http.ResponseWriter, msg string, err error) {
	response := APIResponse{Status: 400, Message: msg + ": " + err.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

func notFound(w http.ResponseWriter, id int) {
	r := APIResponse{Status: 404, Message: strconv.Itoa(id) + " not found"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(r)
}

func internalServerError(w http.ResponseWriter, err error) {
	fmt.Println(err.Error())
	r := APIResponse{Status: 500, Message: "Internal Server Error"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(r)
}
