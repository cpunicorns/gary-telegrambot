package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var db_host = "34.141.10.192"
var db_port = "5400"
var db_user = os.Getenv("GARYS_TRICKS_DB_USER")
var db_pass = os.Getenv("GARYS_TRICKS_DB_PASS")
var db_name = "garys_tricks"

type Tricks struct {
	trick_name        string
	trick_description string
	trick_progress    string
}

type messageSlice []*Tricks

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_pass, db_name)

	db, err := sql.Open("postgres", psqlconn)
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
			case "Tricks":
				msg.Text = GetAllTricks(db)
			case "NewTrick":
				InsertNewTrick(db, update.Message.Text)
				msg.Text = "New trick added"
			case "DiaryOverview":
				msg.Text = "Here you can see the overview of your diary"
			case "NewDiaryEntry":
				msg.Text = "Here you can add a new diary entry"
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		}
	}
}

func InsertNewTrick(db *sql.DB, message string) {
	log.Println("Message text: ", message)
	trick := strings.Split(message, ",")
	trick_name := trick[0]
	trick_description := trick[1]
	trick_progress := trick[2]
	sqlStatement := `INSERT INTO tricks (trick_name, trick_description, trick_progress) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, trick_name, trick_description, trick_progress)
	if err != nil {
		panic(err)
	}
}

func (message messageSlice) String() string {
	var s []string
	for _, u := range message {
		if u != nil {
			s = append(s, fmt.Sprintf("%s %s", u.trick_name, u.trick_description, u.trick_progress))
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
		_ = rows.Scan(&trick.trick_name, &trick.trick_description, &trick.trick_progress)
		tricks = append(tricks, trick)

	}
	fmt.Println(tricks)
	return (messageSlice(tricks).String())
}
