package handlers

import (
	"github.com/Bostanova/go_final_project/app/repo"
	"github.com/Bostanova/go_final_project/app/tasks"
	"log"
	"net/http"
	"time"
)

func taskDone(w http.ResponseWriter, r *http.Request) {
	errResp := make(map[string]string)

	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		errResp["error"] = "Не указан идентификатор"
		createRespons(errResp, http.StatusBadRequest, w)
		return
	}

	task, idExist, err := repo.DB.GetTask(id)
	if err != nil {
		if !idExist {
			errResp["error"] = "Задача не найдена"
			createRespons(errResp, http.StatusBadRequest, w)
			return
		}
		createRespons("", http.StatusInternalServerError, w)
		return
	}

	if task.Repeat == "" {
		err = repo.DB.DeleteTask(task.ID)
		if err != nil {
			errResp["error"] = err.Error()
			createRespons(errResp, http.StatusInternalServerError, w)
			return
		}
	}
	task.Date, err = tasks.NextDate(time.Now(), task.Date, task.Repeat, true)
	if err != nil {
		log.Println(err)
		createResponse("error", error.Error(err), http.StatusInternalServerError, w)
		return
	}
	err = repo.DB.UpdateTask(task)
	if err != nil {
		log.Println(err)
		createResponse("error", error.Error(err), http.StatusInternalServerError, w)
		return
	}
	createRespons(errResp, http.StatusOK, w)
	return
}

func TaskDelete(w http.ResponseWriter, r *http.Request) {
	errResp := make(map[string]string)

	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		errResp["error"] = "Не указан идентификатор"
		createRespons(errResp, http.StatusBadRequest, w)
		return
	}

	_, idExist, err := repo.DB.GetTask(id)
	if err != nil {
		if !idExist {
			errResp["error"] = "Задача не найдена"
			createRespons(errResp, http.StatusBadRequest, w)
			return
		}
		createRespons("", http.StatusInternalServerError, w)
		return
	}

	err = repo.DB.DeleteTask(id)
	if err != nil {
		errResp["error"] = err.Error()
		createRespons(errResp, http.StatusInternalServerError, w)
		return
	}
	createRespons(map[string]string{}, http.StatusOK, w)
	return
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		taskDone(w, r)
	//case "DELETE":
	//	taskDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
