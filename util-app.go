package api_common

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"strings"
)

func ConfigureApp(allowedDomains []string) *fiber.App {
	app := fiber.New()

	// use default cors config
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(allowedDomains, ","),
	}))

	// generate random request id for each call
	app.Use(requestid.New(requestid.Config{
		Header: HTTP_HEADER_REQUEST_ID,
		Generator: func() string {
			return RandomGenerateUuidWithLength(false, 24)
		},
		ContextKey: CTX_REQUESTID,
	}))

	return app
}
