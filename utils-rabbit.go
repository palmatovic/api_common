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

func PublishToMonitor(response interface{}, c *fiber.Ctx, status int, channel *amqp.Channel, exchange string, key string, source string, sourceType string, uuid *string, url *string) (int, interface{}, error) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return 500, GetErrorResponse(API_CODE_COMMON_INTERNAL_SERVER_ERROR, "api_common", "cannot marshal monitor request"), err
	}
	var uuidStr, urlStr string
	if uuid == nil {
		uuidStr = ""
	} else {
		uuidStr = *uuid
	}
	if url == nil {
		urlStr = ""
	} else {
		urlStr = *url
	}

	base64Response := base64.URLEncoding.EncodeToString(jsonResponse)
	monitorRequest := MonitorRequest{Data: MonitorData{Monitor: Monitor{
		Response:   base64Response,
		Uuid:       TernaryOperator(c != nil, c.Locals(CTX_REQUESTID).(string), uuidStr).(string),
		Source:     source,
		SourceType: sourceType,
		Success:    TernaryOperator(status != 200, false, true).(bool),
		Status:     status,
		Endpoint:   TernaryOperator(c != nil, c.OriginalURL(), urlStr).(string),
	}}}

	var monitorJson []byte
	monitorJson, err = json.Marshal(monitorRequest)
	if err != nil {
		return 500, GetErrorResponse(API_CODE_COMMON_INTERNAL_SERVER_ERROR, "api_common", "cannot marshal monitor request"), err
	}
	err = PublishMessage(channel, exchange, key, monitorJson)
	if err != nil {
		return 500, GetErrorResponse(API_CODE_COMMON_INTERNAL_SERVER_ERROR, "api_common", "cannot publish message to monitor queue"), err
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

func PublishToErmes(email string, template string, parameters *[]string, exchange string, queue string, key string, userId string, channel *amqp.Channel) (int, interface{}, error) {
	var status int
	var response interface{}
	var err error
	var jsn []byte

	jsn, err = json.Marshal(ErmesQueue{
		Data: &ErmesQueueData{
			Error: nil,
			ErmesInfo: ErmesInfo{
				To:         email,
				Template:   template,
				Parameters: parameters,
			},
			RabbitReply: RabbitReply{
				Exchange: exchange,
				Queue:    queue,
				Key:      key,
			},
			UserID: &userId,
		},
	})

	if err != nil {
		return 500, GetErrorResponse(API_CODE_COMMON_INTERNAL_SERVER_ERROR, "create user", "cannot marshal message for ermes"), err
	}
	err = PublishMessage(channel, exchange, key, jsn)
	if err != nil {
		return 500, GetErrorResponse(API_CODE_COMMON_INTERNAL_SERVER_ERROR, "create user", "cannot publish to ermes"), err
	}
	return status, response, err
}
