package handlers

import (
	"context"
	graphqlschema "line/src/graphql"

	"github.com/gofiber/fiber/v2"
	gql "github.com/graphql-go/graphql"
)

func GraphQLHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody struct {
			Query string `json:"query"`
		}
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		token := c.Get("Authorization")
		c.Locals("Authorization", token)

		result := gql.Do(gql.Params{
			Schema:        graphqlschema.Schema,
			RequestString: requestBody.Query,
			Context:       context.WithValue(c.Context(), "fiberCtx", c),
		})

		if len(result.Errors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": result.Errors,
			})
		}

		return c.JSON(result)
	}
}
