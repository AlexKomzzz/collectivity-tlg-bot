package telegram

import (
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	commandStart = "start"
	commandDebt  = "debt"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandDebt:
		return b.handleGetDebt(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

// функция при получении команды /start
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {

	// получение JWT из БД по chatID
	_, err := b.getAccessToken(message.Chat.ID)

	// если JWT не удалось получить, значит необходимо его запросить
	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	// при успешном получении токена из БД  - генерация и отправка сообщения, что все ок, продолжайте пользоваться ботом
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.AlreadyAuthorized)
	_, err = b.bot.Send(msg)
	return err
}

// ответ на незнакомую команду
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.UnknownCommand)
	_, err := b.bot.Send(msg)
	return err
}

// получение debt
func (b *Bot) handleGetDebt(message *tgbotapi.Message) error {

	// получение токена из БД
	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	// отправка запроса в АПИ с токеном в URL, получение данных
	//

	// отправка полученныъ данных по debt клиенту в тлг
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.ResultDebt)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) validateURL(text string) error {
	_, err := url.ParseRequestURI(text)
	return err
}
