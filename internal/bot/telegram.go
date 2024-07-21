package bot

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chromedp/cdproto/page"
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/miriam-samuels/telegram-bot/internal/api"
	"github.com/miriam-samuels/telegram-bot/internal/helper"
	types "github.com/miriam-samuels/telegram-bot/internal/repository"
	"github.com/miriam-samuels/telegram-bot/internal/template"
)

// available commands
const (
	CommandStart                 = "/start"
	CommandNews                  = "/news"
	CommandSpaces                = "/spaces"
	CommandTopVolumeSolanaNFT    = "/topvolumesol"
	CommandTopTrendingSolanaNFT  = "/toptrendingsol"
	CommandTopGainersSolanaNFT   = "/topgainerssol"
	CommandTopMarketCapSolanaNFT = "/topmktcapsol"
	CommandSearchCollection      = "/searchc "
	CommandCollectionFloor       = "/floorprice "
	CommandCollectionListings    = "/listings "
	CommandCollectionVolume      = "/volume "
	CommandCollectionHolders     = "/holders "
	CommandCollectionLoans       = "/loans "
	CommandCollectionRaffles     = "/raffles "
)

// availabel links
const (
	LinkExploreCollections = "https://kyzzen.io/explore"
	LinkTable              = "https://pr-1576.ddv7k8ml5gut2.amplifyapp.com/telegram-table"
	LinkGraph              = "https://pr-1576.ddv7k8ml5gut2.amplifyapp.com/telegram-graph/"
)

var kyzzenRedirect = tgApi.NewInlineKeyboardMarkup(
	tgApi.NewInlineKeyboardRow(
		tgApi.NewInlineKeyboardButtonURL("Explore More NFT Collections", LinkExploreCollections),
	),
)

type TelegramBot struct {
	Api *tgApi.BotAPI
}

type TelegramMessage struct {
	Text        string
	User        int64
	ParseMode   string
	ReplyMarkup interface{}
	File        tgApi.FileBytes
}

var Telegram *TelegramBot

func NewTelegramBot(TOKEN string) (*TelegramBot, error) {
	bot, err := tgApi.NewBotAPI(TOKEN)
	if err != nil {
		return nil, err
	}

	// bot.Debug = true // allow debugging

	tBot := TelegramBot{Api: bot}

	commands := []tgApi.BotCommand{
		{Command: CommandStart, Description: "Start the bot and see list of commands"},
		{Command: CommandNews, Description: "Get the latest NFT news"},
		{Command: CommandSpaces, Description: "View upcoming X Spaces on NFTs"},
		{Command: CommandTopVolumeSolanaNFT, Description: "View top volume NFTs on solana"},
		{Command: CommandTopTrendingSolanaNFT, Description: "View trending NFTs on solana"},
		{Command: CommandTopGainersSolanaNFT, Description: "View top gainers on solana"},
		{Command: CommandTopMarketCapSolanaNFT, Description: "View NFTs with highest market cap on solana"},
		{Command: CommandSearchCollection, Description: "Search for an NFT collection e.g /searchc smb gen2"},
	}

	// Create a custom request to set bot commands
	setCommandsConfig := tgApi.NewSetMyCommands(commands...)
	_, err = bot.Request(setCommandsConfig)
	if err != nil {
		log.Panic(err)
	}

	Telegram = &tBot

	return &tBot, err
}

func (bot *TelegramBot) SendChannelMessage(message string) error {
	var telegramChannelID = os.Getenv("TELEGRAM_CHANNEL_ID")
	if os.Getenv("ENV") == "dev" {
		telegramChannelID = os.Getenv("TELEGRAM_CHANNEL_ID_DEV")
	}
	//  create a new message
	msg := tgApi.NewMessageToChannel(telegramChannelID, message)
	msg.ParseMode = "HTML"

	// Send the message
	if _, err := bot.Api.Send(msg); err != nil {
		return err
	}

	log.Println("Message sent successfully")
	return nil
}

func (bot *TelegramBot) SendUserMessage(message TelegramMessage) error {

	mode := "HTML"
	if message.ParseMode != "" {
		mode = message.ParseMode
	}
	//  create a new message
	msg := tgApi.NewMessage(message.User, message.Text)
	msg.ParseMode = mode
	msg.ReplyMarkup = message.ReplyMarkup

	// Send the message
	if _, err := bot.Api.Send(msg); err != nil {
		return err
	}

	log.Println("Message sent successfully")
	return nil
}

