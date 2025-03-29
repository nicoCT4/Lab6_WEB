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
	ID        int    `json:"id"`
	HomeTeam  string `json:"homeTeam"`
	AwayTeam  string `json:"awayTeam"`
	ScoreA    int    `json:"scoreA"`
	ScoreB    int    `json:"scoreB"`
	MatchDate string `json:"matchDate"`
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
		http.Error(w, "Error al obtener los partidos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		err := rows.Scan(&match.ID, &match.HomeTeam, &match.AwayTeam, &match.ScoreA, &match.ScoreB, &match.MatchDate)
		if err != nil {
			http.Error(w, "Error al escanear partido: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		 // Formatear fecha de manera segura
		if len(match.MatchDate) >= 10 {
			match.MatchDate = match.MatchDate[:10]
		}
		
		matches = append(matches, match)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error al leer filas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}


//obtener un partido por su id
func getMatchByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID de partido inválido", http.StatusBadRequest)
		return
	}

	var match Match
	err = db.QueryRow(
		"SELECT id, home_team, away_team, score_a, score_b, match_date FROM matches WHERE id = $1", 
		id,
	).Scan(&match.ID, &match.HomeTeam, &match.AwayTeam, &match.ScoreA, &match.ScoreB, &match.MatchDate)

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Error al obtener el partido: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Formatear fecha de manera segura
	if len(match.MatchDate) >= 10 {
		match.MatchDate = match.MatchDate[:10]
	}

	w.Header().Set("Content-Type", "application/json")
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

func registerGoal(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	matchID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID de partido inválido", http.StatusBadRequest)
		return
	}

	// Decodificar el cuerpo de la petición
	var request struct {
		 TeamID int `json:"teamId"` // ID del equipo que anotó
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Determinar si es gol del equipo local o visitante
	var isHomeTeam bool
	err = db.QueryRow(
		"SELECT $1 = home_team FROM matches WHERE id = $2", 
		request.TeamID, matchID,
	).Scan(&isHomeTeam)
	if err != nil {
		http.Error(w, "Error al verificar equipos", http.StatusInternalServerError)
		return
}

	// Actualizar el marcador
	var updateQuery string
	if isHomeTeam {
		updateQuery = "UPDATE matches SET score_a = score_a + 1 WHERE id = $1"
	} else {
		updateQuery = "UPDATE matches SET score_b = score_b + 1 WHERE id = $1"
	}

	_, err = db.Exec(updateQuery, matchID)
	if err != nil {
		http.Error(w, "Error al actualizar marcador", http.StatusInternalServerError)
		return
	}

	// Registrar el gol en la tabla goals
	_, err = db.Exec(
		"INSERT INTO goals (match_id, team_id, goals) VALUES ($1, $2, 1)",
		matchID, request.TeamID,
	)
	if err != nil {
		http.Error(w, "Error al registrar gol", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Gol registrado correctamente"})
}

func registerCard(w http.ResponseWriter, r *http.Request, cardType string) {
	params := mux.Vars(r)
	matchID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID de partido inválido", http.StatusBadRequest)
		return
	}

	var request struct {
		PlayerID int `json:"playerId"`
		Minute   int `json:"minute"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Registrar la tarjeta
	_, err = db.Exec(
		"INSERT INTO cards (player_id, match_id, card_type, card_time) VALUES ($1, $2, $3, $4)",
		request.PlayerID, matchID, cardType, request.Minute,
	)
	if err != nil {
		http.Error(w, "Error al registrar tarjeta", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Tarjeta %s registrada", cardType)})
}

// Handlers específicos
func registerYellowCard(w http.ResponseWriter, r *http.Request) {
	registerCard(w, r, "yellow")
}

func registerRedCard(w http.ResponseWriter, r *http.Request) {
	registerCard(w, r, "red")
}

func setExtraTime(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	matchID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID de partido inválido", http.StatusBadRequest)
		return
	}

	var request struct {
		Minutes int `json:"minutes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Actualizar en ambas tablas por consistencia
	_, err = db.Exec("UPDATE matches SET extra_time = $1 WHERE id = $2", request.Minutes, matchID)
	if err != nil {
		http.Error(w, "Error al actualizar partido", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(
		"INSERT INTO extra_time (match_id, time) VALUES ($1, $2) ON CONFLICT (match_id) DO UPDATE SET time = $2",
		matchID, request.Minutes,
	)
	if err != nil {
		http.Error(w, "Error al registrar tiempo extra", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Tiempo extra establecido a %d minutos", request.Minutes)})
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
	router.HandleFunc("/api/matches/{id}/goals", registerGoal).Methods("PATCH")
	router.HandleFunc("/api/matches/{id}/yellowcards", registerYellowCard).Methods("PATCH")
	router.HandleFunc("/api/matches/{id}/redcards", registerRedCard).Methods("PATCH")
	router.HandleFunc("/api/matches/{id}/extratime", setExtraTime).Methods("PATCH")

	// Imprimir mensaje en consola
	fmt.Println("Servidor corriendo en el puerto 8080")
	// Iniciar el servidor
	log.Fatal(http.ListenAndServe(":8080", enableCors(router)))
}
