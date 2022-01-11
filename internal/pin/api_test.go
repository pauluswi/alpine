package pin

import (
	"net/http"
	"testing"
	"time"

	"github.com/pauluswi/alpine/internal/auth"
	"github.com/pauluswi/alpine/internal/entity"
	"github.com/pauluswi/alpine/internal/test"
	"github.com/pauluswi/alpine/pkg/log"
	uuid "github.com/satori/go.uuid"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	repo := &mockRepository{items: []entity.Pin{
		{uuid.NewV4().String(), "6281100099", 1, "123456", time.Now(), time.Now(), entity.Metadata{time.Now().UTC()}},
	}}
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{"create ok", "POST", "/create", `{"customer_id": "6281111111","status":1,"credential":"123456"}`, header, http.StatusCreated, ""},
		{"create auth error", "POST", "/create", `{"customer_id":"62811000991","status":1,"credential":"123456"}`, nil, http.StatusUnauthorized, ""},
		{"create input error", "POST", "/create", `"customer_id":"62811000991"}`, header, http.StatusBadRequest, ""},
		{"get ok", "GET", "/get/6281100099", "", header, http.StatusOK, `*"Credential":"123456"`},
		{"get unknown", "GET", "/get/paytokens/62811000991", "", header, http.StatusNotFound, ""},
		{"validate ok", "POST", "/validate", `{"customer_id": "6281111111","status":1,"credential":"123456"}`, header, http.StatusOK, "true"},
		{"change pin ok", "PATCH", "/change", `{"customer_id": "6281111111","status":1,"credential":"222222"}`, header, http.StatusOK, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
