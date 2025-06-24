package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Esteban21T/Sistemas-de-Gesti-n-empresarial.git/pkg/models"
	"github.com/Esteban21T/Sistemas-de-Gesti-n-empresarial.git/pkg/service"
	"github.com/Esteban21T/Sistemas-de-Gesti-n-empresarial.git/pkg/storage"
	"github.com/gin-gonic/gin"
)

const dataFile = "data/books.json"

var (
	books []models.Book
	mu    sync.Mutex
)

// loadData carga los libros desde el archivo JSON.
func loadData() error {
	b, err := storage.LoadBooks(dataFile)
	if err != nil {
		return err
	}
	books = b
	return nil
}

// saveData guarda los libros en el archivo JSON.
func saveData() error {
	return storage.SaveBooks(dataFile, books)
}

func main() {
	// Cargar datos al iniciar
	if err := loadData(); err != nil {
		fmt.Println("Error al cargar datos:", err)
		return
	}

	// Ejecutar servidor web en segundo plano
	go runWebServer()

	// CLI para gestión interactiva
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n> Opciones: register | search | list | exit")
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		cmd := strings.TrimSpace(input)

		switch cmd {
		case "register":
			var b models.Book
			fmt.Print("Título: ")
			b.Title, _ = reader.ReadString('\n')
			fmt.Print("Autor: ")
			b.Author, _ = reader.ReadString('\n')
			fmt.Print("Género: ")
			b.Genre, _ = reader.ReadString('\n')
			fmt.Print("Año: ")
			fmt.Fscan(reader, &b.Year)
			reader.ReadString('\n') // limpiar buffer

			// Evitar shadowing
			var err error
			books, err = service.RegisterBook(books, b)
			if err != nil {
				fmt.Println("❌", err)
			} else {
				saveData()
				fmt.Println("✅ Libro registrado.")
			}

		case "search":
			fmt.Print("Campo (title/author/genre): ")
			field, _ := reader.ReadString('\n')
			field = strings.TrimSpace(field)
			fmt.Print("Término: ")
			term, _ := reader.ReadString('\n')
			term = strings.TrimSpace(term)
			var resultados []models.Book
			switch field {
			case "title":
				resultados = service.SearchByTitle(books, term)
			case "author":
				resultados = service.SearchByAuthor(books, term)
			case "genre":
				resultados = service.SearchByGenre(books, term)
			default:
				fmt.Println("Campo inválido.")
				continue
			}
			service.ListBooks(resultados)

		case "list":
			service.ListBooks(books)

		case "exit":
			fmt.Println("Saliendo…")
			return

		default:
			fmt.Println("Comando no reconocido.")
		}
	}
}

// runWebServer configura y arranca el servidor HTTP con servicios REST.
func runWebServer() {
	router := gin.Default()

	// 1. Salud del servicio
	router.GET("/health", healthHandler)
	// 2. Crear libro
	router.POST("/books", createBookHandler)
	// 3. Listar todos los libros
	router.GET("/books", getBooksHandler)
	// 4. Obtener libro por ID
	router.GET("/books/:id", getBookByIDHandler)
	// 5. Actualizar libro
	router.PUT("/books/:id", updateBookHandler)
	// 6. Eliminar libro
	router.DELETE("/books/:id", deleteBookHandler)
	// 7. Buscar por título (query param)
	router.GET("/books/search", searchBooksHandler)
	// 8. Filtrar por género
	router.GET("/books/genre/:genre", genreBooksHandler)
	// 9. Filtrar por autor
	router.GET("/books/author/:author", authorBooksHandler)
	// 10. Listar géneros disponibles
	router.GET("/genres", genresHandler)

	// Escucha en el puerto 8080
	router.Run(":8080")
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func createBookHandler(c *gin.Context) {
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	var err error
	books, err = service.RegisterBook(books, newBook)
	mu.Unlock()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	saveData()
	c.JSON(http.StatusCreated, newBook)
}

func getBooksHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, books)
}

func getBookByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for _, b := range books {
		if b.ID == id {
			c.JSON(http.StatusOK, b)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Libro no encontrado"})
}

func updateBookHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var upd models.Book
	if err := c.ShouldBindJSON(&upd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for i, b := range books {
		if b.ID == id {
			if upd.Title != "" {
				books[i].Title = upd.Title
			}
			if upd.Author != "" {
				books[i].Author = upd.Author
			}
			if upd.Genre != "" {
				books[i].Genre = upd.Genre
			}
			if upd.Year != 0 {
				books[i].Year = upd.Year
			}
			saveData()
			c.JSON(http.StatusOK, books[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Libro no encontrado"})
}

func deleteBookHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			saveData()
			c.JSON(http.StatusOK, gin.H{"deleted": true})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Libro no encontrado"})
}

func searchBooksHandler(c *gin.Context) {
	title := c.Query("title")
	mu.Lock()
	defer mu.Unlock()
	results := service.SearchByTitle(books, title)
	c.JSON(http.StatusOK, results)
}

func genreBooksHandler(c *gin.Context) {
	genre := c.Param("genre")
	mu.Lock()
	defer mu.Unlock()
	results := service.SearchByGenre(books, genre)
	c.JSON(http.StatusOK, results)
}

func authorBooksHandler(c *gin.Context) {
	author := c.Param("author")
	mu.Lock()
	defer mu.Unlock()
	results := service.SearchByAuthor(books, author)
	c.JSON(http.StatusOK, results)
}

func genresHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	set := make(map[string]bool)
	for _, b := range books {
		set[b.Genre] = true
	}
	var list []string
	for g := range set {
		list = append(list, g)
	}
	c.JSON(http.StatusOK, list)
}
