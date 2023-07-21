package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var db_host = os.Getenv("GARYBOT_DB_HOST")
var db_port = "3306"
var db_user = os.Getenv("GARYBOT_DB_USER")
var db_pass = os.Getenv("GARYBOT_DB_PASS")
var db_name = os.Getenv("GARYBOT_DB_NAME")

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
			case "cheatsheet":
				msg.Text = GetAllCheatsheetInfos(db)
			case "neueInfo":
				InsertNewCheatsheetInfo(db, update.Message.Text)
				msg.Text = "Neue Info wurde hinzugefügt!"
			default:
				msg.Text = "Wuff, ich weiß nicht, was du von mir willst."
			}
			bot.Send(msg)
		}
	}
}
