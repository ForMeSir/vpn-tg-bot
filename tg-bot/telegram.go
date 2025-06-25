package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Client struct {
	bot *tgbotapi.BotAPI
}

func New(apiKey string) *Client {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Panic(err)
	}
	return &Client{
		bot: bot,
	}
}

func (c *Client) SendMessage(text string, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "Markdown"
	_, err := c.bot.Send(msg)
	return err
}
