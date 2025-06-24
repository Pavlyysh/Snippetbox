package main

import (
	"errors"
	"fmt"
	"net/http"
	"pavlyysh/snippetbox/pkg/models"
	"strconv"
)

// Меняем сигнатуры обработчика home, чтобы он определялся как метод
// структуры *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) // Использование помощника notFound()
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	// // Инициализируем срез содержащий пути к двум файлам. Обратите внимание, что
	// // файл home.page.tmpl должен быть *первым* файлом в срезе.
	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }

	// // Используем функцию template.ParseFiles() для чтения файла шаблона.
	// // Если возникла ошибка, мы запишем детальное сообщение ошибки и
	// // используя функцию http.Error() мы отправим пользователю
	// // ответ: 500 Internal Server Error (Внутренняя ошибка на сервере)
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// Поскольку обработчик home теперь является методом структуры application
	// 	// он может получить доступ к логгерам из структуры.
	// 	// Используем их вместо стандартного логгера от Go.
	// 	app.serverError(w, err) // Использование помощника serverError()
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// // Затем мы используем метод Execute() для записи содержимого
	// // шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	// // возможность отправки динамических данных в шаблон.
	// err = ts.Execute(w, nil)
	// if err != nil {
	// 	// Обновляем код для использования логгера-ошибок
	// 	// из структуры application.
	// 	app.serverError(w, err) // Использование помощника serverError()
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }

}

// Меняем сигнатуру обработчика showSnippet, чтобы он был определен как метод
// структуры *application
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id")) // извлекаем значение id из URL и преобразуем в int
	if err != nil || id < 1 {
		app.notFound(w) // Использование помощника notFound()
		return
	}

	// Вызываем метода Get из модели Snipping для извлечения данных для
	// конкретной записи на основе её ID. Если подходящей записи не найдено,
	// то возвращается ответ 404 Not Found (Страница не найдена).
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%v", s)
}

// Меняем сигнатуру обработчика createSnippet, чтобы он определялся как метод
// структуры *application.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Используем r.Method для проверки, использует ли запрос метод POST или нет. Обратите внимание,
	// что http.MethodPost является строкой и содержит текст "POST".
	if r.Method != http.MethodPost {
		// Используем метод Header().Set() для добавления заголовка 'Allow: POST' в
		// карту HTTP-заголовков. Первый параметр - название заголовка, а
		// второй параметр - значение заголовка.
		w.Header().Set("Allow", http.MethodPost)

		// Используем функцию http.Error() для отправки кода состояния 405 с соответствующим сообщением.
		app.clientError(w, http.StatusMethodNotAllowed) // Используем помощник clientError()
		return
	}

	title := "A story about shark"
	content := "Shark says baby,\nbaby says shark,\nshark again said baby"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
