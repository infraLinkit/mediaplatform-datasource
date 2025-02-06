package model

import (
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const TIME_QUERY_EXEC = 72000

type (
	BaseModel struct {
		Config *config.Cfg
		Logs   *logrus.Logger
		DB     *gorm.DB
		R0     *rueidis.Storage
		R1     *rueidis.Storage
	}
)

func NewBaseModel(obj BaseModel) *BaseModel {

	return &BaseModel{
		Config: obj.Config,
		Logs:   obj.Logs,
		DB:     obj.DB,
		R0:     obj.R0,
		R1:     obj.R1,
	}
}
