Sistema de Gestión de Libros Electrónicos

Fecha: 24 de junio de 2025

Objetivo

Este programa de consola está diseñado para facilitar el manejo de una colección de libros electrónicos. Permite almacenar, buscar y consultar información básica de cada libro (título, autor, género y año) en un archivo JSON local, sirviendo como prototipo de un sistema de gestión bibliotecaria ligero.

# Gestión de Libros Electrónicos

**Última actualización:** 29 de junio de 2025

## 📋 Descripción del proyecto

Esta aplicación en Go permite gestionar una “biblioteca” de libros electrónicos en dos modos:

1. **CLI** (línea de comandos)  
2. **UI web + API REST**  

Los datos se persisten en un archivo JSON (`data/books.json`).

---

## ⚙️ Funcionalidades principales

### 1. CLI (`cmd/gestion-libros`)
- **register**: agrega un libro (título, autor, género, año)  
- **search**: busca por título, autor o género (insensible a mayúsculas)  
- **list**: muestra todos los libros  
- **exit**: sale de la CLI  

### 2. API REST (`cmd/gestion-libros-server`)
- **GET /health** → `{ "status": "ok" }`  
- **POST /books** → crea un libro  
- **GET /books** → lista todos los libros  
- **GET /books/:id** → obtiene un libro por ID  
- **PUT /books/:id** → actualiza un libro  
- **DELETE /books/:id** → elimina un libro  
- **GET /books/search?title=…** → busca por título  
- **GET /books/genre/:genre** → filtra por género  
- **GET /books/author/:author** → filtra por autor  
- **GET /genres** → lista géneros únicos  

### 3. UI web (`web/templates/index.html`)
- Formulario **Registrar Nuevo Libro**  
- Formulario **Buscar Libros** (título, autor, género)  
- Tabla interactiva con los resultados  

La UI consume la API mediante `fetch()`.

## 📂 Estructura del repositorio
Gestion-libros-go/
├─ cmd/
│ ├─ gestion-libros/ CLI
│ │ └─ main.go
│ └─ gestion-libros-server/ Servidor web + API REST
│ └─ main.go
├─ data/
│ └─ books.json Base de datos JSON
├─ pkg/
│ ├─ models/book.go Definición de Book
│ ├─ service/ Lógica (register, search, list)
│ └─ storage/json_store.go Lectura/escritura de books.json
├─ web/
│ ├─ static/ CSS/JS (opcional)
│ └─ templates/index.html Interfaz HTML
├─ go.mod
└─ go.sum


---

## 🚀 Cómo ejecutar

### Prerrequisitos

- Go 1.20+ instalado

### A. UI web + API REST

```bash
cd /ruta/a/Gestion-libros-go
go mod tidy
go run ./cmd/gestion-libros-server
