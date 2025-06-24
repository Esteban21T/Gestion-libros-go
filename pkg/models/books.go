package models

// Book representa la estructura de un libro electrónico
type Book struct {
	ID     int    `json:"id"`     // Identificador único
	Title  string `json:"title"`  // Título del libro
	Author string `json:"author"` // Autor o autores
	Genre  string `json:"genre"`  // Género o categoría
	Year   int    `json:"year"`   // Año de publicación
}
