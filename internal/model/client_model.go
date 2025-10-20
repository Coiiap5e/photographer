package model

import "time"

type Client struct {
	Id               int
	FirstName        string
	LastName         string
	Phone            string
	SocialNetworkUrl string
	CreatedAt        time.Time
}
