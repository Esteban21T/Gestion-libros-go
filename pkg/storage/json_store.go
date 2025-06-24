package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Esteban21T/Sistemas-de-Gesti-n-empresarial.git/pkg/models"
)

// LoadBooks lee libros desde el archivo JSON indicado.
func LoadBooks(path string) ([]models.Book, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Si no existe, devolvemos slice vacío
		return []models.Book{}, nil
	}
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	var books []models.Book
	if err := json.Unmarshal(data, &books); err != nil {
		return nil, err
	}
	return books, nil
}

// SaveBooks escribe el slice de libros en el archivo JSON.
func SaveBooks(path string, books []models.Book) error {
	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}
