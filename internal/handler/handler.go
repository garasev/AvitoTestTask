package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/garasev/AvitoTestTask/internal/service"
	"github.com/go-chi/chi"
)

type Handler struct {
	Service service.Service
	Logger  slog.Logger
}

func NewHandler(service service.Service, logger slog.Logger) *Handler {
	return &Handler{
		Service: service,
		Logger:  logger}
}

func (h *Handler) GetSlug(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		h.Logger.Error(err.Error())
	}
	slug, err := h.Service.Repo.GetSlug(id)
	if err != nil {
		h.Logger.Error(err.Error())
	}

	jsonOrder, err := json.Marshal(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Установка заголовка Content-Type на application/json
	w.Header().Set("Content-Type", "application/json")

	// Отправка JSON в ответ
	w.Write(jsonOrder)
}
