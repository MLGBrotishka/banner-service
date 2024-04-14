package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"my_app/internal/cache"
	"my_app/internal/db"
	"my_app/internal/models"
	"my_app/internal/server"
)

type userBannerGetRequest struct {
	FeatureId string
	TagId     string
}

type userBannerGetOutput struct {
	Code int
	Body map[string]interface{}
}

type userBannerGetTestsuite struct {
	Name    string
	Request userBannerGetRequest
	Output  userBannerGetOutput
}

func TestUserBannerGet(t *testing.T) {
	banner := models.BannerNoId{
		TagIds:    []int32{0, 1},
		FeatureId: 0,
		Content: map[string]interface{}{
			"sucsess": "true",
		},
		IsActive: true,
	}
	testsuite := []userBannerGetTestsuite{
		{
			Name: "Not found",
			Request: userBannerGetRequest{
				FeatureId: "1",
				TagId:     "2",
			},
			Output: userBannerGetOutput{
				Code: http.StatusNotFound,
			},
		},
		{
			Name: "Bad request",
			Request: userBannerGetRequest{
				FeatureId: "1s",
				TagId:     "2",
			},
			Output: userBannerGetOutput{
				Code: http.StatusBadRequest,
				Body: map[string]interface{}{
					"error": "tag_id and feature_id are required",
				},
			},
		},
		{
			Name: "Not enought var",
			Request: userBannerGetRequest{
				FeatureId: "",
				TagId:     "2",
			},
			Output: userBannerGetOutput{
				Code: http.StatusBadRequest,
				Body: map[string]interface{}{
					"error": "tag_id and feature_id are required",
				},
			},
		},
		{
			Name: "Correct",
			Request: userBannerGetRequest{
				FeatureId: "0",
				TagId:     "0",
			},
			Output: userBannerGetOutput{
				Code: http.StatusOK,
				Body: map[string]interface{}{
					"sucsess": "true",
				},
			},
		},
	}

	// Инициализация тестовой базы данных и кэша
	db.InitDB()
	t.Log("Сonnected to db")
	defer db.CloseDB()
	cache.InitCache()
	t.Log("Сonnected to cache")
	defer cache.CloseCache()
	db.CreateBanner(banner)

	for _, curTest := range testsuite {
		curUrl := "/user_banner?feature_id=" + curTest.Request.FeatureId + "&tag_id=" + curTest.Request.TagId
		req := httptest.NewRequest(http.MethodGet, curUrl, nil)
		w := httptest.NewRecorder()
		server.UserBannerGet(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != curTest.Output.Code {
			t.Fatalf("Expected status %v; got %v", curTest.Output.Code, res.StatusCode)
		}
		var body map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&body)
		switch {
		case err == io.EOF && curTest.Output.Body == nil:
			t.Log(curTest.Name, ": Pass")
			continue
		case err != nil:
			t.Fatalf("Failed to decode response body: %v", err)
		}

		// Сравнение тела ответа с ожидаемым значением
		if !reflect.DeepEqual(body, curTest.Output.Body) {
			t.Fatalf("Expected body %v; got %v", curTest.Output.Body, body)
		}
		t.Log(curTest.Name, ": Pass")
	}
}
