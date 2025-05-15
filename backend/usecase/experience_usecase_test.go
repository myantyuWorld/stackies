package usecase_test

import (
	"errors"
	"testing"

	"stackies/backend/domain/model"
	"stackies/backend/domain/repository/mock"
	"stackies/backend/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExperienceUsecase_Create(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		setupMock func(*mock.MockExperienceRepository)
		wantErr   error
	}{
		{
			name:  "正常系: 体験作成に成功",
			title: "テスト体験",
			setupMock: func(m *mock.MockExperienceRepository) {
				m.EXPECT().Create(model.Experience{Title: "テスト体験"}).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:  "異常系: repository.Createがエラーを返す",
			title: "テスト体験",
			setupMock: func(m *mock.MockExperienceRepository) {
				m.EXPECT().Create(model.Experience{Title: "テスト体験"}).Return(errors.New("DB error"))
			},
			wantErr: errors.New("DB error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockExperienceRepository(ctrl)
			tt.setupMock(mockRepo)

			uc := usecase.NewExperienceUsecase(mockRepo)
			err := uc.Create(tt.title)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExperienceUsecase_GetAll(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mock.MockExperienceRepository)
		want      []usecase.ExperienceDto
		wantErr   error
	}{
		{
			name: "正常系: 体験一覧を取得",
			setupMock: func(m *mock.MockExperienceRepository) {
				m.EXPECT().GetAll().Return([]model.Experience{
					{ID: 1, Title: "体験1"},
					{ID: 2, Title: "体験2"},
				}, nil)
			},
			want: []usecase.ExperienceDto{
				{ID: 1, Title: "体験1"},
				{ID: 2, Title: "体験2"},
			},
			wantErr: nil,
		},
		{
			name: "正常系: 体験が0件",
			setupMock: func(m *mock.MockExperienceRepository) {
				m.EXPECT().GetAll().Return([]model.Experience{}, nil)
			},
			want:    []usecase.ExperienceDto{},
			wantErr: nil,
		},
		{
			name: "異常系: repository.GetAllがエラーを返す",
			setupMock: func(m *mock.MockExperienceRepository) {
				m.EXPECT().GetAll().Return(nil, errors.New("DB error"))
			},
			want:    nil,
			wantErr: errors.New("DB error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockExperienceRepository(ctrl)
			tt.setupMock(mockRepo)

			uc := usecase.NewExperienceUsecase(mockRepo)
			got, err := uc.GetAll()

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
