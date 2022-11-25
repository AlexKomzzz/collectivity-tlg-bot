package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/AlexKomzzz/collectivity-tlg-bot/pkg/config"
	"github.com/AlexKomzzz/collectivity-tlg-bot/pkg/storage"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type AuthServer struct {
	server *http.Server
	logger *zap.Logger

	storage storage.TokenStorage

	config *config.Config
}

type dataClient struct {
	Debt        string `json:"debt"`
	AccessToken string `json:"token"`
}

func NewAuthServer(storage storage.TokenStorage, config *config.Config) *AuthServer {
	return &AuthServer{
		storage: storage,
		config:  config,
	}
}

func (s *AuthServer) Start() error {

	s.server = &http.Server{
		Handler: s,
		Addr:    s.config.ServPort,
	}

	logger, _ := zap.NewDevelopment(zap.Fields(
		zap.String("app", "authorization server")))
	defer logger.Sync()

	s.logger = logger

	return s.server.ListenAndServe()
}

func (s *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.logger.Debug("received unavailable HTTP method request",
			zap.String("method", r.Method))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// получение chatId из URL
	chatIDQuery := r.URL.Query().Get("chat_id")
	if chatIDQuery == "" {
		s.logger.Debug("received empty chat_id query param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// конвертация из строки в инт
	chatID, err := strconv.ParseInt(chatIDQuery, 10, 64)
	if err != nil {
		s.logger.Debug("received invalid chat_id query param",
			zap.String("chat_id", chatIDQuery))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// var dataBodyReq []byte
	dataClient := &dataClient{}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Debug("invalid body req",
			zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(bodyBytes, dataClient)
	if err != nil {
		s.logger.Debug("error Unmarshal body request",
			zap.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("data client: ", dataClient)
	// сохранение token в БД по chatID
	if err := s.saveTokenInDB(dataClient.AccessToken, chatID); err != nil {
		s.logger.Debug("failed to save access token",
			zap.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// сохранение debt в БД по chatID
	if err := s.saveDebtInDB(dataClient.Debt, chatID); err != nil {
		s.logger.Debug("failed to save debt",
			zap.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Location", s.config.BotURL)
	// w.WriteHeader(http.StatusBadRequest)
	// log.Println("ссылка на бота отправлена: ", s.config.BotURL)
	w.WriteHeader(http.StatusOK)

}

// сохранение debt в БД по chatID
func (s *AuthServer) saveDebtInDB(debt string, chatID int64) error {

	if err := s.storage.Save(chatID, debt, storage.Debt); err != nil {
		return errors.WithMessage(err, "failed to save access token to storage")
	}

	return nil
}

// сохранение token в БД по chatID
func (s *AuthServer) saveTokenInDB(token string, chatID int64) error {

	if err := s.storage.Save(chatID, token, storage.AccessTokens); err != nil {
		return errors.WithMessage(err, "failed to save access token to storage")
	}

	return nil
}
