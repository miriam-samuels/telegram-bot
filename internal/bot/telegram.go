package bot

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/miriam-samuels/telegram-bot/internal/api"
)

type TelegramBot struct {
	Api *tgApi.BotAPI
}

var Telegram *TelegramBot

func NewTelegramBot(TOKEN string) (*TelegramBot, error) {
	bot, err := tgApi.NewBotAPI(TOKEN)
	if err != nil {
		return nil, err
	}

	// bot.Debug = true // allow debugging

	log.Printf("Authorized on account %s", bot.Self.UserName)

	tBot := TelegramBot{Api: bot}

	Telegram = &tBot

	return &tBot, err
}

func (bot *TelegramBot) SendChannelMessage(message string) {
	//  create a new message
	msg := tgApi.NewMessageToChannel(os.Getenv("TELEGRAM_CHANNEL_ID"), message)
	msg.ParseMode = "HTML"

	// Send the message
	if _, err := bot.Api.Send(msg); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Println("Message sent successfully")

}

func (bot *TelegramBot) SendUserMessage(message string, user int64) {
	//  create a new message
	msg := tgApi.NewMessage(user, message)
	msg.ParseMode = "HTML"

	// Send the message
	if _, err := bot.Api.Send(msg); err != nil {
		log.Fatalf("Failed to send message:: %v", err)
	}

	log.Println("Message sent successfully")

}

func (bot *TelegramBot) ListenForUpdate() {
	u := tgApi.NewUpdate(0) // create new update listener
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	updates := bot.Api.GetUpdatesChan(u)

	go receiveUpdates(ctx, updates)

	// Tell the user the bot is online
	log.Println("Start listening for updates. Press enter to stop")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}

func receiveUpdates(ctx context.Context, updates tgApi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(update)
		}
	}

}

func handleUpdate(update tgApi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)
	// Handle button clicks
	case update.CallbackQuery != nil:
		// handleButton(update.CallbackQuery)
		break
	}
}

func handleMessage(message *tgApi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.UserName, text)

	var (
		err error
		msg string
	)
	if strings.HasPrefix(text, "/") {
		switch text {
		case "/news":
			msg, err = api.GetNftNews()
			Telegram.SendUserMessage(msg, user.ID)
		case "/spaces":
			msg, err = api.GetSpaces()
			Telegram.SendUserMessage(msg, user.ID)
			// case "/menu":
			// 	err = sendMenu(message.Chat.ID)
		}
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// func sendMenu(chatId int64) error {
// 	fmt.Println(chatId)
// 	msg := tgApi.NewMessage(chatId, firstMenu)
// 	msg.ParseMode = tgApi.ModeHTML
// 	msg.ReplyMarkup = firstMenuMarkup
// 	_, err := bot.Send(msg)
// 	return nil
// }
