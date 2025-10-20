package testutils

import (
	"time"

	model2 "github.com/Coiiap5e/photographer/internal/model"
)

func CreateTestClient() *model2.Client {
	return &model2.Client{
		FirstName:        "Ivan",
		LastName:         "Ivanov",
		Phone:            "+7(900)000-00-00",
		SocialNetworkUrl: "",
	}
}

func CreateTestClientWithOptions(option ...func(client *model2.Client)) *model2.Client {
	client := CreateTestClient()
	for _, opt := range option {
		opt(client)
	}
	return client
}

func CreateTestShoot(clientID int) *model2.Shoot {
	return &model2.Shoot{
		ClientId:      clientID,
		ShootDate:     time.Now().AddDate(0, 0, 30),
		StartTime:     time.Date(0, 0, 0, 15, 0, 0, 0, time.UTC),
		EndTime:       time.Date(0, 0, 0, 16, 0, 0, 0, time.UTC),
		ShootPrice:    1000,
		ShootLocation: "Pushkin blvd",
		ShootType:     "love story",
		Notes:         "take an umbrella",
	}

}

func CreateTestShootWithOptions(clientID int, option ...func(shoot *model2.Shoot)) *model2.Shoot {
	shoot := CreateTestShoot(clientID)
	for _, opt := range option {
		opt(shoot)
	}
	return shoot
}
