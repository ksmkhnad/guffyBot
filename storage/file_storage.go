package storage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ksmkhnad/guffyBot/models"
)

// readSubscribersFromFile reads subscriber data from a JSON file.
func ReadSubscribersFromFile(filename string) ([]models.FamilyMember, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data models.SubscriptionData
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return nil, err
	}

	return data.FamilyMembers, nil
}

// writeSubscribersToFile writes subscriber data to a JSON file.
func WriteSubscribersToFile(filename string, familyMembers []models.FamilyMember) error {
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
