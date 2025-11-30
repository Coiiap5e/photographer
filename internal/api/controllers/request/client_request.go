package request

type CreateClientRequest struct {
	FirstName        string `json:"first_name" binding:"required"`
	LastName         string `json:"last_name" binding:"required"`
	Phone            string `json:"phone"`
	SocialNetworkUrl string `json:"social_network_url"`
}
