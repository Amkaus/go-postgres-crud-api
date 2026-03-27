package main

import (
	"go-postgres-crud-api/internal/config"
	"go-postgres-crud-api/internal/handler"
	"go-postgres-crud-api/internal/models"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	dsn := cfg.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.Task{}); err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	handler := &handler.Handler{DB: db}

	http.HandleFunc("POST /tasks", handler.CreateTask)
	http.HandleFunc("GET /tasks", handler.ListTasks)
	http.HandleFunc("GET /tasks/{id}", handler.GetTask)
	http.HandleFunc("PUT /tasks/{id}", handler.UpdateTask)
	http.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)

	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
