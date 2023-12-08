package bot

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ksmkhnad/guffyBot/models"
	"io/ioutil"
	"log"
)

// handleCommand handles incoming commands.
func HandleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, familyMembers []models.FamilyMember, exchangeRate float64) {
	command := update.Message.Command()
	switch command {
	case "subscribe":
		HandleSubscribeCommand(bot, update, familyMembers, exchangeRate)
	default:
		// Handle other commands or do nothing for non-commands
	}
}

// handleCallbackQuery handles the callback queries for service subscription.
func HandleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update, familyMembers []models.FamilyMember, exchangeRate float64) {
	// Find the FamilyMember associated with the Telegram ID
	telegramID := update.CallbackQuery.From.ID
	name := update.CallbackQuery.From.FirstName

	var currentMember *models.FamilyMember
	for i := range familyMembers {
		if familyMembers[i].TelegramID == int64(telegramID) {
			currentMember = &familyMembers[i]
			break
		}
	}
	// Determine the service based on the callback data
	var service string
	switch update.CallbackQuery.Data {
	case "subscribe_spotify":
		service = "Spotify"
	case "subscribe_apple":
		service = "Apple"
	default:
		// Handle other services or do nothing for unknown services
		return
	}

	if currentMember != nil && currentMember.Subscriptions[service] {
		// If already subscribed, you may choose to update details or send a message indicating the existing subscription.
		sendMessageToTelegram(bot, update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("You are already subscribed to %s.", service))
		return
	}
	// If the user is not found, create a new FamilyMember
	if currentMember == nil {
		currentMember = &models.FamilyMember{
			TelegramID: int64(telegramID),
			Name:       name,
			Subscriptions: map[string]bool{
				service: true,
			},
		}
		familyMembers = append(familyMembers, *currentMember)
	}

	// Handle the chosen service from the callback query
	switch update.CallbackQuery.Data {
	case "subscribe_spotify":
		SubscribeToService(bot, update, currentMember, "Spotify", familyMembers)
	case "subscribe_apple":
		SubscribeToService(bot, update, currentMember, "Apple", familyMembers)
	}
}

// handleSubscribeCommand handles the subscribe command and prompts the user to choose a service.
func HandleSubscribeCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, familyMembers []models.FamilyMember, exchangeRate float64) {
	// Prompt the user to choose a service
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose a service to subscribe:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Spotify", "subscribe_spotify"),
			tgbotapi.NewInlineKeyboardButtonData("Apple", "subscribe_apple"),
		),
	)

	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending message to Telegram:", err)
	}
}
func writeSubscribersToFile(filename string, familyMembers []models.FamilyMember) error {
	data := models.SubscriptionData{
		FamilyMembers: familyMembers,
	}

	fileContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, fileContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

// subscribeToService handles the subscription to a specific service.
func SubscribeToService(bot *tgbotapi.BotAPI, update tgbotapi.Update, currentMember *models.FamilyMember, service string, familyMembers []models.FamilyMember) {
	if currentMember != nil {
		currentMember.Subscriptions[service] = true
	}

	// Save the updated family members to the JSON file
	err := writeSubscribersToFile("subscribers.json", familyMembers)
	if err != nil {
		log.Println("Error writing subscribers to file:", err)
	}

	// Send a confirmation message to the user
	sendMessageToTelegram(bot, update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("You have been subscribed to %s successfully!", service))

	// If user already had subscriptions, update the existing data
	if len(currentMember.Subscriptions) > 1 {
		// Calculate the total due for each subscription including the new one
		var totalSpotify, totalApple float64
		for s, subscribed := range currentMember.Subscriptions {
			if subscribed {
				// Set the appropriate price for the chosen service
				var p float64
				switch s {
				case "Spotify":
					p = 19.10
				case "Apple":
					p = 18.10
				}

				// Accumulate the total amounts
				if s == "Spotify" {
					totalSpotify += p
				} else if s == "Apple" {
					totalApple += p
				}
			}
		}
	}
}
