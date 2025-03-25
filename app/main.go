package main

import (
	"github.com/Bostanova/go_final_project/app/handlers"
	"github.com/Bostanova/go_final_project/app/repo"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	// загрузка данных из .env в переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// получение переменных окружения
	port := os.Getenv("TODO_PORT")
	port = ":" + port
	dbFile := os.Getenv("TODO_DBFILE")

	// создание БД
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile = filepath.Join(filepath.Dir(appPath), dbFile)

	repo.DB, err = repo.NewDB(dbFile)
	if err != nil {
		log.Fatalf("DB creating/opening failed: %v", err)
	}
	defer repo.DB.Close()

	// запуск сервера
	r := chi.NewRouter()

	// обработчики
	r.Handle("/*", http.FileServer(http.Dir("./web")))

	r.Get("/api/nextdate", handlers.NextDateHandler)
	r.Get("/api/tasks", handlers.GetTasksHandler)
	r.Post("/api/task/done", handlers.TaskDoneHandler)
	r.HandleFunc("/api/task", handlers.TaskHandler)

	log.Print("Starting the server...")
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("server starting failed: %v", err)
	}
}
