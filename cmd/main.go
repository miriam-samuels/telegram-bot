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
	"github.com/robfig/cron"
	"github.com/rs/cors"
)

func init() {
	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func main() {

	// startServer()

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

// connection port and host for local environment
const (
	CONN_PORT = "6000"
)

func startServer() {
	// Get port if it exists in env file
	port := os.Getenv("PORT")

	// check if port exists in env file else use constant
	if port == "" {
		port = CONN_PORT
	}

	// create new router
	router := mux.NewRouter().StrictSlash(true)

	//  cross origin
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "OPTIONS"},
		// Debug:            true,
	}).Handler(router)

	// add more configurations to server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	// start server
	fmt.Println("starting server on port :: " + port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
