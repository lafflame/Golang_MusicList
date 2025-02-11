package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Выводим все треки в начале программы и потом предлагаем выбор
func main() {
	allTracks()
	vibor()
}

// Выбор пользователем действий
func vibor() {
	//Главное меню
	fmt.Println("\nВыберите действие:\n1.Добавить трек\n2.Выдать случайный трек\n3.Удалить трек\n4.Вывести все треки\n5.Выход\n")
	var choice int
	fmt.Scan(&choice)

	//Очистка буфера
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	//Пользователь выбирает действие
	switch choice {
	case 1:
		zapis() //Добавить трек
	case 2:
		random() //Выдать случайный трек
	case 3:
		fmt.Println("В разработке...") //Удалить трек
		time.Sleep(2 * time.Second)
		vibor()
	case 4:
		allTracks() //Вывести все треки
	case 5:
		fmt.Println("Хорошего дня!") //Выход
		return
	case 6:
		parsing()
	default:
		fmt.Println("Некорректный выбор")
		vibor()
	}
	prodolzhenie() //Если пользователь хочет выбрать ещё что-то
}

// Вопрос пользователю, чтобы не запускать программу заново после каждого действия
func prodolzhenie() {
	fmt.Println("Вы хотите продолжить? (y/n)")
	var choice string
	fmt.Scan(&choice)
	if choice == "y" {
		vibor()
	} else {
		fmt.Println("Хорошего дня!")
		time.Sleep(2 * time.Second)
		return
	}
}

// Вывод всех треков на экран
func allTracks() {
	//Стандартная работа с файлом и отложенное закрытие
	f, err := os.Open("Tracks.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Println("Список всех аудиозаписей в медиатеке: ")
	time.Sleep(2 * time.Second) //Имитация подключения к файлу и раздумья...
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break //Если не получается считать новую строку - значит все они закончились и цикл прекращается
		}
		fmt.Print(line)
	}
}

// Запись пользовательского трека в файл по стандарту "№ : "NameOfArtist - NameOfSong"
func zapis() {
	//Стандартная работа с файлом и отложенное закрытие
	filePath := "Tracks.txt"
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	//Читаем имя исполнителя
	fmt.Println("Введите имя исполнителя или название группы: ")
	nameOfArtist, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	nameOfArtist = strings.TrimSpace(nameOfArtist)

	//Читаем название трек
	fmt.Println("Введите название трека: ")
	trackName, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	trackName = strings.TrimSpace(trackName)

	//Добавляем номер трека в начало, тире между треком и названием группы
	lastTrackNumber, err := lastTrackNumber()
	data := strconv.Itoa(lastTrackNumber+1) + ": " + nameOfArtist + " - " + trackName
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Добавляем текст в конец файла
	_, err = f.WriteString(data + "\n")
	if err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		fmt.Println("Трек был успешно добавлен!\n")
	}
}

// Вывод случайного трека
func random() {
	//Стандартная работа с файлом и отложенное закрытие
	filePath := "Tracks.txt"
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var tracks []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		tracks = append(tracks, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}
	if len(tracks) == 0 {
		fmt.Println("Медиатека пуста. Добавьте треки перед использованием функции случайного выбора.")
		return
	}
	randomIndex := rand.Intn(len(tracks))
	fmt.Println("\nСлучайный трек из медиатеки:")
	fmt.Println(tracks[randomIndex])
}

// Определение последнего числа (кол-во строк) для работы с добавлением треков и определения кол-ва треков
func lastTrackNumber() (int, error) {
	//Стандартная работа с файлом и отложенное закрытие
	f, err := os.Open("Tracks.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	maxNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2) // Разделяем строку на номер и название
		if len(parts) < 2 {
			continue // Если длина строки меньше 2 - пропускаем
		}

		// Пробуем преобразовать номер в число
		number, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err == nil && number > maxNumber {
			maxNumber = number
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err // Возвращаем ошибку, если она произошла при чтении
	}

	return maxNumber, nil
}

func gettingInfo(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Ошибка при выполнении запроса:", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Ошибка: статус код %d", resp.StatusCode)
	}

	// Парсим HTML-документ
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Ошибка при парсинге HTML:", err)
	}

	// Извлекаем имя исполнителя
	artistName := doc.Find("h1.page-artist__title").First().Text()

	// Выводим результат
	fmt.Printf("Имя исполнителя: %s\n", artistName)
}

// Получение данных с библиотеки пользователя
func parsing() {
	var url string
	fmt.Println("Введите URL с аудиозаписью (!URL должен открываться без авторизации!): ")
	fmt.Scan(&url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("При получении данных произошла ошибка."+
			"\nУбедитесь, что медиатека общедоступна", err)
	}
	defer resp.Body.Close()

}

//package main
//
//import (
//"fmt"
//"log"
//"net/http"
//
//"github.com/PuerkitoBio/goquery"
//)
//
//func main() {
//	// URL страницы с исполнителем на Яндекс.Музыке
//	url := "https://music.yandex.ru/artist/3121" // Замените на нужный URL
//
//	// Выполняем HTTP GET запрос
//	resp, err := http.Get(url)
//	if err != nil {
//		log.Fatal("Ошибка при выполнении запроса:", err)
//	}
//	defer resp.Body.Close()
//
//	// Проверяем статус ответа
//	if resp.StatusCode != http.StatusOK {
//		log.Fatalf("Ошибка: статус код %d", resp.StatusCode)
//	}
//
//	// Парсим HTML-документ
//	doc, err := goquery.NewDocumentFromReader(resp.Body)
//	if err != nil {
//		log.Fatal("Ошибка при парсинге HTML:", err)
//	}
//
//	// Извлекаем имя исполнителя
//	artistName := doc.Find("h1.page-artist__title").First().Text()
//
//	// Выводим результат
//	fmt.Printf("Имя исполнителя: %s\n", artistName)
//}
