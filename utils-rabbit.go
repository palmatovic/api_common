package api_common

import (
	"errors"
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
