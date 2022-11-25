package main

import (
	"log"

	"github.com/AlexKomzzz/collectivity-tlg-bot/pkg/config"
	"github.com/AlexKomzzz/collectivity-tlg-bot/pkg/server"
	"github.com/AlexKomzzz/collectivity-tlg-bot/pkg/storage"
	"github.com/AlexKomzzz/collectivity-tlg-bot/pkg/storage/boltdb"
	"github.com/AlexKomzzz/collectivity-tlg-bot/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// инициализация конфиг и env
	cfg, err := config.Init()
	if err != nil {
		log.Fatalln("ошибка при инициализации файлов config и env: ", err)
	}

	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalln("ошибка при создании бота: ", err)
	}
	botApi.Debug = true

	db, err := initBolt()
	if err != nil {
		log.Fatalln("Ошибка при создании БД: ", err)
	}
	storage := boltdb.NewTokenStorage(db)

	bot := telegram.NewBot(botApi, storage, cfg.Messages)

	redirectServer := server.NewAuthServer(storage, cfg)

	go func() {
		if err := redirectServer.Start(); err != nil {
			log.Fatalln("Ошибка при запуске сервера: ", err)
		}
		log.Println("ServBot START!")

	}()
	if err := bot.Start(); err != nil {
		log.Fatalln("Ошибка при запуске бота: ", err)
	}
}

func initBolt() (*bolt.DB, error) {
	db, err := bolt.Open("/bolt_db/bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(storage.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(storage.Debt))
		return err
	}); err != nil {
		return nil, err
	}

	return db, nil
}
