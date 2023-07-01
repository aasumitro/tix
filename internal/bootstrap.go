package internal

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/forms/v1"
	"gopkg.in/gomail.v2"
)

type boostrap struct {
	engine     *gin.Engine
	db         *sql.DB
	cache      *redis.Client
	mailer     *gomail.Dialer
	googleForm *forms.Service
}

type BoostrapOption func(*boostrap)

func WithEngine(engine *gin.Engine) BoostrapOption {
	return func(boostrap *boostrap) {
		boostrap.engine = engine
	}
}

func WithPostgreDatabase(db *sql.DB) BoostrapOption {
	return func(boostrap *boostrap) {
		boostrap.db = db
	}
}

func WithRedisCache(cache *redis.Client) BoostrapOption {
	return func(boostrap *boostrap) {
		boostrap.cache = cache
	}
}

func WithMailer(mailer *gomail.Dialer) BoostrapOption {
	return func(boostrap *boostrap) {
		boostrap.mailer = mailer
	}
}

func WithGoogleFormService(googleFromService *forms.Service) BoostrapOption {
	return func(boostrap *boostrap) {
		boostrap.googleForm = googleFromService
	}
}

func RunApp(options ...BoostrapOption) {
	boot := &boostrap{}
	for _, option := range options {
		option(boot)
	}
	boot.newPublicAPIProvider()
	boot.newTixAPIProvider()
}
