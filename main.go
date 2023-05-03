package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var token string

type Task struct {
	id    int
	name  string
	owner int
}

var tasks = []Task{
	{
		id:    1,
		name:  "create tg bot",
		owner: 12312,
	},
	{
		id:    2,
		name:  "go eat",
		owner: 12312,
	},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token = os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Could not create bot")
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "help":
			msg.Text = "This is task manager bot\n type /tasks"
		case "tasks":
			for _, task := range tasks {
				msg.Text += task.name + "\n"
			}
		default:
			msg.Text = "I don't understand this command"
		}
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

}
