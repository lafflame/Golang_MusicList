package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	f, err := os.Open("Tracks.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Println("Список всех аудиозаписей в медиатеке: ")
	time.Sleep(1 * time.Second) //имитация подключения к файлу и раздумья...
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		processLine(line)
	}
	vibor()
}

func vibor() {
	fmt.Println("\nВыберите действие:\n1.Добавить трек\n2.Выдать случайный трек\n3.Удалить трек\n4.Выход")
	var choice int
	fmt.Scan(&choice)

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	switch choice {
	case 1:
		zapis()
	case 2:
		fmt.Println("В разработке...\n")
		time.Sleep(2 * time.Second)
		vibor()
	case 3:
		fmt.Println("В разработке...\n")
		time.Sleep(2 * time.Second)
		vibor()
	case 4:
		fmt.Println("Хорошего дня!")
		time.Sleep(2 * time.Second)
		return
	default:
		fmt.Println("Некорректный выбор\n")
		vibor()
	}
}

func zapis() {
	filePath := "Tracks.txt"
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Println("Введите имя исполнителя или название группы: ")
	nameOfArtist, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	nameOfArtist = strings.TrimSpace(nameOfArtist)
	fmt.Println("Введите название трека: ")
	trackName, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	trackName = strings.TrimSpace(trackName)
	lastTrackNumber, err := lastTrackNumber(filePath)
	data := strconv.Itoa(lastTrackNumber+1) + ": " + nameOfArtist + " - " + trackName
	f.WriteString(data)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close() // Закрываем файл в конце работы

	// Добавляем текст в конец файла
	_, err = file.WriteString(data + "\n") // Обязательно добавляем символ новой строки
	if err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		fmt.Println("Трек был успешно добавлен!\n")
	}
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

func processLine(line string) {
	fmt.Print(line) // Здесь можно добавить любую обработку строки
}

func lastTrackNumber(file string) (int, error) {
	f, err := os.Open("Tracks.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	maxNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		// Строка имеет вид "1: ..."
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
