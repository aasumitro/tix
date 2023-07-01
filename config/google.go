package config

import (
	"context"
	"fmt"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
	"log"
)

func (cfg *Config) InitGoogleFormConn() {
	log.Println("Trying to init google form service conn . . . .")
	googleFormsSingleton.Do(func() {
		ctx := context.Background()
		formsService, err := forms.NewService(ctx, option.WithCredentialsFile(cfg.GoogleCredentialPath))
		if err != nil {
			panic(fmt.Sprintf("GOOGLE_FORM_ERROR, error create new service: %s", err.Error()))
		}
		GoogleForm = formsService
		log.Println("Google form service connection created . . . .")
	})
}

// FormsServiceWrapper is a custom implementation of IFormsService that wraps *forms.FormsService.
type FormsServiceWrapper struct {
	Service *forms.FormsService
}

// Get implements the Get method of IFormsService.
func (w *FormsServiceWrapper) Get(formID string) *forms.FormsGetCall {
	return w.Service.Get(formID)
}

// Responses implements the Responses method of IFormsService.
func (w *FormsServiceWrapper) Responses() *forms.FormsResponsesService {
	return w.Service.Responses
}
