package handler

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/model"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"
)

type (
	IncomingHandler struct {
		Config *config.Cfg
		Logs   *logrus.Logger
		DB     *gorm.DB
		Rmqp   rmqp.AMQP
		R      *rueidis.Storage
		RCP    *redis.Client
		DS     *model.BaseModel
		GS     *sheets.Service
	}
)

func NewIncomingHandler(obj IncomingHandler) *IncomingHandler {

	b := model.NewBaseModel(model.BaseModel{
		Config: obj.Config,
		Logs:   obj.Logs,
		DB:     obj.DB,
		R:      obj.R,
		RCP:    obj.RCP,
	})

	return &IncomingHandler{
		Config: obj.Config,
		Logs:   obj.Logs,
		DB:     obj.DB,
		R:      obj.R,
		RCP:    obj.RCP,
		Rmqp:   obj.Rmqp,
		DS:     b,
		GS:     obj.GS,
	}
}
