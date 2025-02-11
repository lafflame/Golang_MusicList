package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Выводим все треки в начале программы и потом предлагаем выбор
func main() {
	allTracks()
	vibor()
}

// Выбор пользователем действий
func vibor() {
	//Главное меню
	fmt.Println("\nВыберите действие:\n1.Добавить трек\n2.Выдать случайный трек" +
		"\n3.Удалить трек\n4.Вывести все треки\n5.Поиск трека\n6.Редактирование трека\n" +
		"7.Поиск клипа на Youtube\n8.Статистика в медиатеке\n9.Выход")
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
		deleteTrack() // Удаление трека
	case 4:
		allTracks() //Вывести все треки
	case 5:
		searchTrack() // Поиск трека
	case 6:
		editTrack() // Редактирование трека
	case 7:
		playYouTubeClip() // Открытие клипа
	case 8:
		showStatistics() // Статистика в медиатеке
	case 9:
		gettingInfo() // Секретная фича - парсинг (который блокируется сайтом)
	case 10:
		fmt.Println("Хорошего дня!")
		return
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

// Парсинг
func gettingInfo() {
	url := "https://music.yandex.ru/artist/1426524"
	// Выполняем HTTP-запрос
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
	// Извлекаем название трека
	trackName := doc.Find("h1.page-artist__title").Text()
	fmt.Println(trackName)

}

// Открытие браузера и ссылки
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}

// Поиск трека
func searchTrack() {
	fmt.Println("Введите название трека или имя исполнителя для поиска:")
	reader := bufio.NewReader(os.Stdin)
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	// Открываем файл с треками
	filePath := "Tracks.txt"
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer f.Close()

	// Сканируем файл построчно
	scanner := bufio.NewScanner(f)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
			fmt.Println(line)
			found = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	if !found {
		fmt.Println("Трек не найден.")
	}
}

func deleteTrack() {
	fmt.Println("Введите номер трека для удаления:")
	var trackNumber int
	fmt.Scan(&trackNumber)

	// Открываем файл с треками
	filePath := "Tracks.txt"
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer f.Close()

	// Читаем все строки из файла
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	// Ищем трек с указанным номером
	found := false
	var updatedLines []string
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		// Пробуем преобразовать номер трека в число
		number, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}

		// Если номер трека совпадает, пропускаем его (удаляем)
		if number == trackNumber {
			found = true
			continue
		}

		// Добавляем строку в обновлённый список
		updatedLines = append(updatedLines, line)
	}

	if !found {
		fmt.Println("Трек с таким номером не найден.")
		return
	}

	// Добавляем пустую строку в конец
	updatedLines = append(updatedLines, "")

	// Перезаписываем файл с обновлённым списком треков
	err = os.WriteFile(filePath, []byte(strings.Join(updatedLines, "\n")), 0644)
	if err != nil {
		fmt.Println("Ошибка при записи файла:", err)
		return
	}

	fmt.Println("Трек успешно удалён.")
}

// Изменение названия трека
func editTrack() {
	fmt.Println("Введите номер трека для редактирования:")
	var trackNumber int
	fmt.Scan(&trackNumber)

	// Очистка буфера после `fmt.Scan`
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	// Открываем файл с треками
	filePath := "Tracks.txt"
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer f.Close()

	// Читаем все строки из файла
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	// Ищем трек с указанным номером
	found := false
	var updatedLines []string
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		// Пробуем преобразовать номер трека в число
		number, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}

		// Если номер трека совпадает, редактируем его
		if number == trackNumber {
			found = true
			fmt.Println("Текущий трек:", line)
			fmt.Println("Введите новое имя исполнителя:")
			nameOfArtist, _ := reader.ReadString('\n')
			nameOfArtist = strings.TrimSpace(nameOfArtist)

			fmt.Println("Введите новое название трека:")
			trackName, _ := reader.ReadString('\n')
			trackName = strings.TrimSpace(trackName)

			// Формируем новую строку для трека
			newLine := strconv.Itoa(number) + ": " + nameOfArtist + " - " + trackName
			updatedLines = append(updatedLines, newLine)
		} else {
			// Добавляем строку в обновлённый список без изменений
			updatedLines = append(updatedLines, line)
		}
	}

	if !found {
		fmt.Println("Трек с таким номером не найден.")
		return
	}

	// Перезаписываем файл с обновлённым списком треков
	err = os.WriteFile(filePath, []byte(strings.Join(updatedLines, "\n")), 0644)
	if err != nil {
		fmt.Println("Ошибка при записи файла:", err)
		return
	}

	fmt.Println("Трек успешно отредактирован.")
}

// Поиск трека на ютубе
func playYouTubeClip() {
	fmt.Println("Введите название трека или исполнителя для поиска на YouTube:")
	reader := bufio.NewReader(os.Stdin)
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	if query == "" {
		fmt.Println("Запрос не может быть пустым!")
		return
	}

	// Формируем URL для поиска на YouTube
	searchURL := "https://www.youtube.com/results?search_query=" + strings.ReplaceAll(query, " ", "+")

	fmt.Println("Открываю YouTube...")
	err := openBrowser(searchURL)
	if err != nil {
		fmt.Println("Ошибка при открытии браузера:", err)
	}
}

// Функция для вывода статистики по трекам
func showStatistics() {
	// Открываем файл с треками
	filePath := "Tracks.txt"
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer f.Close()

	// Сканируем файл построчно
	scanner := bufio.NewScanner(f)
	trackCount := 0
	artistCounts := make(map[string]int) // Мапа для подсчета треков по исполнителям

	for scanner.Scan() {
		line := scanner.Text()
		trackCount++

		// Разделяем строку на части: номер, исполнитель и название трека
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		// Извлекаем имя исполнителя
		trackInfo := strings.SplitN(parts[1], "-", 2)
		if len(trackInfo) < 2 {
			continue
		}
		artist := strings.TrimSpace(trackInfo[0])

		// Увеличиваем счетчик для текущего исполнителя
		artistCounts[artist]++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	// Выводим общую статистику
	fmt.Println("\nСтатистика медиатеки:")
	fmt.Printf("Общее количество треков: %d\n", trackCount)

	// Выводим количество треков по каждому исполнителю
	fmt.Println("\nКоличество треков по исполнителям:")
	for artist, count := range artistCounts {
		fmt.Printf("%s: %d треков\n", artist, count)
	}

	// Находим самого популярного исполнителя
	maxCount := 0
	popularArtist := ""
	for artist, count := range artistCounts {
		if count > maxCount {
			maxCount = count
			popularArtist = artist
		}
	}

	if popularArtist != "" {
		fmt.Printf("\nСамый популярный исполнитель: %s (%d треков)\n", popularArtist, maxCount)
	} else {
		fmt.Println("\nНет данных о самом популярном исполнителе.")
	}
}
