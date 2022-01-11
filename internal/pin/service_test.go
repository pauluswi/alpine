package pin

import (
	"context"
	"testing"
	"time"

	"github.com/pauluswi/alpine/internal/entity"
	"github.com/pauluswi/alpine/pkg/log"
	"github.com/stretchr/testify/assert"

	uuid "github.com/satori/go.uuid"
)

func Test_service_TokenCycle(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(&mockRepository{}, logger)

	ctx := context.Background()

	// create pin
	err := s.Create(ctx, entity.Input{CustomerID: "6281100099", Status: 1, Credential: "123456"})
	assert.Nil(t, err)

	pin, err := s.Get(ctx, "6281100099")
	assert.Nil(t, err)
	assert.NotEmpty(t, pin.Credential)

	// pin validation
	val, err := s.Validate(ctx, entity.Input{CustomerID: "6281100099", Status: 1, Credential: "123456"})
	assert.Nil(t, err)
	assert.Equal(t, true, val)

	// change pin
	err = s.Change(ctx, entity.Input{CustomerID: "6281100099", Status: 1, Credential: "123433"})
	assert.Nil(t, err)
}

type mockRepository struct {
	items []entity.Pin
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Pin, error) {
	var pin entity.Pin
	pin.ID = uuid.NewV4().String()
	pin.CustomerID = "6281100099"
	pin.Status = 1
	pin.Credential = "123456"
	pin.CreatedAt = time.Now()
	pin.UpdatedAt = time.Now()
	return pin, nil
}

func (m mockRepository) Save(ctx context.Context, pin entity.Pin) error {
	return nil
}

func (m mockRepository) Update(ctx context.Context, pin entity.Pin) error {
	return nil
}
