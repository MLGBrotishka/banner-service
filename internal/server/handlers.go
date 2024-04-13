package server

import (
	"encoding/json"
	"fmt"
	"my_app/internal/db"
	"my_app/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func ValidateInt(param string) (*int, error) {
	if param == "" {
		return nil, fmt.Errorf("value required")
	}
	value, err := strconv.Atoi(param)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func UserBannerGet(w http.ResponseWriter, r *http.Request) {
	// Получение параметров запроса
	useLastRevision := r.URL.Query().Get("use_last_revision") == "true"
	featureId, featureErr := ValidateInt(r.URL.Query().Get("feature_id"))
	tagId, tagErr := ValidateInt(r.URL.Query().Get("tag_id"))
	// Проверка наличия параметров
	var errorResponse models.ErrorResponse
	if featureErr != nil || tagErr != nil {
		errorResponse.Error = "tag_id and feature_id are required"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	isAdmin := false
	role, ok := r.Context().Value(RoleKey).(Role)
	if ok {
		isAdmin = role == AdminRole
	}
	// Получение баннера из базы данных
	bannerContent, err := db.GetBannerForUser(featureId, tagId, useLastRevision, isAdmin)
	if err != nil {
		if strings.Contains(err.Error(), "no banner found") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			// Обработка других типов ошибок
			errorResponse.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bannerContent)
}

func BannersGet(w http.ResponseWriter, r *http.Request) {
	// Получение параметров запроса
	limit, err := ValidateInt(r.URL.Query().Get("limit"))
	if err != nil && !strings.Contains(err.Error(), "value required") {
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
		return
	}
	offset, err := ValidateInt(r.URL.Query().Get("offset"))
	if err != nil && !strings.Contains(err.Error(), "value required") {
		http.Error(w, "Invalid offset value", http.StatusBadRequest)
		return
	}
	featureId, featureErr := ValidateInt(r.URL.Query().Get("feature_id"))
	tagId, tagErr := ValidateInt(r.URL.Query().Get("tag_id"))
	// Проверка наличия параметров
	var errorResponse models.ErrorResponse
	if featureErr != nil && tagErr != nil {
		errorResponse.Error = "At least one of feature_id or tag_id must be provided"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Получение баннеров из базы данных
	banners, err := db.GetBanners(featureId, tagId, limit, offset)
	if err != nil {
		errorResponse.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	if len(banners) == 0 {
		banners = []models.BannerExpanded{}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(banners)
}

func BannerPost(w http.ResponseWriter, r *http.Request) {
	// Получение параметров запроса
	var banner models.BannerNoId
	var errorResponse models.ErrorResponse
	err := json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		errorResponse.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	var response models.IdResponse
	// Создание баннера в базе данных
	response.BannerId, err = db.CreateBanner(banner)
	if err != nil {
		errorResponse.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func BannerIdPatch(w http.ResponseWriter, r *http.Request) {
	// Получение параметров запроса
	vars := mux.Vars(r)
	idStr := vars["id"]
	var errorResponse models.ErrorResponse
	if idStr == "" {
		errorResponse.Error = "Banner Id is required"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	id, err := ValidateInt(idStr)
	if err != nil {
		errorResponse.Error = "Invalid banner Id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var banner models.BannerNoId
	err = json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		errorResponse.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	exist, _ := db.BannerExists(*id)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Обновление баннера в базе данных
	err = db.UpdateBanner(*id, banner)
	if err != nil {
		errorResponse.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func BannerIdDelete(w http.ResponseWriter, r *http.Request) {
	// Получение параметров запроса
	vars := mux.Vars(r)
	idStr := vars["id"]
	var errorResponse models.ErrorResponse
	if idStr == "" {
		errorResponse.Error = "Banner Id is required"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	id, err := ValidateInt(idStr)
	if err != nil {
		errorResponse.Error = "Invalid banner Id"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	exist, _ := db.BannerExists(*id)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Удаление баннера из базы данных
	err = db.DeleteBanner(*id)
	if err != nil {
		errorResponse.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
