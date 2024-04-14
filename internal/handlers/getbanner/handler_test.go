package getbanner

import (
	"AvitoTech/internal/handlers/getbanner/mocks"
	"AvitoTech/internal/models"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

type argsType struct {
	tagID     int64
	limit     int
	offset    int
	featureID int64
}

// TODO: докинуть тест кейсы
func TestHandler_GetBanners(t *testing.T) {
	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2020-09-30 10:00:00")
	updatedAt, _ := time.Parse("2006-01-02 15:04:05", "2020-09-30 14:00:00")

	tests := []struct {
		name               string
		args               argsType
		prepareDBInt       func() DBInt
		prepareUserChecker func() UserChecker
		expectedResult     []models.Banner
		expectedCode       int
	}{
		{
			name: "successful get by tag id",
			args: argsType{tagID: 1, limit: 1, offset: 0, featureID: -1},
			prepareDBInt: func() DBInt {
				mockDB := mocks.NewDBInt(t)

				mockDB.EXPECT().GetBannersByTagID(uint64(1), 1, 0).
					Return([]models.Banner{
						{
							IsActive:  true,
							BannerID:  1,
							TagIDs:    []uint64{1, 2, 3},
							FeatureID: 3,
							Content: pgtype.JSONB{
								Bytes:  []byte(`{"url":"some_url","text":"Проверка доступа","title":"Post Request"}`),
								Status: pgtype.Present,
							},
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						},
					}, nil)

				return mockDB
			},
			prepareUserChecker: func() UserChecker {
				mockUserChecker := mocks.NewUserChecker(t)

				mockUserChecker.EXPECT().IsAdmin(mock.Anything).Return(true, nil)

				return mockUserChecker
			},
			expectedResult: []models.Banner{
				{
					IsActive:  true,
					BannerID:  1,
					TagIDs:    []uint64{1, 2, 3},
					FeatureID: 3,
					Content: pgtype.JSONB{
						Bytes:  []byte(`{"url":"some_url","text":"Проверка доступа","title":"Post Request"}`),
						Status: pgtype.Present,
					},
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := NewHandler(tt.prepareDBInt(), tt.prepareUserChecker())

			app.Get("/api/banner", handler.GetBanners)

			req := httptest.NewRequest(http.MethodGet, buildPath(tt.args), nil)

			resp, _ := app.Test(req, -1)

			var banners []models.Banner

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			json.Unmarshal(body, &banners)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			assert.Equal(t, tt.expectedResult, banners)
		})
	}
}

func buildPath(args argsType) string {
	var builder strings.Builder

	builder.WriteString("/api/banner?")

	if args.featureID != -1 {
		builder.WriteString("feature_id=")
		builder.WriteString(strconv.FormatInt(args.featureID, 10))
		builder.WriteString("&")
	}
	if args.limit != -1 {
		builder.WriteString("limit=")
		builder.WriteString(strconv.Itoa(args.limit))
		builder.WriteString("&")
	}
	if args.offset != -1 {
		builder.WriteString("offset=")
		builder.WriteString(strconv.Itoa(args.offset))
		builder.WriteString("&")
	}
	if args.tagID != -1 {
		builder.WriteString("tag_id=")
		builder.WriteString(strconv.FormatInt(args.tagID, 10))
		builder.WriteString("&")
	}

	return builder.String()
}
