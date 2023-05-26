package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Cheatsheet struct {
	info string
}

type cheatsheetSlice []*Cheatsheet

func InsertNewCheatsheetInfo(db *sql.DB, info string) {
	log.Println("Info text: ", info)
	unfiltered_info := strings.Replace(info, "/neueInfo ", "", -1)
	infos := strings.Split(unfiltered_info, ",")
	info_text := infos[0]
	sqlStatement := `INSERT INTO cheatsheet (info) VALUES (?)`
	_, err := db.Exec(sqlStatement, info_text)
	if err != nil {
		panic(err)
	}
}

func (text cheatsheetSlice) CheatsheetString() string {
	var s []string
	for _, u := range text {
		if u != nil {
			s = append(s, fmt.Sprintf("%s", u.info))
		}
	}
	return strings.Join(s, "\n")
}

func GetAllCheatsheetInfos(db *sql.DB) string {
	sqlStatement := `SELECT * FROM cheatsheet`
	cheatsheets := make([]*Cheatsheet, 0)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println("Error getting all cheatsheets: ", err)
	}

	for rows.Next() {
		cheatsheet := new(Cheatsheet)
		_ = rows.Scan(&cheatsheet.info)
		cheatsheets = append(cheatsheets, cheatsheet)

	}

	fmt.Println(cheatsheets)
	return (cheatsheetSlice(cheatsheets).CheatsheetString())
}
