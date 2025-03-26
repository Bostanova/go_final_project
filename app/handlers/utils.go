package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func createResponse(body any, status int, w http.ResponseWriter) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(bodyJSON)
}
