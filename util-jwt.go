package api_common

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

// GetJwtFromContext returns the jwt object, given the fiber context
func GetJwtFromContext(c *fiber.Ctx) (*jwt.Token, error) {
	user := c.Locals("user").(*jwt.Token)
	if user != nil {
		return user, nil
	} else {
		return nil, fmt.Errorf("cannot find jwt in context")
	}
}

// GetJwtUserId returns the userid in the jwt, given the fiber context
func GetJwtUserId(c *fiber.Ctx) (string, error) {
	userJwt, errGetJwtFromContext := GetJwtFromContext(c)
	if errGetJwtFromContext != nil {
		return "", fmt.Errorf("cannot find jwt in context")
	}
	userClaims := userJwt.Claims.(jwt.MapClaims)
	if userClaims == nil {
		return "", fmt.Errorf("malformed jwt, cannot find any claims")
	}
	if userClaims["sub"] == nil {
		return "", fmt.Errorf("malformed jwt, cannot find sub claim")
	}
	userId := userClaims["sub"].(string)
	if userId == "" {
		return "", fmt.Errorf("malformed jwt, cannot find sub claim")
	}
	return userId, nil
}

func RequiresRefreshToken(serviceConfig MicroserviceConfiguration) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		// get userid from jwt
		token, err := GetJwtFromContext(ctx)
		if err != nil {
			log.WithError(err).Panic("cannot get jwt from context")
			return ctx.Status(401).JSON(GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "cannot get jwt from context", err.Error()))
		}
		claims := token.Claims.(jwt.MapClaims)

		if len(claims) != len(serviceConfig.Application.Jwt.Api.RefreshToken.Claims) {
			log.Errorf("invalid token provided")
			return ctx.Status(401).JSON(GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "invalid token", "invalid token provided"))
		}
		for i, _ := range claims {
			if !StringArrayContains(serviceConfig.Application.Jwt.Api.RefreshToken.Claims, i) {
				log.Errorf("invalid token provided")
				return ctx.Status(401).JSON(GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "invalid token", "invalid token provided"))
			}
		}

		return ctx.Next()
	}
}
