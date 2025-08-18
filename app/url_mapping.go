package app

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/handler"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	_ "github.com/lib/pq"
	"github.com/mikhail-bigun/fiberlogrus"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
	"gorm.io/gorm"
)

type App3rdParty struct {
	Config *config.Cfg
	Logs   *logrus.Logger
	DB     *gorm.DB
	R      *rueidis.Storage
	Rmqp   rmqp.AMQP
}

func MapUrls(obj App3rdParty) *fiber.App {

	f := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		BodyLimit:   100 * 1024 * 1024,
	})

	f.Use(
		fiberlogrus.New(
			fiberlogrus.Config{
				Logger: helper.MakeLogger(
					helper.Setup{
						Env:     obj.Config.LogEnv,
						Logname: obj.Config.LogPath + "/access_log",
						Display: true,
						Level:   obj.Config.LogLevel,
					}),
				Tags: []string{
					fiberlogrus.TagIP,
					fiberlogrus.TagIPs,
					fiberlogrus.TagProtocol,
					fiberlogrus.TagHost,
					fiberlogrus.TagPort,
					fiberlogrus.TagMethod,
					fiberlogrus.TagPath,
					fiberlogrus.TagURL,
					fiberlogrus.TagUA,
					fiberlogrus.TagBody,
					fiberlogrus.TagRoute,
					fiberlogrus.TagQueryStringParams,
					fiberlogrus.TagStatus,
					fiberlogrus.TagPid,
					fiberlogrus.TagReferer,
					fiberlogrus.TagLatency,
				},
			}))

	h := handler.NewIncomingHandler(handler.IncomingHandler{
		Config: obj.Config,
		Logs:   obj.Logs,
		R:      obj.R,
		DB:     obj.DB,
		Rmqp:   obj.Rmqp,
	})

	// V1
	v1 := f.Group("/v1") // v1

	// Postback
	v1.Get("/postback/:urlservicekey/", h.Postback)
	v1.Get("/postback", h.Postback2)

	// Report
	rpt := v1.Group("/report") // Report
	//rpt.Get("/report/", h.Report).Name("Report API")
	rpt.Get("/pinreport", h.DisplayPinReport).Name("Pin Report Summary FE")
	rpt.Get("/datasentapiperformance", h.DisplayPinPerformanceReport).Name("Pin Performance Api Report Summary FE")
	rpt.Get("/cpareportlist", h.AuthMiddleware, h.DisplayCPAReport).Name("Receive Pin CPA Report Transactional")
	rpt.Get("/costreport/:v", h.AuthMiddleware, h.DisplayCostReport).Name("Receive Pin Cost Report / detail Transactional")
	rpt.Get("/conversionlog", h.DisplayConversionLogReport).Name("Conversion Log Report")
	rpt.Get("/campaign-monitoring-summary", h.DisplayCampaignSummary).Name("Campaign Summary")
	rpt.Get("/campaign-monitoring-summary/chart", h.DisplayCampaignSummaryChart).Name("Campaign Summary Chart")
	rpt.Get("/alertreport/:v", h.DisplayAlertReportAll).Name("All Alert Report list/")
	rpt.Get("/trafficreport", h.DisplayTrafficReport).Name("Traffic Report list")
	rpt.Get("/mainstreamreport", h.AuthMiddleware, h.DisplayMainstreamReport).Name("Mainstream Report list")
	rpt.Get("/budgetmonitoring", h.DisplayBudgetMonitoring).Name("Budget Monitoring list")
	rpt.Get("/performance-report", h.DisplayPerformanceReport).Name("Performance Report list")
	rpt.Get("/revenuemonitoring", h.DisplayRevenueMonitoring).Name("Revenue Monitoring list")
	rpt.Get("/revenuemonitoring/chart", h.DisplayRevenueMonitoringChart).Name("Revenue Monitoring chart")
	rpt.Get("/defaultinput/", h.DisplayDefaultInput).Name("Default Input for cpa n mainstream")
	rpt.Get("/redirectiontime", h.DisplayRedirectionTime).Name("Redirection Time")

	// API Internal
	internal := v1.Group("/int") // Internal API
	internal.Put("/setdata/:v/", h.SetData).Name("SetTargetDailyBudget")
	internal.Put("/updatedata/:v/", h.UpdateAgencyFeeAndCostConversion).Name("UpdateAgencyFeeAndCostConversion")
	internal.Put("/updateratio/:v/", h.UpdateRatio).Name("Update Ratio Transactional")
	internal.Put("/updatepostback/:v/", h.UpdatePostback).Name("Update Postback Transactional")
	internal.Put("/updateagencycost/:v", h.UpdateAgencyCost).Name("Update Agency fee and cost per conversion in db")
	internal.Put("/updatestatusalert/:v", h.UpdateStatusAlert).Name("Update Status Alert in db")
	internal.Get("/pinreport/", h.TrxPinReport).Name("Receive Pin Report Transactional")
	internal.Get("/datasentapiperformance/", h.TrxPerformancePinReport).Name("Receive Pin API Performance Report Transactional")
	internal.Get("/exportcpa/", h.AuthMiddleware, h.ExportCpaButton).Name("Export CPA-Report Button")
	internal.Get("/exportcost/", h.ExportCostButton).Name("Export Cost-Report Button")
	internal.Get("/exportcostdetail/", h.ExportCostDetailButton).Name("Export Cost-Report-Detail Button")
	internal.Get("/pinperformance", h.PinPerformance).Name("Receive Pin Performance Report Transactional")
	internal.Post("/uploadexcel", h.UploadExcel).Name("Upload Excal SMS Campaign")
	internal.Put("/updateexcel/:id", h.UpdateExcel).Name("Update Excal SMS Campaign")
	internal.Put("/upsertexcel/", h.UpsertExcel).Name("Upsert Excal SMS Campaign")
	internal.Get("/getdataarpu/", h.GetDataArpu).Name("Get Data ARPU")

	// Management
	management := v1.Group("/management") // Management
	// Campaign
	campaign := management.Group("/campaign") // Campaign
	campaign.Get("/", h.DisplayCampaignManagement).Name("Campaign Management FE")
	campaign.Get("/campaigncounts", h.GetCampaignCounts).Name("Campaign Management Campaign Counts FE")
	campaign.Get("/:v", h.DisplayCampaignManagement).Name("Campaign Management Detail FE")
	campaign.Post("/send", h.SendCampaignHandler).Name("Campaign Management Send FE")
	campaign.Post("/updatestatus", h.UpdateStatusCampaign).Name("Update status campaign on campaign_details")
	campaign.Post("/editcampaign", h.EditCampaign).Name("Edit capping campaign on campaign_details")
	campaign.Post("/delcampaign", h.DelCampaign).Name("Edit capping campaign on campaign_details")
	campaign.Post("/updatekeymainstream", h.UpdateKeyMainstream).Name("Update key mainstream campaign on campaign_details")
	campaign.Post("/updategooglesheet", h.UpdateGoogleSheet).Name("Update google sheet campaign on campaign_details")
	campaign.Post("/editmocappingservices2s", h.EditMOCappingServiceS2S).Name("Update mocappingservices2s campaign on campaign_details")

	// Menu
	menu := management.Group("/menu") // Menu
	menu.Post("/", h.CreateMenu).Name("Menu Management Create FE")
	menu.Get("/", h.GetAllMenus).Name("Menu Management FE")
	menu.Get("/:id", h.GetMenuByID).Name("Menu Management Edit FE")
	menu.Put("/:id", h.UpdateMenu).Name("Menu Management Update FE")
	menu.Delete("/:id", h.DeleteMenu).Name("Menu Management Delete FE")
	// role
	role := management.Group("/role") // role
	role.Post("/", h.CreateRole).Name("Role Management Create FE")
	role.Get("/", h.GetRoleTable).Name("Role Management FE")
	role.Put("/:id", h.UpdateRole).Name("Role Management Update FE")
	role.Delete("/:id", h.DeleteRole).Name("Role Management Delete FE")
	// user
	user := management.Group("/user") // uset
	user.Post("/", h.CreateUser).Name("User Management Create FE")
	user.Get("/", h.GetUserTable).Name("User Management FE")
	user.Get("/usercounts", h.GetUserCounts).Name("User Management User Counts FE")
	user.Put("/:id", h.UpdateUser).Name("User Management Update FE")
	user.Put("/assignservice/:id", h.AssignService).Name("User Management Assign Service & Adnet FE")
	user.Put("/updatestatus/:id", h.UpdateUserStatus).Name("User Management Update Status FE")
	user.Delete("/:id", h.DeleteUser).Name("User Management Delete FE")
	user.Get("/approvalrequest", h.GetUserApplovalRequestTable).Name("User Management Approval Request FE")
	user.Put("/approveuser/:id", h.ApproveUser).Name("User Management Approve User FE")

	// User Log
	userlog := management.Group("/userlog")
	userlog.Get("/ip", h.GetIpAddress).Name("User Log IP")
	userlog.Post("/", h.CreateUserLog).Name(" Save User Log List")
	userlog.Get("/", h.DisplayUserLogList).Name(" Display User Log List")
	userlog.Get("/:id", h.DisplayUserLogHistory).Name(" Display User Log History")

	//  Country and Service Management
	countryService := management.Group("/country-service")
	countryService.Get("/email", h.DisplayEmail).Name("Display Email")
	countryService.Get("/email/:id", h.DisplayEmailByID).Name("Display Email By ID")
	countryService.Post("/email", h.CreateEmail).Name("Create Email")
	countryService.Put("/email/:id", h.UpdateEmail).Name("Update Email")
	countryService.Delete("/email/:id", h.DeleteEmail).Name("Delete Email")
	countryService.Get("/country", h.DisplayCountry).Name("Create Country")
	countryService.Post("/country", h.CreateCountry).Name("Create Country")
	countryService.Put("/country/:id", h.UpdateCountry).Name("Update Country")
	countryService.Delete("/country/:id", h.DeleteCountry).Name("Delete Country")
	countryService.Get("/company", h.DisplayCompany).Name("Create Company")
	countryService.Post("/company", h.CreateCompany).Name("Create Company")
	countryService.Put("/company/:id", h.UpdateCompany).Name("Update Company")
	countryService.Delete("/company/:id", h.DeleteCompany).Name("Delete Company")
	countryService.Get("/domain", h.DisplayDomain).Name("Create Domain")
	countryService.Post("/domain", h.CreateDomain).Name("Create Domain")
	countryService.Put("/domain/:id", h.UpdateDomain).Name("Update Domain")
	countryService.Delete("/domain/:id", h.DeleteDomain).Name("Delete Domain")
	countryService.Get("/operator", h.DisplayOperator).Name("Create Operator")
	countryService.Post("/operator", h.CreateOperator).Name("Create Operator")
	countryService.Put("/operator/:id", h.UpdateOperator).Name("Update Operator")
	countryService.Delete("/operator/:id", h.DeleteOperator).Name("Delete Operator")
	countryService.Get("/partner", h.DisplayPartner).Name("Create Partner")
	countryService.Post("/partner", h.CreatePartner).Name("Create Partner")
	countryService.Put("/partner/:id", h.UpdatePartner).Name("Update Partner")
	countryService.Delete("/partner/:id", h.DeletePartner).Name("Delete Partner")
	countryService.Get("/service", h.DisplayService).Name("Create Service")
	countryService.Post("/service", h.CreateService).Name("Create Service")
	countryService.Put("/service/:id", h.UpdateService).Name("Update Service")
	countryService.Delete("/service/:id", h.DeleteService).Name("Delete Service")
	countryService.Get("/adnet-list", h.DisplayAdnetList).Name("Create AdnetList")
	countryService.Post("/adnet-list", h.CreateAdnetList).Name("Create AdnetList")
	countryService.Put("/adnet-list/:id", h.UpdateAdnetList).Name("Update AdnetList")
	countryService.Delete("/adnet-list/:id", h.DeleteAdnetList).Name("Delete AdnetList")
	countryService.Get("/agency", h.DisplayAgency).Name("Show Agency")
	countryService.Post("/agency", h.CreateAgency).Name("Create Agency")
	countryService.Put("/agency/:id", h.UpdateAgency).Name("Update Agency")
	countryService.Delete("/agency/:id", h.DeleteAgency).Name("Delete Agency")
	countryService.Get("/channel", h.DisplayChannel).Name("Show Channel")
	countryService.Post("/channel", h.CreateChannel).Name("Create Channel")
	countryService.Put("/channel/:id", h.UpdateChannel).Name("Update Channel")
	countryService.Delete("/channel/:id", h.DeleteChannel).Name("Delete Channel")
	countryService.Get("/mainstream-group", h.DisplayMainstreamGroup).Name("Show MainStreamGroup")
	countryService.Post("/mainstream-group", h.CreateMainstreamGroup).Name("Create MainStreamGroup")
	countryService.Put("/mainstream-group/:id", h.UpdateMainstreamGroup).Name("Update MainStreamGroup")
	countryService.Delete("/mainstream-group/:id", h.DeleteMainstreamGroup).Name("Delete MainStreamGroup")

	ipRange := management.Group("/ipranges")
	ipRange.Get("/", h.GetIPRangeFiles).Name(" Display IP Ranges List List")
	ipRange.Post("/upload", h.UploadIPRangeRows).Name("Upload IP Ranges CSV")
	ipRange.Post("/implement", h.ImplementIPRange).Name("Implement IP Ranges")
	ipRange.Post("/download", h.DownloadIPRangeCSV).Name("Download IP Ranges")

	// API External
	v1.Group("/ext") // External API

	return f
}
