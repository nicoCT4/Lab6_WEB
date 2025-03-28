# 📊 LaLiga Tracker - Parte 1 (Backend + Frontend)

Este proyecto es un sistema de seguimiento de partidos de La Liga, desarrollado en Go (Golang) con una API REST y una interfaz web sencilla en HTML + JS.  
Permite crear, consultar, actualizar y eliminar partidos.

---

## 📁 Estructura del Proyecto
LAB6_WEB/ ├── backend/ │ ├── Dockerfile │ ├── docker-compose.yml.example │ ├── go.mod │ ├── go.sum │ └── main.go ├── LaLigaTracker.html └── CapturaFront.png

## 🖼️ Vista del Frontend funcionando

![Vista del frontend](CapturaFront.png)

---

## 🧪 Pruebas con curl

#Crear un partido:
curl -X POST http://localhost:8080/api/matches \
  -H "Content-Type: application/json" \
  -d '{"homeTeam":"Barcelona","awayTeam":"Real Madrid","scoreA":2,"scoreB":1,"matchDate":"2025-04-01"}'

#Listar partidos:
curl http://localhost:8080/api/matches

#Obtener partido por ID:
curl http://localhost:8080/api/matches/1

#Actualizar partido:
curl -X PUT http://localhost:8080/api/matches/1 \
  -H "Content-Type: application/json" \
  -d '{"homeTeam":"Atletico","awayTeam":"Sevilla","scoreA":1,"scoreB":1,"matchDate":"2025-04-05"}'

#Eliminar partido:
curl -X DELETE http://localhost:8080/api/matches/1


## 🌐 Cómo levantar el frontend
Abre el archivo LaLigaTracker.html en tu navegador.

O bien...

Usa Live Server en VSCode:

Haz clic derecho en LaLigaTracker.html

Selecciona "Open with Live Server"

Asegúrate de que el backend esté corriendo en http://localhost:8080 ya que es ahí donde el HTML busca la API.

## Autor
Nicolás Concuá UVG-2025
