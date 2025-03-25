package handlers

import (
	"github.com/Bostanova/go_final_project/app/repo"
	"github.com/Bostanova/go_final_project/app/tasks"
	"log"
	"net/http"
	"time"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]string)

	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		body["error"] = "Не указан идентификатор"
		createResponse(body, http.StatusBadRequest, w)
		return
	}

	task, idExist, err := repo.DB.GetTask(id)
	if err != nil {
		if !idExist {
			body["error"] = "Задача не найдена"
			createResponse(body, http.StatusBadRequest, w)
			return
		}
		createResponse("", http.StatusInternalServerError, w)
		return
	}

	if task.Repeat == "" {
		err = repo.DB.DeleteTask(task.ID)
		if err != nil {
			body["error"] = err.Error()
			createResponse(body, http.StatusInternalServerError, w)
			return
		}
	}
	task.Date, err = tasks.NextDate(time.Now(), task.Date, task.Repeat, true)
	if err != nil {
		log.Println(err)
		body["error"] = err.Error()
		createResponse(body, http.StatusInternalServerError, w)
		return
	}

	err = repo.DB.UpdateTask(task)
	if err != nil {
		log.Println(err)
		body["error"] = err.Error()
		createResponse(body, http.StatusInternalServerError, w)
		return
	}

	createResponse(body, http.StatusOK, w)
	return
}
