package config

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"google.golang.org/api/forms/v1"
	"gopkg.in/gomail.v2"
	"log"
	"sync"
)

var (
	cfgSingleton         sync.Once
	postgreSingleton     sync.Once
	redisSingleton       sync.Once
	mailerSingleton      sync.Once
	engineSingleton      sync.Once
	googleFormsSingleton sync.Once

	Instance   *Config
	Postgre    *sql.DB
	Redis      *redis.Client
	Mailer     *gomail.Dialer
	Engine     *gin.Engine
	GoogleForm *forms.Service
)

type Config struct {
	AppName        string `mapstructure:"APP_NAME"`
	AppDescription string `mapstructure:"APP_DESC"`
	AppDebug       bool   `mapstructure:"APP_DEBUG"`
	AppURL         string `mapstructure:"APP_URL"`

	SupabaseProjectURL string `mapstructure:"SUPABASE_PROJECT_URL"`
	SupabaseAPIKey     string `mapstructure:"SUPABASE_API_KEY"`
	SupabaseAPIKeyRoot string `mapstructure:"SUPABASE_API_KEY_ROOT"`
	SupabaseJWTSecret  string `mapstructure:"SUPABASE_JWT_SECRET"`

	PostgreDsnURL string `mapstructure:"POSTGRE_DSN_URL"`
	RedisDsnURL   string `mapstructure:"REDIS_DSN_URL"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	SentryDsnURL  string `mapstructure:"SENTRY_DSN_URL"`

	MailHost     string `mapstructure:"MAIL_HOST"`
	MailPort     int    `mapstructure:"MAIL_PORT"`
	MailUsername string `mapstructure:"MAIL_USERNAME"`
	MailPassword string `mapstructure:"MAIL_PASSWORD"`

	GoogleCredentialPath string `mapstructure:"GOOGLE_CREDENTIAL_PATH"`
}

func LoadEnv() {
	// notify that app try to load config file
	log.Println("Load configuration file . . . .")
	cfgSingleton.Do(func() {
		// find environment file
		viper.AutomaticEnv()
		// error handling for specific case
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				panic(".env file not found!, please copy .env.example and paste as .env")
			}
			panic(fmt.Sprintf("ENV_ERROR: %s", err.Error()))
		}
		// notify that config file is ready
		log.Println("configuration file: ready")
		// extract config to struct
		if err := viper.Unmarshal(&Instance); err != nil {
			panic(fmt.Sprintf("ENV_ERROR: %s", err.Error()))
		}
	})
}
