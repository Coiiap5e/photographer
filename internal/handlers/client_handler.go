package handlers

import "github.com/Coiiap5e/photographer/internal/api/controllers"

type ClientHandler struct {
	clientController *controllers.ClientController
}

func NewClientHandler(clientController *controllers.ClientController) *ClientHandler {
	return &ClientHandler{
		clientController: clientController,
	}
}
