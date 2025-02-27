package engine

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type DumpData struct {
	Name      string
	BookBytes []byte
}

// DumpBookData получает декодированный ряд байтов книги и ее название
func DumpBookData(bookIdList []int, username, password string) (resArray []DumpData) {

	for _, bookId := range bookIdList {
		resValue, err := dumpData(bookId, username, password)

		if err != nil {
			log.Fatal(err)
			continue
		}

		resArray = append(resArray, resValue)
	}

	return resArray
}

func dumpData(bookId int, username, password string) (DumpData, error) {
	// создаем авторизованного клиента
	client := Auth(username, password) // Передаем логин и пароль

	// ссылка на зашифрованный контент книги
	link := "https://www.iprbookshop.ru/pdfstream.php?publicationId=" + strconv.Itoa(bookId) + "&part=null"

	requestModel, err := http.NewRequest("GET", link, nil)

	if err != nil {
		return DumpData{}, errors.New("Сайт недоступен!")
	}

	// делаем запрос на сайт
	response, err := client.Do(requestModel)

	if err != nil {
		return DumpData{}, errors.New("Сайт недоступен!")
	}

	defer response.Body.Close()

	bodyText, err := io.ReadAll(response.Body)

	if err != nil {
		return DumpData{}, errors.New("Сайт недоступен!")
	}

	if len(bodyText) == 25462 {
		return DumpData{}, errors.New("Книги не существует!")
	}

	return DumpData{Name: strconv.Itoa(bookId), BookBytes: DecodeBytes(bodyText)}, nil
}

// Min поиск минимального значения в массиве
func Min(arr []int) int {
	min := arr[0]

	for _, val := range arr {
		if val < min {
			min = val
		}
	}
	return min
}

// Auth возвращает авторизованный клиент
func Auth(username, password string) http.Client {
	authUrl := "https://www.iprbookshop.ru/95835"

	// данные для авторизации
	data := url.Values{}
	data.Set("action", "login")
	data.Set("username", username)
	data.Set("password", password)
	data.Set("rememberme", "1")

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, authUrl, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// запрос на авторизацию
	authReq, _ := client.Do(r)

	// создаем контейнер для куки авторизации
	cookieJar, _ := cookiejar.New(nil)
	url, _ := url.Parse(authUrl)

	// устанавливаем куки авторизации
	cookieJar.SetCookies(url, authReq.Cookies())

	// создаем клиент с авторизацией
	client = &http.Client{Jar: cookieJar}

	return *client
}

// DecodeBytes декодирует набор байтов
func DecodeBytes(b []byte) []byte {

	for i := 0; i < len(b); i += 2048 {
		for j := i; j < Min([]int{i + 100, len(b) - 1}); j += 2 {
			b[j], b[j+1] = b[j+1], b[j]
		}
	}
	return b
}

// Name контейнер для имени книги
type Name struct {
	name string
}

// Получает название книги с сайта
func GetBookName(bookId int) string {
	link := "https://www.iprbookshop.ru/" + strconv.Itoa(bookId) + ".html"

	Name := Name{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.iprbookshop.ru", "iprbookshop.ru"),
	)

	c.OnHTML("h4.header-orange", func(e *colly.HTMLElement) {
		Name.name = e.Text
	})

	c.Visit(link)

	return Name.name
}
