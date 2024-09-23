package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Router struct {
	mux   *http.ServeMux
	nodes [][]string // Список Storage узлов
}

func NewRouter(mux *http.ServeMux, nodes [][]string) *Router {
	r := &Router{mux: mux, nodes: nodes}

	// Регистрация эндпоинтов для роутера
	r.mux.Handle("/", http.FileServer(http.Dir("../front/dist")))
	r.mux.Handle("/select", http.RedirectHandler("/storage/select", http.StatusTemporaryRedirect))
	r.mux.Handle("/insert", http.RedirectHandler("/storage/insert", http.StatusTemporaryRedirect))
	r.mux.Handle("/replace", http.RedirectHandler("/storage/replace", http.StatusTemporaryRedirect))
	r.mux.Handle("/delete", http.RedirectHandler("/storage/delete", http.StatusTemporaryRedirect))

	return r
}

func (r *Router) Run() {
	slog.Info("Router is running")
}

func (r *Router) Stop() {
	slog.Info("Stopping Router")
}

type Storage struct {
	mux      *http.ServeMux
	name     string
	replicas []string
	leader   bool
}

func NewStorage(mux *http.ServeMux, name string, replicas []string, leader bool) *Storage {
	s := &Storage{mux: mux, name: name, replicas: replicas, leader: leader}

	s.mux.HandleFunc("/"+name+"/select", func(w http.ResponseWriter, r *http.Request) {
		// Вернуть geojson объекты в формате feature collection
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"type": "FeatureCollection", "features": []}`))
	})

	s.mux.HandleFunc("/"+name+"/insert", func(w http.ResponseWriter, r *http.Request) {
		// Логика для вставки нового geojson объекта
		w.Write([]byte("Inserted"))
	})

	s.mux.HandleFunc("/"+name+"/replace", func(w http.ResponseWriter, r *http.Request) {
		// Логика для замены geojson объекта
		w.Write([]byte("Replaced"))
	})

	s.mux.HandleFunc("/"+name+"/delete", func(w http.ResponseWriter, r *http.Request) {
		// Логика для удаления geojson объекта
		w.Write([]byte("Deleted"))
	})

	return s
}

func (r *Storage) Run() {
	slog.Info("Storage is running")
}

func (r *Storage) Stop() {
	slog.Info("Stopping Storage")
}

func main() {
	mux := http.NewServeMux()

	// Создаем экземпляры Router и Storage
	nodes := [][]string{{"storage-node-1"}} // Пример списка узлов
	router := NewRouter(mux, nodes)
	storage := NewStorage(mux, "storage", nil, true)

	// Запускаем сервисы в отдельных горутинах
	go router.Run()
	go storage.Run()

	// Запуск HTTP сервера
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Info("Error starting server", "err", err)
		}
	}()

	// Ожидание сигнала завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	slog.Info("Received signal", "signal", sig)

	// Остановка сервисов в обратном порядке
	server.Shutdown(context.Background())
	storage.Stop()
	router.Stop()
}
