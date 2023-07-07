package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Tricks struct {
	trick_name        string
	trick_description string
	trick_difficulty  string
}

type messageSlice []*Tricks

func InsertNewTrick(db *sql.DB, message string) {
	log.Println("Message text: ", message)
	unfiltered_message := strings.Replace(message, "/neuerTrick ", "", -1)
	trick := strings.Split(unfiltered_message, ",")
	trick_name := trick[0]
	trick_description := trick[1]
	trick_difficulty := trick[2]
	sqlStatement := `INSERT INTO tricks (name, description, difficulty) VALUES (?, ?, ?)`
	_, err := db.Exec(sqlStatement, trick_name, trick_description, trick_difficulty)
	if err != nil {
		panic(err)
	}
}

func (message messageSlice) String() string {
	var s []string
	for _, u := range message {
		if u != nil {
			s = append(s, fmt.Sprintf("Trick: %s Beschreibung: %s Schwierigkeit: %s", u.trick_name, u.trick_description, u.trick_difficulty))
		}
	}
	return strings.Join(s, "\n")
}

func GetAllTricks(db *sql.DB) string {
	sqlStatement := `SELECT * FROM tricks`
	tricks := make([]*Tricks, 0)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println("Error getting all tricks: ", err)
	}

	for rows.Next() {
		trick := new(Tricks)
		_ = rows.Scan(&trick.trick_name, &trick.trick_description, &trick.trick_difficulty)
		tricks = append(tricks, trick)

	}

	fmt.Println(tricks)
	return (messageSlice(tricks).String())
}
