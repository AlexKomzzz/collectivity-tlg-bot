package telegram

import "github.com/AlexKomzzz/collectivity-tlg-bot/pkg/storage"

func (b *Bot) getDebtUser(chatID int64) (string, error) {
	return b.storage.Get(chatID, storage.Debt)
}
