package api_common

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strconv"
)

// Elog add extra info on every log
func Elog(c *fiber.Ctx) *log.Entry {
	actor, org, role, hierarchy, _ := GetJwtUser(c)
	ips := append([]string{c.IP()}, c.IPs()...)
	reqId := c.Locals(CTX_REQUESTID).(string)
	return log.WithFields(log.Fields{
		"actor":     actor,
		"org":       org,
		"role":      role,
		"hierarchy": hierarchy,
		"ips":       ips,
		"uuid":      reqId,
	})
}

// GetJwtFromContext returns the jwt object, given the fiber context
func GetJwtFromContext(c *fiber.Ctx) (*jwt.Token, error) {
	user := c.Locals("user")
	if user != nil && user.(*jwt.Token) != nil {
		return user.(*jwt.Token), nil
	} else {
		return nil, fmt.Errorf("cannot find jwt in context")
	}
}

// GetJwtUser returns the userid in the jwt, given the fiber context
func GetJwtUser(c *fiber.Ctx) (string, string, string, string, error) {
	userJwt, errGetJwtFromContext := GetJwtFromContext(c)
	if errGetJwtFromContext != nil {
		return "", "", "", "", fmt.Errorf("cannot find jwt in context")
	}
	userClaims := userJwt.Claims.(jwt.MapClaims)
	if userClaims == nil {
		return "", "", "", "", fmt.Errorf("malformed jwt, cannot find any claims")
	}
	if userClaims["sub"] == nil {
		return "", "", "", "", fmt.Errorf("malformed jwt, cannot find sub claim")
	}
	userId := userClaims["sub"].(string)
	if userId == "" {
		return "", "", "", "", fmt.Errorf("malformed jwt, cannot find sub claim")
	}
	if userClaims["org"] == nil {
		return "", "", "", "", fmt.Errorf("malformed jwt, cannot find sub claim")
	}
	org := userClaims["org"].(string)
	if org == "" {
		return "", "", "", "", fmt.Errorf("malformed jwt, cannot find sub claim")
	}
	if userClaims["role"] == nil {
		return "", "", "", "", fmt.Errorf("malformed jwt, cannot find sub claim")
	}
	role := userClaims["role"].(string)
	if role == "" {
		return "", "", "", "", fmt.Errorf("malformed jwt, cannot find sub claim")
	}
	if userClaims["hierarchy"] == nil {
		return "", "", "", "", fmt.Errorf("malformed jwt, cannot find sub claim")
	}
	hierarchy := userClaims["hierarchy"].(string)
	if hierarchy == "" {
		return "", "", "", "", fmt.Errorf("malformed jwt, cannot find sub claim")
	}
	return userId, org, role, hierarchy, nil
}

func RequiresRefreshToken(serviceConfig MicroserviceConfiguration, channel *amqp.Channel, source string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		// get userid from jwt
		var response interface{}
		token, err := GetJwtFromContext(ctx)
		if err != nil {
			Elog(ctx).WithError(err).Panic("cannot get jwt from context")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires refresh token", err.Error())
			err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, source, "rest", nil, nil)
			if err != nil {
				Elog(ctx).WithError(err).Errorf("cannot send message to monitor")
			} else {
				Elog(ctx).Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		claims := token.Claims.(jwt.MapClaims)

		if len(claims) != len(serviceConfig.Application.Jwt.Api.RefreshToken.Claims) {
			Elog(ctx).Errorf("invalid token provided")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires refresh token", "invalid token provided")
			err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, source, "rest", nil, nil)
			if err != nil {
				Elog(ctx).WithError(err).Errorf("cannot send message to monitor")
			} else {
				Elog(ctx).Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		for i, _ := range claims {
			if !StringArrayContains(serviceConfig.Application.Jwt.Api.RefreshToken.Claims, i) {
				Elog(ctx).Errorf("invalid token provided")
				response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires refresh token", "invalid token provided")
				err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, source, "rest", nil, nil)
				if err != nil {
					Elog(ctx).WithError(err).Errorf("cannot send message to monitor")
				} else {
					Elog(ctx).Infof("successfully sent message to monitor")
				}
				return ctx.Status(401).JSON(response)
			}
		}
		return ctx.Next()
	}
}

