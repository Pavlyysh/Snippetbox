package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Создаем новый флаг командной строки, значение по умолчанию: ":4000".
	// Добавляем небольшую справку, объясняющая, что содержит данный флаг.
	// Значение флага будет сохранено в переменной addr.
	addr := flag.String("addr", ":4000", "server address")

	// Мы вызываем функцию flag.Parse() для извлечения флага из командной строки.
	// Она считывает значение флага из командной строки и присваивает его содержимое
	// переменной. Вам нужно вызвать ее *до* использования переменной addr
	// иначе она всегда будет содержать значение по умолчанию ":4000".
	// Если есть ошибки во время извлечения данных - приложение будет остановлено.
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Инициализируем FileServer, он будет обрабатывать
	// HTTP-запросы к статическим файлам из папки "./ui/static".
	// Обратите внимание, что переданный в функцию http.Dir путь
	// является относительным корневой папке проекта
	fileServer := http.FileServer(neutredFileSystem{http.Dir("./static")})

	// Используем функцию mux.Handle() для регистрации обработчика для
	// всех запросов, которые начинаются с "/static/". Мы убираем
	// префикс "/static" перед тем как запрос достигнет http.FileServer
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Значение, возвращаемое функцией flag.String(), является указателем на значение
	// из флага, а не самим значением. Нам нужно убрать ссылку на указатель
	// то есть перед использованием добавьте к нему префикс *. Обратите внимание, что мы используем
	// функцию log.Printf() для записи логов в журнал работы нашего приложения.
	log.Printf("start web-server on %s\n", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

type neutredFileSystem struct {
	fs http.FileSystem
}

func (nfs neutredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
