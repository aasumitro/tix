package service

import (
	"github.com/aasumitro/tix/internal/domain"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
	"sync"
)

type tixService struct {
	mu                      sync.Mutex
	redisCache              *redis.Client
	googleServiceRepository domain.IGoogleServiceRepository
	authRESTRepository      domain.IAuthRESTRepository
	postgreSQLRepository    domain.IPostgreSQLRepository
	mailer                  *gomail.Dialer
}

type TixOptions func(*tixService)

func WithGoogleServiceRepository(
	googleFormService domain.IGoogleServiceRepository,
) TixOptions {
	return func(service *tixService) {
		service.googleServiceRepository = googleFormService
	}
}

func WithRedisCache(
	redisCache *redis.Client,
) TixOptions {
	return func(service *tixService) {
		service.redisCache = redisCache
	}
}

func WithAuthRESTRepository(
	authRESTRepository domain.IAuthRESTRepository,
) TixOptions {
	return func(service *tixService) {
		service.authRESTRepository = authRESTRepository
	}
}

func WithPostgreSQLRepository(
	postgreSQLRepository domain.IPostgreSQLRepository,
) TixOptions {
	return func(service *tixService) {
		service.postgreSQLRepository = postgreSQLRepository
	}
}

func WithMailer(mailer *gomail.Dialer) TixOptions {
	return func(service *tixService) {
		service.mailer = mailer
	}
}

func NewTixService(
	options ...TixOptions,
) domain.ITixService {
	service := &tixService{}
	for _, option := range options {
		option(service)
	}
	return service
}