func RequiresAccessToken(applicationClaims []string, channel *amqp.Channel, serviceConfig MicroserviceConfiguration, source string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var response interface{}
		token, err := GetJwtFromContext(ctx)
		if err != nil {
			Elog(ctx).WithError(err).Panic("cannot get jwt from context")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires access token", err.Error())
			err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, source, "rest", nil, nil)
			if err != nil {
				Elog(ctx).WithError(err).Errorf("cannot send message to monitor")
			} else {
				Elog(ctx).Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		claims := token.Claims.(jwt.MapClaims)

		if len(claims) != len(applicationClaims) {
			Elog(ctx).Errorf("invalid token provided")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires access token", "invalid token provided")
			err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, source, "rest", nil, nil)
			if err != nil {
				Elog(ctx).WithError(err).Errorf("cannot send message to monitor")
			} else {
				Elog(ctx).Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		for i, _ := range claims {
			if !StringArrayContains(applicationClaims, i) {
				Elog(ctx).Errorf("invalid token provided")
				response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires access token", "invalid token provided")
				err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, source, "rest", nil, nil)
				if err != nil {
					Elog(ctx).WithError(err).Errorf("cannot send message to monitor")
				} else {
					Elog(ctx).Infof("successfully sent message to monitor")
				}
				return ctx.Status(401).JSON(response)
			}
		}

		return ctx.Next()
	}
}

func RequiresHierarchy(hierarchies []int, channel *amqp.Channel, serviceConfig MicroserviceConfiguration, source string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var response interface{}
		token, err := GetJwtFromContext(ctx)
		if err != nil {
			Elog(ctx).WithError(err).Panic("cannot get jwt from context")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires hierarchy", err.Error())
			err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, source, "rest", nil, nil)
			if err != nil {
				Elog(ctx).WithError(err).Errorf("cannot send message to monitor")
			} else {
				Elog(ctx).Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		claims := token.Claims.(jwt.MapClaims)
		jwtHierarchy := int(claims["hierarchy"].(float64))
		if !IntArrayContains(hierarchies, jwtHierarchy) {
			Elog(ctx).Errorf("Unauthorized user hierarchy: %d, with role %s", jwtHierarchy, claims["role"].(string))
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires hierarchy", fmt.Sprintf("Unauthorized user hierarchy: %d, with role %s", jwtHierarchy, claims["role"].(string)))
			err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, source, "rest", nil, nil)
			if err != nil {
				Elog(ctx).WithError(err).Errorf("cannot send message to monitor")
			} else {
				Elog(ctx).Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		return ctx.Next()
	}
}

func RequiresFirstLogin(isRequired bool, channel *amqp.Channel, serviceConfig MicroserviceConfiguration, source string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var response interface{}
		token, err := GetJwtFromContext(ctx)
		if err != nil {
			Elog(ctx).WithError(err).Panic("cannot get jwt from context")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires first login", "cannot get jwt from context")
			err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, source, "rest", nil, nil)
			if err != nil {
				Elog(ctx).WithError(err).Errorf("cannot send message to monitor")
			} else {
				Elog(ctx).Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		claims := token.Claims.(jwt.MapClaims)
		var firstLogin bool
		firstLogin, err = strconv.ParseBool(claims["first_login"].(string))
		if err != nil {
			Elog(ctx).WithError(err).Panic("cannot get first_login claim")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires first login", "cannot get first_login claim")
			err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, source, "rest", nil, nil)
			if err != nil {
				Elog(ctx).WithError(err).Errorf("cannot send message to monitor")
			} else {
				Elog(ctx).Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}
		if firstLogin != isRequired {
			Elog(ctx).Errorf("invalid token provided")
			response = GetErrorResponse(API_CODE_COMMON_UNAUTHORIZED, "requires first login", "invalid token provided")
			err = PublishToMonitor(response, ctx, 401, channel, serviceConfig.Infrastructure.Rabbit.Monitor.Exchange, serviceConfig.Infrastructure.Rabbit.Monitor.Key, source, "rest", nil, nil)
			if err != nil {
				Elog(ctx).WithError(err).Errorf("cannot send message to monitor")
			} else {
				Elog(ctx).Infof("successfully sent message to monitor")
			}
			return ctx.Status(401).JSON(response)
		}

		return ctx.Next()
	}
}
