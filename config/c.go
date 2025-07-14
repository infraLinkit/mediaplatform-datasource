package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var APP_PATH = "/Users/wiliewahyuhidayat/Documents/GO/mediaplatform/cores/" // local

const (
	//APP_PATH       = "/app"
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
		Host1                                  string
		Host2                                  string
		Host3                                  string
		Host4                                  string
		Host5                                  string
		Host6                                  string
		Host7                                  string
		Host8                                  string
		Host9                                  string
		Host10                                 string
		Port1                                  string
		Port2                                  string
		Port3                                  string
		Port4                                  string
		Port5                                  string
		Port6                                  string
		Port7                                  string
		Port8                                  string
		Port9                                  string
		Port10                                 string
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
		ARPUUsername                           string
		ARPUPassword                           string
		APILINKITDashboard                     string
		GSType                                 string
		GSProjectID                            string
		GSPrivateKeyID                         string
		GSPrivateKey                           string
		GSClientEmail                          string
		GSClientID                             string
		GSAuthURI                              string
		GSTokenURI                             string
		GSAuthProvider                         string
		GSClient                               string
		GSUniversalDomain                      string
		CronResetCapping                       string
	}

	Setup struct {
		Config *Cfg
		Logs   *logrus.Logger
		R      *rueidis.Storage
		DB     *gorm.DB
		Rmqp   rmqp.AMQP
		GS     *sheets.Service
	}
)

func InitCfg() *Cfg {

	if os.Getenv("APPPATH") != "" {
		APP_PATH = os.Getenv("APPPATH")
	}

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
		Host1:                                  os.Getenv("HOST1"),
		Host2:                                  os.Getenv("HOST2"),
		Host3:                                  os.Getenv("HOST3"),
		Host4:                                  os.Getenv("HOST4"),
		Host5:                                  os.Getenv("HOST5"),
		Host6:                                  os.Getenv("HOST6"),
		Host7:                                  os.Getenv("HOST7"),
		Host8:                                  os.Getenv("HOST8"),
		Host9:                                  os.Getenv("HOST9"),
		Host10:                                 os.Getenv("HOST10"),
		Port1:                                  os.Getenv("PORT1"),
		Port2:                                  os.Getenv("PORT2"),
		Port3:                                  os.Getenv("PORT3"),
		Port4:                                  os.Getenv("PORT4"),
		Port5:                                  os.Getenv("PORT5"),
		Port6:                                  os.Getenv("PORT6"),
		Port7:                                  os.Getenv("PORT7"),
		Port8:                                  os.Getenv("PORT8"),
		Port9:                                  os.Getenv("PORT9"),
		Port10:                                 os.Getenv("PORT10"),
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
		ARPUUsername:                           os.Getenv("ARPU_USERNAME"),
		ARPUPassword:                           os.Getenv("ARPU_PASSWORD"),
		APILINKITDashboard:                     os.Getenv("APILINKIT_DASHBOARD"),
		GSType:                                 os.Getenv("GSTYPE"),
		GSProjectID:                            os.Getenv("GSPROJECT_ID"),
		GSPrivateKeyID:                         os.Getenv("GSPRIVATE_KEY_ID"),
		GSPrivateKey:                           os.Getenv("GSPRIVATE_KEY"),
		GSClientEmail:                          os.Getenv("GSCLIENT_EMAIL"),
		GSClientID:                             os.Getenv("GSCLIENT_ID"),
		GSAuthURI:                              os.Getenv("GSAUTH_URI"),
		GSTokenURI:                             os.Getenv("GSTOKEN_URI"),
		GSAuthProvider:                         os.Getenv("GSAUTH_PROVIDER"),
		GSClient:                               os.Getenv("GSCLIENT"),
		GSUniversalDomain:                      os.Getenv("GSUNIVERSAL_DOMAIN"),
		CronResetCapping:                       os.Getenv("CRONRESETCAPPING"),
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
		GS:     c.InitGoogleSheet(l),
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

func (c *Cfg) InitGoogleSheet(l *logrus.Logger) *sheets.Service {
	type KeyGS struct {
		Type            string `json:"type"`
		ProjectID       string `json:"project_id"`
		PrivateKeyID    string `json:"private_key_id"`
		PrivateKey      string `json:"private_key"`
		ClientEmail     string `json:"client_email"`
		ClientID        string `json:"client_id"`
		AuthURI         string `json:"auth_uri"`
		TokenURI        string `json:"token_uri"`
		AuthProvider    string `json:"auth_provider_x509_cert_url"`
		Client          string `json:"client_x509_cert_url"`
		UniversalDomain string `json:"universal_domain"`
	}

	var (
		sheetKey = KeyGS{
			Type:            c.GSType,
			ProjectID:       c.GSProjectID,
			PrivateKeyID:    c.GSPrivateKeyID,
			PrivateKey:      strings.ReplaceAll(c.GSPrivateKey, `\n`, "\n"),
			ClientEmail:     c.GSClientEmail,
			ClientID:        c.GSClientID,
			AuthURI:         c.GSAuthURI,
			TokenURI:        c.GSTokenURI,
			AuthProvider:    c.GSAuthProvider,
			Client:          c.GSClient,
			UniversalDomain: c.GSUniversalDomain,
		}
	)

	credential, err := json.Marshal(sheetKey)
	if err != nil {
		l.Fatalf("Google Sheet Failed to get key: %v", err)
	}

	srv, err := sheets.NewService(context.Background(), option.WithCredentialsJSON(credential))
	if err != nil {
		l.Fatalf("Unable to retrieve Google Sheets client: %v", err)
	}

	l.Info("[v] Google sheet connection successful established\n")

	return srv
}
