package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var token string
var idCounter int

const fileName = "data.json"

type Task struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	OwnerId           int    `json:"owner_id"`
	OwnerUsername     string `json:"owner_username"`
	PerformerId       int    `json:"performer_id"`
	PerformerUsername string `json:"performed_username"`
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

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Something went wrong reading file", fileName)
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		log.Fatal("Something went wrong unmarshalling data")
	}

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch cmd := update.Message.Command(); {
		case cmd == "get_my_id":
			msg.Text = fmt.Sprint(update.Message.From.ID)
		case cmd == "help":
			msg.Text = "This is task manager bot\n type /tasks"
		case cmd == "tasks":
			for _, task := range tasks {
				msg.Text += fmt.Sprintf("%d) %s\n", task.Id, task.Name)
			}
		case cmd == "new":
			newTaskName := update.Message.CommandArguments()
			task := Task{
				Id:            idCounter,
				Name:          newTaskName,
				OwnerId:       int(update.Message.From.ID),
				OwnerUsername: update.Message.From.UserName,
			}
			tasks = append(tasks, task)
			rawData, err := json.Marshal(tasks)
			if err != nil {
				log.Fatal("Error marshaling data")
			}
			err = os.WriteFile(fileName, rawData, 0644)
			if err != nil {
				log.Fatal("Error writing to file")
			}
			idCounter++

		case strings.HasPrefix(cmd, "assign_"):
			taskId, err := strconv.Atoi(strings.Split(cmd, "_")[1])
			if err != nil {
				log.Fatal("Error parsing task id")
			}
			found := false
			for _, task := range tasks {
				if task.Id == taskId {
					found = true
					task.PerformerId = int(update.Message.From.ID)
					task.PerformerUsername = update.Message.From.UserName
					msg.Text = fmt.Sprintf("Task '%s' assigned to you", task.Name)
					break
				}
			}
			if !found {
				msg.Text = "Did not found task with given id"
			}
			// TODO: overwrite file
		case cmd == "my":
			found := false
			for _, task := range tasks {
				if task.PerformerId == int(update.Message.From.ID) {
					found = true
					msg.Text += fmt.Sprintf("%d) %s by @&s", task.Id, task.Name)
				}
			}
			if !found {
				msg.Text = "No tasks assigned"
			}
		case cmd == "owner":
			found := false
			for _, task := range tasks {
				if task.OwnerId == int(update.Message.From.ID) {
					found = true
					msg.Text += fmt.Sprintf("%d) '%s' by @%s\n", task.Id, task.Name, task.OwnerUsername)
				}
			}
			if !found {
				msg.Text = "No tasks owned by you"
			}
		default:
			msg.Text = "I don't understand this command"
		}
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func updateFile() {

}
