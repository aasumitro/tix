package rest

import (
	"context"
	"github.com/aasumitro/tix/internal/domain"
	"google.golang.org/api/forms/v1"
)

type IFormsService interface {
	Get(formID string) *forms.FormsGetCall
	Responses() *forms.FormsResponsesService
}

type googleServiceRepository struct {
	googleFormService IFormsService
}

func (repository *googleServiceRepository) GetEvent(
	ctx context.Context,
	formID string,
) (*forms.Form, error) {
	data, err := repository.googleFormService.
		Get(formID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repository *googleServiceRepository) GetResponses(
	ctx context.Context,
	formID string,
) (*forms.ListFormResponsesResponse, error) {
	data, err := repository.googleFormService.Responses().
		List(formID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func NewGoogleServiceRepository(
	googleFormService IFormsService,
) domain.IGoogleServiceRepository {
	return &googleServiceRepository{googleFormService}
}
