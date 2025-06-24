package service

import (
	"strings"

	"github.com/Esteban21T/Sistemas-de-Gesti-n-empresarial.git/pkg/models"
)

// SearchByTitle busca libros cuyo título contenga el query (insensible a mayúsculas).
func SearchByTitle(books []models.Book, query string) []models.Book {
	var result []models.Book
	lowerQ := strings.ToLower(query)
	for _, b := range books {
		if strings.Contains(strings.ToLower(b.Title), lowerQ) {
			result = append(result, b)
		}
	}
	return result
}

// SearchByAuthor busca libros cuyo autor contenga el query (insensible a mayúsculas).
func SearchByAuthor(books []models.Book, query string) []models.Book {
	var result []models.Book
	lowerQ := strings.ToLower(query)
	for _, b := range books {
		if strings.Contains(strings.ToLower(b.Author), lowerQ) {
			result = append(result, b)
		}
	}
	return result
}

// SearchByGenre devuelve todos los libros cuyo género coincida exactamente (ignora mayúsculas/minúsculas).
func SearchByGenre(books []models.Book, genre string) []models.Book {
	var result []models.Book
	for _, b := range books {
		if strings.EqualFold(b.Genre, genre) {
			result = append(result, b)
		}
	}
	return result
}
