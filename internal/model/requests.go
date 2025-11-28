package model

import "time"

type CreateClientRequest struct {
	FirstName        string `json:"first_name" binding:"required"`
	LastName         string `json:"last_name" binding:"required"`
	Phone            string `json:"phone"`
	SocialNetworkUrl string `json:"social_network_url"`
}

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
