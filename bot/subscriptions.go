package bot

import (
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ksmkhnad/guffyBot/models"
)

// SubscriptionsHandler handles sending subscription details to a specific family member.
func SubscriptionsHandler(bot *tgbotapi.BotAPI, chatID int64, familyMembers []models.FamilyMember, exchangeRate float64) {
	currentMember := findFamilyMemberByTelegramID(chatID, familyMembers)

	if currentMember == nil {
		log.Printf("No family member found with Telegram ID: %d\n", chatID)
		return
	}

	totalSpotifyKZT, totalAppleKZT := calculateTotalDue(currentMember)
	membersCount := 6
	totalSpotifyKZT /= float64(membersCount)
	totalAppleKZT /= float64(membersCount)

	/*// Convert to Euro and USD
	totalSpotifyEUR := totalSpotifyKZT * exchangeRate
	totalAppleEUR := totalAppleKZT * exchangeRate*/

	// Send messages indicating the amount each family member needs to pay for each subscribed service
	if totalSpotifyKZT > 0 {
		messageSpotify := fmt.Sprintf("%s: For Spotify - Pay %.2f KZT", currentMember.Name, totalSpotifyKZT)
		sendMessageToTelegram(bot, chatID, messageSpotify)
	}

	if totalAppleKZT > 0 {
		messageApple := fmt.Sprintf("%s: For Apple - Pay %.2f KZT", currentMember.Name, totalAppleKZT)
		sendMessageToTelegram(bot, chatID, messageApple)
	}
}

// findFamilyMemberByTelegramID finds the FamilyMember associated with the given Telegram ID.
func findFamilyMemberByTelegramID(telegramID int64, familyMembers []models.FamilyMember) *models.FamilyMember {
	for _, member := range familyMembers {
		if member.TelegramID == telegramID {
			return &member
		}
	}
	return nil
}

// calculateTotalDue calculates the total amount due for each subscription.
func calculateTotalDue(currentMember *models.FamilyMember) (float64, float64) {
	var totalSpotifyKZT, totalAppleKZT float64

	for s, subscribed := range currentMember.Subscriptions {
		if subscribed {
			// Set the appropriate price for the chosen service
			var price float64
			switch s {
			case "Spotify":
				price = 3900
			case "Apple":
				price = 1680
			}

			// Accumulate the total amounts
			if s == "Spotify" {
				totalSpotifyKZT += price
			} else if s == "Apple" {
				totalAppleKZT += price
			}
		}
	}

	return totalSpotifyKZT, totalAppleKZT
}
