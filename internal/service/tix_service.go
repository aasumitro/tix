package service

import (
	"github.com/aasumitro/tix/internal/domain"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/forms/v1"
	"gopkg.in/gomail.v2"
	"sync"
)

type tixService struct {
	mu                   sync.Mutex
	googleFormService    *forms.Service
	redisCache           *redis.Client
	authRESTRepository   domain.IAuthRESTRepository
	postgreSQLRepository domain.IPostgreSQLRepository
	mailer               *gomail.Dialer
}

type TixOptions func(*tixService)

func WithGoogleFormService(
	googleFormService *forms.Service,
) TixOptions {
	return func(service *tixService) {
		service.googleFormService = googleFormService
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
