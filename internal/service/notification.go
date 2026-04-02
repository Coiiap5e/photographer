package service

import "github.com/Coiiap5e/photographer/internal/model"

type Notifier interface {
	Notify(shoot model.Shoot) error
	NotifyMessage(message string) error
}
