package validators

import (
	"strings"

	"github.com/Coiiap5e/photographer/internal/api/controllers/dto"
	"github.com/Coiiap5e/photographer/internal/errors"
)

func ValidateCreateClient(req *dto.CreateClientRequest) error {
	if strings.TrimSpace(req.FirstName) == "" && strings.TrimSpace(req.LastName) == "" {
		return errors.New(errors.ErrCodeValidation, "client name is required")
	}

	return nil
}
