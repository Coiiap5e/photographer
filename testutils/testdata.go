package testutils

import (
	"time"

	"github.com/Coiiap5e/photographer/model"
)

func CreateTestClient() *model.Client {
	return &model.Client{
		FirstName:        "Ivan",
		LastName:         "Ivanov",
		Phone:            "+7(900)000-00-00",
		SocialNetworkUrl: "",
	}
}

func CreateTestClientWithOptions(option ...func(client *model.Client)) *model.Client {
	client := CreateTestClient()
	for _, opt := range option {
		opt(client)
	}
	return client
}

func CreateTestShoot(clientID int) *model.Shoot {
	return &model.Shoot{
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

func CreateTestShootWithOptions(clientID int, option ...func(shoot *model.Shoot)) *model.Shoot {
	shoot := CreateTestShoot(clientID)
	for _, opt := range option {
		opt(shoot)
	}
	return shoot
}
