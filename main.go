package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var db_host = "34.107.46.5"
var db_port = "3306"
var db_user = os.Getenv("GARYS_TRICKS_DB_USER")
var db_pass = os.Getenv("GARYS_TRICKS_DB_PASS")
var db_name = "garybot"

type Tricks struct {
	trick_name        string
	trick_description string
	trick_difficulty  string
}

type DiaryEntry struct {
	text string
	date string
}

type Tasks struct {
	task    string
	weekday string
}

type messageSlice []*Tricks
type diarySlice []*DiaryEntry
type taskSlice []*Tasks

func main() {
	db, err := sql.Open("mysql", db_user+":"+db_pass+"@tcp("+db_host+":"+db_port+")/"+db_name)
	if err != nil {
		log.Println("Connection to database failed:", err)
		panic(err)
	}
	defer db.Close()

	telegramToken := os.Getenv("GARYS_TRICKS_TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Fatal("GARYS_TRICKS_TELEGRAM_TOKEN environment variable not set")
	}
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			switch update.Message.Command() {
			case "tricks":
				msg.Text = GetAllTricks(db)
			case "neuerTrick":
				InsertNewTrick(db, update.Message.Text)
				msg.Text = "Neuer Trick wurde hinzufügt!"
			case "tagebuch":
				msg.Text = GetAllDiaryEntries(db)
			case "neuerEintrag":
				InsertNewDiaryEntry(db, update.Message.Text)
				msg.Text = "Neuer Tagebucheintrag wurde hinzugefügt!"
			case "wochenplan":
				msg.Text = GetAllTasks(db)
			case "neueAufgabe":
				InsertNewTask(db, update.Message.Text)
				msg.Text = "Neue Aufgabe wurde hinzugefügt!"
			default:
				msg.Text = "Wuff, ich weiß nicht, was du von mir willst."
			}
			bot.Send(msg)
		}
	}
}

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
			s = append(s, fmt.Sprintf("%s %s %s", u.trick_name, u.trick_description, u.trick_difficulty))
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
