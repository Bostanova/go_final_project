package handlers

import (
	"encoding/json"
	"github.com/Bostanova/go_final_project/app/repo"
	"log"
	"net/http"
)

func createBody(body map[string]interface{}, w http.ResponseWriter) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bodyJson)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]interface{})

	tasks, err := repo.DB.GetTasks()
	if err != nil {
		body["error"] = err.Error()
		createBody(body, w)
		return
	}

	if tasks == nil {
		body["tasks"] = []string{}
		createBody(body, w)
		return
	}

	body["tasks"] = tasks
	createBody(body, w)
}
