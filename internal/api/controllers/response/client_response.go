package response

import (
	"time"

	"github.com/Coiiap5e/photographer/internal/model"
)

type ClientResponse struct {
	ID               int       `json:"id"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	Phone            string    `json:"phone"`
	SocialNetworkUrl string    `json:"socialNetworkUrl"`
	CreatedAt        time.Time `json:"createdAt"`
}

func ToClientResponse(client *model.Client) *ClientResponse {
	return &ClientResponse{
		ID:               client.Id,
		FirstName:        client.FirstName,
		LastName:         client.LastName,
		Phone:            client.Phone,
		SocialNetworkUrl: client.SocialNetworkUrl,
		CreatedAt:        client.CreatedAt,
	}
}

func ToClientsResponse(clients []*model.Client) []ClientResponse {
	responses := make([]ClientResponse, len(clients))
	for i, client := range clients {
		responses[i] = *ToClientResponse(client)
	}
	return responses
}
