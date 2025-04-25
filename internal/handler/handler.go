// internal/handler/handler.go
package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/mail" // Для валидации email
	"regexp"  // Для валидации телефона
	"runtime/debug"
	"strings" // Для TrimSpace
	"time"   // Для валидации даты
	"unicode/utf8" // Для длины пароля

	"shanraq.xyz/internal/config"
	"shanraq.xyz/internal/view"
)
type Handler struct {
	Logger        *slog.Logger
	TemplateCache view.TemplateCache
	Config        config.Config
}

func NewHandler(l *slog.Logger, tc view.TemplateCache, cfg config.Config) *Handler {
	return &Handler{
		Logger:        l,
		TemplateCache: tc,
		Config:        cfg,
	}
}

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

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" { h.notFound(w, r); return }
	http.Redirect(w, r, "/welcome", http.StatusSeeOther)
}

func (h *Handler) Welcome(w http.ResponseWriter, r *http.Request) {
	err := view.Render(w, r, http.StatusOK, "welcome.html", h.TemplateCache, h.Config, nil)
	if err != nil { h.serverError(w, r, err) }
}

// func (h *Handler) LoginGet(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "/welcome", http.StatusSeeOther) // Или просто "/"
// }

func (h *Handler) LoginGet(w http.ResponseWriter, r *http.Request) {
 	data := &view.TemplateData{ Form: view.RegisterForm{Errors: map[string]string{}} } // Передаем пустую форму
 	err := view.Render(w, r, http.StatusOK, "login.html", h.TemplateCache, h.Config, data)
 	if err != nil { h.serverError(w, r, err) }
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	err := view.Render(w, r, http.StatusOK, "dashboard.html", h.TemplateCache, h.Config, nil)
	if err != nil { h.serverError(w, r, err) }
}

// RegisterPost обрабатывает отправку формы регистрации.
func (h *Handler) RegisterPost(w http.ResponseWriter, r *http.Request) {
	// Парсим форму
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, r, http.StatusBadRequest)
		return
	}

	// Извлекаем данные и убираем лишние пробелы
	form := view.RegisterForm{
		Gender:      strings.TrimSpace(r.PostForm.Get("gender")),
		DOB:         strings.TrimSpace(r.PostForm.Get("dob")),
		FirstName:   strings.TrimSpace(r.PostForm.Get("first_name")),
		LastName:    strings.TrimSpace(r.PostForm.Get("last_name")),
		MiddleName:  strings.TrimSpace(r.PostForm.Get("middle_name")), // Необязательное
		Email:       strings.TrimSpace(r.PostForm.Get("email")),
		PhoneNumber: strings.TrimSpace(r.PostForm.Get("phone_number")),
		Password:    r.PostForm.Get("password"), // Пароль не тримим
		Errors:      map[string]string{},
	}

	if form.Gender == "" {
		form.Errors["Gender"] = "Поле 'Пол' обязательно для заполнения"
	} else if form.Gender != "male" && form.Gender != "female" {
		form.Errors["Gender"] = "Недопустимое значение для поля 'Пол'"
	}

	if form.DOB == "" {
		form.Errors["DOB"] = "Поле 'Дата рождения' обязательно для заполнения"
	} else {
		_, err := time.Parse("2006-01-02", form.DOB)
		if err != nil {
			form.Errors["DOB"] = "Неверный формат даты (ожидается ГГГГ-ММ-ДД)"
		}
		// Можно добавить проверку на возраст, если нужно
	}

	if form.FirstName == "" {
		form.Errors["FirstName"] = "Поле 'Имя' обязательно для заполнения"
	}

	if form.LastName == "" {
		form.Errors["LastName"] = "Поле 'Фамилия' обязательно для заполнения"
	}

	if form.Email == "" {
		form.Errors["Email"] = "Поле 'Email' обязательно для заполнения"
	} else {
		_, err := mail.ParseAddress(form.Email) // Используем стандартный пакет для валидации email
		if err != nil {
			form.Errors["Email"] = "Неверный формат Email адреса"
		}
	}

	// Пример простой валидации для формата +7XXXXXXXXXX (можно усложнить)
	phoneRegex := `^\+7\d{10}$`
	if form.PhoneNumber == "" {
		form.Errors["PhoneNumber"] = "Поле 'Телефон' обязательно для заполнения"
	} else if matched, _ := regexp.MatchString(phoneRegex, form.PhoneNumber); !matched {
		form.Errors["PhoneNumber"] = "Неверный формат телефона (ожидается +7XXXXXXXXXX)"
	}

	if form.Password == "" {
		form.Errors["Password"] = "Поле 'Пароль' обязательно для заполнения"
	} else if utf8.RuneCountInString(form.Password) < 8 {
		form.Errors["Password"] = "Пароль должен содержать не менее 8 символов"
	}

	// Если есть ошибки, рендерим форму снова с ошибками
	if len(form.Errors) > 0 {
		h.Logger.Info("Ошибки валидации при регистрации", "errors", form.Errors)
		data := &view.TemplateData{Form: form} // Передаем форму с ошибками
		// Используем статус 422 Unprocessable Entity
		err = view.Render(w, r, http.StatusUnprocessableEntity, "login.html", h.TemplateCache, h.Config, data)
		if err != nil {
			h.serverError(w, r, err)
		}
		return
	}

	// На данном этапе просто логируем
	h.Logger.Info("Успешная валидация регистрации",
		"gender", form.Gender,
		"dob", form.DOB,
		"firstName", form.FirstName,
		"lastName", form.LastName,
		"email", form.Email,
		"phone", form.PhoneNumber,
	)

	// TODO: Хешировать пароль перед сохранением!
	// TODO: Сохранить пользователя в БД.
	// TODO: Добавить Flash сообщение об успехе.
	// TODO: Решить, куда перенаправлять пользователя (на /login или сразу в /dashboard?).
	//       Пока просто вернем ответ OK или перенаправим на страницу входа.

	// Устанавливаем Flash сообщение (пока просто пример, нужна реализация сессий)
	// app.sessionManager.Put(r.Context(), "flash", "Регистрация прошла успешно! Теперь вы можете войти.")

	// Перенаправляем на страницу входа (или главную)
	http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther) // Добавим параметр для возможного сообщения
	// Или можно просто: fmt.Fprintln(w, "Регистрация успешна (данные в логе)")

}

// TODO: Добавить LoginPost, LogoutPost