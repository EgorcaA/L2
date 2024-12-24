package main

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"html"
)

func downloadFile(url string, filename string) error {
	// Открытие URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ошибка при скачивании %s: %v", url, err)
	}
	defer resp.Body.Close()

	// Открытие файла для записи
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("не удалось создать файл %s: %v", filename, err)
	}
	defer file.Close()

	// Копирование содержимого страницы в файл
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл %s: %v", filename, err)
	}

	return nil
}

func getLinks(url string) ([]string, error) {
	// Отправляем GET запрос
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при скачивании %s: %v", url, err)
	}
	defer resp.Body.Close()

	// Парсим HTML-страницу
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге HTML: %v", err)
	}

	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link := a.Val
					if strings.HasPrefix(link, "http") || strings.HasPrefix(link, "https") {
						links = append(links, link)
					} else if strings.HasPrefix(link, "/") {
						links = append(links, url+link)
					}
				}
			}
		}
		// Рекурсивно обрабатываем все узлы
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links, nil
}

func saveURLContent(url, baseDir string) error {
	// Скачиваем страницу
	fileName := filepath.Join(baseDir, "index.html")
	err := downloadFile(url, fileName)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении страницы %s: %v", url, err)
	}

	// Получаем ссылки на другие ресурсы на странице
	links, err := getLinks(url)
	if err != nil {
		return fmt.Errorf("ошибка при получении ссылок из %s: %v", url, err)
	}

	// Скачиваем каждый ресурс
	for _, link := range links {
		// Преобразуем URL в локальный путь
		parsedURL, err := url.Parse(link)
		if err != nil {
			return fmt.Errorf("ошибка при разборе URL %s: %v", link, err)
		}

		// Генерация пути для сохранения
		fileName := filepath.Join(baseDir, parsedURL.Path)
		if strings.HasSuffix(fileName, "/") {
			fileName += "index.html"
		}

		// Создание директории, если она не существует
		dir := filepath.Dir(fileName)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("ошибка при создании каталога %s: %v", dir, err)
		}

		// Скачиваем файл
		err = downloadFile(link, fileName)
		if err != nil {
			fmt.Printf("Ошибка при скачивании %s: %v\n", link, err)
		}
	}

	return nil
}

func main() {
	// Получаем URL от пользователя
	fmt.Print("Введите URL: ")
	var url string
	fmt.Scanln(&url)

	// Устанавливаем базовую директорию для сохранения
	baseDir := "downloaded_site"
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Ошибка при создании директории: %v\n", err)
		return
	}

	// Скачиваем сайт
	start := time.Now()
	err = saveURLContent(url, baseDir)
	if err != nil {
		fmt.Printf("Ошибка при скачивании сайта: %v\n", err)
		return
	}

	fmt.Printf("Скачивание завершено за %v\n", time.Since(start))
}
