package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/miriam-samuels/telegram-bot/internal/api"
	"github.com/miriam-samuels/telegram-bot/internal/bot"
	"github.com/robfig/cron/v3"
)

var (
	TelegramToken string
	env           = os.Getenv("ENV")
)

func init() {
	// Find .env.yaml file
	err := godotenv.Load(".env.yaml")
	if err != nil {
		log.Fatalf("Error loading .env.yaml file: %s", err)
	}
	TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	if env == "dev" {
		TelegramToken = os.Getenv("TELEGRAM_TOKEN_DEV")
	}
	if TelegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}
}

func main() {
	// initate telegram bot
	tBot, err := bot.NewTelegramBot(TelegramToken)
	if err != nil {
		log.Fatal("error getting bot :: ", err)
	}

	// handle cron for daily post
	handleCron(tBot)

	// handle fetching collection details
	api.Collections, _ = api.GetAllCollections() // first fetch
	go handleCollectionsTicker()

	// handle listening for tg update
	bot.Telegram.ListenForUpdate()
	// handleRouter()

	select {}
}

func handleCron(tBot *bot.TelegramBot) {
	// Initialize cron scheduler
	c := cron.New(cron.WithSeconds())
	// Define the task to be run at a specific time every day
	_, err := c.AddFunc("0 0 16 * * *", func() {
		message, err := api.GetNftNews()
		if err != nil {
			log.Println(err)
		}
		if err := tBot.SendChannelMessage(message); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	_, err = c.AddFunc("0 0 12 * * *", func() {
		message, err := api.GetSpaces()
		if err != nil {
			log.Println(err)
		}
		if err := tBot.SendChannelMessage(message); err != nil {
			log.Println(err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	// Start the cron scheduler
	c.Start()
}

func handleCollectionsTicker() {
	var collectionError error

	const tickRate = 1 * time.Hour // update collection list every hour

	ticker := time.NewTicker(tickRate).C

	for range ticker {
		api.Collections, collectionError = api.GetAllCollections()
		if collectionError != nil {
			log.Printf("Failled to update colllection:: %v", collectionError)
		}
	}
}

func handleRouter() {
	port := ":8080"
	if env == "dev" {
		port = ":8000"
	}
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})
	router.HandleFunc("/webhook", bot.Webhook)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})
	log.Printf("Listening on %v", port)
	go http.ListenAndServe(port, router)
}
