// internal/handler/handler.go
package handler

import (
	"database/sql"
	"fmt"
	"html/template" 
	"log/slog"
	"net/http"
	"net/mail"
	"regexp"
	"runtime/debug"
	"strings"
	"time"
	"unicode/utf8"

	"shanraq.xyz/internal/config"
	"shanraq.xyz/internal/view"
)

type Handler struct {
	Logger        *slog.Logger
	TemplateCache map[string]*template.Template // Теперь тип известен
	Config        config.Config
	DB            *sql.DB
}

func NewHandler(l *slog.Logger, tc view.TemplateCache, cfg config.Config, db *sql.DB) *Handler {
	return &Handler{
		Logger:        l,
		TemplateCache: tc,
		Config:        cfg,
		DB:            db,
	}
}

// --- Хелперы для ошибок ---
func (h *Handler) serverError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.Logger.Error(trace, "method", r.Method, "uri", r.URL.RequestURI())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *Handler) clientError(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	h.clientError(w, r, http.StatusNotFound)
}

// --- Основные обработчики ---
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.notFound(w, r)
		return
	}
	http.Redirect(w, r, "/welcome", http.StatusSeeOther)
}

func (h *Handler) Welcome(w http.ResponseWriter, r *http.Request) {
	err := view.Render(w, r, http.StatusOK, "welcome.html", h.TemplateCache, h.Config, nil)
	if err != nil {
		h.serverError(w, r, err)
	}
}

// LoginGet рендерит шаблон "login.html"
func (h *Handler) LoginGet(w http.ResponseWriter, r *http.Request) {
	data := &view.TemplateData{Form: view.RegisterForm{Errors: map[string]string{}}}
	err := view.Render(w, r, http.StatusOK, "login.html", h.TemplateCache, h.Config, data)
	if err != nil {
		h.serverError(w, r, err)
	}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// TODO: Добавить проверку аутентификации
	err := view.Render(w, r, http.StatusOK, "dashboard.html", h.TemplateCache, h.Config, nil)
	if err != nil {
		h.serverError(w, r, err)
	}
}

// RegisterPost обрабатывает отправку формы регистрации.
func (h *Handler) RegisterPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, r, http.StatusBadRequest)
		return
	}
	form := view.RegisterForm{ /* ... извлечение данных ... */
		Gender:      strings.TrimSpace(r.PostForm.Get("gender")),
		DOB:         strings.TrimSpace(r.PostForm.Get("dob")),
		FirstName:   strings.TrimSpace(r.PostForm.Get("first_name")),
		LastName:    strings.TrimSpace(r.PostForm.Get("last_name")),
		MiddleName:  strings.TrimSpace(r.PostForm.Get("middle_name")),
		Email:       strings.TrimSpace(r.PostForm.Get("email")),
		PhoneNumber: strings.TrimSpace(r.PostForm.Get("phone_number")),
		Password:    r.PostForm.Get("password"),
		Errors:      map[string]string{},
	}
	// --- Валидация ---
	if form.Gender == "" { form.Errors["Gender"] = "Поле 'Пол' обязательно для заполнения" } else if form.Gender != "male" && form.Gender != "female" { form.Errors["Gender"] = "Недопустимое значение для поля 'Пол'" }
	if form.DOB == "" { form.Errors["DOB"] = "Поле 'Дата рождения' обязательно для заполнения" } else { _, err := time.Parse("2006-01-02", form.DOB); if err != nil { form.Errors["DOB"] = "Неверный формат даты (ожидается ГГГГ-ММ-ДД)" } }
	if form.FirstName == "" { form.Errors["FirstName"] = "Поле 'Имя' обязательно для заполнения" }
	if form.LastName == "" { form.Errors["LastName"] = "Поле 'Фамилия' обязательно для заполнения" }
	if form.Email == "" { form.Errors["Email"] = "Поле 'Email' обязательно для заполнения" } else { _, err := mail.ParseAddress(form.Email); if err != nil { form.Errors["Email"] = "Неверный формат Email адреса" } }
	phoneRegex := `^\+7\d{10}$`
	if form.PhoneNumber == "" { form.Errors["PhoneNumber"] = "Поле 'Телефон' обязательно для заполнения" } else if matched, _ := regexp.MatchString(phoneRegex, form.PhoneNumber); !matched { form.Errors["PhoneNumber"] = "Неверный формат телефона (ожидается +7XXXXXXXXXX)" }
	if form.Password == "" { form.Errors["Password"] = "Поле 'Пароль' обязательно для заполнения" } else if utf8.RuneCountInString(form.Password) < 8 { form.Errors["Password"] = "Пароль должен содержать не менее 8 символов" }
	// --- Конец Валидации ---
	if len(form.Errors) > 0 {
		h.Logger.Info("Ошибки валидации при регистрации", "errors", form.Errors)
		data := &view.TemplateData{Form: form}
		err = view.Render(w, r, http.StatusUnprocessableEntity, "login.html", h.TemplateCache, h.Config, data)
		if err != nil { h.serverError(w, r, err) }
		return
	}
	h.Logger.Info("Успешная валидация регистрации", "email", form.Email)
	// TODO: Хешировать пароль, Сохранить пользователя в БД, Добавить Flash
	http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
}

// LoginPost обрабатывает отправку формы входа.
func (h *Handler) LoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, r, http.StatusBadRequest)
		return
	}
	form := view.LoginForm{
		Email:    strings.TrimSpace(r.PostForm.Get("email")),
		Password: r.PostForm.Get("password"),
		Errors:   map[string]string{},
	}
	if form.Email == "" { form.Errors["Email"] = "Email не может быть пустым" }
	if form.Password == "" { form.Errors["Password"] = "Пароль не может быть пустым" }

	if len(form.Errors) > 0 {
		h.Logger.Info("Ошибки валидации при входе", "email", form.Email, "errors", form.Errors)
		data := &view.TemplateData{Form: form}
		data.Flash = "Ошибка входа: проверьте email и пароль."
		err = view.Render(w, r, http.StatusUnprocessableEntity, "login.html", h.TemplateCache, h.Config, data)
		if err != nil { h.serverError(w, r, err) }
		return
	}

	h.Logger.Info("Попытка входа", "email", form.Email)
	// --- Заглушка для проверки пользователя (TODO) ---
	// ... (комментарии TODO как были) ...
	// --- Конец Заглушки ---
	h.Logger.Info("Успешный вход (заглушка)", "email", form.Email)
	// TODO: Создать сессию пользователя.
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// TODO: Добавить LogoutPost