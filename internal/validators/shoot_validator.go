package validators

import (
	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
)

type ShootValidator struct {
	clock *clock.Clock
}

func NewShootValidator(clock *clock.Clock) *ShootValidator {
	return &ShootValidator{clock: clock}
}

func (v *ShootValidator) ValidateCreateShoot(req *model.CreateShootRequest) error {
	if req.ShootDate.Before(v.clock.Now()) {
		return errors.New(errors.ErrCodeValidation, "shoot date must be in the future")
	}

	if len(req.Clients) == 0 {
		return errors.New(errors.ErrCodeValidation, "at least one client is required")
	}

	if err := v.validateShootPrice(req.ShootPrice); err != nil {
		return err
	}

	return nil
}

func (v *ShootValidator) validateShootPrice(price int) error {
	if price < 0 {
		return errors.New(errors.ErrCodeValidation, "price must be greater than zero")
	}

	return nil
}
