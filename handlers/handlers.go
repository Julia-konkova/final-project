package handlers

import (
	"encoding/json"
	"final-project/pkg/api"
	"final-project/pkg/db"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type Handlers struct {
	WebDir string
}

func New(webDir string) *Handlers {
	return &Handlers{WebDir: webDir}
}

type Response struct {
	ID  string `json:"id,omitempty"`
	Err string `json:"error,omitempty"`
}

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("error encoding JSON:", err)
	}
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, errorMsg string) {
	response := Response{Err: errorMsg}
	writeJSONResponse(w, statusCode, response)
}

func (h *Handlers) ServeWebFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" || path == "/index.html" {
		path = "index.html"
	}

	fullPath := filepath.Join(h.WebDir, path)
	http.ServeFile(w, r, fullPath)
}

func (h *Handlers) TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		log.Printf("task output error: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка вывода списка задач")
		return
	}

	writeJSONResponse(w, http.StatusOK, TasksResp{Tasks: tasks})
}

func (h *Handlers) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("id is required")
		writeErrorResponse(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("task not found: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Задача не найдена")
		return
	}

	writeJSONResponse(w, http.StatusOK, task)
}

func (h *Handlers) DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("id is required")
		writeErrorResponse(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("task not found: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Задача не найдена")
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(task.ID)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "Не удалось удалить задачу")
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
		return
	}

	now := time.Now()
	next, err := api.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		log.Printf("error calling nextdate: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Не удалось изменить дату задачи")
		return
	}

	err = db.UpdateDate(next, id)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Не удалось сохранить дату задачи")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
}

func (h *Handlers) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("error decoding request body: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Ошибка декодирования тела запроса")
		return
	}

	if task.Title == "" {
		log.Println("empty title of task")
		writeErrorResponse(w, http.StatusBadRequest, "Отсутствует заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		log.Printf("invalid type of date: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный формат даты")
		return
	}

	err := db.UpdateTask(&task)
	if err != nil {
		log.Printf("error to updating task: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка при редактировании задачи")
		return
	}

	writeJSONResponse(w, http.StatusOK, db.Task{})
}

func (h *Handlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("id is required")
		writeErrorResponse(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("task not found: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Задача не найдена")
		return
	}

	err = db.DeleteTask(task.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Не удалось удалить задачу")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{})
}

func (h *Handlers) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("error decoding request body: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Ошибка декодирования тела запроса")
		return
	}

	if task.Title == "" {
		log.Println("empty title of task")
		writeErrorResponse(w, http.StatusBadRequest, "Отсутствует заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		log.Printf("invalid type of date: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный формат даты")
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("error to creating task: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка при создании задачи")
		return
	}

	response := Response{ID: strconv.FormatInt(id, 10)}
	writeJSONResponse(w, http.StatusOK, response)
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(api.Layout)
		return nil
	}

	t, err := time.Parse(api.Layout, task.Date)
	if err != nil {
		return err
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	taskDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if task.Repeat != "" {
		if !taskDay.Before(today) {
			return nil
		}

		next, err := api.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		if !api.AfterNow(t, now) {
			task.Date = next
		}
	} else {
		if t.Year() != now.Year() || t.Month() != now.Month() || t.Day() != now.Day() {
			task.Date = now.Format(api.Layout)
		}
	}

	return nil
}
