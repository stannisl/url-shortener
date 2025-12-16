package service

import (
	"context"
	"fmt"
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

func TestUrlService_CreateShortenUrl(t *testing.T) {
	type mockBehavior func(r *mocks.MockUrlRepository, url string)
	type args struct {
		ctx       context.Context
		originUrl string
	}
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		want         *domain.ShortenUrlModel
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			mockBehavior: func(r *mocks.MockUrlRepository, url string) {
				r.On("CreateShortUrl", mock.Anything, url).
					Return(&domain.ShortenUrlModel{
						ID:          1,
						OriginalUrl: url,
						ShortCode:   "abcdefg",
					}, nil)
			},
			args: args{
				ctx:       context.Background(),
				originUrl: "https://google.com",
			},
			want: &domain.ShortenUrlModel{
				ID:          1,
				OriginalUrl: "https://google.com",
				ShortCode:   "abcdefg",
			},
			wantErr: assert.NoError,
		},
		{
			name: "Fail",
			mockBehavior: func(r *mocks.MockUrlRepository, url string) {
				r.On("CreateShortUrl", mock.Anything, url).
					Return(nil, domain.ErrNotUnique)
			},
			args: args{
				ctx:       context.Background(),
				originUrl: "https://google.com",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUrlRepo := new(mocks.MockUrlRepository)
			tt.mockBehavior(mockUrlRepo, tt.args.originUrl)
			us := &UrlService{
				urlRepository: mockUrlRepo,
			}
			got, err := us.CreateShortenUrl(tt.args.ctx, tt.args.originUrl)
			if !tt.wantErr(t, err, fmt.Sprintf("CreateShortenUrl(%v, %v)", tt.args.ctx, tt.args.originUrl)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CreateShortenUrl(%v, %v)", tt.args.ctx, tt.args.originUrl)
		})
	}
}
