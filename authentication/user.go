package authentication

// // CreateUser creates a firebase user, and returns the user model, and error
// func CreateUser(data map[string]interface{}) (*models.User, error) {
// 	params := (&auth.UserToCreate{}).
// 		Email(data["email"].(string)).
// 		EmailVerified(false).
// 		DisplayName(data["username"].(string)).
// 		Password(data["password"].(string))
// 	u, err := authClient.CreateUser(context.Background(), params)
// 	if err != nil {
// 		log.Println("Error creating user\n", err)
// 		return nil, err
// 	}

// 	return &models.User{UID: u.UID, Email: u.Email, DisplayName: u.DisplayName}, err
// }

// func UpdateUser()
