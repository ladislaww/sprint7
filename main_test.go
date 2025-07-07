package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {

	handler := http.HandlerFunc(mainHandle)
	city := "moscow"
	total := len(cafeList[city])

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, total},
		{count: -1, want: 0},
	}

	for _, r := range requests {
		url := fmt.Sprintf("/cafe?city=%s&count=%d", city, r.count)
		request := httptest.NewRequest("GET", url, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		require.Equal(t, http.StatusOK, response.Code)

		body := response.Body.String()

		var cafes []string
		if strings.TrimSpace(body) != "" {
			cafes = strings.Split(body, ",")
		}

		assert.Equal(t, r.want, len(cafes))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "moscow"

	tests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},     // ничего не найдёт
		{"кофе", 2},       // частичное совпадение
		{"вилка", 1},      // точное совпадение
		{"КОФЕ", 2},       // проверка на регистр
		{"", len(cafeList[city])}, // без параметра поиска — вернуть всё
	}

	for _, tt := range tests {
		url := fmt.Sprintf("/cafe?city=%s&search=%s", city, tt.search)
		req := httptest.NewRequest("GET", url, nil)
		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		body := resp.Body.String()
		var cafes []string
		if strings.TrimSpace(body) != "" {
			cafes = strings.Split(body, ",")
		}

		assert.Equal(t, tt.wantCount, len(cafes))

		if tt.search != "" {
			for _, name := range cafes {
				assert.Contains(t, strings.ToLower(name), strings.ToLower(tt.search))
			}
		}
	}
}
