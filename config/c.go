package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	EVENT_TRAFFIC  = "Traffic"
	EVENT_LANDING  = "Landing"
	EVENT_CLICK    = "Click Landing"
	EVENT_REDIRECT = "Redirect"

	OK_DESC          = "OK"
	BAD_REQUEST_DESC = "Bad Request"
)

type (
	Cfg struct {
		AppDockerVer                           string
		AppHost                                string
		AppHostPort                            string
		AppApi                                 string
		AppApiPort                             string
		RedisHost                              string
		RedisPort                              int
		RedisDBIndex                           int
		RedisPwd                               string
		RedisKeyExpiration                     int64
		PSQLHost                               string
		PSQLUsername                           string
		PSQLPassword                           string
		PSQLPort                               string
		PSQLDB                                 string
		RabbitMQHost                           string
		RabbitMQPort                           int
		RabbitMQUsername                       string
		RabbitMQPassword                       string
		RabbitMQVHost                          string
		RabbitMQDataType                       string
		RabbitMQPixelStorageExchangeName       string
		RabbitMQPixelStorageQueueName          string
		RabbitMQRedisCounterExchangeName       string
		RabbitMQRedisCounterQueueName          string
		RabbitMQRatioExchangeName              string
		RabbitMQRatioQueueName                 string
		RabbitMQRatioQueueThreshold            int
		RabbitMQPostbackAdnetExchangeName      string
		RabbitMQPostbackAdnetQueueName         string
		RabbitMQCampaignManagementExchangeName string
		RabbitMQCampaignManagementQueueName    string
		RabbitMQAlertManagementExchangeName    string
		RabbitMQAlertManagementQueueName       string
		LogEnv                                 string
		LogPath                                string
		LogLevel                               string
		TZ                                     *time.Location
		APIARPU                                string
	}

	Setup struct {
		Config *Cfg
		Logs   *logrus.Logger
		R      *rueidis.Storage
		DB     *gorm.DB
		Rmqp   rmqp.AMQP
	}
)

func InitCfg() *Cfg {

	loc, _ := time.LoadLocation(os.Getenv("TZ"))

	rabbitmq_port, _ := strconv.Atoi(os.Getenv("RABBITMQPORT"))
	redis_dbindex, _ := strconv.Atoi(os.Getenv("REDISDBINDEX"))
	redis_port, _ := strconv.Atoi(os.Getenv("REDISPORT"))
	redis_exp, _ := strconv.Atoi(os.Getenv("REDISKEYEXPIRE"))
	ratio_queue_threshold, _ := strconv.Atoi(os.Getenv("RABBITMQRATIOQUEUETHRESHOLD"))

	cfg := &Cfg{
		AppDockerVer:                           os.Getenv("APP_DOCKER_VER"),
		AppHost:                                os.Getenv("APPHOST"),
		AppHostPort:                            os.Getenv("APPHOSTPORT"),
		AppApi:                                 os.Getenv("APPAPI"),
		AppApiPort:                             os.Getenv("APPAPIPORT"),
		RedisHost:                              os.Getenv("REDISHOST"),
		RedisPort:                              redis_port,
		RedisDBIndex:                           redis_dbindex,
		RedisPwd:                               os.Getenv("REDISPASSWORD"),
		RedisKeyExpiration:                     int64(redis_exp),
		PSQLHost:                               os.Getenv("DB_HOST"),
		PSQLUsername:                           os.Getenv("DB_USERNAME"),
		PSQLPassword:                           os.Getenv("DB_PASSWORD"),
		PSQLPort:                               os.Getenv("DB_PORT"),
		PSQLDB:                                 os.Getenv("DB_DATABASE"),
		RabbitMQHost:                           os.Getenv("RABBITMQHOST"),
		RabbitMQPort:                           rabbitmq_port,
		RabbitMQUsername:                       os.Getenv("RABBITMQUSERNAME"),
		RabbitMQPassword:                       os.Getenv("RABBITMQPASSWORD"),
		RabbitMQVHost:                          os.Getenv("RABBITMQVHOST"),
		RabbitMQDataType:                       "application/json",
		RabbitMQPixelStorageExchangeName:       os.Getenv("RABBITMQPIXELSTORAGEEXCHANGENAME"),
		RabbitMQPixelStorageQueueName:          os.Getenv("RABBITMQPIXELSTORAGEQUEUENAME"),
		RabbitMQRedisCounterExchangeName:       os.Getenv("RABBITMQREDISCOUNTEREXCHANGENAME"),
		RabbitMQRedisCounterQueueName:          os.Getenv("RABBITMQREDISCOUNTERQUEUENAME"),
		RabbitMQRatioExchangeName:              os.Getenv("RABBITMQRATIOEXCHANGENAME"),
		RabbitMQRatioQueueName:                 os.Getenv("RABBITMQRATIOQUEUENAME"),
		RabbitMQRatioQueueThreshold:            ratio_queue_threshold,
		RabbitMQPostbackAdnetExchangeName:      os.Getenv("RABBITMQPOSTBACKADNETEXCHANGENAME"),
		RabbitMQPostbackAdnetQueueName:         os.Getenv("RABBITMQPOSTBACKADNETQUEUENAME"),
		RabbitMQCampaignManagementExchangeName: os.Getenv("RABBITMQCAMPAIGNMANAGEMENTEXCHANGENAME"),
		RabbitMQCampaignManagementQueueName:    os.Getenv("RABBITMQCAMPAIGNMANAGEMENTQUEUENAME"),
		RabbitMQAlertManagementExchangeName:    os.Getenv("RABBITMQALERTMANAGEMENTEXCHANGENAME"),
		RabbitMQAlertManagementQueueName:       os.Getenv("RABBITMQALERTMANAGEMENTQUEUENAME"),
		LogEnv:                                 os.Getenv("LOGENV"),
		LogPath:                                os.Getenv("LOGPATH"),
		LogLevel:                               os.Getenv("LOGLEVEL"),
		TZ:                                     loc,
		APIARPU:                                os.Getenv("APIARPU"),
	}

	return cfg
}

