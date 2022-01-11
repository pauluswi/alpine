package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

// pin struct is defined here
type Pin struct {
	ID         string    `db:"id" validate:"required,uuid"`
	CustomerID string    `db:"customer_id" validate:"required"`
	Status     int       `db:"status" validate:"required"`
	Credential string    `db:"credential" validate:"required"`
	CreatedAt  time.Time `db:"created_at" validate:"required"`
	UpdatedAt  time.Time `db:"updated_at" validate:"required"`
	Metadata   Metadata  `db:"metadata" validate:"-"`
}

type Metadata struct {
	ValidatedAt time.Time `json:"validated_at"`
}

func NewPin() *Pin {
	return &Pin{
		ID: uuid.NewV4().String(),
	}
}

// Value returns m as a value.  This does a validating unmarshal into another
// RawMessage.  If m is invalid json, it returns an error.
func (m Metadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan stores the src in *m.  No validation is done.
func (m *Metadata) Scan(src interface{}) error {
	return json.Unmarshal([]byte(fmt.Sprintf("%s", src)), &m)
}

// --- I/O for Service function

// Input
type Input struct {
	CustomerID string `json:"customer_id" validate:"required,numeric,startswith=62,min=10"`
	Status     int    `json:"status" validate:"required"`
	Credential string `json:"credential" validate:"required,numeric,min=6,max=6"`
}

// OutGenerate .
type Output struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Status     int    `json:"status"`
	Credential string `json:"credential"`
}

// --- list of Postgres Error Code
// --- complete error codes see: https://www.postgresql.org/docs/13/errcodes-appendix.html

const (
	PGErrCodeUniqueViolation = "23505"
)

// --- list of constraint name, for list constraint see migrations/postgres sql schema

const (
	PGConstraintUniqueTokenAndTokenDate = "idx_unq_tokens_token_token_date"
)

var (
	ErrInputValidation = fmt.Errorf("validation input error")
)
