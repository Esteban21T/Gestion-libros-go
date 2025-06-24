package service

import (
	"fmt"

	"github.com/Esteban21T/Sistemas-de-Gesti-n-empresarial.git/pkg/models"
)

// ListBooks imprime en consola cada libro con formato legible.
func ListBooks(books []models.Book) {
	for _, b := range books {
		fmt.Printf("%d: %s - %s (%s, %d)\n",
			b.ID, b.Title, b.Author, b.Genre, b.Year)
	}
}
