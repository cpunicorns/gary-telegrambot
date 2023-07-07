package main

/* import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Reminders struct {
	date string
}

type remindersSlice []*Reminders

func (text remindersSlice) RemindersString() string {
	var s []string
	for _, u := range text {
		if u != nil {
			s = append(s, fmt.Sprintf("%s", u.date))
		}
	}
	return strings.Join(s, "\n")
}

func CreateDBEvent(db *sql.DB) {
	sqlStatement := `CREATE EVENT IF NOT EXISTS wurmkur ON SCHEDULE AT \
				(SELECT date from diary where text like '%Wurmkur%') DO INSERT INTO reminders VALUES (NOW());`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}

//func SendReminder(db *sql.DB, reminders []*Reminders) string {
//	loc, _ := time.LoadLocation("Europe/Berlin")
//	today := time.Now().In(loc)
//	lessThanTwoWeeks := today.Add(14 * time.Weekday)
//	for _, reminder := range reminders {
//		if reminder != nil {
//			reminder_date := time.Parse("2006-01-02 15:04:05", reminder.date)
//			if lessThanTwoWeeks.After(reminder_date) {
//				message := fmt.Sprintf("Die letzte Wurmkur war am %s", reminder.date)
//				return message
//			}
//		}
//	}
//}

func ReadReminder(db *sql.DB) []*Reminders {
	sqlStatement := `SELECT * FROM reminders`
	reminders := make([]*Reminders, 0)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println("Error getting all wurmkur entries: ", err)
	}
	if rows != nil {
		for rows.Next() {
			reminder := new(Reminders)
			_ = rows.Scan(&reminder.date)
			reminders = append(reminders, reminder)
		}
	}
	return (remindersSlice(reminders).RemindersString())
}
*/
