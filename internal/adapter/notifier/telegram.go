package notifier

import (
	"fmt"
	"strconv"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	bot       *tgbotapi.BotAPI
	channelID int64
}

func NewTelegramNotifier(token string, channelID string) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeTelegramBotInit, "failed to initialize telegram bot")
	}
	id, err := strconv.ParseInt(channelID, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInvalidInput, "invalid telegram channel id")
	}

	return &TelegramNotifier{
		bot:       bot,
		channelID: id,
	}, nil
}

func (n *TelegramNotifier) Notify(shoot model.Shoot) error {
	message := fmt.Sprintf(
		"<b>Upcoming Shoot!</b>\n\n<b>ID:</b> %d\n<b>Date:</b> %s\n<b>Time:</b> %s\n<b>Type:</b> %s\n<b>Location:</b> %s\n<b>Price:</b> %.2f USD",
		shoot.Id,
		shoot.ShootDate.Format("2006-01-02"),
		shoot.StartTime.Format("15:04"),
		shoot.ShootType,
		shoot.ShootLocation,
		shoot.PriceUSD,
	)
	msg := tgbotapi.NewMessage(n.channelID, message)
	msg.ParseMode = tgbotapi.ModeHTML
	if _, err := n.bot.Send(msg); err != nil {
		return errors.Wrap(err, errors.ErrCodeTelegramBotSend, "failed to send message to telegram")
	}
	return nil
}

func (n *TelegramNotifier) NotifyMessage(message string) error {
	msg := tgbotapi.NewMessage(n.channelID, message)
	msg.ParseMode = tgbotapi.ModeHTML
	if _, err := n.bot.Send(msg); err != nil {
		return errors.Wrap(err, errors.ErrCodeTelegramBotSend, "failed to send message to telegram")
	}
	return nil
}
