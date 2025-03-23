package repo

import (
	"database/sql"
	"github.com/Bostanova/go_final_project/app/models"
	_ "modernc.org/sqlite"
	"os"
)

const (
	createTableQuery = `
	CREATE TABLE IF NOT EXISTS scheduler (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	date CHAR(8) NOT NULL DEFAULT "",
    	title VARCHAR(128) NOT NULL DEFAULT "",
    	comment TEXT,
    	repeat VARCHAR(128) NOT NULL DEFAULT ""
	);
	`
	createIndexDateQuery = `CREATE INDEX IF NOT EXISTS index_dates ON scheduler (date);`
)

type DataBase struct {
	db *sql.DB
}

var DB *DataBase

func NewDB(dbFile string) (*DataBase, error) {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	dbConn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	db := &DataBase{db: dbConn}

	if install {
		err := db.createDB()
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (db *DataBase) createDB() error {
	_, err := db.db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	_, err = db.db.Exec(createIndexDateQuery)
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBase) Close() error {
	return db.db.Close()
}

// AddTask добавляет задачу в БД и возвращает id добавленной задачи либо ошибку
func (db *DataBase) AddTask(task models.Task) (int, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"
	res, err := db.db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	insertID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(insertID), nil
}

// UpdateTask обновляет задачу
func (db *DataBase) UpdateTask(task models.Task) error {
	_, err := db.db.Exec("UPDATE scheduler SET date= :date, title= :title, comment= :comment, repeat= :repeat WHERE id= :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	return nil
}

// DeleteTask удаляет задачу
func (db *DataBase) DeleteTask(id string) error {
	_, err := db.db.Exec("DELETE FROM scheduler WHERE id= :id", sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}

// GetTask принимает id задачи в качестве аргумента и возвращает информацию об этой задаче
func (db *DataBase) GetTask(id string) (models.Task, bool, error) {
	var task models.Task
	row := db.db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, false, err
	}

	return task, true, nil
}

// GetTasks возвращает упорядоченный по дате список всех задач
func (db *DataBase) GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	rows, err := db.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 20")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