func (bot *TelegramBot) SendImage2User(message TelegramMessage) error {
	//  create a new message
	msg := tgApi.NewPhoto(message.User, message.File)
	msg.ParseMode = "HTML"
	msg.Caption = message.Text
	msg.ReplyMarkup = message.ReplyMarkup

	// Send the message
	if _, err := bot.Api.Send(msg); err != nil {
		return err
	}

	log.Println("Image sent successfully")
	return nil
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
	}
	var message struct {
		UpdateID int           `json:"update_id"`
		Message  tgApi.Message `json:"message"`
	}
	if err := json.Unmarshal(body, &message); err != nil {
		log.Printf("failed to unmarshal body: %v", err)
		return
	}
	err = handleMessage(&message.Message)
	if err != nil {
		log.Printf("failed to send message: %v", err)
	}
}

func (bot *TelegramBot) ListenForUpdate() {
	u := tgApi.NewUpdate(0) // create new update listener
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	updates := bot.Api.GetUpdatesChan(u)

	go receiveUpdates(ctx, updates)

	// Tell the user the bot is online
	log.Println("Start listening for updates.")

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
		err := handleMessage(update.Message)
		if err != nil {
			Telegram.SendUserMessage(TelegramMessage{Text: template.ErrorMessage, User: int64(update.Message.Chat.ID)})
		}
	// Handle button clicks
	case update.CallbackQuery != nil:
		err := handleButton(update.CallbackQuery)
		if err != nil {
			Telegram.SendUserMessage(TelegramMessage{Text: template.ErrorMessage, User: int64(update.Message.Chat.ID)})
		}
	}
}

func handleMessage(message *tgApi.Message) error {
	var buf []byte // Store image bytes

	user := message.From
	text := helper.DivideFromSymbol(message.Text, "@")
	chat := message.Chat

	log.Printf("%s wrote %s", user.UserName, text)
	if strings.HasPrefix(text, "/") { //  for general commands
		switch text {
		case CommandStart:
			return Telegram.SendUserMessage(TelegramMessage{Text: template.WelcomeMessage, User: int64(chat.ID)})

		case CommandNews:
			msg, err := api.GetNftNews()
			if err != nil {
				return err
			}
			return Telegram.SendUserMessage(TelegramMessage{Text: msg, User: int64(chat.ID)})

		case CommandSpaces:
			msg, err := api.GetSpaces()
			if err != nil {
				return err
			}
			return Telegram.SendUserMessage(TelegramMessage{Text: msg, User: int64(chat.ID)})

		case CommandTopVolumeSolanaNFT:
			s := helper.Screenshot{
				Url:        LinkTable + "?cmd=topvolumesolananft",
				ImageBytes: &buf,
				Selector:   "div.explore-collections",
			}

			err := s.TakeScreenshot()
			if err != nil || len(buf) == 0 {
				return err
			}

			file := tgApi.FileBytes{
				Name:  "file",
				Bytes: buf,
			}

			return Telegram.SendImage2User(TelegramMessage{User: int64(chat.ID), File: file, ReplyMarkup: kyzzenRedirect})

		case CommandTopTrendingSolanaNFT:
			s := helper.Screenshot{
				Url:        LinkTable + "?cmd=toptrendingsolananft",
				ImageBytes: &buf,
				Selector:   "div.explore-collections",
			}

			err := s.TakeScreenshot()
			if err != nil || len(buf) == 0 {
				return err
			}

			file := tgApi.FileBytes{
				Name:  "file",
				Bytes: buf,
			}

			return Telegram.SendImage2User(TelegramMessage{User: int64(chat.ID), File: file, ReplyMarkup: kyzzenRedirect})

		case CommandTopGainersSolanaNFT:
			s := helper.Screenshot{
				Url:        LinkTable + "?cmd=topgainerssolananft",
				ImageBytes: &buf,
				Selector:   "div.explore-collections",
			}

			err := s.TakeScreenshot()
			if err != nil || len(buf) == 0 {
				return err
			}

			file := tgApi.FileBytes{
				Name:  "file",
				Bytes: buf,
			}

			return Telegram.SendImage2User(TelegramMessage{User: int64(chat.ID), File: file, ReplyMarkup: kyzzenRedirect})

		case CommandTopMarketCapSolanaNFT:
			s := helper.Screenshot{
				Url:        LinkTable + "?cmd=topmarketcapsolananft",
				ImageBytes: &buf,
				Selector:   "div.explore-collections",
			}

			err := s.TakeScreenshot()
			if err != nil || len(buf) == 0 {
				return err
			}

			file := tgApi.FileBytes{
				Name:  "file",
				Bytes: buf,
			}

			return Telegram.SendImage2User(TelegramMessage{User: int64(chat.ID), File: file, ReplyMarkup: kyzzenRedirect})

		default:
			//  for specific commands
			if strings.Contains(text, CommandSearchCollection) {
				return handleCollectionSearch(text, chat)
			} else { //list of available commands
				return Telegram.SendUserMessage(TelegramMessage{Text: template.WelcomeMessage, User: int64(chat.ID)})
			}
		}

	}
	return nil
}

