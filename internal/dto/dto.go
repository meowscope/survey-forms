package dto

import (
	"example.com/m/internal/models"
)

type RequestCreateSurvey struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Questions_list []models.Question `json:"questions_list"`
}

func InternalCreateSurvey(survey models.Survey) RequestCreateSurvey {
	return RequestCreateSurvey{
		Name:           survey.Name,
		Description:    survey.Description,
		Questions_list: survey.Questions_list,
	}
}
