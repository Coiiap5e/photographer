package request

import (
	"time"

	"github.com/Coiiap5e/photographer/internal/model"
)

type CreateClientRequest struct {
	FirstName        string `json:"first_name" binding:"required"`
	LastName         string `json:"last_name" binding:"required"`
	Phone            string `json:"phone"`
	SocialNetworkUrl string `json:"social_network_url"`
}

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