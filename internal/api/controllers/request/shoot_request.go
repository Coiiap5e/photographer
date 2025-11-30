package request

import (
	"time"

	"github.com/Coiiap5e/photographer/internal/model"
)

type CreateShootRequest struct {
	ShootDate     time.Time            `json:"shootDate" binding:"required"`
	StartTime     time.Time            `json:"startTime" binding:"required"`
	EndTime       time.Time            `json:"endTime" binding:"required"`
	ShootPrice    int                  `json:"shootPrice"`
	ShootLocation string               `json:"shootLocation"`
	ShootType     string               `json:"shootType"`
	Notes         string               `json:"notes"`
	Clients       []ShootClientRequest `json:"clients" binding:"required, min=1"`
}

type ShootClientRequest struct {
	ClientID         int    `json:"clientID" binding:"required"`
	IsMainClient     bool   `json:"isMainClient" binding:"required"`
	RelationshipType string `json:"relationshipType" binding:"required"`
}

func ToShootClientDomain(clientRequests []ShootClientRequest) []*model.ShootClient {
	clients := make([]*model.ShootClient, len(clientRequests))
	for i, clientRequest := range clientRequests {
		clients[i] = &model.ShootClient{
			ClientID:         clientRequest.ClientID,
			IsMainClient:     clientRequest.IsMainClient,
			RelationshipType: clientRequest.RelationshipType,
		}
	}
	return clients
}
