package services

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"main.go/schemas"
)

func (r *Service) ProcessWithML(c *fiber.Ctx, u string) error {
	agent := fiber.Get(u)
	_, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	var attrs *schemas.Attribute
	err := json.Unmarshal(body, &attrs)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal attrs")
	}

	return nil
}
