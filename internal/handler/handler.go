// internal/handler/handler.go
package handler

import (
	"fmt" 
	"log/slog"
	"net/http"
	"runtime/debug" 

	"shanraq.xyz/internal/config"
	"shanraq.xyz/internal/view"
)

// Handler содержит зависимости для обработчиков HTTP.
type Handler struct {
	Logger        *slog.Logger
	TemplateCache view.TemplateCache
	Config        config.Config
}

// NewHandler создает новый экземпляр Handler.
func NewHandler(l *slog.Logger, tc view.TemplateCache, cfg config.Config) *Handler {
	return &Handler{
		Logger:        l,
		TemplateCache: tc,
		Config:        cfg,
	}
}

// serverError логирует детальную ошибку со стектрейсом и отправляет пользователю 500.
func (h *Handler) serverError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.Logger.Error(trace, "method", r.Method, "uri", r.URL.RequestURI()) // Логируем стектрейс
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError отправляет пользователю специфичный статус-код и сообщение.
// Теперь принимает 'r *http.Request' для потенциального логирования.
func (h *Handler) clientError(w http.ResponseWriter, r *http.Request, status int) {
	// Можно добавить логирование клиентских ошибок при необходимости
	// h.Logger.Warn("Client error", "status", status, "method", r.Method, "uri", r.URL.RequestURI())
	http.Error(w, http.StatusText(status), status)
}

// notFound - обработчик для 404 Not Found.
// Теперь принимает 'r *http.Request', чтобы соответствовать http.HandlerFunc.
func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	h.clientError(w, r, http.StatusNotFound) // Передаем 'r' в clientError
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.notFound(w, r) // Передаем 'r'
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

func (h *Handler) LoginGet(w http.ResponseWriter, r *http.Request) {
	data := &view.TemplateData{
		Form: struct{}{},
	}
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

// TODO: Добавить LoginPost, RegisterPost, LogoutPost
