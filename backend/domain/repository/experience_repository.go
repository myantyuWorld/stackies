//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock
package repository

import (
	"stackies/backend/infra/repository/model"
)

type ExperienceRepository interface {
	GetAll() ([]model.Experience, error)
	Create(experience model.Experience) error
}
