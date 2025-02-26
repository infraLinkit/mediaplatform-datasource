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
	rpt.Get("/conversionlog", h.DisplayConversionLogReport).Name("Conversion Log Report")
	rpt.Get("/campaign-monitoring-summary", h.DisplayCampaignSummary).Name("Campaign Summary")

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

	// Management
	management := v1.Group("/management") // Management
	// Campaign
	campaign := management.Group("/campaign") // Campaign
	campaign.Get("/", h.DisplayCampaignManagement).Name("Campaign Management FE")
	campaign.Get("/:v", h.DisplayCampaignManagement).Name("Campaign Management Detail FE")
	campaign.Post("/send", h.SendCampaignHandler).Name("Campaign Management Send FE")
	// Menu
	menu := management.Group("/menu") // Menu
	menu.Post("/", h.CreateMenu).Name("Menu Management Create FE")
	menu.Get("/", h.GetAllMenus).Name("Menu Management FE")
	menu.Get("/:id", h.GetMenuByID).Name("Menu Management Edit FE")
	menu.Put("/:id", h.UpdateMenu).Name("Menu Management Update FE")
	menu.Delete("/:id", h.DeleteMenu).Name("Menu Management Delete FE")

	//  Country and Service Management
	countryService := management.Group("/country-service")
	countryService.Get("/country", h.DisplayCountry).Name("Create Country")
	countryService.Post("/country", h.CreateCountry).Name("Create Country")
	countryService.Put("/country/:id", h.UpdateCountry).Name("Update Country")
	countryService.Get("/company", h.DisplayCompany).Name("Create Company")
	countryService.Post("/company", h.CreateCompany).Name("Create Company")
	countryService.Put("/company/:id", h.UpdateCompany).Name("Update Company")
	countryService.Get("/domain", h.DisplayDomain).Name("Create Domain")
	countryService.Post("/domain", h.CreateDomain).Name("Create Domain")
	countryService.Put("/domain/:id", h.UpdateDomain).Name("Update Domain")
	countryService.Get("/operator", h.DisplayOperator).Name("Create Operator")
	countryService.Post("/operator", h.CreateOperator).Name("Create Operator")
	countryService.Put("/operator/:id", h.UpdateOperator).Name("Update Operator")
	countryService.Get("/partner", h.DisplayPartner).Name("Create Partner")
	countryService.Post("/partner", h.CreatePartner).Name("Create Partner")
	countryService.Put("/partner/:id", h.UpdatePartner).Name("Update Partner")
	countryService.Get("/service", h.DisplayService).Name("Create Service")
	countryService.Post("/service", h.CreateService).Name("Create Service")
	countryService.Put("/service/:id", h.UpdateService).Name("Update Service")
	countryService.Get("/adnet-list", h.DisplayAdnetList).Name("Create AdnetList")
	countryService.Post("/adnet-list", h.CreateAdnetList).Name("Create AdnetList")
	countryService.Put("/adnet-list/:id", h.UpdateAdnetList).Name("Update AdnetList")

	// API External
	v1.Group("/ext") // External API

	return f
}
