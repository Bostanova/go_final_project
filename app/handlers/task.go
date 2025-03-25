package handlers

import (
	"encoding/json"
	"github.com/Bostanova/go_final_project/app/models"
	"github.com/Bostanova/go_final_project/app/repo"
	"github.com/Bostanova/go_final_project/app/tasks"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	body := make(map[string]interface{})

	bodyReq, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyReq, &task)
	if err != nil {
		body["error"] = "ошибка десериализации JSON"
		createResponse(body, http.StatusBadRequest, w)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if len(task.Title) == 0 {
		body["error"] = "не указан заголовок задачи"
		createResponse(body, http.StatusBadRequest, w)
		return
	}

	if len(task.Date) == 0 {
		task.Date = time.Now().Format(models.TimeFormat)
	}

	task.Date, err = tasks.NextDate(time.Now(), task.Date, task.Repeat, false)
	if err != nil {
		log.Println(err)
		body["error"] = err.Error()
		createResponse(body, http.StatusInternalServerError, w)
		return
	}

	id, err := repo.DB.AddTask(task)
	if err != nil {
		log.Println(err)
		body["error"] = err.Error()
		createResponse(body, http.StatusInternalServerError, w)
		return
	}

	idStr := strconv.Itoa(id)
	body["id"] = idStr
	createResponse(body, http.StatusCreated, w)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]interface{})

	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		body["error"] = "Задача не найдена"
		createResponse(body, http.StatusBadRequest, w)
		return
	}

	task, idExist, err := repo.DB.GetTask(id)
	if err != nil {
		if !idExist {
			body["error"] = "Не указан идентификатор"
			createResponse(body, http.StatusBadRequest, w)
			return
		}
		createResponse(body, http.StatusInternalServerError, w)
		return
	}
	createResponse(task, http.StatusOK, w)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	body := make(map[string]interface{})

	bodyReq, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyReq, &task)
	if err != nil {
		body["error"] = "ошибка десериализации JSON"
		createResponse(body, http.StatusBadRequest, w)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if len(task.Title) == 0 {
		body["error"] = "не указан заголовок задачи"
		createResponse(body, http.StatusBadRequest, w)
		return
	}
	if len(task.Date) == 0 {
		task.Date = time.Now().Format(models.TimeFormat)
	}

	task.Date, err = tasks.NextDate(time.Now(), task.Date, task.Repeat, false)
	if err != nil {
		log.Println(err)
		body["error"] = err.Error()
		createResponse(body, http.StatusInternalServerError, w)
		return
	}

	_, idExist, err := repo.DB.GetTask(task.ID)
	if err != nil {
		body["error"] = "Задача не найдена"
		createResponse(body, http.StatusInternalServerError, w)
		return
	}
	if !idExist {
		body["error"] = "Задача не найдена"
		createResponse(body, http.StatusBadRequest, w)
		return
	}
	err = repo.DB.UpdateTask(task)
	if err != nil {
		body["error"] = "Задача не найдена"
		createResponse(body, http.StatusInternalServerError, w)
		return
	}
	createResponse(body, http.StatusOK, w)
	return
}

func TaskDelete(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]string)

	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		body["error"] = "Не указан идентификатор"
		createResponse(body, http.StatusBadRequest, w)
		return
	}

	_, idExist, err := repo.DB.GetTask(id)
	if err != nil {
		if !idExist {
			body["error"] = "Задача не найдена"
			createResponse(body, http.StatusBadRequest, w)
			return
		}
		createResponse("", http.StatusInternalServerError, w)
		return
	}

	err = repo.DB.DeleteTask(id)
	if err != nil {
		body["error"] = err.Error()
		createResponse(body, http.StatusInternalServerError, w)
		return
	}
	createResponse(map[string]string{}, http.StatusOK, w)
	return
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTaskHandler(w, r)
	case "POST":
		addTaskHandler(w, r)
	case "PUT":
		putTaskHandler(w, r)
	case "DELETE":
		TaskDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
