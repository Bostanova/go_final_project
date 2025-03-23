package handlers

import (
	"fmt"
	"github.com/Bostanova/go_final_project/app/models"
	"github.com/Bostanova/go_final_project/app/tasks"
	"log"
	"net/http"
	"time"
)

//const TimeFormat = "20060102"

func checkParam(r *http.Request, param string) (string, error) {
	resParam := r.URL.Query().Get(param)
	if len(resParam) == 0 {
		return resParam, fmt.Errorf("параметр %s не указан", param)
	}
	if v := r.Header.Get(param); len(v) > 0 {
		return fmt.Sprintf("%s", v), nil
	}
	return resParam, nil
}
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now, err := checkParam(r, "now")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	date, err := checkParam(r, "date")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	repeat, err := checkParam(r, "repeat")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	nowFormat, err := time.Parse(models.TimeFormat, now)
	if err != nil {
		log.Println("неверный формат даты")
		w.Write([]byte(err.Error()))
	}

	nextDate, err := tasks.NextDate(nowFormat, date, repeat, false)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(nextDate))
}
