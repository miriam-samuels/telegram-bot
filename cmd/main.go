package main

import (
	"github.com/joho/godotenv"
	"github.com/miriam-samuels/telegram-bot/internal/api"
	"github.com/miriam-samuels/telegram-bot/internal/bot"
	"github.com/robfig/cron"
	"log"
	"os"
)

var (
	telegramToken string
)

func init() {
	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	telegramToken = os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}
}

func main() {
	// initate telegram bot
	tBot, err := bot.NewTelegramBot(telegramToken)
	if err != nil {
		log.Fatal("error getting bot :: ", err)
	}

	// Initialize cron scheduler
	c := cron.New()
	// Define the task to be run at a specific time every day
	err = c.AddFunc("0 16 * * *", func() {
		message, err := api.GetNftNews()
		if err != nil {
			log.Println(err)
		}
		tBot.SendChannelMessage(message)
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	err = c.AddFunc("0 12 * * *", func() {
		message, err := api.GetSpaces()
		if err != nil {
			log.Println(err)
		}
		tBot.SendChannelMessage(message)
	})

	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	// Start the cron scheduler
	c.Start()

	bot.Telegram.ListenForUpdate()

	// Keep the application running
	select {}
}
