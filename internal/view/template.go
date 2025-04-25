// internal/view/template.go
package view

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"shanraq.xyz/internal/config"
)

type RegisterForm struct {
    Gender      string
    DOB         string
    FirstName   string
    LastName    string
    MiddleName  string
    Email       string
    PhoneNumber string
    Password    string
    Errors      map[string]string
}

type TemplateData struct {
	CurrentYear     int
	Config          config.Config
	CurrentPath     string
	Flash           string
	IsAuthenticated bool
	Form            any
	Data            any
}

type TemplateCache map[string]*template.Template

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache() (TemplateCache, error) {
	cache := map[string]*template.Template{}
	templateDirFS := os.DirFS("./static/tmpl")

	log.Println("Поиск шаблонов страниц в './static/tmpl/pages/*.html' на диске...")
	pages, err := fs.Glob(templateDirFS, "pages/*.html")
	if err != nil { return nil, err }
	if len(pages) == 0 {
		log.Println("ОШИБКА: Не найдено ни одного файла по шаблону 'pages/*.html' в ./static/tmpl/")
		// ... (логирование содержимого директорий) ...
		return nil, fmt.Errorf("не найдены файлы шаблонов страниц в ./static/tmpl/pages")
	} else {
		log.Printf("Найдено страниц: %d (%v)\n", len(pages), pages)
	}

	for _, page := range pages {
		name := filepath.Base(page)

		log.Printf("Парсинг шаблонов для страницы: %s (путь: %s)", name, page)
    	ts, err := template.New("base.html").Funcs(functions).ParseFS(templateDirFS, "base.html", "parts/*.html", page)
    	if err != nil { return nil, fmt.Errorf("ошибка парсинга шаблонов для %s: %w", name, err) }
    	cache[name] = ts
	}
	return cache, nil
}

func Render(w http.ResponseWriter, r *http.Request, status int, page string, cache map[string]*template.Template, cfg config.Config, data *TemplateData) error {
	ts, ok := cache[page]
	if !ok {
		return fmt.Errorf("шаблон %s не существует в кэше", page)
	}
	buf := new(bytes.Buffer)
	if data == nil {
		data = &TemplateData{}
	}
	data.CurrentYear = time.Now().Year()
	data.Config = cfg
	data.CurrentPath = r.URL.Path
	if data.Form == nil {
		data.Form = RegisterForm{Errors: map[string]string{}}
	}

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	_, err = buf.WriteTo(w)
	return err
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}