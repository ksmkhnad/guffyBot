// utils/file.go
package utils

import (
	"encoding/json"
	"github.com/ksmkhnad/guffyBot/models"
	"io/ioutil"
)

// WriteSubscribersToFile writes subscriber data to a JSON file.
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
