package bot

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ksmkhnad/guffyBot/storage"
	"github.com/ksmkhnad/guffyBot/utils"
)

func RunBot() {
	// Create a new bot instance
	bot, err := tgbotapi.NewBotAPI("tokenBot") // Replace with your actual bot token
	if err != nil {
		log.Fatal(err)
	}

	exchangeRate, err := utils.GetExchangeRate()
	if err != nil {
		log.Println("Error fetching exchange rate:", err)
		// Set a default exchange rate or handle the error accordingly
		exchangeRate = 460.15
	}

	// Read existing family members from the JSON file
	familyMembers, err := storage.ReadSubscribersFromFile("subscribers.json")
	if err != nil {
		log.Println("Error reading subscribers from file:", err)
	}

	for i := range familyMembers {
		{
			chatID := familyMembers[i].TelegramID
			SubscriptionsHandler(bot, chatID, familyMembers, exchangeRate)
		}
	}

	// Set up updates channel
	u := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(u)

	// Process incoming messages and commands
	for update := range updates {
		// Handle callback queries
		if update.CallbackQuery != nil {
			HandleCallbackQuery(bot, update, familyMembers, exchangeRate)
			continue
		}

		if update.Message == nil {
			continue
		}

		// Handle commands
		if update.Message.IsCommand() {
			HandleCommand(bot, update, familyMembers, exchangeRate)
			continue
		}

		// Handle other types of messages as needed
	}
}
