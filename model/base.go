package model

import (
	"database/sql"

	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const TIME_QUERY_EXEC = 72000

type (
	BaseModel struct {
		Config    *entity.Cfg
		Logs      *logrus.Logger
		DBPostgre *sql.DB
		R         *rueidis.Storage
	}
)

func NewBaseModel(obj BaseModel) *BaseModel {

	return &BaseModel{
		Config:    obj.Config,
		Logs:      obj.Logs,
		DBPostgre: obj.DBPostgre,
		R:         obj.R,
	}
}
