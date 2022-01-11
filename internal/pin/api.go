package pin

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/pauluswi/alpine/internal/entity"
	"github.com/pauluswi/alpine/internal/errors"
	"github.com/pauluswi/alpine/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}
	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Get("/get/<id>", res.get)
	r.Post("/validate", res.validate)
	r.Post("/create", res.create)
	r.Patch("/change", res.change)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	pin, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(pin)
}

func (r resource) validate(c *routing.Context) error {
	var input entity.Input
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("Bad Request")
	}

	valid, err := r.service.Validate(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(valid, http.StatusOK)
}

func (r resource) create(c *routing.Context) error {
	var input entity.Input
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("Bad Request")
	}

	err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(nil, http.StatusCreated)
}

func (r resource) change(c *routing.Context) error {
	var input entity.Input
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("Bad Request")
	}
	err := r.service.Change(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(nil, http.StatusOK)
}
