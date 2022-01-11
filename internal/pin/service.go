package pin

import (
	"context"
	"fmt"
	"time"

	"github.com/pauluswi/alpine/internal/entity"
	"github.com/pauluswi/alpine/pkg/log"
	"github.com/pauluswi/alpine/pkg/validator"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	Get(ctx context.Context, customer_id string) (entity.Pin, error)
	Validate(ctx context.Context, req entity.Input) (valid bool, err error)
	Create(ctx context.Context, req entity.Input) (err error)
	Change(ctx context.Context, req entity.Input) (err error)
}

// // pin represents the data about an payment token.
// type pin struct {
// 	entity.Pin
// }

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new payment token service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// --- list of error and constants
var (
	ErrValidation  = fmt.Errorf("validation error")
	ErrCreatePin   = fmt.Errorf("pin creation error")
	ErrDBPersist   = fmt.Errorf("persist to database error")
	ErrPinNotFound = fmt.Errorf("pin not found in database")
)

// Get returns a pin belongs to a customer
func (s service) Get(ctx context.Context, id string) (out entity.Pin, err error) {
	pin, err := s.repo.Get(ctx, id)
	if err != nil {
		return pin, err
	}
	return pin, nil
}

// Validate will check pin information
func (s service) Validate(ctx context.Context, req entity.Input) (valid bool, err error) {
	defer func() {
		if err != nil {
			s.logger.Error(ctx, err.Error())
		}
	}()

	err = validator.ValidateWithOpts(req, validator.Opts{Mode: validator.ModeVerbose})
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrValidation, err)
		return false, err
	}

	pin, err := s.repo.Get(ctx, req.CustomerID)
	if err != nil {
		return false, err
	}

	// validate pin
	// First generate random 16 byte salt
	var salt = GenerateRandomSalt(saltSize)
	// Hash password using the salt
	var hashedPassword = HashPassword(req.Credential, salt)
	valid = DoPasswordsMatch(hashedPassword, pin.Credential, salt)
	if !valid {
		return false, err
	}
	return true, nil
}

// Save creates a pin information
func (s service) Create(ctx context.Context, req entity.Input) (err error) {
	defer func() {
		if err != nil {
			s.logger.Error(ctx, err.Error())
		}
	}()

	err = validator.ValidateWithOpts(req, validator.Opts{Mode: validator.ModeVerbose})
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrValidation, err)
		return err
	}
	// ** I use a very simple algorithm to encrypt pin / credential
	// ** in real world the algorithm must be more details and secure
	// First generate random 16 byte salt
	var salt = GenerateRandomSalt(saltSize)
	// Hash password using the salt
	var hashedPassword = HashPassword(req.Credential, salt)

	// build new pin
	pin := entity.NewPin()
	pin.CustomerID = req.CustomerID
	pin.Status = req.Status
	pin.Credential = hashedPassword
	pin.CreatedAt = time.Now()
	pin.UpdatedAt = time.Now()

	err = validator.ValidateWithOpts(pin, validator.Opts{Mode: validator.ModeVerbose})
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrValidation, err)
		return err
	}

	err = s.repo.Save(ctx, *pin)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrDBPersist, err)
		return err
	}

	return err
}

// Save creates a pin information
func (s service) Change(ctx context.Context, req entity.Input) (err error) {
	defer func() {
		if err != nil {
			s.logger.Error(ctx, err.Error())
		}
	}()

	err = validator.ValidateWithOpts(req, validator.Opts{Mode: validator.ModeVerbose})
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrValidation, err)
		return err
	}

	currPin, err := s.repo.Get(ctx, req.CustomerID)
	if err != nil {
		return err
	}

	// build new pin
	var pin entity.Pin
	pin.ID = currPin.ID
	pin.Status = req.Status
	pin.Credential = req.Credential
	pin.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, pin)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrDBPersist, err)
		return err
	}

	return err
}