func handleButton(callback *tgApi.CallbackQuery) error {
	var s helper.Screenshot
	var buf []byte
	var caption string
	var err error
	var inlineKeyboard tgApi.InlineKeyboardMarkup
	chatId := int64(callback.Message.Chat.ID)

	if strings.Contains(callback.Data, CommandSearchCollection) {
		handleCollectionSearch(callback.Data, callback.Message.Chat)
		return nil
	}
	if strings.Contains(callback.Data, CommandCollectionRaffles) {
		collectionName := strings.Replace(callback.Data, CommandCollectionRaffles, "", 1)
		msg, err := api.GetRaffles(collectionName)
		if err != nil {
			return err
		}
		inlineKeyboard = tgApi.NewInlineKeyboardMarkup(
			tgApi.NewInlineKeyboardRow(
				tgApi.NewInlineKeyboardButtonURL("Visit Kyzzen for more raffles", "https://kyzzen.io/raffles-aggregator"),
			),
		)
		return Telegram.SendUserMessage(TelegramMessage{Text: msg, User: chatId, ReplyMarkup: inlineKeyboard})
	}

	if strings.Contains(callback.Data, CommandCollectionFloor) { // get floor price details of a collection
		collectionName := strings.Replace(callback.Data, CommandCollectionFloor, "", 1)
		caption, err = api.GetACollection(collectionName, template.CollectionFloor)
		if err != nil {
			return err
		}

		s = helper.Screenshot{
			Url:        LinkGraph + helper.FormatLink(collectionName) + "?activeTab=floor",
			ImageBytes: &buf,
			Selector:   "div#floor-trend",
		}
		inlineKeyboard = tgApi.NewInlineKeyboardMarkup(
			tgApi.NewInlineKeyboardRow(
				tgApi.NewInlineKeyboardButtonURL("Visit Kyzzen for advanced analytics", "https://kyzzen.io/collections/"+helper.FormatLink(collectionName)+"?activeTab=analytics"),
			),
		)
	} else if strings.Contains(callback.Data, CommandCollectionListings) { // get floor price details of a collection
		collectionName := strings.Replace(callback.Data, CommandCollectionListings, "", 1)
		caption, err = api.GetACollection(collectionName, template.CollectionFloor)
		if err != nil {
			return err
		}

		s = helper.Screenshot{
			Url:        LinkGraph + helper.FormatLink(collectionName) + "?activeTab=listings",
			ImageBytes: &buf,
			Selector:   "div#no_listings",
		}
		inlineKeyboard = tgApi.NewInlineKeyboardMarkup(
			tgApi.NewInlineKeyboardRow(
				tgApi.NewInlineKeyboardButtonURL("Visit Kyzzen for advanced analytics", "https://kyzzen.io/collections/"+helper.FormatLink(collectionName)+"?activeTab=analytics"),
			),
		)
	} else if strings.Contains(callback.Data, CommandCollectionVolume) { // get volume details of a collection
		collectionName := strings.Replace(callback.Data, CommandCollectionVolume, "", 1)
		caption, err = api.GetACollection(collectionName, template.CollectionVol)
		if err != nil {
			return err
		}

		s = helper.Screenshot{
			Url:        LinkGraph + helper.FormatLink(collectionName) + "?activeTab=sales",
			ImageBytes: &buf,
			Selector:   "div#volume-sales",
		}

		inlineKeyboard = tgApi.NewInlineKeyboardMarkup(
			tgApi.NewInlineKeyboardRow(
				tgApi.NewInlineKeyboardButtonURL("Visit Kyzzen for advanced analytics", "https://kyzzen.io/collections/"+helper.FormatLink(collectionName)+"?activeTab=analytics"),
			),
		)
	} else if strings.Contains(callback.Data, CommandCollectionHolders) { // get holderss details of a collection
		collectionName := strings.Replace(callback.Data, CommandCollectionHolders, "", 1)
		caption, err = api.GetACollection(collectionName, template.CollectionHolders)
		if err != nil {
			return err
		}

		s = helper.Screenshot{
			Url:        LinkGraph + helper.FormatLink(collectionName) + "?activeTab=holders",
			ImageBytes: &buf,
			Selector:   "div#holders",
		}

		inlineKeyboard = tgApi.NewInlineKeyboardMarkup(
			tgApi.NewInlineKeyboardRow(
				tgApi.NewInlineKeyboardButtonURL("Visit Kyzzen for advanced analytics", "https://kyzzen.io/collections/"+helper.FormatLink(collectionName)+"?activeTab=analytics"),
			),
		)

	} else if strings.Contains(callback.Data, CommandCollectionLoans) {
		collectionName := strings.Replace(callback.Data, CommandCollectionLoans, "", 1)
		caption, err = api.GetLoansData(collectionName)
		if err != nil {
			return err
		}

		s = helper.Screenshot{
			Url:        "https://kyzzen.io/collections/" + helper.FormatLink(collectionName) + "?activeTab=loans",
			ImageBytes: &buf,
			Selector:   "div.collection-loans",
			Viewport:   page.Viewport{Width: 820, Height: 510},
		}

		inlineKeyboard = tgApi.NewInlineKeyboardMarkup(
			tgApi.NewInlineKeyboardRow(
				tgApi.NewInlineKeyboardButtonURL("Visit Kyzzen for advanced analytics", "https://kyzzen.io/collections/"+helper.FormatLink(collectionName)+"?activeTab=loans"),
			),
		)
	}

	log.Println(s.Url)

	err = s.TakeScreenshot()
	if err != nil || len(buf) == 0 {
		return err
	}

	file := tgApi.FileBytes{
		Name:  "file",
		Bytes: buf,
	}

	return Telegram.SendImage2User(TelegramMessage{Text: caption, User: chatId, File: file, ReplyMarkup: inlineKeyboard})
}

