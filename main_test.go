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

func TestCafeSearc(t *testing.T) {

	handler := http.HandlerFunc(mainHandle)
	city := "moscow"

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, r := range requests {
		url := fmt.Sprintf("/cafe?city=%s&search=%s", city, r.search)
		req := httptest.NewRequest("GET", url, nil)
		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp)

		body := resp.Body.String()

		var cafe []string
		if strings.TrimSpace(body) != "" {
			cafe = strings.Split(body, ",")
		}

		assert.Equal(t, r.wantCount, len(cafe))

		for _, name := range cafe {
			nameLower := strings.ToLower(name)
			searchLower := strings.ToLower(r.search)

			assert.Contains(t, nameLower, searchLower)
		}

	}
}
