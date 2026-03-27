package handler

import (
	"encoding/json"
	"go-postgres-crud-api/internal/models"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
	result := h.DB.Create(&task)
	if result.Error != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task

	result := h.DB.Find(&tasks)

	if result.Error != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)

}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	result := h.DB.First(&task, id)
	if result.Error != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if result := h.DB.First(&task, id); result.Error != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	var updating_task models.Task
	if err := json.NewDecoder(r.Body).Decode(&updating_task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updating_task.Title != "" {
		task.Title = updating_task.Title
	}
	if updating_task.Description != "" {
		task.Description = updating_task.Description
	}

	h.DB.Save(&task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	result := h.DB.Delete(&models.Task{}, id)
	if result.Error != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
