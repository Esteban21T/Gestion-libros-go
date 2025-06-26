// cmd/gestion-libros-server/main.go
package main

import (
	"html/template"
	"net/http"
	"path/filepath"
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

// loadData carga los libros desde el archivo JSON.
func loadData() error {
	b, err := storage.LoadBooks(dataFile)
	if err != nil {
		return err
	}
	books = b
	return nil
}

// saveData guarda el slice de libros en el archivo JSON.
func saveData() error {
	return storage.SaveBooks(dataFile, books)
}

func main() {
	// 1. Cargar datos
	if err := loadData(); err != nil {
		panic("Error al cargar datos: " + err.Error())
	}

	// 2. Configurar router
	router := gin.Default()

	// 3. Cargar plantillas HTML desde web/templates
	router.SetHTMLTemplate(template.Must(
		template.ParseGlob(filepath.Join("web", "templates", "*.html")),
	))
	// 4. Servir archivos estáticos
	router.Static("/static", filepath.Join("web", "static"))
	// 5. Ruta principal que muestra la UI
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 6. Endpoints REST
	router.GET("/health", healthHandler)
	router.POST("/books", createBookHandler)
	router.GET("/books", getBooksHandler)
	router.GET("/books/:id", getBookByIDHandler)
	router.PUT("/books/:id", updateBookHandler)
	router.DELETE("/books/:id", deleteBookHandler)
	router.GET("/books/search", searchBooksHandler)
	router.GET("/books/genre/:genre", genreBooksHandler)
	router.GET("/books/author/:author", authorBooksHandler)
	router.GET("/genres", genresHandler)

	// 7. Iniciar servidor
	if err := router.Run(":8080"); err != nil {
		panic("Error iniciando servidor: " + err.Error())
	}
}

// healthHandler comprueba el estado de salud del microservicio.
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// createBookHandler recibe un JSON y registra un nuevo libro.
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
	if err := saveData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newBook)
}

// getBooksHandler lista todos los libros.
func getBooksHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, books)
}

// getBookByIDHandler obtiene un libro por su ID.
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

// updateBookHandler actualiza los campos de un libro existente.
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
			if err := saveData(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, books[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Libro no encontrado"})
}

// deleteBookHandler elimina un libro por ID.
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
			if err := saveData(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"deleted": true})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Libro no encontrado"})
}

// searchBooksHandler busca libros por título.
func searchBooksHandler(c *gin.Context) {
	title := c.Query("title")
	mu.Lock()
	defer mu.Unlock()
	results := service.SearchByTitle(books, title)
	c.JSON(http.StatusOK, results)
}

// genreBooksHandler filtra libros por género.
func genreBooksHandler(c *gin.Context) {
	genre := c.Param("genre")
	mu.Lock()
	defer mu.Unlock()
	results := service.SearchByGenre(books, genre)
	c.JSON(http.StatusOK, results)
}

// authorBooksHandler filtra libros por autor.
func authorBooksHandler(c *gin.Context) {
	author := c.Param("author")
	mu.Lock()
	defer mu.Unlock()
	results := service.SearchByAuthor(books, author)
	c.JSON(http.StatusOK, results)
}

// genresHandler lista todos los géneros disponibles.
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
