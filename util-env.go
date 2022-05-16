package api_common

import (
	"crypto/rsa"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type Env struct {
	DB                    *gorm.DB
	ServiceConfig         MicroserviceConfiguration
	DatabaseEncryptionKey string
	AuthJwtPrivateKey     *rsa.PrivateKey
	AuthJwtPublicKey      *rsa.PublicKey
	RabbitChannel         *amqp.Channel
}
