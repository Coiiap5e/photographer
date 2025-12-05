package response

import (
	"time"

	"github.com/Coiiap5e/photographer/internal/model"
)

type ShootResponse struct {
	ID            int                       `json:"id"`
	ShootDate     time.Time                 `json:"shotDate"`
	StartTime     time.Time                 `json:"startTime"`
	EndTime       time.Time                 `json:"endTime"`
	ShootPrice    int                       `json:"shotPrice"`
	ShootLocation string                    `json:"shotLocation"`
	ShootType     string                    `json:"shotType"`
	Notes         string                    `json:"notes"`
	CreatedAt     time.Time                 `json:"createdAt"`
	UpdatedAt     time.Time                 `json:"updatedAt"`
	Clients       []ShootClientInfoResponse `json:"clients"`
}

type ShootClientInfoResponse struct {
	ClientID         int    `json:"clientID"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Phone            string `json:"phone"`
	IsMainClient     bool   `json:"isMainClient"`
	RelationshipType string `json:"relationshipType"`
}

func ToShootClientInfoResponse(clients []model.ShootClientInfo) []ShootClientInfoResponse {
	clientsResponse := make([]ShootClientInfoResponse, len(clients))
	for i, clientInfo := range clients {
		clientsResponse[i] = ShootClientInfoResponse{
			ClientID:         clientInfo.ClientID,
			FirstName:        clientInfo.FirstName,
			LastName:         clientInfo.LastName,
			Phone:            clientInfo.Phone,
			IsMainClient:     clientInfo.IsMainClient,
			RelationshipType: clientInfo.RelationshipType,
		}
	}

	return clientsResponse
}

func ToShootResponse(shoot *model.Shoot) *ShootResponse {
	clients := ToShootClientInfoResponse(shoot.Clients)
	return &ShootResponse{
		ID:            shoot.Id,
		ShootDate:     shoot.ShootDate,
		StartTime:     shoot.StartTime,
		EndTime:       shoot.EndTime,
		ShootPrice:    shoot.ShootPrice,
		ShootLocation: shoot.ShootLocation,
		ShootType:     shoot.ShootType,
		Notes:         shoot.Notes,
		CreatedAt:     shoot.CreatedAt,
		UpdatedAt:     shoot.UpdatedAt,
		Clients:       clients,
	}
}
