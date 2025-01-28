package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	//Выводим все треки в начале программы и потом предлагаем выбор
	allTracks()
	vibor()
}

func vibor() {
	//Главное меню
	fmt.Println("\nВыберите действие:\n1.Добавить трек\n2.Выдать случайный трек\n3.Удалить трек\n4.Вывести все треки\n5.Выход\n")
	var choice int
	fmt.Scan(&choice)

	//Очистка буфера
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	//Пользователь выбирает, что ему интересно
	switch choice {
	case 1:
		zapis()
		prodolzhenie()
	case 2:
		random()
		prodolzhenie()
	case 3:
		fmt.Println("В разработке...")
		time.Sleep(2 * time.Second)
		vibor()
	case 4:
		allTracks()
		prodolzhenie()
	case 5:
		fmt.Println("Хорошего дня!")
		time.Sleep(2 * time.Second)
		return
	default:
		fmt.Println("Некорректный выбор")
		vibor()
	}
}

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

func zapis() {
	//Стандартная работа с файлом и отложенное закрытие
	filePath := "Tracks.txt"
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//Читаем пользовательский ввод - название группы
	fmt.Println("Введите имя исполнителя или название группы: ")
	nameOfArtist, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	nameOfArtist = strings.TrimSpace(nameOfArtist)

	//Читаем пользовательский ввод - трек
	fmt.Println("Введите название трека: ")
	trackName, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	trackName = strings.TrimSpace(trackName)

	//Добавляем номер трека в начало, тире между треком и названием группы
	lastTrackNumber, err := lastTrackNumber(filePath)
	data := strconv.Itoa(lastTrackNumber+1) + ": " + nameOfArtist + " - " + trackName
	f.WriteString(data)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Добавляем текст в конец файла
	_, err = file.WriteString(data + "\n")
	if err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		fmt.Println("Трек был успешно добавлен!\n")
	}
}

func lastTrackNumber(file string) (int, error) {
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
			continue // Если строка не содержит разделителя ":", пропускаем
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
