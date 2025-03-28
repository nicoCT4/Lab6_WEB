package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Match struct {
	ID     int    `json:"id"`
	TeamA  string `json:"teamA"`
	TeamB  string `json:"teamB"`
	ScoreA int    `json:"scoreA"`
	ScoreB int    `json:"scoreB"`
	Date   string `json:"date"`
}

var matches = []Match{}
var nextID = 1

func getMatches(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(matches)
}

func getMatchByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for _, match := range matches {
		if match.ID == id {
			json.NewEncoder(w).Encode(match)
			return
		}
	}
	http.NotFound(w, r)
}

func createMatch(w http.ResponseWriter, r *http.Request) {
	var match Match
	json.NewDecoder(r.Body).Decode(&match)
	match.ID = nextID
	nextID++
	matches = append(matches, match)
	json.NewEncoder(w).Encode(match)
}

func updateMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, match := range matches {
		if match.ID == id {
			json.NewDecoder(r.Body).Decode(&matches[i])
			matches[i].ID = id
			json.NewEncoder(w).Encode(matches[i])
			return
		}
	}
	http.NotFound(w, r)
}

func deleteMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, match := range matches {
		if match.ID == id {
			matches = append(matches[:i], matches[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/matches", getMatches).Methods("GET")
	router.HandleFunc("/api/matches/{id}", getMatchByID).Methods("GET")
	router.HandleFunc("/api/matches", createMatch).Methods("POST")
	router.HandleFunc("/api/matches/{id}", updateMatch).Methods("PUT")
	router.HandleFunc("/api/matches/{id}", deleteMatch).Methods("DELETE")

	fmt.Println("Servidor corriendo en el puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
