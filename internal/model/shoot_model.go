package model

import "time"

type Shoot struct {
	Id            int
	ShootDate     time.Time
	StartTime     time.Time
	EndTime       time.Time
	ShootPrice    int
	ShootLocation string
	ShootType     string
	Notes         string
	CreatedAt     time.Time

	Clients []ShootClientInfo
}
