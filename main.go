package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var db_host = "34.107.46.5"
var db_port = "3306"
var db_user = os.Getenv("GARYS_TRICKS_DB_USER")
var db_pass = os.Getenv("GARYS_TRICKS_DB_PASS")
var db_name = "garys_tricks"

type Tricks struct {
	trick_name        string
	trick_description string
	trick_difficulty  string
}

type messageSlice []*Tricks

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

	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Garybot is up and running!")
}

func InsertNewTrick(db *sql.DB, message string) {
	log.Println("Message text: ", message)
	unfiltered_message := strings.Replace(message, "/NewTrick ", "", -1)
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
			s = append(s, fmt.Sprintf("%s %s", u.trick_name, u.trick_description, u.trick_difficulty))
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
