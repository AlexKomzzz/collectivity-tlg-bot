package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	commandStart = "start"
	commandDebt  = "debt"
	commandHelp  = "help"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandDebt:
		return b.handleGetDebt(message)
	case commandHelp:
		return b.handleHelp(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

// функция при получении команды /start
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	log.Println("Получена команда /start")
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

// ответ на запрос справки /help
func (b *Bot) handleHelp(message *tgbotapi.Message) error {
	log.Println("Получена команда /help")

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.Help)
	_, err := b.bot.Send(msg)
	return err
}

// получение debt
func (b *Bot) handleGetDebt(message *tgbotapi.Message) error {

	log.Println("Получена команда /debt")

	// получение debt из БД
	debtUser, err := b.getDebtUser(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	// отправка полученных данных debt клиенту в тлг
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(b.messages.Responses.ResultDebt, debtUser))
	_, err = b.bot.Send(msg)
	return err
}

// ответ на незнакомую команду
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	log.Println("Получена неизвестная команда")

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.UnknownCommand)
	_, err := b.bot.Send(msg)
	return err
}
