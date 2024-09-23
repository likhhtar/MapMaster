package main

import (
	"bytes"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func TestAPI(t *testing.T) {
	mux := http.NewServeMux()

	// Создаем нано-сервисы Storage и Router
	storage := NewStorage(mux, "test", []string{}, true)
	go func() { storage.Run() }()
	t.Cleanup(storage.Stop)

	router := NewRouter(mux, [][]string{{"test"}})
	go func() { router.Run() }()
	t.Cleanup(router.Stop)

	// Ожидаем, что сервисы запустились
	time.Sleep(100 * time.Millisecond)

	// Создаем тестовые данные
	feature := geojson.NewFeature(orb.Point{rand.Float64(), rand.Float64()})
	body, err := feature.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	// Табличный шаблон для разных запросов
	tests := []struct {
		method string
		url    string
		body   []byte
		want   int
	}{
		{
			method: "POST",
			url:    "/test/insert",
			body:   body,
			want:   http.StatusOK,
		},
		{
			method: "POST",
			url:    "/test/replace",
			body:   body,
			want:   http.StatusOK,
		},
		{
			method: "POST",
			url:    "/test/delete",
			body:   body,
			want:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.url, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, bytes.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			// Проверка на редирект
			if rr.Code == http.StatusTemporaryRedirect {
				req, err := http.NewRequest(tt.method, rr.Header().Get("Location"), bytes.NewReader(tt.body))
				if err != nil {
					t.Fatal(err)
				}
				rr = httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
			}

			// Проверяем финальный статус-код
			if rr.Code != tt.want {
				t.Errorf("handler returned wrong status code: got %v, want %v", rr.Code, tt.want)
			}
		})
	}
}
