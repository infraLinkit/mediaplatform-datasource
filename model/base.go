package model

import (
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"
)

const TIME_QUERY_EXEC = 72000

type (
	BaseModel struct {
		Config *config.Cfg
		Logs   *logrus.Logger
		DB     *gorm.DB
		R      *rueidis.Storage
		GS     *sheets.Service
	}
)

func NewBaseModel(obj BaseModel) *BaseModel {

	return &BaseModel{
		Config: obj.Config,
		Logs:   obj.Logs,
		DB:     obj.DB,
		R:      obj.R,
		GS:     obj.GS,
	}
}
