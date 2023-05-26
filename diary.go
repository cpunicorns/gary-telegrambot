package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type DiaryEntry struct {
	text string
	date string
}

type diarySlice []*DiaryEntry

func (text diarySlice) DiaryString() string {
	var s []string
	for _, u := range text {
		if u != nil {
			s = append(s, fmt.Sprintf("%s %s", u.text, u.date))
		}
	}
	return strings.Join(s, "\n")
}

func GetAllDiaryEntries(db *sql.DB) string {
	sqlStatement := `SELECT * FROM diary`
	entries := make([]*DiaryEntry, 0)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println("Error getting all diary entries: ", err)
	}

	for rows.Next() {
		entry := new(DiaryEntry)
		_ = rows.Scan(&entry.text, &entry.date)
		entries = append(entries, entry)

	}

	fmt.Println(entries)
	return (diarySlice(entries).DiaryString())
}

func InsertNewDiaryEntry(db *sql.DB, message string) {
	log.Println("Message text: ", message)
	unfiltered_message := strings.Replace(message, "/neuerEintrag ", "", -1)
	entry := strings.Split(unfiltered_message, ",")
	entry_text := entry[0]
	entry_date := entry[1]
	sqlStatement := `INSERT INTO diary (text, date) VALUES (?, ?)`
	_, err := db.Exec(sqlStatement, entry_text, entry_date)
	if err != nil {
		panic(err)
	}
}
