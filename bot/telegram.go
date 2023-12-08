package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

// sendMessageToTelegram sends a message to a Telegram chat.
func sendMessageToTelegram(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending message to Telegram:", err)
	}
}
