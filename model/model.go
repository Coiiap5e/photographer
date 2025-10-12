package model

import "time"

type Client struct {
	Id               int       `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Phone            string    `json:"phone"`
	SocialNetworkUrl string    `json:"social_network_url"`
	CreatedAt        time.Time `json:"created_at"`
}

type Shoot struct {
	Id            int       `json:"id"`
	ClientId      int       `json:"client_id"`
	ShootDate     time.Time `json:"shoot_date"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	ShootPrice    int       `json:"shot_price"`
	ShootLocation string    `json:"shot_location"`
	ClientName    string    `json:"client_name"`
	ShootType     string    `json:"shot_type"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}
