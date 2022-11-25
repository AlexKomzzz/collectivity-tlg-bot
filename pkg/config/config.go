package config

import (
	"log"

	"github.com/spf13/viper"
)

type Messages struct {
	Responses
	Errors
}

type Responses struct {
	Start             string `mapstructure:"start"`
	RedirectURL       string `mapstructure:"redirect_url"`
	RegisterURL       string `mapstructure:"register_url"`
	AuthLink          string `mapstructure:"auth_link"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	UnknownCommand    string `mapstructure:"unknown_command"`
	ResultDebt        string `mapstructure:"result_debt"`
	GetDebtLink       string `mapstructure:"get_debt_link"`
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidURL   string `mapstructure:"invalid_url"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

type Config struct {
	TelegramToken string

	BotURL     string `mapstructure:"bot_url"`
	BoltDBFile string `mapstructure:"db_file"`
	ServPort   string `mapstructure:"serv_port"`

	Messages Messages
}

func Init() (*Config, error) {
	if err := setUpViper(); err != nil {
		log.Println("Ошибка при инициализации файла конфиг")
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		log.Println("Ошибка при парсинге даных из конфига в структуру")
		return nil, err
	}

	if err := fromEnv(&cfg); err != nil {
		log.Println("Ошибка при загрузке .env файла")
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("messages.response", &cfg.Messages.Responses); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("messages.error", &cfg.Messages.Errors); err != nil {
		return err
	}

	return nil
}

func fromEnv(cfg *Config) error {
	if err := viper.BindEnv("token_bot"); err != nil {
		log.Println("Ошибка при загрузке .env файла")
		return err
	}
	cfg.TelegramToken = viper.GetString("token_bot")

	return nil
}

func setUpViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	return viper.ReadInConfig()
}
