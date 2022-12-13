package boltdb

import (
	"errors"
	"strconv"

	"github.com/AlexKomzzz/collectivity-tlg-bot/pkg/storage"
	"github.com/boltdb/bolt"
)

type TokenStorage struct {
	db *bolt.DB
}

func NewTokenStorage(db *bolt.DB) *TokenStorage {
	return &TokenStorage{db: db}
}

func (s *TokenStorage) Save(chatID int64, debt string, bucket storage.Bucket) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatID), []byte(debt))
	})
}

// определение JWT из БД по chatID
func (s *TokenStorage) Get(chatID int64, bucket storage.Bucket) (string, error) {
	var token string

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		token = string(b.Get(intToBytes(chatID)))
		return nil
	})

	if token == "" {
		return "", errors.New("not found")
	}

	return token, err
}

// удаление токена из БД
func (s *TokenStorage) Delete(chatID int64, bucket storage.Bucket) error {
	return s.db.Update(func(tx *bolt.Tx) error {

		// удаление данных о задожденности из ведра Debt
		b := tx.Bucket([]byte(storage.Debt))
		err := b.Delete(intToBytes(chatID))
		if err != nil {
			return err
		}

		// удаление токена из ведра AccessTokens
		bTwo := tx.Bucket([]byte(storage.AccessTokens))
		return bTwo.Delete(intToBytes(chatID))
	})
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
