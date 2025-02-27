package main

import (
	"IprbooksDumper/engine"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Ввод логина
	fmt.Print("Введите ваш логин -> ")
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login) // Убираем лишние пробелы и переносы строки

	// Ввод пароля
	fmt.Print("Введите ваш пароль -> ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Авторизация
	engine.Auth(login, password)

	// Ввод ID книг
	fmt.Print("Введите ID вашей книги, если книг несколько, введите их ID через запятую -> ")
	text, _ := reader.ReadString('\n')

	num := strings.Replace(text, "\n", "", -1)
	idList := strings.Split(num, ",")

	var idListRes []int

	// Проверка введенных ID
	for _, val := range idList {
		convertId, err := strconv.Atoi(strings.TrimSpace(val))

		if err != nil {
			log.Println("Не валидный ID: ", val)
			continue
		}

		idListRes = append(idListRes, convertId)
	}

	resInfoList := engine.DumpBookData(idListRes, login, password) // Передаем логин и пароль

	if len(resInfoList) == 0 {
		panic("Все ID не валидные.")
	}

	// Сохранение книг
	for _, dumpBook := range resInfoList {
		bookID := dumpBook.Name // Теперь Name — это ID
		engine.SaveToFile(bookID, dumpBook.BookBytes)
		fmt.Println("Файл записан.")
	}
}
