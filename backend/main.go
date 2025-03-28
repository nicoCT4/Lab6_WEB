package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"   //manejar peticiones
	"strconv"    //convertir strings a int

	"github.com/gorilla/mux"  //manejar rutas
)

//Estructura de un partido 
type Match struct {
	ID     int    `json:"id"`   
	TeamA  string `json:"teamA"`
	TeamB  string `json:"teamB"`
	ScoreA int    `json:"scoreA"`
	ScoreB int    `json:"scoreB"`
	Date   string `json:"date"`
}

//Lista de partidos
var matches = []Match{}
//control de los id
var currentID int = 1

//obtener todos los partidos
func getMatches(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(matches)  //enviar la lista de partidos en formato json
}

//obtener un partido por su id
func getMatchByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)  //obtener los parametros de la url
	id, _ := strconv.Atoi(params["id"]) //convertir el id a entero
	for _, match := range matches {  //recorrer la lista de partidos
		if match.ID == id {
			json.NewEncoder(w).Encode(match)  //enviar el partido en formato json
			return
		}
	}
	http.NotFound(w, r)  //si no se encuentra el partido, enviar un error 404
}

//crear un partido
func createMatch(w http.ResponseWriter, r *http.Request) {
	var match Match
	json.NewDecoder(r.Body).Decode(&match)  //convertir el json a un objeto Match
	match.ID = currentID  //asignar un id al partido
	currentID++  //incrementar el id
	matches = append(matches, match)
	json.NewEncoder(w).Encode(match)  //enviar el partido en formato json
}

//actualizar un partido
func updateMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)  //obtener los parametros de la url
	id, _ := strconv.Atoi(params["id"])  //convertir el id a entero
	for i, match := range matches {  //recorrer la lista de partidos
		if match.ID == id {  //si se encuentra el partido
			json.NewDecoder(r.Body).Decode(&matches[i])  //convertir el json a un objeto Match
			matches[i].ID = id  //asignar el id al partido
			json.NewEncoder(w).Encode(matches[i])  //enviar el partido en formato json
			return
		}
	}
	http.NotFound(w, r)
}

//eliminar un partido
func deleteMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)  //obtener los parametros de la url
	id, _ := strconv.Atoi(params["id"])  
	for i, match := range matches {
		if match.ID == id {
			matches = append(matches[:i], matches[i+1:]...)  //eliminar el partido de la lista
			w.Header().Set("Content-Type", "application/json")  //enviar un mensaje en formato json
			w.WriteHeader(http.StatusOK)  //enviar un estado 200
			json.NewEncoder(w).Encode(map[string]string{"message": "Partido eliminado"})  //enviar un mensaje en formato json
			return
		}
	}
	http.NotFound(w, r)
}

//funcion principal
func main() {
	router := mux.NewRouter()

	//rutas que maneja el backend
	router.HandleFunc("/api/matches", getMatches).Methods("GET")
	router.HandleFunc("/api/matches/{id}", getMatchByID).Methods("GET")
	router.HandleFunc("/api/matches", createMatch).Methods("POST")
	router.HandleFunc("/api/matches/{id}", updateMatch).Methods("PUT")
	router.HandleFunc("/api/matches/{id}", deleteMatch).Methods("DELETE")

	//imprimir mensaje en consola
	fmt.Println("Servidor corriendo en el puerto 9090")
	//iniciar el servidor
	log.Fatal(http.ListenAndServe(":9090", router))
}
