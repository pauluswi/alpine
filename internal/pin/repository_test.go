package pin

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/pauluswi/alpine/internal/entity"
	"github.com/pauluswi/alpine/internal/test"
	"github.com/pauluswi/alpine/pkg/log"
	"github.com/stretchr/testify/assert"

	uuid "github.com/satori/go.uuid"
)

func TestRepository(t *testing.T) {
	logger, _ := log.NewForTest()
	db := test.DB(t)
	test.ResetTables(t, db, "pinstore")
	repo := NewRepository(db, logger)

	ctx := context.Background()

	// create
	err := repo.Save(ctx, entity.Pin{
		ID:         uuid.NewV4().String(),
		CustomerID: "62811000001",
		Status:     1,
		Credential: "123456",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	assert.Nil(t, err)

	// get
	pin, err := repo.Get(ctx, "62811000001")
	assert.Nil(t, err)
	assert.Equal(t, "123456", pin.Credential)
	_, err = repo.Get(ctx, "999990")
	assert.Equal(t, sql.ErrNoRows, err)

	// // update
	// err = repo.Update(ctx, entity.Pin{
	// 	ID:         pin.ID,
	// 	CustomerID: pin.CustomerID,
	// 	Status:     2,
	// 	Metadata:   entity.Metadata{time.Now().UTC()},
	// 	UpdatedAt:  time.Now(),
	// })
	// assert.Nil(t, err)

	// // get after update
	// updatedpin, err := repo.Get(ctx, "62811000001")
	// assert.Nil(t, err)
	// assert.Equal(t, 2, updatedpin.Status)

	// // update
	// err = repo.Update(ctx, entity.Pin{
	// 	ID:         pin.ID,
	// 	CustomerID: pin.CustomerID,
	// 	Status:     2,
	// 	Metadata:   entity.Metadata{time.Now().UTC()},
	// 	UpdatedAt:  time.Now(),
	// })
	// assert.Nil(t, err)

	// // get after delete
	// _, err = repo.Get(ctx, pin.ID)
	// assert.Equal(t, sql.ErrNoRows, err)
}
