package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain/entity"
	"github.com/aasumitro/tix/internal/domain/response"
	"strings"
	"time"
)

func (service *tixService) FetchForms(
	ctx context.Context,
	formID string,
) (items []*response.GoogleFormQuestion, err error) {
	cacheKey := fmt.Sprintf("google-form-%s", formID)

	data, err := func() (any, error) {
		if valueCache, errCache := service.redisCache.Get(
			ctx, cacheKey,
		).Result(); errCache == nil {
			return valueCache, nil
		}

		data, err := service.googleFormService.Forms.
			Get(formID).Context(ctx).Do()
		if err != nil {
			return nil, err
		}

		for _, q := range data.Items {
			items = append(items, &response.GoogleFormQuestion{
				ID:    q.QuestionItem.Question.QuestionId,
				Title: q.Title,
			})
		}

		if jsonData, err := json.Marshal(items); err == nil {
			service.redisCache.Set(ctx, cacheKey, jsonData, common.GoogleFormCacheTimeDuration)
		}

		return data, nil
	}()
	if err != nil {
		return nil, err
	}

	if data, ok := data.([]*response.GoogleFormQuestion); ok {
		items = data
	}

	if data, ok := data.(string); ok {
		var form []*response.GoogleFormQuestion
		if err := json.Unmarshal([]byte(data), &form); err != nil {
			return nil, err
		}
		items = form
	}

	return items, nil
}

func (service *tixService) FetchResponds(
	ctx context.Context,
	formID string,
) (items []*response.GoogleFormRespond, err error) {
	questions, err := service.FetchForms(ctx, formID)
	if err != nil {
		return nil, err
	}

	data, err := service.googleFormService.Forms.Responses.
		List(formID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	for _, resp := range data.Responses {
		var answer response.GoogleFormRespondAnswer
		for _, responseAnswer := range resp.Answers {
			for _, q := range questions {
				if q.ID == responseAnswer.QuestionId {
					key := strings.Replace(strings.ToLower(q.Title), " ", "_", 1)
					if responseAnswer.TextAnswers != nil {
						switch key {
						case "pekerjaan":
							answer.Job = responseAnswer.TextAnswers.Answers[0].Value
						case "tanggal_lahir":
							answer.DoB = responseAnswer.TextAnswers.Answers[0].Value
						case "email":
							answer.Email = responseAnswer.TextAnswers.Answers[0].Value
						case "nama":
							answer.Name = responseAnswer.TextAnswers.Answers[0].Value
						case "nomor_telepon":
							answer.Phone = responseAnswer.TextAnswers.Answers[0].Value
						}
					}
					if responseAnswer.FileUploadAnswers != nil && key == "bukti_transfer" {
						answer.PoP = fmt.Sprintf(
							"https://drive.google.com/uc?export=view&id=%s",
							responseAnswer.FileUploadAnswers.Answers[0].FileId)
					}
					break
				}
			}
		}

		items = append(items, &response.GoogleFormRespond{
			RespondID:         resp.ResponseId,
			LastSubmittedTime: resp.LastSubmittedTime,
			CreateTime:        resp.CreateTime,
			Answer:            &answer,
		})
	}

	return items, nil
}

func (service *tixService) SyncRespondData(
	ctx context.Context,
	formID string,
) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	event, err := service.postgreSQLRepository.GetEventByGoogleFormID(ctx, formID)
	if err != nil {
		return err
	}

	respondents, err := service.FetchResponds(ctx, formID)
	if err != nil {
		return err
	}

	var newParticipant []*entity.Participant
	for _, respond := range respondents {
		data, err := service.postgreSQLRepository.GetParticipantByEmailAndEventID(
			ctx, respond.Answer.Email, event.ID)
		if (err == nil || errors.Is(err, sql.ErrNoRows)) && data == nil {
			newParticipant = append(newParticipant, &entity.Participant{
				EventID: event.ID,
				Name:    respond.Answer.Name,
				Email:   respond.Answer.Email,
				Phone:   respond.Answer.Phone,
				Job:     respond.Answer.Job,
				PoP:     respond.Answer.PoP,
				DoB:     respond.Answer.DoB,
			})
		}
	}

	cacheKey := fmt.Sprintf("participants-%s", formID)
	service.redisCache.Del(ctx, cacheKey)
	return service.postgreSQLRepository.InsertManyParticipants(
		ctx, newParticipant, time.Now().Unix())
}
