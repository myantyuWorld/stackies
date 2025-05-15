//go:generate mockgen -source=experience_usecase.go -destination=mock/mock_$GOFILE -package=mock
package usecase

import (
	"stackies/backend/domain/model"
	"stackies/backend/domain/repository"
)

type ExperienceDto struct {
	ID    int
	Title string
}

type experienceUsecase struct {
	experienceRepository repository.ExperienceRepository
}

// Create implements ExperienceUsecase.
func (e *experienceUsecase) Create(title string) error {
	return e.experienceRepository.Create(*model.NewExperience(title))
}

// GetAll implements ExperienceUsecase.
func (e *experienceUsecase) GetAll() ([]ExperienceDto, error) {
	experiences, err := e.experienceRepository.GetAll()
	if err != nil {
		return nil, err
	}
	experienceDtos := make([]ExperienceDto, len(experiences))
	for i, experience := range experiences {
		experienceDtos[i] = ExperienceDto{
			ID:    experience.ID,
			Title: experience.Title,
		}
	}
	return experienceDtos, nil
}

type ExperienceUsecase interface {
	Create(title string) error
	GetAll() ([]ExperienceDto, error)
}

func NewExperienceUsecase(experienceRepository repository.ExperienceRepository) ExperienceUsecase {
	return &experienceUsecase{experienceRepository: experienceRepository}
}
