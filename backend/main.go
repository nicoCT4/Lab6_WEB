package main

import (
	"encoding/json"
	"fmt"
	"time"      
	"log"
	"net/http"   //manejar peticiones
	"strconv"   //convertir tipos
	"github.com/gorilla/mux"  //manejar rutas
	"os" //manejar variables de entorno
	"database/sql" //base de datos
	_ "github.com/lib/pq" // driver de postgres
)

var db *sql.DB


func connectToDB() {
	var err error
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	// Intenta conectarse varias veces
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
					fmt.Println("¡Conectado a PostgreSQL!")
					return
			}
		}
		log.Printf("Intento %d: Error de conexión, reintentando en 5 segundos...", i+1)
		time.Sleep(5 * time.Second)
	}
	log.Fatal("No se pudo conectar a PostgreSQL después de 5 intentos")
}


//Estructura de un partido 
type Match struct {
	ID        int    
	HomeTeam  string 
	AwayTeam  string 
	ScoreA    int    
	ScoreB    int    
	MatchDate string 
}

//crear un partido
func createMatch(w http.ResponseWriter, r *http.Request) {
	var match Match
	_ = json.NewDecoder(r.Body).Decode(&match)

	// Insertar el partido en la base de datos
	sqlStatement := `
		INSERT INTO matches (home_team, away_team, score_a, score_b, match_date)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	id := 0
	err := db.QueryRow(sqlStatement, match.HomeTeam, match.AwayTeam, match.ScoreA, match.ScoreB, match.MatchDate).Scan(&id)
	if err != nil {
		log.Fatal("Error al insertar el partido: ", err)
	}

	match.ID = id
	json.NewEncoder(w).Encode(match)  // Enviar el partido con ID asignado
}

//obtener todos los partidos
func getMatches(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, home_team, away_team, score_a, score_b, match_date FROM matches")
	if err != nil {
		log.Fatal("Error al obtener los partidos:", err)
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		err := rows.Scan(&match.ID, &match.HomeTeam, &match.AwayTeam, &match.ScoreA, &match.ScoreB, &match.MatchDate)
		if err != nil {
			log.Fatal("Error al escanear partido:", err)
		}
		matches = append(matches, match)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error al leer filas:", err)
	}

	json.NewEncoder(w).Encode(matches)
}


//obtener un partido por su id
func getMatchByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var match Match
	err := db.QueryRow("SELECT id, home_team, away_team, score_a, score_b, match_date FROM matches WHERE id = $1", id).
		Scan(&match.ID, &match.HomeTeam, &match.AwayTeam, &match.ScoreA, &match.ScoreB, &match.MatchDate)

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Fatal("Error al obtener el partido:", err)
	}

	json.NewEncoder(w).Encode(match)
}


//actualizar un partido
func updateMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var match Match
	_ = json.NewDecoder(r.Body).Decode(&match)

	// Actualizar el partido en la base de datos
	sqlStatement := `
		UPDATE matches SET home_team=$1, away_team=$2, score_a=$3, score_b=$4, match_date=$5
		WHERE id=$6`
	_, err := db.Exec(sqlStatement, match.HomeTeam, match.AwayTeam, match.ScoreA, match.ScoreB, match.MatchDate, id)
	if err != nil {
		log.Fatal("Error al actualizar el partido:", err)
	}

	match.ID = id
	json.NewEncoder(w).Encode(match)
}


//eliminar un partido
func deleteMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	// Eliminar el partido de la base de datos
	sqlStatement := `DELETE FROM matches WHERE id = $1`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal("Error al eliminar el partido:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Partido eliminado"})
}



// Middleware para permitir CORS (evita el error de fetch desde otro puerto como el 3000)
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

//funcion principal
func main() {
	connectToDB() // Conectar a la base de datos

	router := mux.NewRouter()

	// Servir el HTML (ruta relativa al ejecutable)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "LaLigaTracker.html")
	})

	// Rutas que maneja el backend
	router.HandleFunc("/api/matches", getMatches).Methods("GET")
	router.HandleFunc("/api/matches/{id}", getMatchByID).Methods("GET")
	router.HandleFunc("/api/matches", createMatch).Methods("POST")
	router.HandleFunc("/api/matches/{id}", updateMatch).Methods("PUT")
	router.HandleFunc("/api/matches/{id}", deleteMatch).Methods("DELETE")

	// Imprimir mensaje en consola
	fmt.Println("Servidor corriendo en el puerto 8080")
	// Iniciar el servidor
	log.Fatal(http.ListenAndServe(":8080", enableCors(router)))
}
