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

func GetRabbitQueue(ch *amqp.Channel, ex string, q string) error {
	err := ch.ExchangeDeclare(
		ex,       // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		q,     // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

func GetRabbitConsumer(ch *amqp.Channel, exchange string, queue string, key string) (<-chan amqp.Delivery, error) {
	var err error
	var msgs <-chan amqp.Delivery

	err = ch.QueueBind(queue, key, exchange, false, nil)

	if err != nil {
		return nil, TernaryOperator(ch.Close() != nil, err, errors.New("cannot close channel")).(error)
	}
	msgs, err = ch.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)
	if err != nil {
		return nil, TernaryOperator(ch.Close() != nil, err, errors.New("cannot close channel")).(error)
	}
	return msgs, nil
}
