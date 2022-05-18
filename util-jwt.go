package api_common

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strconv"
)

// GetJwtFromContext returns the jwt object, given the fiber context
func GetJwtFromContext(c *fiber.Ctx) (*jwt.Token, error) {
	user := c.Locals("user")
	if user != nil && user.(*jwt.Token) != nil {
		return user.(*jwt.Token), nil
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

func RequiresRefreshToken(serviceConfig MicroserviceConfiguration, channel *amqp.Channel) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		// get userid from jwt
		var response interface{}
		token, err := GetJwtFromContext(ctx)
		if err != nil {
			log.WithError(err).Panic("cannot get jwt from context")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires refresh token", err.Error())
			_, _, err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, "auth", "rest", nil, nil)
			if err != nil {
				log.WithError(err).Errorf("cannot send message to monitor")
			} else {
				log.Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		claims := token.Claims.(jwt.MapClaims)

		if len(claims) != len(serviceConfig.Application.Jwt.Api.RefreshToken.Claims) {
			log.Errorf("invalid token provided")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires refresh token", "invalid token provided")
			_, _, err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, "auth", "rest", nil, nil)
			if err != nil {
				log.WithError(err).Errorf("cannot send message to monitor")
			} else {
				log.Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		for i, _ := range claims {
			if !StringArrayContains(serviceConfig.Application.Jwt.Api.RefreshToken.Claims, i) {
				log.Errorf("invalid token provided")
				response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires refresh token", "invalid token provided")
				_, _, err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, "auth", "rest", nil, nil)
				if err != nil {
					log.WithError(err).Errorf("cannot send message to monitor")
				} else {
					log.Infof("successfully sent message to monitor")
				}
				return ctx.Status(401).JSON(response)
			}
		}
		return ctx.Next()
	}
}

func RequiresAccessToken(applicationClaims []string, channel *amqp.Channel, serviceConfig MicroserviceConfiguration) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var response interface{}
		token, err := GetJwtFromContext(ctx)
		if err != nil {
			log.WithError(err).Panic("cannot get jwt from context")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires access token", err.Error())
			_, _, err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, "auth", "rest", nil, nil)
			if err != nil {
				log.WithError(err).Errorf("cannot send message to monitor")
			} else {
				log.Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		claims := token.Claims.(jwt.MapClaims)

		if len(claims) != len(applicationClaims) {
			log.Errorf("invalid token provided")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires access token", "invalid token provided")
			_, _, err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, "auth", "rest", nil, nil)
			if err != nil {
				log.WithError(err).Errorf("cannot send message to monitor")
			} else {
				log.Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		for i, _ := range claims {
			if !StringArrayContains(applicationClaims, i) {
				log.Errorf("invalid token provided")
				response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires access token", "invalid token provided")
				_, _, err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, "auth", "rest", nil, nil)
				if err != nil {
					log.WithError(err).Errorf("cannot send message to monitor")
				} else {
					log.Infof("successfully sent message to monitor")
				}
				return ctx.Status(401).JSON(response)
			}
		}

		return ctx.Next()
	}
}

func RequiresHierarchy(hierarchies []int, channel *amqp.Channel, serviceConfig MicroserviceConfiguration) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var response interface{}
		token, err := GetJwtFromContext(ctx)
		if err != nil {
			log.WithError(err).Panic("cannot get jwt from context")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires hierarchy", err.Error())
			_, _, err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, "auth", "rest", nil, nil)
			if err != nil {
				log.WithError(err).Errorf("cannot send message to monitor")
			} else {
				log.Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		claims := token.Claims.(jwt.MapClaims)
		jwtHierarchy := int(claims["hierarchy"].(float64))
		if !IntArrayContains(hierarchies, jwtHierarchy) {
			log.Errorf("Unauthorized user hierarchy: %d, with role %s", jwtHierarchy, claims["role"].(string))
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires hierarchy", fmt.Sprintf("Unauthorized user hierarchy: %d, with role %s", jwtHierarchy, claims["role"].(string)))
			_, _, err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, "auth", "rest", nil, nil)
			if err != nil {
				log.WithError(err).Errorf("cannot send message to monitor")
			} else {
				log.Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		return ctx.Next()
	}
}

func RequiresFirstLogin(isRequired bool, channel *amqp.Channel, serviceConfig MicroserviceConfiguration) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var response interface{}
		token, err := GetJwtFromContext(ctx)
		if err != nil {
			log.WithError(err).Panic("cannot get jwt from context")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires first login", "cannot get jwt from context")
			_, _, err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, "auth", "rest", nil, nil)
			if err != nil {
				log.WithError(err).Errorf("cannot send message to monitor")
			} else {
				log.Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		claims := token.Claims.(jwt.MapClaims)
		var firstLogin bool
		firstLogin, err = strconv.ParseBool(claims["first_login"].(string))
		if err != nil {
			log.WithError(err).Panic("cannot get first_login claim")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires first login", "cannot get first_login claim")
			_, _, err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, "auth", "rest", nil, nil)
			if err != nil {
				log.WithError(err).Errorf("cannot send message to monitor")
			} else {
				log.Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		if firstLogin != isRequired {
			log.Errorf("invalid token provided")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires first login", "invalid token provided")
			_, _, err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, "auth", "rest", nil, nil)
			if err != nil {
				log.WithError(err).Errorf("cannot send message to monitor")
			} else {
				log.Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}

		return ctx.Next()
	}
}
