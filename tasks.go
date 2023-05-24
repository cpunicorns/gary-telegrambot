package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Tasks struct {
	task    string
	weekday string
}

type taskSlice []*Tasks

func (text taskSlice) TaskString() string {
	var s []string
	for _, u := range text {
		if u != nil {
			s = append(s, fmt.Sprintf("%s %s", u.task, u.weekday))
		}
	}
	return strings.Join(s, "\n")
}

func GetAllTasks(db *sql.DB) string {
	sqlStatement := `SELECT * FROM tasks`
	tasks := make([]*Tasks, 0)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println("Error getting all tasks: ", err)
	}

	for rows.Next() {
		task := new(Tasks)
		_ = rows.Scan(&task.task, &task.weekday)
		tasks = append(tasks, task)

	}
	fmt.Println(tasks)
	return (taskSlice(tasks).TaskString())
}

func InsertNewTask(db *sql.DB, message string) {
	log.Println("Message text: ", message)
	unfiltered_message := strings.Replace(message, "/neueAufgabe ", "", -1)
	task := strings.Split(unfiltered_message, ",")
	task_summary := task[0]
	task_weekday := task[1]
	sqlStatement := `INSERT INTO tasks (task, weekday) VALUES (?, ?)`
	_, err := db.Exec(sqlStatement, task_summary, task_weekday)
	if err != nil {
		panic(err)
	}
}
