package router

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
	if featureErr != nil || tagErr != nil {
		http.Error(w, "tag_id and feature_id are required", http.StatusBadRequest)
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
			http.Error(w, "No banner found", http.StatusNotFound)
		} else {
			// Обработка других типов ошибок
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
	if featureErr != nil && tagErr != nil {
		http.Error(w, "At least one of feature_id or tag_id must be provided", http.StatusBadRequest)
		return
	}

	// Получение баннеров из базы данных
	banners, err := db.GetBanners(featureId, tagId, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(banners) == 0 {
		http.Error(w, "No banner found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(banners)
}

func BannerPost(w http.ResponseWriter, r *http.Request) {
	// Получение параметров запроса
	var banner models.BannerNoId
	err := json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var response models.InlineResponse201
	// Создание баннера в базе данных
	response.BannerId, err = db.CreateBanner(banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	if idStr == "" {
		http.Error(w, "Banner Id is required", http.StatusBadRequest)
		return
	}

	id, err := ValidateInt(idStr)
	if err != nil {
		http.Error(w, "Invalid banner Id", http.StatusBadRequest)
		return
	}

	var banner models.BannerNoId
	err = json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exist, _ := db.BannerExists(*id)
	if !exist {
		http.Error(w, "Banner Id is not found", http.StatusNotFound)
		return
	}

	// Обновление баннера в базе данных
	err = db.UpdateBanner(*id, banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func BannerIdDelete(w http.ResponseWriter, r *http.Request) {
	// Получение параметров запроса
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		http.Error(w, "Banner Id is required", http.StatusBadRequest)
		return
	}

	id, err := ValidateInt(idStr)
	if err != nil {
		http.Error(w, "Invalid banner Id", http.StatusBadRequest)
		return
	}

	exist, _ := db.BannerExists(*id)
	if !exist {
		http.Error(w, "Banner Id is not found", http.StatusNotFound)
		return
	}

	// Удаление баннера из базы данных
	err = db.DeleteBanner(*id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
