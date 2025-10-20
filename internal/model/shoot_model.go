package model

import "time"

type Shoot struct {
	Id              int
	ClientId        int
	ShootDate       time.Time
	StartTime       time.Time
	EndTime         time.Time
	ShootPrice      int
	ShootLocation   string
	ClientFirstName string
	ClientLastName  string
	ShootType       string
	Notes           string
	CreatedAt       time.Time
}
