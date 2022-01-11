package pin

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/pauluswi/alpine/internal/entity"
	"github.com/pauluswi/alpine/pkg/dbcontext"
	"github.com/pauluswi/alpine/pkg/log"
)

// Repository encapsulates the logic to access pin from the data source.
type Repository interface {
	// Get returns the customer's pin information with the specified token string.
	Get(ctx context.Context, customerid string) (entity.Pin, error)
	// Save will store a pin information into data source.
	Save(ctx context.Context, pin entity.Pin) error
	// Update will store an updated pin information into data source.
	Update(ctx context.Context, pin entity.Pin) error
	// // Delete will remove a pin information from data source.
	// Delete(ctx context.Context, customerid string) error
}

// repository persists pin in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new pin repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get returns the customer's token information with the specified token string.
func (r repository) Get(ctx context.Context, customerid string) (entity.Pin, error) {
	var pin entity.Pin
	err := r.db.With(ctx).Select("id", "customer_id", "status", "credential", "metadata", "created_at", "updated_at").
		From("pinstore").
		Where(dbx.HashExp{"customer_id": customerid}).
		One(&pin)
	return pin, err
}

// Save will store a token information into data source.
func (r repository) Save(ctx context.Context, pin entity.Pin) error {
	_, err := r.db.With(ctx).Insert("pinstore", dbx.Params{
		"id":          pin.ID,
		"customer_id": pin.CustomerID,
		"status":      pin.Status,
		"credential":  pin.Credential,
		"metadata":    "{}",
		"created_at":  pin.CreatedAt,
		"updated_at":  pin.UpdatedAt,
	}).Execute()
	return err
}

// Update will store an updated token information into data source.
func (r repository) Update(ctx context.Context, pin entity.Pin) error {
	_, err := r.db.With(ctx).Update("pinstore", dbx.Params{
		"status":     pin.Status,
		"credential": pin.Credential,
		"metadata":   pin.Metadata,
		"updated_at": pin.UpdatedAt,
	}, dbx.HashExp{"id": pin.ID}).Execute()
	return err
}

// // Update will store an updated token information into data source.
// func (r repository) Delete(ctx context.Context, id string) error {
// 	_, err := r.db.With(ctx).Delete("pinstore", dbx.HashExp{"ID": id}).Execute()
// 	return err
// }
