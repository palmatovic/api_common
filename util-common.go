package api_common

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"regexp"
	"strconv"
)

// ParseBool parses a string to extract a boolean or returns
// a default value
func ParseBool(input string, defaultValue bool) bool {
	result, err := strconv.ParseBool(input)
	if err != nil {
		return defaultValue
	}
	return result
}

// StringArrayContains returns true if stringArray contains the provided
// element, false otherwise
func StringArrayContains(stringArray []string, element string) bool {
	for _, a := range stringArray {
		if a == element {
			return true
		}
	}
	return false
}

// IntArrayContains returns true if intArray contains the provided
// element, false otherwise
func IntArrayContains(stringArray []int, element int) bool {
	for _, a := range stringArray {
		if a == element {
			return true
		}
	}
	return false
}

// CountCharTypes returns the number of each char type
func CountCharTypes(input string) (lower int, upper int, numeric int, special int) {
	lower, upper, numeric, special = 0, 0, 0, 0
	for i := 0; i < len(input); i++ {
		if input[i] >= 65 && input[i] <= 90 {
			upper++
		} else if input[i] >= 97 && input[i] <= 122 {
			lower++
		} else if input[i] >= 48 && input[i] <= 57 {
			numeric++
		} else {
			special++
		}
	}
	return lower, upper, numeric, special
}

//CheckRegex checks if string respect regex rule
func CheckRegex(regexRule string, stringToCheck string) (bool, error) {
	result, err := regexp.MatchString(regexRule, stringToCheck)
	return result, err
}

// TernaryOperator reproduces the ternary operator construct missing in golang
func TernaryOperator(condition bool, ifTrue interface{}, ifFalse interface{}) interface{} {
	if condition {
		return ifTrue
	} else {
		return ifFalse
	}
}

func GetErrorResponse(code string, reason string, detail string) interface{} {
	return fiber.Map{
		"status": false,
		"data": ErrorData{Error: Error{
			ErrorCode: code,
			Reason:    reason,
			Detail:    detail,
		}},
	}
}

func GetSuccessResponse(data interface{}) interface{} {
	return fiber.Map{
		"status": true,
		"data":   data,
	}
}

func Response(c *fiber.Ctx, response interface{}, status int, channel *amqp.Channel, exchange string, key string, source string) error {
	var err error
	status, response, err = PublishToMonitor(response, c, status, channel, exchange, key, source, "rest", nil, nil)
	if err != nil {
		log.WithFields(log.Fields{"uuid": c.Locals(CTX_REQUESTID).(string)}).WithError(err).Errorf("cannot send message to monitor queue")
	} else {
		log.WithFields(log.Fields{"uuid": c.Locals(CTX_REQUESTID).(string)}).Infof("sent message to monitor queue")
	}
	return c.Status(status).JSON(response)
}
