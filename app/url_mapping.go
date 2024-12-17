package app

import (
	"database/sql"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/handler"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
)

type App3rdParty struct {
	Config *config.Cfg
	Logs   *logrus.Logger
	PS     *sql.DB
	R      *rueidis.Storage
	Rmqp   rmqp.AMQP
}

func MapUrls(obj App3rdParty) *fiber.App {

	f := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	h := handler.NewIncomingHandler(handler.IncomingHandler{
		Config: obj.Config,
		Logs:   obj.Logs,
		R:      obj.R,
		PS:     obj.PS,
		Rmqp:   obj.Rmqp,
	})

	// Landing Page
	f.Get("/v1/postback/", h.Postback).Name("Postback from messaging")
	f.Get("/v1/report/", h.Report).Name("Report API")
	f.Put("/v1/int/setdata/:v/", h.SetData).Name("SetTargetDailyBudget")

	return f
}
