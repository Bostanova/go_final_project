package handlers

import (
	"github.com/Bostanova/go_final_project/app/repo"
	"net/http"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]interface{})

	tasks, err := repo.DB.GetTasks()
	if err != nil {
		body["error"] = err.Error()
		createResponse(body, http.StatusInternalServerError, w)
		return
	}

	if tasks == nil {
		body["tasks"] = []string{}
		createResponse(body, http.StatusOK, w)
		return
	}

	body["tasks"] = tasks
	createResponse(body, http.StatusOK, w)
}
