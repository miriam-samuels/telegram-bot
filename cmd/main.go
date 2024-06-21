package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/miriam-samuels/telegram-bot/internal/api"
	"github.com/miriam-samuels/telegram-bot/internal/bot"
	"github.com/robfig/cron"
)

func init() {
	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func main() {

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}
	// initate telegram bot
	tBot, err := bot.NewTelegramBot(token)
	if err != nil {
		log.Fatal("error getting bot :: ", err)
	}

	// Initialize cron scheduler
	c := cron.New()
	// Define the task to be run at a specific time every day
	err = c.AddFunc("0 16 * * *", func() {
		message := api.GetNftNews()
		tBot.SendChannelMessage(message)
	})
	c.AddFunc("0 12 * * *", func() {
		message := api.GetSpaces()
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
