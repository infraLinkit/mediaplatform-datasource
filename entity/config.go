package entity

import (
	"database/sql"
	"time"

	"github.com/gofiber/storage/rueidis"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
)

type (
	Cfg struct {
		AppHost              string
		RedisHost            string
		RedisPort            int
		RedisPwd             string
		RedisKeyExpiration   int64
		PSQLHost             string
		PSQLUsername         string
		PSQLPassword         string
		PSQLPort             string
		PSQLDB               string
		RabbitMQHost         string
		RabbitMQPort         int
		RabbitMQUsername     string
		RabbitMQPassword     string
		RabbitMQVHost        string
		RabbitMQExchangeName string
		RabbitMQQueueName    string
		RabbitMQDataType     string
		LogPath              string
		LogLevel             string
		TZ                   *time.Location
	}

	Setup struct {
		Config *Cfg
		Logs   *logrus.Logger
		R      *rueidis.Storage
		DB     *sql.DB
		Rmqp   rmqp.AMQP
	}
)
