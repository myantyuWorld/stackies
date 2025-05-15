package model

import "stackies/backend/infra/repository/model"

type Experience struct {
	ID    int
	Title string
}

func NewExperience(title string) *model.Experience {
	return &model.Experience{
		Title: title,
	}
}

func (e *Experience) ConvertToEntity() *model.Experience {
	return &model.Experience{
		ID:    e.ID,
		Title: e.Title,
	}
}
