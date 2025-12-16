package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/wb-go/wbf/dbpg"
)

func TestUrlRepository_CreateShortUrl(t *testing.T) {
	type mockBehavior func(mock sqlmock.Sqlmock, originUrl string)

	type args struct {
		ctx       context.Context
		originUrl string
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		want         *domain.ShortenUrlModel
		wantErr      bool
		uniqueErr    bool
	}{
		{ // 1
			name: "Success - No Collisions",
			mockBehavior: func(mock sqlmock.Sqlmock, originUrl string) {
				mock.ExpectQuery("INSERT INTO urls").
					WithArgs(originUrl, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"original_url", "short_code"}).
						AddRow(originUrl, "abcdef1222"))
			},
			args: args{
				ctx:       context.Background(),
				originUrl: "https://google.com",
			},
			want: &domain.ShortenUrlModel{
				ID:          0,
				OriginalUrl: "https://google.com",
				ShortCode:   "abcdef1222",
				CreatedAt:   time.Time{},
			},
			wantErr:   false,
			uniqueErr: false,
		},
		{ // 2
			name: "Success - After 2 Collisions",
			mockBehavior: func(mock sqlmock.Sqlmock, originUrl string) {
				err := pq.Error{Code: "23505"}
				for range 2 {
					mock.ExpectQuery("INSERT INTO urls").
						WithArgs(originUrl, sqlmock.AnyArg()).
						WillReturnError(&err)
				}
				mock.ExpectQuery("INSERT INTO urls").
					WithArgs(originUrl, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"original_url", "short_code"}).
						AddRow(originUrl, "abcdef1222"))
			},
			args: args{
				ctx:       context.Background(),
				originUrl: "https://google.com",
			},
			want: &domain.ShortenUrlModel{
				ID:          0,
				OriginalUrl: "https://google.com",
				ShortCode:   "abcdef1222",
				CreatedAt:   time.Time{},
			},
			wantErr:   false,
			uniqueErr: false,
		},
		{ // 3
			name: "Failure - DB Error (Not Collision)",
			mockBehavior: func(mock sqlmock.Sqlmock, originUrl string) {
				mock.ExpectQuery("INSERT INTO urls").
					WithArgs(originUrl, sqlmock.AnyArg()).
					WillReturnError(errors.New("db connection error"))
			},
			args: args{
				ctx:       context.Background(),
				originUrl: "https://fail.com",
			},
			want:      nil,
			wantErr:   true,
			uniqueErr: false,
		},
		{ // 4
			name: "Failure - INSERT ERROR",
			mockBehavior: func(mock sqlmock.Sqlmock, originUrl string) {
				err := pq.Error{Code: "23505"}
				for range 5 {
					mock.ExpectQuery("INSERT INTO urls").
						WithArgs(originUrl, sqlmock.AnyArg()).
						WillReturnError(&err)
				}
			},
			args: args{
				ctx:       context.Background(),
				originUrl: "https://fail.com",
			},
			want:      nil,
			wantErr:   true,
			uniqueErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			tt.mockBehavior(mock, tt.args.originUrl)

			wrappedDB := &dbpg.DB{Master: mockDB}

			r := &UrlRepository{
				db: wrappedDB,
			}

			got, err := r.CreateShortUrl(tt.args.ctx, tt.args.originUrl)

			if tt.wantErr {
				assert.Error(t, err)

				if tt.uniqueErr {
					assert.True(t, errors.Is(err, domain.ErrNotUnique))
				}
			} else {
				assert.NoError(t, err)
				if got != nil {
					assert.Equal(t, tt.args.originUrl, got.OriginalUrl)
					assert.NotEmpty(t, got.ShortCode)
					assert.True(t, len(got.ShortCode) <= 10)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_generateRandomUrl(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1",
			args: args{3},
			want: "dsads",
		},
		{
			name: "Test 2",
			args: args{5},
			want: "dsads",
		},
		{
			name: "Test 3",
			args: args{0},
			want: "dsdsd",
		},
		{
			name: "Test 4",
			args: args{10},
			want: "dsadsddsds",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, len(tt.want), len(generateRandomUrl(tt.args.length)), "generateRandomUrl(%v)", tt.args.length)
		})
	}
}
