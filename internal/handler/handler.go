package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/garasev/AvitoTestTask/internal/models"
	"github.com/garasev/AvitoTestTask/internal/service"
	"github.com/go-chi/chi"
)

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

type Handler struct {
	service service.Service
	logger  slog.Logger
}

func NewHandler(service service.Service, logger slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger}
}

func (h *Handler) GetSlug(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}

	slug, err := h.service.GetSlug(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonSlug, err := json.Marshal(slug)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonSlug)
}

func (h *Handler) GetSlugs(w http.ResponseWriter, r *http.Request) {
	slugs, err := h.service.GetSlugs()
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonSlugs, err := json.Marshal(slugs)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonSlugs)
}

func (h *Handler) AddSlug(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		h.logger.Error("Content Type is not application/json")
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var slug models.Slug
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&slug)
	if err != nil {
		h.logger.Error(err.Error())
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	id, err := h.service.AddSlug(slug)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Slug with the same name already exists", http.StatusConflict)
		return
	}
	h.logger.Info("Success: new slug with id=" + strconv.Itoa(id))
	errorResponse(w, "Success: new slug with id="+strconv.Itoa(id), http.StatusOK)
	return
}

func (h *Handler) DeleteSlug(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		h.logger.Error("Content Type is not application/json")
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var slug models.Slug
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&slug)
	if err != nil {
		h.logger.Error(err.Error())
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	err = h.service.DeleteSlug(slug)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Slug with the same name doesn't exists", http.StatusNoContent)
		return
	}
	h.logger.Info("Successful deletion")
	errorResponse(w, "Successful deletion", http.StatusNoContent)
	return
}
