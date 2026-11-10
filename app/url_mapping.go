package app

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/infraLinkit/mediaplatform-datasource/handler"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	_ "github.com/lib/pq"
	"github.com/mikhail-bigun/fiberlogrus"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"
)

type App3rdParty struct {
	Config *config.Cfg
	Logs   *logrus.Logger
	DB     *gorm.DB
	R      *rueidis.Storage
	RCP    *redis.Client
	Rmqp   rmqp.AMQP
	RM     *helper.RabbitManager
	GS     *sheets.Service
}

func MapUrls(obj App3rdParty) *fiber.App {

	f := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		BodyLimit:   100 * 1024 * 1024,
	})

	// Only needed when the React SPA is served from another origin.
	if origins := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")); origins != "" {
		f.Use(cors.New(cors.Config{
			AllowOrigins: origins,
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
			AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			// The IP Ranges download names its file in this header.
			ExposeHeaders: "Content-Disposition",
		}))
	}

	f.Static("/storage", "./storage")

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
		RCP:    obj.RCP,
		DB:     obj.DB,
		Rmqp:   obj.Rmqp,
		RM:     obj.RM,
		GS:     obj.GS,
	})

	// V1
	v1 := f.Group("/v1") // v1

	// Auth
	auth := v1.Group("/auth")
	authTooMany := func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too many attempts. Please try again in a minute."})
	}
	// Same as LoginRequest in Laravel: 5 failures per login name and IP. Counters are per instance.
	loginLimiter := limiter.New(limiter.Config{
		Max:                    5,
		Expiration:             time.Minute,
		SkipSuccessfulRequests: true,
		LimitReached:           authTooMany,
		KeyGenerator: func(c *fiber.Ctx) string {
			var body handler.LoginRequest
			_ = c.BodyParser(&body)
			name := body.Email
			if name == "" {
				name = body.Username
			}
			return "login|" + strings.ToLower(name) + "|" + c.IP()
		},
	})
	publicAuthLimiter := limiter.New(limiter.Config{Max: 10, Expiration: time.Minute, LimitReached: authTooMany})
	auth.Get("/config", h.GetAuthConfig)
	auth.Post("/login", loginLimiter, h.Login)
	auth.Post("/register", publicAuthLimiter, h.Register)
	auth.Post("/forgot-password", publicAuthLimiter, h.ForgotPassword)
	auth.Post("/reset-password", publicAuthLimiter, h.ResetPassword)
	auth.Post("/refresh", h.RefreshToken)
	auth.Get("/me", h.AuthMiddleware, h.GetMe)
	auth.Get("/profile", h.AuthMiddleware, h.GetProfile)
	auth.Put("/profile", h.AuthMiddleware, h.UpdateProfile)

	dashboard := f.Group("/dashboard", h.AuthGate)
	dashboard.Get("get-data", h.AuthMiddleware, h.DisplayDashboardData)
	dashboard.Get("get-top-campaign", h.AuthMiddleware, h.DisplayDashboardTopCampaign)
	dashboard.Get("get-top-partner-spend", h.AuthMiddleware, h.DisplayDashboardTopPartnerSpend)
	dashboard.Get("get-report-list", h.AuthMiddleware, h.DisplayDashboardReport)
	dashboard.Get("get-country-stats", h.AuthMiddleware, h.DisplayCountryStats)
	dashboard.Get("get-ops-stats", h.AuthMiddleware, h.DisplayOpsStats)
	dashboard.Get("get-alerts", h.AuthMiddleware, h.DisplayAlerts)
	dashboard.Get("get-rollup", h.AuthMiddleware, h.DisplayRollup)
	dashboard.Get("get-campaign-hierarchy", h.AuthMiddleware, h.DisplayCampaignHierarchy)
	dashboard.Get("get-adnet-stats", h.AuthMiddleware, h.DisplayAdnetStats)
	dashboard.Get("get-heatmap", h.AuthMiddleware, h.DisplayHeatmap)
	dashboard.Get("get-campaign-daily", h.AuthMiddleware, h.DisplayCampaignDaily)
	dashboard.Get("get-service-daily", h.AuthMiddleware, h.DisplayServiceDaily)
	dashboard.Get("get-filter-options", h.AuthMiddleware, h.DisplayFilterOptions)

	// Postback
	v1.Get("/postback/:urlservicekey/", h.Postback)
	v1.Get("/postback", h.PostbackV3)
	v1.Get("/postback_billed", h.PostbackBilled)
	v1.Get("/postback_sync", h.PostbackDirectReply)
	v1.Get("/inquire/campid", h.InquiryCampID)
	v1.Get("/inquire/api-campid", h.InquiryAPICampID)

	// Report
	rpt := v1.Group("/report", h.AuthGate) // Report
	//rpt.Get("/report/", h.Report).Name("Report API")
	rpt.Get("/pinreport", h.AuthMiddleware, h.DisplayPinReport).Name("Pin Report Summary FE")
	rpt.Get("/datasentapiperformance", h.DisplayPinPerformanceReport).Name("Pin Performance Api Report Summary FE")
	rpt.Get("/cpareportlist", h.AuthMiddleware, h.DisplayCPAReport).Name("Receive Pin CPA Report Transactional")
	rpt.Get("/costreport/:v", h.AuthMiddleware, h.DisplayCostReport).Name("Receive Pin Cost Report / detail Transactional")
	rpt.Get("/conversionlog", h.DisplayConversionLogReport).Name("Conversion Log Report")
	rpt.Get("/campaign-monitoring-summary", h.DisplayCampaignSummary).Name("Campaign Summary")
	rpt.Get("/campaign-monitoring-summary/chart", h.DisplayCampaignSummaryChart).Name("Campaign Summary Chart")
	rpt.Get("/alertreport/:v", h.DisplayAlertReportAll).Name("All Alert Report list/")
	rpt.Get("/alert-dashboard", h.DisplayAlertDashboard).Name("Alert Dashboard list")
	rpt.Post("/alert-dashboard/change-status", h.AuthMiddleware, h.ChangeAlertStatus).Name("Alert Dashboard change status")
	rpt.Get("/trafficreport", h.DisplayTrafficReport).Name("Traffic Report list")
	rpt.Get("/trafficreport/chart", h.GetTrafficReportChart).Name("Traffic Report chart")
	rpt.Get("/mainstreamreport", h.AuthMiddleware, h.DisplayMainstreamReport).Name("Mainstream Report list")
	rpt.Get("/google-traffic-report", h.AuthMiddleware, h.DisplayGoogleTrafficReport).Name("Google Traffic Report")
	rpt.Get("/budgetmonitoring", h.DisplayBudgetMonitoring).Name("Budget Monitoring list")
	rpt.Get("/performance-report", h.DisplayPerformanceReport).Name("Performance Report list")
	rpt.Get("/revenuemonitoring", h.DisplayRevenueMonitoring).Name("Revenue Monitoring list")
	rpt.Get("/revenuemonitoring/chart", h.DisplayRevenueMonitoringChart).Name("Revenue Monitoring chart")
	rpt.Get("/defaultinput/", h.DisplayDefaultInput).Name("Default Input for cpa n mainstream")
	rpt.Get("/redirectiontime", h.DisplayRedirectionTime).Name("Redirection Time")
	rpt.Post("/resend-data", h.ResendData).Name("Resend Data")
	rpt.Post("/resend-data-apireport", h.ResendDataAPIReport).Name("Resend Data API Report")
	rpt.Post("/campaign-monitoring-summary/edit-target-budget", h.EditTargetBudgetLevel).Name("Edit Target Budget")
	rpt.Post("/campaign-monitoring-summary/edit-target-budget-batch", h.EditTargetBudgetBatch).Name("Edit Target Budget Batch")
	rpt.Put("/campaign-monitoring-summary/tdb", h.AuthMiddleware, h.SetTargetDailyBudget).Name("Set Target Daily Budget")
	rpt.Get("/ioreport", h.DisplaySummaryBudgetIO).Name("IO Report")
	rpt.Put("/ioreport/update", h.UpdateSummaryBudgetIO).Name("IO Report Update")

	rpt.Get("/campaign-spending-channel", h.DisplayCampaignSpendingChannel).
		Name("Campaign Spending Channel")

	rpt.Get("/campaign-spending-channel/country-children", h.DisplayCampaignSpendingChannelCountryChildren).
		Name("Campaign Spending Channel Country Children")

	// API Internal
	internal := v1.Group("/int", h.AuthGate) // Internal API
	internal.Put("/setdata/:v/", h.SetData).Name("SetTargetDailyBudget")
	internal.Put("/updatedata/:v/", h.UpdateAgencyFeeAndCostConversion).Name("UpdateAgencyFeeAndCostConversion")
	internal.Put("/updateratio/:v/", h.UpdateRatio).Name("Update Ratio Transactional")
	internal.Put("/updatepostback/:v/", h.UpdatePostback).Name("Update Postback Transactional")
	internal.Put("/updateagencycost/:v", h.UpdateAgencyCost).Name("Update Agency fee and cost per conversion in db")
	internal.Put("/updatestatusalert/:v", h.UpdateStatusAlert).Name("Update Status Alert in db")
	internal.Get("/datasentapipinreport/", h.TrxPinReport).Name("Receive Pin Report Transactional")
	internal.Get("/uploadsmsapipinreport/", h.TrxUploadSmsApiPinReport).Name("Receive Upload SMS API Pin Report (requires existing campaign)")
	internal.Post("/pinreport/editpayout", h.EditPayoutAPIReport).Name("Edit payout api report")
	internal.Get("/datasentapiperformance/", h.TrxPerformancePinReport).Name("Receive Pin API Performance Report Transactional")
	internal.Post("/pinperformance/editcpa", h.EditCpaAPIPerformanceReport).Name("Edit cpa api performance report")
	internal.Post("/pinperformance/editarpu", h.EditArpuAPIPerformanceReport).Name("Edit arpu api performance report")
	internal.Get("/exportcpa/", h.AuthMiddleware, h.ExportCpaButton).Name("Export CPA-Report Button")
	internal.Get("/exportcost/", h.ExportCostButton).Name("Export Cost-Report Button")
	internal.Get("/exportcostdetail/", h.ExportCostDetailButton).Name("Export Cost-Report-Detail Button")
	internal.Get("/pinperformance", h.PinPerformance).Name("Receive Pin Performance Report Transactional")
	internal.Post("/uploadexcel", h.UploadExcel).Name("Upload Excal SMS Campaign")
	internal.Put("/updateexcel/:id", h.UpdateExcel).Name("Update Excal SMS Campaign")
	internal.Post("/uploadsmsexcel", h.AuthMiddleware, h.UploadSMSExcel).Name("Upload SMS Excel Rows")
	internal.Put("/upsertexcel/", h.UpsertExcel).Name("Upsert Excal SMS Campaign")
	internal.Get("/getdataarpu/", h.GetDataArpu).Name("Get Data ARPU")
	internal.Get("/get_urlservice_in_summarylanding", h.GetURLServiceInSummaryLanding).Name("Get URL service for total load time")
	internal.Put("/update_response_url_service_in_summarylanding", h.UpdateResponseURLServiceInSummaryLanding).Name("Get URL service for total load time")
	internal.Put("/updatepayout", h.AuthMiddleware, h.UpdateRealPayout).Name("Update Real Payout")

	// Management
	management := v1.Group("/management", h.AuthGate) // Management
	// Campaign
	campaign := management.Group("/campaign") // Campaign
	campaign.Post("/bulk-update-status", h.RequirePermission("campaign-management.bulk-update-status"), h.BulkUpdateCampaignStatus).Name("Campaign Management Bulk Status")
	campaign.Delete("/bulk-delete", h.RequirePermission("campaign-management.bulk-delete"), h.BulkDeleteCampaigns).Name("Campaign Management Bulk Delete")
	campaign.Get("/", h.RequirePermission("campaign-management.index"), h.DisplayCampaignManagement).Name("Campaign Management FE")
	campaign.Get("/campaigncounts", h.RequirePermission("campaign-management.index"), h.GetCampaignCounts).Name("Campaign Management Campaign Counts FE")
	campaign.Get("/options", h.AuthMiddleware, h.GetCampaignOptions).Name("Campaign Management Options FE")
	campaign.Get("/lp-design-types", h.AuthMiddleware, h.GetLPDesignTypes).Name("Get LP Design Types")
	campaign.Get("/email-templates", h.RequirePermission("campaign-management.show"), h.CampaignEmailTemplates).Name("Campaign Management Email Templates")
	campaign.Post("/send-email", h.RequirePermission("campaign-management.show"), h.SendCampaignEmail).Name("Campaign Management Send Email")
	campaign.Get("/:v", h.RequirePermission("campaign-management.show", "campaign-management.edit"), h.DisplayCampaignManagement).Name("Campaign Management Detail FE")
	campaign.Get("/edit-data/:campaign_id", h.AuthMiddleware, h.RequirePermission("campaign-management.edit"), h.CampaignEditData).Name("Campaign Edit Data")
	campaign.Post("/update", h.AuthMiddleware, h.RequirePermission("campaign-management.edit"), h.UpdateCampaignForm).Name("Campaign Update Form")
	campaign.Post("/update-upload-sms", h.AuthMiddleware, h.RequirePermission("campaign-management.edit"), h.UpdateUploadSmsCampaign).Name("Campaign Update Upload SMS")
	campaign.Post("/update-single-url-s2s", h.AuthMiddleware, h.RequirePermission("campaign-management.edit"), h.UpdateSingleUrlS2S).Name("Campaign Update Single URL S2S")
	campaign.Post("/mainstream/store", h.AuthMiddleware, h.RequirePermission("campaign-mainstream.create"), h.StoreCampaignMainstream).Name("Campaign Mainstream Store")
	campaign.Post("/mainstream/store-single-url", h.AuthMiddleware, h.RequirePermission("campaign-mainstream.create"), h.StoreCampaignMainstreamSingleUrl).Name("Campaign Mainstream Store Single URL")
	campaign.Post("/mainstream/update", h.AuthMiddleware, h.RequirePermission("campaign-management.edit"), h.UpdateCampaignMainstream).Name("Campaign Mainstream Update")
	campaign.Post("/mainstream/update-single-url", h.AuthMiddleware, h.RequirePermission("campaign-management.edit"), h.UpdateCampaignMainstreamSingleUrl).Name("Campaign Mainstream Update Single URL")
	campaign.Post("/store", h.AuthMiddleware, h.RequirePermission("campaign-management.create"), h.StoreCampaign).Name("Campaign Management Store")
	campaign.Post("/store-single-url-s2s", h.AuthMiddleware, h.RequirePermission("campaign-management.create"), h.StoreSingleUrlS2S).Name("Campaign Management Store Single URL S2S")
	campaign.Post("/send", h.RequirePermission("campaign-management.create", "campaign-mainstream.create"), h.SendCampaignHandler).Name("Campaign Management Send FE")
	campaign.Post("/updatestatus", h.RequirePermission("campaign-management.bulk-update-status", "campaign-management.edit"), h.UpdateStatusCampaign).Name("Update status campaign on campaign_details")
	campaign.Post("/editcampaign", h.EditCampaign).Name("Edit capping campaign on campaign_details")
	campaign.Post("/editmocapping", h.EditCampaignMOCapping).Name("Edit mocapping campaign on campaign_details")
	campaign.Post("/editratio", h.EditCampaignRatio).Name("Edit ratio campaign on campaign_details")
	campaign.Post("/editpo", h.EditCampaignPO).Name("Edit postback campaign on campaign_details")
	campaign.Post("/delcampaign", h.RequirePermission("campaign-management.bulk-delete"), h.DelCampaign).Name("Edit capping campaign on campaign_details")
	campaign.Post("/updatekeymainstream", h.UpdateKeyMainstream).Name("Update key mainstream campaign on campaign_details")
	campaign.Post("/updategooglesheet", h.UpdateGoogleSheet).Name("Update google sheet campaign on campaign_details")
	campaign.Post("/updategooglesheetbillable", h.UpdateGoogleSheetBillable).Name("Update google sheet billable campaign on campaign_details")
	campaign.Post("/editmocappingservices2s", h.RequirePermission("setting-campaign.store"), h.EditMOCappingServiceS2S).Name("Update mocappingservices2s campaign on campaign_details")
	campaign.Post("/editpoaf", h.EditPOAF).Name("Edit poaf")
	campaign.Post("/editcampaignmanagementdetail", h.RequirePermission("setting-campaign.store", "campaign-management.edit"), h.EditCampaignManagementDetail).Name("Edit campaign on campaign_details")
	campaign.Post("/updatecampaignmanagement", h.UpdateCampaign).Name("Edit campaign on form")
	campaign.Post("/store-lp-design", h.AuthMiddleware, h.RequireRoleCode("super_admin"), h.StoreLPDesign).Name("Store LP Design Type")
	campaign.Delete("/lp-design-types/:id", h.AuthMiddleware, h.RequireRoleCode("super_admin"), h.DeleteLPDesignType).Name("Delete LP Design Type")

	// Campaign Setting
	campaignSetting := management.Group("/campaign-setting") // Campaign Setting
	campaignSetting.Get("/adnet-level", h.RequirePermission("setting-campaign.index"), h.CampaignSettingAdnetLevel).Name("Campaign Setting Adnet Level")
	campaignSetting.Get("/service-level", h.RequirePermission("setting-campaign.index"), h.CampaignSettingServiceLevel).Name("Campaign Setting Service Level")
	campaignSetting.Post("/editratio", h.RequirePermission("setting-campaign.store", "campaign-management.edit", "cpa-report.edit", "mainstream-report.edit"), h.EditCampaignSettingRatio).Name("Edit ratio campaign on campaign_setting")
	campaignSetting.Post("/editpo", h.RequirePermission("setting-campaign.store", "campaign-management.edit", "cpa-report.edit", "mainstream-report.edit"), h.EditCampaignSettingPO).Name("Edit po campaign on campaign_setting")
	campaignSetting.Post("/editmocapping", h.RequirePermission("setting-campaign.store", "campaign-management.edit", "cpa-report.edit", "mainstream-report.edit"), h.EditCampaignSettingMOCapping).Name("Edit mocapping campaign on campaign_setting")

	// Campaign Mainstream Setting
	campaignMainstreamSetting := management.Group("/campaign-setting-mainstream", h.AuthMiddleware)
	campaignMainstreamSetting.Get("/list", h.RequirePermission("setting-campaign-mainstream.index"), h.CampaignSettingMainstream).Name("Campaign Setting Mainstream List")
	campaignMainstreamSetting.Get("/lp-design-types", h.GetLPDesignTypes).Name("Get LP Design Types Mainstream")
	campaignMainstreamSetting.Post("/store-lp-design", h.RequireRoleCode("super_admin"), h.StoreLPDesign).Name("Store LP Design Type Mainstream")
	campaignMainstreamSetting.Post("/store", h.RequirePermission("setting-campaign-mainstream.index"), h.UpdateKeyMainstream).Name("Store Key Mainstream")
	campaignMainstreamSetting.Post("/store-google-sheet", h.RequirePermission("setting-campaign-mainstream.index"), h.UpdateGoogleSheet).Name("Store Google Sheet Mainstream")
	campaignMainstreamSetting.Post("/store-google-sheet-billable", h.RequirePermission("setting-campaign-mainstream.index"), h.UpdateGoogleSheetBillable).Name("Store Google Sheet Billable Mainstream")

	// Generic Upload
	v1.Post("/upload", h.AuthMiddleware, h.UploadFile).Name("Generic File Upload V1")
	management.Post("/upload", h.AuthMiddleware, h.UploadFile).Name("Generic File Upload Management")

	// Menu
	menu := management.Group("/menu") // Menu
	menu.Get("/sidebar", h.AuthMiddleware, h.GetSidebarMenu).Name("Menu Sidebar Dynamic FE")
	menu.Post("/", h.RequirePermission("menu-management.create"), h.CreateMenu).Name("Menu Management Create FE")
	menu.Get("/", h.RequirePermission("menu-management.index", "menu-management.create", "menu-management.edit", "role-management.create", "role-management.edit"), h.GetAllMenus).Name("Menu Management FE")
	menu.Get("/:id", h.RequirePermission("menu-management.edit"), h.GetMenuByID).Name("Menu Management Edit FE")
	menu.Put("/:id", h.RequirePermission("menu-management.edit"), h.UpdateMenu).Name("Menu Management Update FE")
	menu.Delete("/:id", h.RequirePermission("menu-management.delete"), h.DeleteMenu).Name("Menu Management Delete FE")
	// role
	role := management.Group("/role") // role
	role.Post("/", h.RequirePermission("role-management.create"), h.CreateRole).Name("Role Management Create FE")
	role.Get("/", h.RequirePermission("role-management.index"), h.GetRoleTable).Name("Role Management FE")
	role.Get("/:id", h.AuthMiddleware, h.RequirePermission("role-management.edit"), h.GetRoleByID).Name("Role Management Detail FE")
	role.Put("/:id", h.RequirePermission("role-management.edit"), h.UpdateRole).Name("Role Management Update FE")
	role.Delete("/:id", h.RequirePermission("role-management.delete"), h.DeleteRole).Name("Role Management Delete FE")
	// user
	user := management.Group("/user") // uset
	user.Post("/", h.RequirePermission("user-management.create"), h.CreateUser).Name("User Management Create FE")
	user.Get("/", h.RequirePermission("user-management.index"), h.GetUserTable).Name("User Management FE")
	user.Get("/usercounts", h.RequirePermission("user-management.index"), h.GetUserCounts).Name("User Management User Counts FE")
	user.Get("/options", h.AuthMiddleware, h.RequirePermission("user-management.create", "user-management.edit", "user-management.show", "user-management.assign-service", "user-management.approval-request"), h.GetUserOptions).Name("User Management Options Dropdown FE")
	user.Get("/approvalrequest", h.RequirePermission("user-management.approval-request"), h.GetUserApplovalRequestTable).Name("User Management Approval Request FE")
	user.Get("/:id", h.AuthMiddleware, h.RequirePermission("user-management.edit", "user-management.show", "user-management.assign-service", "activity-log.history"), h.GetUserDetail).Name("User Management Detail FE")
	user.Put("/:id", h.RequirePermission("user-management.edit"), h.UpdateUser).Name("User Management Update FE")
	user.Put("/assignservice/:id", h.RequirePermission("user-management.assign-service"), h.AssignService).Name("User Management Assign Service & Adnet FE")
	user.Put("/updatestatus/:id", h.RequirePermission("user-management.index"), h.UpdateUserStatus).Name("User Management Update Status FE")
	user.Delete("/:id", h.RequirePermission("user-management.delete"), h.DeleteUser).Name("User Management Delete FE")
	user.Put("/approveuser/:id", h.RequirePermission("user-management.approval-request"), h.ApproveUser).Name("User Management Approve User FE")

	// User Log
	userlog := management.Group("/userlog")
	userlog.Get("/ip", h.GetIpAddress).Name("User Log IP")
	userlog.Post("/", h.CreateUserLog).Name(" Save User Log List")
	userlog.Get("/", h.RequirePermission("activity-log.index"), h.DisplayUserLogList).Name(" Display User Log List")
	userlog.Get("/:id", h.RequirePermission("activity-log.history"), h.DisplayUserLogHistory).Name(" Display User Log History")

	budgetio := management.Group("/budget-io")
	budgetio.Post("/", h.AuthMiddleware, h.RequirePermission("budget-io.create"), h.CreateBudgetIORequest).Name("Budget IO Create")
	budgetio.Get("/options", h.AuthMiddleware, h.RequirePermission("budget-io.create"), h.GetBudgetIOOptions).Name("Budget IO Options")
	budgetio.Get("/budgetiolist", h.AuthMiddleware, h.RequirePermission("budget-io.index"), h.ListBudgetIORequests).Name("Budget IO List")
	budgetio.Get("/budgetiolistall", h.AuthMiddleware, h.RequirePermission("budget-io.view-all-list"), h.ListAllBudgetIORequests).Name("Budget IO List All")
	budgetio.Get("/budgetioapproved", h.AuthMiddleware, h.ListApprovedBudgetIORequests).Name("Budget IO Approved")
	budgetio.Get("/budgetioapprovedall", h.AuthMiddleware, h.RequirePermission("budget-io.view-all-list"), h.ListAllApprovedBudgetIORequests).Name("Budget IO Approved All")
	budgetio.Get("/:id", h.AuthMiddleware, h.GetBudgetIORequest).Name("Budget IO Detail")
	budgetio.Put("/:id", h.AuthMiddleware, h.RequirePermission("budget-io.revise"), h.UpdateBudgetIORequest).Name("Budget IO Update")
	budgetio.Post("/:id/approve", h.AuthMiddleware, h.RequirePermission("budget-io.approve-reject"), h.ApproveBudgetIORequest).Name("Budget IO Approve")
	budgetio.Post("/:id/reject", h.AuthMiddleware, h.RequirePermission("budget-io.approve-reject"), h.RejectBudgetIORequest).Name("Budget IO Reject")

	//  Country and Service Management
	countryService := management.Group("/country-service")
	countryService.Get("/email", h.DisplayEmail).Name("Display Email")
	countryService.Get("/email/:id", h.DisplayEmailByID).Name("Display Email By ID")
	countryService.Post("/email", h.CreateEmail).Name("Create Email")
	countryService.Put("/email/:id", h.UpdateEmail).Name("Update Email")
	countryService.Delete("/email/:id", h.DeleteEmail).Name("Delete Email")
	countryService.Get("/country", h.DisplayCountry).Name("Create Country")
	countryService.Get("/country/:code", h.DisplayCountryInfo).Name("Country Information")
	countryService.Post("/country", h.RequireUnique(&entity.Country{}, "name", "code", "code2"), h.CreateCountry).Name("Create Country")
	countryService.Put("/country/:id", h.RequireUnique(&entity.Country{}, "name", "code", "code2"), h.UpdateCountry).Name("Update Country")
	countryService.Delete("/country/:id", h.DeleteCountry).Name("Delete Country")
	countryService.Get("/continent", h.DisplayCountry).Name("Display Continent")
	countryService.Get("/company/maskedid", h.NextCompanyMaskedID).Name("Next Company Masked ID")
	countryService.Get("/company-group/maskedid", h.NextCompanyGroupMaskedID).Name("Next Company Group Masked ID")
	countryService.Get("/company", h.DisplayCompany).Name("Create Company")
	countryService.Post("/company", h.RequireUnique(&entity.Company{}, "name"), h.CreateCompany).Name("Create Company")
	countryService.Put("/company/:id", h.RequireUnique(&entity.Company{}, "name"), h.UpdateCompany).Name("Update Company")
	countryService.Delete("/company/:id", h.DeleteCompany).Name("Delete Company")
	countryService.Get("/domain", h.DisplayDomain).Name("Create Domain")
	countryService.Post("/domain", h.RequireUnique(&entity.Domain{}, "url"), h.CreateDomain).Name("Create Domain")
	countryService.Put("/domain/:id", h.RequireUnique(&entity.Domain{}, "url"), h.UpdateDomain).Name("Update Domain")
	countryService.Delete("/domain/:id", h.DeleteDomain).Name("Delete Domain")
	countryService.Get("/operator", h.DisplayOperator).Name("Create Operator")
	countryService.Post("/operator", h.CreateOperator).Name("Create Operator")
	countryService.Put("/operator/:id", h.UpdateOperator).Name("Update Operator")
	countryService.Delete("/operator/:id", h.DeleteOperator).Name("Delete Operator")
	countryService.Get("/api-operator-list", h.DisplayAPIOperatorList).Name("Show API Operator List")
	countryService.Get("/operator-alias", h.DisplayOperatorAlias).Name("Display Operator Alias")
	countryService.Post("/operator-alias", h.CreateOperatorAlias).Name("Create Operator Alias")
	countryService.Put("/operator-alias/:id", h.UpdateOperatorAlias).Name("Update Operator Alias")
	countryService.Delete("/operator-alias/:id", h.DeleteOperatorAlias).Name("Delete Operator Alias")
	countryService.Get("/partner", h.DisplayPartner).Name("Create Partner")
	countryService.Post("/partner", h.CreatePartner).Name("Create Partner")
	countryService.Put("/partner/:id", h.UpdatePartner).Name("Update Partner")
	countryService.Delete("/partner/:id", h.DeletePartner).Name("Delete Partner")
	countryService.Get("/service", h.DisplayService).Name("Create Service")
	countryService.Post("/service", h.CreateService).Name("Create Service")
	countryService.Put("/service/:id", h.UpdateService).Name("Update Service")
	countryService.Delete("/service/:id", h.DeleteService).Name("Delete Service")
	countryService.Get("/api-service-list", h.DisplayAPIServiceList).Name("Show API Service List")
	countryService.Get("/adnet-list", h.DisplayAdnetList).Name("Create AdnetList")
	countryService.Get("/api-adnet-list", h.DisplayAPIAdnetList).Name("Show API AdnetList")
	countryService.Post("/adnet-list", h.RequireUnique(&entity.AdnetList{}, "code"), h.CreateAdnetList).Name("Create AdnetList")
	countryService.Put("/adnet-list/:id", h.UpdateAdnetList).Name("Update AdnetList")
	countryService.Delete("/adnet-list/:id", h.DeleteAdnetList).Name("Delete AdnetList")
	countryService.Put("/adnet-list/:id/full", h.AuthMiddleware, h.UpdateAdnetAndCampaigns).Name("Update Adnet And Campaigns")
	countryService.Get("/agency", h.DisplayAgency).Name("Show Agency")
	countryService.Post("/agency", h.RequireUnique(&entity.Agency{}, "name", "code"), h.CreateAgency).Name("Create Agency")
	countryService.Put("/agency/:id", h.RequireUnique(&entity.Agency{}, "name", "code"), h.UpdateAgency).Name("Update Agency")
	countryService.Delete("/agency/:id", h.DeleteAgency).Name("Delete Agency")
	countryService.Get("/channel", h.DisplayChannel).Name("Show Channel")
	countryService.Post("/channel", h.RequireUnique(&entity.Channel{}, "name"), h.CreateChannel).Name("Create Channel")
	countryService.Put("/channel/:id", h.RequireUnique(&entity.Channel{}, "name"), h.UpdateChannel).Name("Update Channel")
	countryService.Delete("/channel/:id", h.DeleteChannel).Name("Delete Channel")
	countryService.Get("/mainstream-group", h.DisplayMainstreamGroup).Name("Show MainStreamGroup")
	countryService.Post("/mainstream-group", h.RequireUnique(&entity.MainstreamGroup{}, "name", "unique_domain"), h.CreateMainstreamGroup).Name("Create MainStreamGroup")
	countryService.Put("/mainstream-group/:id", h.RequireUnique(&entity.MainstreamGroup{}, "name", "unique_domain"), h.UpdateMainstreamGroup).Name("Update MainStreamGroup")
	countryService.Delete("/mainstream-group/:id", h.DeleteMainstreamGroup).Name("Delete MainStreamGroup")
	countryService.Get("/domain-service", h.DisplayDomainService).Name("Show Domain Service")

	countryService.Post("/edit-adnet-dsp-status", h.UpdateDSPAdnetStatus).Name("Edit DSP Status")

	countryService.Get("/company-group", h.DisplayCompanyGroup).Name("Display Company Group")
	countryService.Post("/company-group", h.RequireUnique(&entity.CompanyGroup{}, "name"), h.CreateCompanyGroup).Name("Create Company Group")
	countryService.Put("/company-group/:id", h.RequireUnique(&entity.CompanyGroup{}, "name"), h.UpdateCompanyGroup).Name("Update Company Group")
	countryService.Delete("/company-group/:id", h.DeleteCompanyGroup).Name("Delete Company Group")

	ipRange := management.Group("/ipranges")
	ipRange.Get("/", h.GetIPRangeFiles).Name(" Display IP Ranges List List")
	ipRange.Post("/upload", h.UploadIPRangeRows).Name("Upload IP Ranges CSV")
	ipRange.Post("/implement", h.ImplementIPRange).Name("Implement IP Ranges")
	ipRange.Post("/download", h.DownloadIPRangeCSV).Name("Download IP Ranges")

	// API External
	ext := v1.Group("/ext", h.ExtAPIKeyMiddleware) // External API
	ext.Get("/campaign-details", h.GetExtCampaignDetails).Name("External Campaign Summary (campaign_id, service, shortcode, adnet, cost)")

	return f
}
