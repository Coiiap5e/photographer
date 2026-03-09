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
	PriceUSD      float64
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Clients []ShootClientInfo
}

type ShootDateTimePatch struct {
	ShootDate time.Time
	StartTime time.Time
	EndTime   time.Time
}
