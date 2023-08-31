package handler

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	var addSlug models.AddSlug
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&addSlug)
	if err != nil {
		h.logger.Error(err.Error())
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	if addSlug.Percent < 0 || addSlug.Percent > 100 {
		h.logger.Error("Percent incorrect")
		errorResponse(w, "Percent incorrect", http.StatusBadRequest)
		return
	}
	users, err := h.service.AddSlug(
		models.Slug{Name: addSlug.Name},
		addSlug.Percent,
	)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusConflict)
		return
	}
	if users != nil {
		result := strings.Join(func() []string {
			userStrings := make([]string, len(users))
			for i, user := range users {
				userStrings[i] = strconv.Itoa(user)
			}
			return userStrings
		}(), ", ")
		h.logger.Info("Success: new slug was added for users with id: " + result)
		errorResponse(w, "Success: new slug was added for users with id: "+result, http.StatusOK)
		return
	}
	h.logger.Info("Success: new slug was added")
	errorResponse(w, "Success: new slug was added", http.StatusOK)
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
	errorResponse(w, "Successful deletion", http.StatusOK)
	return
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonUsers, err := json.Marshal(users)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUsers)
}

func (h *Handler) GetUserSlugs(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.CheckUser(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !user {
		h.logger.Error("User with id = " + idStr + " does'n exists")
		errorResponse(w, "User with id = "+idStr+" does'n exists", http.StatusBadRequest)
		return
	}

	slugs, err := h.service.GetUserSlugs(id)
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

func (h *Handler) AddUsers(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		h.logger.Error("Content Type is not application/json")
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var users models.AddUsers
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&users)
	if err != nil {
		h.logger.Error(err.Error())
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	if users.Count <= 0 {
		h.logger.Error("The number of adding users must be greater than 0.")
		errorResponse(w, "The number of adding users must be greater than 0.", http.StatusBadRequest)
		return
	}

	err = h.service.AddUsers(users.Count)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	errorResponse(w, "Users added:"+strconv.Itoa(users.Count), http.StatusOK)
	return
}

func (h *Handler) AddUserSlugs(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.CheckUser(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !user {
		h.logger.Error("User with id = " + idStr + " does'n exists")
		errorResponse(w, "User with id = "+idStr+" does'n exists", http.StatusBadRequest)
		return
	}

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		h.logger.Error("Content Type is not application/json")
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var slugs models.AddDeleteSlugs
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&slugs)
	if err != nil {
		h.logger.Error(err.Error())
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	addSlugs, err := h.service.CheckSlugs(slugs.AddSlugs)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !addSlugs {
		h.logger.Error("Bad Request. One or more slugs in add_slugs don't exists.")
		errorResponse(w, "Bad Request. One or more slugs in add_slugs don't exists.", http.StatusBadRequest)
		return
	}
	deleteSlugs, err := h.service.CheckSlugs(slugs.DeleteSlugs)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !deleteSlugs {
		h.logger.Error("Bad Request. One or more slugs in delete_slugs don't exists.")
		errorResponse(w, "Bad Request. One or more slugs in delete_slugs don't exists.", http.StatusBadRequest)
		return
	}

	err = h.service.AddSlugsUser(id, slugs.AddSlugs, time.Duration(slugs.SlugDuration)*time.Minute)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.service.DeleteSlugsUser(id, slugs.DeleteSlugs)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	errorResponse(w, "Success: All slugs have been added and removed accordingly for the user", http.StatusOK)
	return
}

func (h *Handler) GetUserArchive(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.CheckUser(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !user {
		h.logger.Error("User with id = " + idStr + " does'n exists")
		errorResponse(w, "User with id = "+idStr+" does'n exists", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if year > time.Now().Year() {
		h.logger.Error("Bad Request: year incorrect")
		errorResponse(w, "Bad Request: year incorrect", http.StatusBadRequest)
		return
	}
	month, err := strconv.Atoi(r.URL.Query().Get("month"))
	data := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	archives, err := h.service.GetUserArchive(id, data)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=file.csv")

	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)

	for _, arch := range archives {
		err = wr.Write(arch.Write())
		if err != nil {
			h.logger.Error(err.Error())
			errorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	wr.Flush()

	if err := wr.Error(); err != nil {
		h.logger.Error("Error flushing CSV writer:" + err.Error())
		errorResponse(w, "Error flushing CSV writer:"+err.Error(), http.StatusInternalServerError)
		return
	}
}
