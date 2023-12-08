package models

// Subscription represents a subscription for a service.
type Subscription struct {
	Service string
	Price   float64
}

// FamilyMember represents a family member with their subscriptions.
type FamilyMember struct {
	Name          string
	TelegramID    int64
	Subscriptions map[string]bool
}

// SubscriptionData represents the structure of the data to be stored in the JSON file.
type SubscriptionData struct {
	FamilyMembers []FamilyMember
}
