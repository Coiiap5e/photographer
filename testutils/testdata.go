package testutils

import (
	"time"

	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
)

func CreateTestClient() *model.Client {
	return &model.Client{
		FirstName:        "Ivan",
		LastName:         "Ivanov",
		Phone:            "+7(900)000-00-00",
		SocialNetworkUrl: "",
	}
}

func CreateTestShootClients(id int) []*model.ShootClient {
	return []*model.ShootClient{
		{
			ClientID:         id,
			IsMainClient:     true,
			RelationshipType: "Main",
		},
	}
}

func CreateTestClientWithOptions(option ...func(client *model.Client)) *model.Client {
	client := CreateTestClient()
	for _, opt := range option {
		opt(client)
	}
	return client
}

func CreateTestShoot(clientID int, baseDate *clock.Clock) *model.Shoot {
	//TODO: доработать тест после многие ко многим (использовать id)
	baseTime := baseDate.Now()
	return &model.Shoot{
		ShootDate:     baseTime,
		StartTime:     baseTime.Add(15 * time.Hour),
		EndTime:       baseTime.Add(16 * time.Hour),
		ShootPrice:    1000,
		ShootLocation: "Pushkin blvd",
		ShootType:     "love story",
		Notes:         "take an umbrella",
	}

}

func CreateTestShootWithOptions(clientID int, baseDate *clock.Clock, option ...func(shoot *model.Shoot)) *model.Shoot {
	shoot := CreateTestShoot(clientID, baseDate)
	for _, opt := range option {
		opt(shoot)
	}
	return shoot
}
