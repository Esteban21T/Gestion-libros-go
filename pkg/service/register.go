package service

import (
	"errors"

	"github.com/Esteban21T/Sistemas-de-Gesti-n-empresarial.git/pkg/models"
)

// RegisterBook valida y añade un libro, generando un ID incremental.
func RegisterBook(books []models.Book, newBook models.Book) ([]models.Book, error) {
	if newBook.Title == "" || newBook.Author == "" {
		return books, errors.New("Título y autor son obligatorios")
	}
	maxID := 0
	for _, b := range books {
		if b.ID > maxID {
			maxID = b.ID
		}
	}
	newBook.ID = maxID + 1
	return append(books, newBook), nil
}
