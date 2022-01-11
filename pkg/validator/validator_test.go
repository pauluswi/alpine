package validator_test

import (
	"errors"
	"fmt"
	"time"

	ovalidator "github.com/go-playground/validator/v10"

	"github.com/pauluswi/alpine/pkg/validator"
)

func Example_basic() {
	// UserForm struct
	type UserForm struct {
		Name     string    `validate:"required,gte=7"`
		Email    string    `validate:"email"`
		Age      int       `validate:"required,numeric,min=1,max=99"`
		CreateAt int       `validate:"gte=1"`
		Safe     int       `validate:"-"`
		UpdateAt time.Time `validate:"required"`
		Code     string    `validate:"required"`
	}

	err := validator.Validate(UserForm{})
	fmt.Println(err.Error())

	// Output:
	// UserForm.Name must be required; UserForm.Email must be email; UserForm.Age must be required, actual value is 0; UserForm.CreateAt must be gte{1}, actual value is 0; UserForm.UpdateAt must be required, actual value is 0001-01-01 00:00:00 +0000 UTC; UserForm.Code must be required
}

func Example_optional_mode_compact() {
	// UserForm struct
	type UserForm struct {
		Name     string    `validate:"required,gte=7"`
		Email    string    `validate:"email"`
		Age      int       `validate:"required,numeric,min=1,max=99"`
		CreateAt int       `validate:"gte=1"`
		Safe     int       `validate:"-"`
		UpdateAt time.Time `validate:"required"`
		Code     string    `validate:"required"`
	}

	err := validator.ValidateWithOpts(UserForm{}, validator.Opts{Mode: validator.ModeCompact})
	fmt.Println(err.Error())

	// Output:
	// invalid fields: Name, Email, Age, CreateAt, UpdateAt, Code
}

func Example_optional_mode_verbose() {
	// UserForm struct
	type UserForm struct {
		Name     string    `validate:"required,gte=7"`
		Email    string    `validate:"email"`
		Age      int       `validate:"required,numeric,min=1,max=99"`
		CreateAt int       `validate:"gte=1"`
		Safe     int       `validate:"-"`
		UpdateAt time.Time `validate:"required"`
		Code     string    `validate:"required"`
	}

	err := validator.ValidateWithOpts(UserForm{}, validator.Opts{Mode: validator.ModeVerbose})
	fmt.Println(err.Error())

	// (SKIP) Output:
	// invalid fields at /Users/avreghlybarra/Developer/linkaja/neogodam/pkg/validator/example_test.go:61: UserForm.Name must be required; UserForm.Email must be email; UserForm.Age must be required, actual value is 0; UserForm.CreateAt must be gte{1}, actual value is 0; UserForm.UpdateAt must be required, actual value is 0001-01-01 00:00:00 +0000 UTC; UserForm.Code must be required
}

func Example_unwrap_to_original_error() {
	// UserForm struct
	type UserForm struct {
		Name     string    `validate:"required,gte=7"`
		Email    string    `validate:"email"`
		Age      int       `validate:"required,numeric,min=1,max=99"`
		CreateAt int       `validate:"gte=1"`
		Safe     int       `validate:"-"`
		UpdateAt time.Time `validate:"required"`
		Code     string    `validate:"required"`
	}

	err := validator.ValidateWithOpts(UserForm{}, validator.Opts{Mode: validator.ModeVerbose})

	oerr := ovalidator.ValidationErrors{}
	_ = errors.As(err, &oerr)
	fmt.Println(oerr.Error())

	// Output:
	// Key: 'UserForm.Name' Error:Field validation for 'Name' failed on the 'required' tag
	// Key: 'UserForm.Email' Error:Field validation for 'Email' failed on the 'email' tag
	// Key: 'UserForm.Age' Error:Field validation for 'Age' failed on the 'required' tag
	// Key: 'UserForm.CreateAt' Error:Field validation for 'CreateAt' failed on the 'gte' tag
	// Key: 'UserForm.UpdateAt' Error:Field validation for 'UpdateAt' failed on the 'required' tag
	// Key: 'UserForm.Code' Error:Field validation for 'Code' failed on the 'required' tag
}

func Example_change_validator() {
	// UserForm struct
	type UserForm struct {
		Name     string    `validate:"required,gte=7"`
		Email    string    `validate:"email"`
		Age      int       `validate:"required,numeric,min=1,max=99"`
		CreateAt int       `validate:"gte=1"`
		Safe     int       `validate:"-"`
		UpdateAt time.Time `validate:"required"`
		Code     string    `validate:"required"`
	}

	custval := ovalidator.New()

	validator.SetGlobal(custval)
	err := validator.Validate(UserForm{})
	fmt.Println(err)

	// Output:
	// UserForm.Name must be required; UserForm.Email must be email; UserForm.Age must be required, actual value is 0; UserForm.CreateAt must be gte{1}, actual value is 0; UserForm.UpdateAt must be required, actual value is 0001-01-01 00:00:00 +0000 UTC; UserForm.Code must be required
}
