package telegram

import (
	"fmt"

	"github.com/AlexKomzzz/collectivity-tlg-bot/pkg/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// формирование сообщения о необходимости прохождения авторизации
func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {

	// формирование url для регистрации с данными о ссылке редиректа
	authLink := fmt.Sprintf("%s?redirest_url=%s/?chat_id=%d", b.messages.Responses.AuthLink, b.messages.Responses.RedirectURL, message.Chat.ID)

	msgText := fmt.Sprintf(b.messages.Responses.Start, authLink, b.messages.Responses.RegisterURL)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.storage.Get(chatID, storage.AccessTokens)
}
