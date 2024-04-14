package get_user_banner

import (
	"AvitoTech/integration_tests"
	"github.com/steinfletcher/apitest"
	"net/http"
	"testing"
)

func TestGetUserBanner_NotAdminSuccess(t *testing.T) {
	apitest.New().
		EnableNetworking().
		Get("http://localhost:8080/api/user_banner").
		Header("Authorization", integrationtests.GenerateToken(2)).
		QueryParams(map[string]string{
			"feature_id": "2",
			"tag_id":     "1",
		}).
		Expect(t).
		Status(http.StatusOK).
		BodyFromFile("not_admin_success.json").
		End()
}

func TestGetUserBanner_AdminSuccess(t *testing.T) {
	apitest.New().
		EnableNetworking().
		Get("http://localhost:8080/api/user_banner").
		Header("Authorization", integrationtests.GenerateToken(1)).
		QueryParams(map[string]string{
			"feature_id": "3",
			"tag_id":     "1",
		}).
		Expect(t).
		Status(http.StatusOK).
		BodyFromFile("admin_success.json").
		End()
}

func TestGetUserBanner_NotAdminFailure(t *testing.T) {
	apitest.New().
		EnableNetworking().
		Get("http://localhost:8080/api/user_banner").
		Header("Authorization", integrationtests.GenerateToken(2)).
		QueryParams(map[string]string{
			"feature_id": "3",
			"tag_id":     "1",
		}).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
