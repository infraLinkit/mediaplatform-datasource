package app

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/handler"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
	"gorm.io/gorm"
)

type App3rdParty struct {
	Config *config.Cfg
	Logs   *logrus.Logger
	DB     *gorm.DB
	R0     *rueidis.Storage
	R1     *rueidis.Storage
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
		R0:     obj.R0,
		R1:     obj.R1,
		DB:     obj.DB,
		Rmqp:   obj.Rmqp,
	})

	// V1
	v1 := f.Group("/v1") // v1

	// Postback
	v1.Get("/postback/:urlservicekey/", h.Postback)

	// Report
	rpt := v1.Group("/report") // Report
	//rpt.Get("/report/", h.Report).Name("Report API")
	rpt.Get("/pinreport", h.DisplayPinReport).Name("Pin Report Summary FE")
	rpt.Get("/datasentapiperformance", h.DisplayPinPerformanceReport).Name("Pin Performance Api Report Summary FE")
	rpt.Get("/cpareportlist", h.DisplayCPAReport).Name("Receive Pin CPA Report Transactional")
	rpt.Get("/costreport/:v", h.DisplayCostReport).Name("Receive Pin Cost Report / detail Transactional")

	// API Internal
	internal := v1.Group("/int") // Internal API
	internal.Put("/setdata/:v/", h.SetData).Name("SetTargetDailyBudget")
	internal.Put("/updatedata/:v/", h.UpdateAgencyFeeAndCostConversion).Name("UpdateAgencyFeeAndCostConversion")
	internal.Put("/updateratio/:v/", h.UpdateRatio).Name("Update Ratio Transactional")
	internal.Put("/updatepostback/:v/", h.UpdatePostback).Name("Update Postback Transactional")
	internal.Put("/updateagencycost/:v", h.UpdateAgencyCost).Name("Update Agency fee and cost per conversion in db")
	internal.Get("/pinreport/", h.TrxPinReport).Name("Receive Pin Report Transactional")
	internal.Get("/datasentapiperformance/", h.TrxPerformancePinReport).Name("Receive Pin API Performance Report Transactional")
	internal.Get("/exportcpa/", h.ExportCpaButton).Name("Export CPA-Report Button")
	internal.Get("/exportcost/", h.ExportCostButton).Name("Export Cost-Report Button")
	internal.Get("/exportcostdetail/", h.ExportCostDetailButton).Name("Export Cost-Report-Detail Button")

	// API External
	v1.Group("/ext") // External API

	return f
}
