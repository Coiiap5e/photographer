package model

import "time"

type ClientResponse struct {
	ID               int       `json:"id"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	Phone            string    `json:"phone"`
	SocialNetworkUrl string    `json:"socialNetworkUrl"`
	CreatedAt        time.Time `json:"createdAt"`
}

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
