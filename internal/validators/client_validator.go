package validators

import (
	"strings"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
)

type ClientValidator struct {
}

func NewClientValidator() *ClientValidator {
	return &ClientValidator{}
}

func (v *ClientValidator) ValidateCreateClient(req *model.CreateClientRequest) error {
	if strings.TrimSpace(req.FirstName) == "" && strings.TrimSpace(req.LastName) == "" {
		return errors.New(errors.ErrCodeValidation, "client name is required")
	}

	return nil
}
