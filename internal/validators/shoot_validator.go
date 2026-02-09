package validators

import (
	"github.com/Coiiap5e/photographer/internal/api/controllers/dto"
	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/utils/clock"
)

func ValidateCreateShoot(req *dto.CreateShootRequest, clock *clock.Clock) error {
	if req.ShootDate.Before(clock.Now()) {
		return errors.New(errors.ErrCodeValidation, "shoot date must be in the future")
	}

	if len(req.Clients) == 0 {
		return errors.New(errors.ErrCodeValidation, "at least one client is required")
	}

	if err := validateShootPrice(req.ShootPrice); err != nil {
		return err
	}

	return nil
}

func validateShootPrice(price int) error {
	if price < 0 {
		return errors.New(errors.ErrCodeValidation, "price must be greater than zero")
	}

	return nil
}

func ValidateID(id int) error {
	if id < 0 {
		return errors.New(errors.ErrCodeValidation, "id must be greater than zero")
	}

	return nil
}

func ValidateDate(req *dto.UpdateShootDateTimeRequest, clock *clock.Clock) error {
	if req.ShootDate.Before(clock.Now()) {
		return errors.New(errors.ErrCodeValidation, "shoot date must be in the future")
	}

	return nil
}
