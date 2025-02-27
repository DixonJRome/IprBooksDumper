package engine

import (
	"fmt"
	"os"
)

// SaveToFile создание pdf файла
func SaveToFile(bookID string, data []byte) {
	filename := bookID + ".pdf"

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}

	_, err = f.Write(data)
	if err != nil {
		fmt.Println("Ошибка при записи файла:", err)
	}

	f.Close()
	fmt.Println("Файл сохранён как", filename)
}
