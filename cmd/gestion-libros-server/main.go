// cmd/gestion-libros-server/main.go
package main

import (
	"net/http"
	"strconv"
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

func loadBooks() error {
	b, err := storage.LoadBooks(dataFile)
	if err != nil {
		return err
	}
	books = b
	return nil
}

func saveBooks() error {
	return storage.SaveBooks(dataFile, books)
}

func main() {
	// Carga inicial de datos
	if err := loadBooks(); err != nil {
		panic("Error al cargar libros: " + err.Error())
	}

	r := gin.Default()

	// Endpoints REST (al menos 8 servicios)
	r.GET("/health", healthHandler)                    // 1. Salud del servicio
	r.POST("/books", createBookHandler)                // 2. Crear libro
	r.GET("/books", getBooksHandler)                   // 3. Listar todos los libros
	r.GET("/books/:id", getBookByIDHandler)            // 4. Obtener libro por ID
	r.PUT("/books/:id", updateBookHandler)             // 5. Actualizar libro
	r.DELETE("/books/:id", deleteBookHandler)          // 6. Eliminar libro
	r.GET("/books/search", searchBooksHandler)         // 7. Buscar por título (query param)
	r.GET("/books/genre/:genre", genreBooksHandler)    // 8. Filtrar por género
	r.GET("/books/author/:author", authorBooksHandler) // 9. Filtrar por autor
	r.GET("/genres", genresHandler)                    // 10. Listar géneros disponibles

	r.Run(":8080") // Inicia servidor en el puerto 8080
}
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
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

func createBookHandler(c *gin.Context) {
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var err error
	mu.Lock()
	books, err = service.RegisterBook(books, newBook)
	mu.Unlock()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := saveBooks(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newBook)
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
			if err := saveBooks(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
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
			if err := saveBooks(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
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
	genresMap := make(map[string]bool)
	for _, b := range books {
		genresMap[b.Genre] = true
	}
	var genres []string
	for g := range genresMap {
		genres = append(genres, g)
	}
	c.JSON(http.StatusOK, genres)
}
