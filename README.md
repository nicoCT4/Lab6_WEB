# 📊 LaLiga Tracker - Parte 1 (Backend + Frontend)

Este proyecto es un sistema de seguimiento de partidos de La Liga, desarrollado en Go (Golang) con una API REST y una interfaz web sencilla en HTML + JS.  
Permite crear, consultar, actualizar y eliminar partidos.

---

## 🖼️ Vista del Frontend funcionando

![Vista del frontend](CapturaFront.png)

---

## **API y Frontend para Seguimiento de Partidos**  
Sistema completo para administrar partidos de fútbol con estadísticas en tiempo real.

[Acceder al Sistema »](http://localhost:8080)

---

## **📌 Objetivos**
- **Implementar API REST**: Desarrollar backend en Go con endpoints CRUD
- **Crear interfaz interactiva**: Frontend HTML/JS sin frameworks
- **Gestionar eventos en vivo**: Registrar goles, tarjetas y tiempo extra
- **Documentación profesional**: Incluir Swagger y pruebas Postman

---

## **📖 Descripción del Proyecto**  
Sistema que permite:
- Crear/consultar partidos de LaLiga
- Registrar eventos durante el juego
- Visualizar estadísticas en tiempo real
- Administrar equipos y jugadores

**Tecnologías clave**:  
✔ Golang (Backend)  
✔ PostgreSQL (Base de datos)  
✔ HTML5/CSS/JS (Frontend)  
✔ Docker (Contenedorización)

---

## **📋 Requerimientos Técnicos**

✅ **Backend**:  
- API REST con 10+ endpoints  
- Manejo de CORS  
- Swagger UI integrado  

✅ **Frontend**:  
- Interfaz responsive  
- Funcionalidad sin recargas  
- Compatibilidad con navegadores modernos  

✅ **Base de datos**:  
- Modelo relacional (equipos, partidos, jugadores)  
- Migraciones SQL versionadas  

---

## **🎯 Criterios de Evaluación**

### Backend (50 puntos)
- **20 pts**: Implementación CRUD completa  
- **10 pts**: Documentación Swagger  
- **10 pts**: Manejo de errores  
- **10 pts**: Estructura de código limpia  

### Frontend (30 puntos)
- **15 pts**: Funcionalidad interactiva  
- **10 pts**: Diseño responsive  
- **5 pts**: Consumo correcto de API  

### Documentación (20 puntos)
- **10 pts**: README profesional  
- **5 pts**: Colección Postman  
- **5 pts**: Comentarios en código  

---

## **📂 Estructura del Proyecto**
LaLigaTracker/
├── backend/
│ ├── main.go
│ ├── go.mod
│ ├── go.sum
│ ├── migrations/
│ └── 
├── DockerFile
├── docker-compose.yml
├── LaLigaTracker.html
├── capturaFront.png
└── README.md

## 🧪 Pruebas con cURL

### Operaciones CRUD
```bash
# Crear partido
curl -X POST http://localhost:8080/api/matches \
  -H "Content-Type: application/json" \
  -d '{"homeTeam":"Barcelona","awayTeam":"Real Madrid","matchDate":"2025-06-15"}'

# Listar partidos
curl http://localhost:8080/api/matches

# Obtener partido específico
curl http://localhost:8080/api/matches/1

# Actualizar partido
curl -X PUT http://localhost:8080/api/matches/1 \
  -H "Content-Type: application/json" \
  -d '{"homeTeam":"Atlético Madrid","awayTeam":"Sevilla"}'

# Eliminar partido
curl -X DELETE http://localhost:8080/api/matches/1

## 📬 Contacto

**Desarrollador:**  
👨‍💻 Nicolás Concuá  
📧 [con23197@uvg.edu.gt](mailto:con23197@uvg.edu.gt)  
🔗 [GitHub](https://github.com/tuusuario)  


**Institución:**  
🏫 Universidad del Valle de Guatemala  
📚 Ingeniería en Ciencias de la Computación    
