package tasks

import (
	"fmt"
	"github.com/Bostanova/go_final_project/app/models"
	"strconv"
	"strings"
	"time"
)

//const TimeFormat = "20060102"

func getYear(now time.Time, date time.Time) (time.Time, error) {
	resDate := date.AddDate(1, 0, 0)
	for !resDate.After(now) {
		resDate = resDate.AddDate(1, 0, 0)
	}

	return resDate, nil
}

func getDay(now time.Time, date time.Time, repeatSlice []string) (time.Time, error) {
	if len(repeatSlice) < 2 {
		return time.Time{}, fmt.Errorf("не указан интервал в днях")
	}

	daysNumber, err := strconv.Atoi(repeatSlice[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("недопустимый формат интервала")
	}
	if daysNumber > 400 {
		return time.Time{}, fmt.Errorf("превышен максимально допустимый интервал") // using time.Time{} as zero date. Can check it with (t Time) IsZero() function
	}

	resDate := date.AddDate(0, 0, daysNumber)
	for !resDate.After(now) {
		resDate = resDate.AddDate(0, 0, daysNumber)
	}
	return resDate, nil
}

func NextDate(now time.Time, date string, repeat string, taskExist bool) (string, error) {
	//TODO добавить правило: Если правило не указано, отмеченная выполненной задача будет удаляться из таблицы.
	var resDate time.Time

	d, err := time.Parse(models.TimeFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты")
	}

	repeatSlice := strings.Split(repeat, " ")
	switch repeatSlice[0] {
	case "d":
		if date == now.Format(models.TimeFormat) && !taskExist {
			return date, nil
		}
		resDate, err = getDay(now, d, repeatSlice)
	case "y":
		resDate, err = getYear(now, d)
	case "":
		resDate = time.Now()
	default:
		err = fmt.Errorf("недопустимый символ")
	}

	if resDate.IsZero() {
		return "", err
	}
	return resDate.Format(models.TimeFormat), nil
}