func handleCollectionSearch(text string, chat *tgApi.Chat) error {
	collectionName := strings.Replace(text, CommandSearchCollection, "", 1)

	msg, err := api.GetACollection(collectionName, template.CollectionInfo)
	if err != nil {
		var customError *types.CustomError // stores info about the collection

		if errors.As(err, &customError) {
			var similarCollectionsKeyboardRows [][]tgApi.InlineKeyboardButton

		loop:
			for _, col := range api.Collections {
				if helper.CheckWordSimilarity(customError.Message, col["name"].(string)) {
					row := tgApi.NewInlineKeyboardRow(
						tgApi.NewInlineKeyboardButtonData(col["name"].(string), CommandSearchCollection+col["name"].(string)),
					)
					similarCollectionsKeyboardRows = append(similarCollectionsKeyboardRows, row)
				}

				if len(similarCollectionsKeyboardRows) == 10 {
					break loop
				}
			}

			if len(similarCollectionsKeyboardRows) == 0 {
				msg = "Collection Not Found"
				return Telegram.SendUserMessage(TelegramMessage{Text: msg, User: int64(chat.ID)})
			}

			inlineKeyboard := tgApi.NewInlineKeyboardMarkup(
				similarCollectionsKeyboardRows...,
			)
			return Telegram.SendUserMessage(TelegramMessage{Text: msg, User: int64(chat.ID), ReplyMarkup: inlineKeyboard})
		}
		return err
	}

	inlineKeyboard := tgApi.NewInlineKeyboardMarkup(
		tgApi.NewInlineKeyboardRow(
			tgApi.NewInlineKeyboardButtonData("Floor Price", CommandCollectionFloor+collectionName),
			tgApi.NewInlineKeyboardButtonData("Volume", CommandCollectionVolume+collectionName),
			tgApi.NewInlineKeyboardButtonData("Holders", CommandCollectionHolders+collectionName),
		),
		tgApi.NewInlineKeyboardRow(
			tgApi.NewInlineKeyboardButtonData("Listings", CommandCollectionListings+collectionName),
			tgApi.NewInlineKeyboardButtonData("Loan", CommandCollectionLoans+collectionName),
			tgApi.NewInlineKeyboardButtonData("Raffles", CommandCollectionRaffles+collectionName),
		),
	)
	return Telegram.SendUserMessage(TelegramMessage{Text: msg, User: int64(chat.ID), ReplyMarkup: inlineKeyboard})

}
