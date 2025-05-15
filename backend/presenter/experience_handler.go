package presenter

import (
	"net/http"
	"stackies/backend/usecase"

	"github.com/labstack/echo/v4"
)

type experienceHandler struct {
	experienceUsecase usecase.ExperienceUsecase
}

type CreateExperienceRequest struct {
	Title string `json:"title"`
}

type ExperienceResponse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (e *ExperienceResponse) ConvertToDto(experience usecase.ExperienceDto) {
	e.ID = experience.ID
	e.Title = experience.Title
}

// Create implements ExperienceHandler.
func (e *experienceHandler) Create(c echo.Context) error {
	var request CreateExperienceRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err := e.experienceUsecase.Create(request.Title)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, request)
}

// GetAll implements ExperienceHandler.
func (e *experienceHandler) GetAll(c echo.Context) error {
	experiences, err := e.experienceUsecase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := make([]ExperienceResponse, len(experiences))
	for i, experience := range experiences {
		response[i].ConvertToDto(experience)
	}
	return c.JSON(http.StatusOK, response)
}

type ExperienceHandler interface {
	Create(c echo.Context) error
	GetAll(c echo.Context) error
}

func NewExperienceHandler(experienceUsecase usecase.ExperienceUsecase) ExperienceHandler {
	return &experienceHandler{experienceUsecase: experienceUsecase}
}
