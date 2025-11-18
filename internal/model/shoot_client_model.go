package model

import "time"

type ShootClient struct {
	ShootID          int
	ClientID         int
	IsMainClient     bool
	RelationshipType string
	CreatedAt        time.Time
}

type ShootClientInfo struct {
	ClientID         int
	FirstName        string
	LastName         string
	Phone            string
	IsMainClient     bool
	RelationshipType string
}
