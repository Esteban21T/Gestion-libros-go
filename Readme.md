Sistema de Gestión de Libros Electrónicos

Fecha: 24 de junio de 2025

Objetivo

Este programa de consola está diseñado para facilitar el manejo de una colección de libros electrónicos. Permite almacenar, buscar y consultar información básica de cada libro (título, autor, género y año) en un archivo JSON local, sirviendo como prototipo de un sistema de gestión bibliotecaria ligero.

Funcionalidades Principales

Registrar un libro (register)

Solicita al usuario los datos de un nuevo libro: título, autor, género y año.

Valida que título y autor no estén vacíos.

Asigna automáticamente un ID incremental.

Guarda los datos en data/books.json.

Buscar libros (search)

Ofrece tres modos de búsqueda:

Por título: coincidencia parcial insensible a mayúsculas.

Por autor: coincidencia parcial insensible a mayúsculas.

Por género: coincidencia exacta ignorando mayúsculas.

Muestra los resultados formateados en consola.

Listar todos los libros (list)

Lee y muestra en pantalla todos los libros registrados, con su ID, título, autor, género y año.

Salir del programa (exit)

Termina la ejecución de la aplicación.

Estructura de Carpetas
gestion-libros-go/
├── cmd/               # Ejecutable de consola
│   └── gestion-libros/
│       └── main.go    # Punto de entrada
├── pkg/               # Lógica y modelos
│   ├── models/book.go # Definición de Book
│   ├── storage/       # Lectura/Escritura de JSON
│   └── service/       # Registro, búsqueda y listado
├── data/              # Almacenamiento local
│   └── books.json     # Base de datos JSON (inicializado con `[]`)
├── go.mod             # Módulo Go
└── README.md          # Documentación del proyecto