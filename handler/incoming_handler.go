package handler

import (
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/model"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
	"gorm.io/gorm"
)

type (
	IncomingHandler struct {
		Config *config.Cfg
		Logs   *logrus.Logger
		DB     *gorm.DB
		Rmqp   rmqp.AMQP
		R0     *rueidis.Storage
		R1     *rueidis.Storage
		DS     *model.BaseModel
	}
)

func NewIncomingHandler(obj IncomingHandler) *IncomingHandler {

	b := model.NewBaseModel(model.BaseModel{
		Config: obj.Config,
		Logs:   obj.Logs,
		DB:     obj.DB,
		R0:     obj.R0,
		R1:     obj.R1,
	})

	return &IncomingHandler{
		Config: obj.Config,
		Logs:   obj.Logs,
		DB:     obj.DB,
		R0:     obj.R0,
		R1:     obj.R1,
		Rmqp:   obj.Rmqp,
		DS:     b,
	}
}
