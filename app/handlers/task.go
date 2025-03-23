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

type Response struct {
	Error string `json:"error,omitempty"`
	ID    string `json:"id,omitempty"`
}

func createResponse(key, value string, status int, w http.ResponseWriter) {
	var respBody Response
	if key == "error" {
		respBody.Error = value
	} else {
		respBody.ID = value
	}
	rJson, err := json.Marshal(respBody)
	if err != nil {
		log.Println(err)
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	w.Write(rJson)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &task)
	if err != nil {
		createResponse("error", "ошибка десериализации JSON", http.StatusBadRequest, w)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if len(task.Title) == 0 {
		createResponse("error", "не указан заголовок задачи", http.StatusBadRequest, w)
		return
	}

	if len(task.Date) == 0 {
		task.Date = time.Now().Format(models.TimeFormat)
	}

	task.Date, err = tasks.NextDate(time.Now(), task.Date, task.Repeat, false)
	if err != nil {
		log.Println(err)
		createResponse("error", error.Error(err), http.StatusInternalServerError, w)
		return
	}

	id, err := repo.DB.AddTask(task)
	if err != nil {
		log.Println(err)
		createResponse("error", error.Error(err), http.StatusInternalServerError, w)
		return
	}

	idStr := strconv.Itoa(id)
	createResponse("id", idStr, http.StatusCreated, w)
}

func createRespons(body any, status int, w http.ResponseWriter) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	w.Write(bodyJSON)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	errResp := make(map[string]string)

	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		errResp["error"] = "Задача не найдена"
		createRespons(errResp, http.StatusBadRequest, w)
		return
	}

	task, idExist, err := repo.DB.GetTask(id)
	if err != nil {
		if !idExist {
			errResp["error"] = "Не указан идентификатор"
			createRespons(errResp, http.StatusBadRequest, w)
			return
		}
		createRespons("", http.StatusInternalServerError, w)
		return
	}
	createRespons(task, http.StatusOK, w)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	//var task models.Task
	errResp := make(map[string]string)
	errResp["error"] = "Задача не найдена"

	var task models.Task
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &task)
	if err != nil {
		createResponse("error", "ошибка десериализации JSON", http.StatusBadRequest, w)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if len(task.Title) == 0 {
		createResponse("error", "не указан заголовок задачи", http.StatusBadRequest, w)
		return
	}
	if len(task.Date) == 0 {
		task.Date = time.Now().Format(models.TimeFormat)
	}
	task.Date, err = tasks.NextDate(time.Now(), task.Date, task.Repeat, false)
	if err != nil {
		log.Println(err)
		createResponse("error", error.Error(err), http.StatusInternalServerError, w)
		return
	}
	_, idExist, err := repo.DB.GetTask(task.ID)
	if err != nil {
		createRespons(errResp, http.StatusInternalServerError, w)
		return
	}
	if !idExist {
		createRespons(errResp, http.StatusBadRequest, w)
		return
	}
	err = repo.DB.UpdateTask(task)
	if err != nil {
		createRespons(errResp, http.StatusInternalServerError, w)
		return
	}
	createRespons(map[string]string{}, http.StatusOK, w)
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
