package service

import (
	"context"
	"testing"

	"github.com/stannisl/url-shortener/internal/domain"
	mocks "github.com/stannisl/url-shortener/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUrlService_GetOriginUrl(t *testing.T) {
	type mockBehavior func(r *mocks.MockUrlRepository, shortCode string)

	tests := []struct {
		name          string
		shortCode     string
		mockSetup     mockBehavior
		expectedError bool
	}{
		{
			name:      "Success",
			shortCode: "validCode",
			mockSetup: func(r *mocks.MockUrlRepository, shortCode string) {
				r.On("GetOriginUrl", mock.Anything, shortCode).
					Return(&domain.ShortenUrlModel{
						ID:          1,
						OriginalUrl: "http://go.dev/",
						ShortCode:   shortCode,
					}, nil)
			},
			expectedError: false,
		},
		{
			name:      "Success",
			shortCode: "validCode",
			mockSetup: func(r *mocks.MockUrlRepository, shortCode string) {
				r.On("GetOriginUrl", mock.Anything, shortCode).
					Return(nil, domain.ErrModelNotFound)
			},
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUrlRepository)
			tt.mockSetup(mockRepo, tt.shortCode)
			service := NewUrlService(mockRepo)

			_, err := service.GetOriginUrl(context.Background(), tt.shortCode)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
