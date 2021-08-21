package models

func ValidAPIKey(keyString string) (string, error) {
	apiKey := &APIKey{
		ID: keyString,
	}
	err := db.First(&apiKey).Error
	if err != nil || apiKey.UserID == "" {
		return "", err
	}
	return apiKey.UserID, nil
}

func GenerateAPIKey(userID string) (APIKey, error) {
	apiKey := &APIKey{
		UserID: userID,
	}

	err := db.Create(&apiKey).Error
	return *apiKey, err
}
