package api_common

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

func GetRabbitChannel(url string) (*amqp.Channel, error) {
	var err error
	var conn *amqp.Connection
	var ch *amqp.Channel

	conn, err = amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err = conn.Channel()
	if err != nil {
		return nil, TernaryOperator(conn.Close() != nil, err, errors.New("cannot close connection")).(error)
	}

	return ch, nil
}

func PublishMessage(channel *amqp.Channel, exchange string, key string, json []byte) error {
	err := channel.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{Body: json})
	if err != nil {
		return err
	}
	return nil
}

func PublishToMonitor(response interface{}, c *fiber.Ctx, status int, channel *amqp.Channel, exchange string, key string, source string, sourceType string) (int, interface{}, error) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		status = 500
		response = GetErrorResponse(API_CODE_COMMON_INTERNAL_SERVER_ERROR, "api_common", "cannot marshal monitor request")
	} else {
		base64Response := base64.URLEncoding.EncodeToString(jsonResponse)
		monitorRequest := MonitorRequest{Data: MonitorData{Monitor: Monitor{
			Response:   base64Response,
			Uuid:       c.Locals(CTX_REQUESTID).(string),
			Source:     source,
			SourceType: sourceType,
			Success:    TernaryOperator(status != 200, false, true).(bool),
			Status:     status,
			Endpoint:   c.OriginalURL(),
		}}}

		var monitorJson []byte
		monitorJson, err = json.Marshal(monitorRequest)
		if err != nil {
			status = 500
			response = GetErrorResponse(API_CODE_COMMON_INTERNAL_SERVER_ERROR, "api_common", "cannot marshal monitor request")
		} else {
			err = PublishMessage(channel, exchange, key, monitorJson)
			if err != nil {
				status = 500
				response = GetErrorResponse(API_CODE_COMMON_INTERNAL_SERVER_ERROR, "api_common", "cannot publish message to monitor queue")
			}
		}
	}
	return status, response, err
}
