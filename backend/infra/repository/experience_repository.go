package repository

import (
	"stackies/backend/domain/repository"
	"stackies/backend/infra/repository/model"

	"gorm.io/gorm"
)

type experienceRepository struct {
	db *gorm.DB
}

// GetAll implements repository.ExperienceRepository.
func (e *experienceRepository) GetAll() ([]model.Experience, error) {
	var experiences []model.Experience
	if err := e.db.Find(&experiences).Error; err != nil {
		return nil, err
	}
	return experiences, nil
}

// Create implements repository.ExperienceRepository.
func (e *experienceRepository) Create(experience model.Experience) error {
	if err := e.db.Create(&experience).Error; err != nil {
		return err
	}
	return nil
}

func NewExperienceRepository(db *gorm.DB) repository.ExperienceRepository {
	return &experienceRepository{
		db: db,
	}
}