func (c *Cfg) Initiate(logname string) *Setup {

	//l := helper.MakeLogger(c.LogPath+"/"+logname, true, c.LogLevel)
	l := helper.MakeLogger(
		helper.Setup{Env: c.LogEnv, Logname: c.LogPath + "/" + logname, Display: true, Level: c.LogLevel})
	l.Info(fmt.Sprintf("Config Loaded : %#v\n", c))

	return &Setup{
		Config: c,
		Logs:   l,
		R:      c.InitRedisJSON(l, c.RedisDBIndex),
		DB:     c.InitGormPgx(l),
		Rmqp:   c.InitMessageBroker(),
	}

}

func (c *Cfg) InitRedis(l *logrus.Logger, dbindex int) *redis.Client {

	r := redis.NewClient(&redis.Options{
		Addr:     c.RedisHost + ":" + strconv.Itoa(c.RedisPort),
		Password: c.RedisPwd,
		DB:       dbindex,
	})

	pong, errRedis := r.Ping().Result()

	if errRedis == nil && pong == "PONG" {

		l.Info(fmt.Sprintf("[v] conn successful established of the redis : %s\n", pong))
		return r
	} else {

		l.Info(fmt.Sprintf("[x] An Error occured when establishing of the redis : %#v\n", errRedis))

		panic(errRedis)
	}

}

func (c *Cfg) InitRedisJSON(l *logrus.Logger, dbindex int) *rueidis.Storage {

	port := strconv.Itoa(c.RedisPort)

	r := rueidis.New(rueidis.Config{
		InitAddress: []string{c.RedisHost + ":" + port},
		Username:    "",
		Password:    c.RedisPwd,
		SelectDB:    dbindex,
		Reset:       false,
		TLSConfig:   nil,
	})

	pong := r.Conn().B().Ping()
	l.Info(fmt.Sprintf("[v] Status conn successful established of the redis : %#v\n", pong))

	return r
}

func (c *Cfg) InitPsql(l *logrus.Logger) *sql.DB {

	db, err := sql.Open("postgres", "postgresql://"+c.PSQLUsername+":"+c.PSQLPassword+"@"+c.PSQLHost+":"+c.PSQLPort+"/"+c.PSQLDB+"?sslmode=disable")
	if err != nil {

		// panic the function then hard exit
		l.Info(fmt.Sprintf("[x] An Error occured when establishing of the database : %#v\n", err))

		panic(err)

	} else {

		l.Info("[v] Database successful established\n")
	}

	return db
}

func (c *Cfg) InitGormPgx(l *logrus.Logger) *gorm.DB {

	/* dsn := "host=" + c.PSQLHost + " user=" + c.PSQLUsername + " password=" + c.PSQLPassword + " dbname=" + c.PSQLDB + " port=" + c.PSQLPort + " sslmode=disable TimeZone=" + c.TZ.String()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) */

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=" + c.PSQLHost + " user=" + c.PSQLUsername + " password=" + c.PSQLPassword + " dbname=" + c.PSQLDB + " port=" + c.PSQLPort + " sslmode=disable TimeZone=" + c.TZ.String(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {

		// panic the function then hard exit
		l.Info(fmt.Sprintf("[x] An Error occured when establishing of the database : %#v\n", err))

		panic(err)

	} else {

		l.Info("[v] Database GORM successful established\n")

		sqlDB, _ := db.DB()

		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(10)

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(100)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return db
}

func (c *Cfg) InitMessageBroker() rmqp.AMQP {
	var rb rmqp.AMQP

	rb.SetAmqpURL(c.RabbitMQHost, c.RabbitMQPort, c.RabbitMQUsername, c.RabbitMQPassword, c.RabbitMQVHost)

	rb.SetUpConnectionAmqp()

	//rb.SetUpChannel("direct", true, c.RabbitMQExchangeName, true, c.RabbitMQQueueName)
	return rb
}
