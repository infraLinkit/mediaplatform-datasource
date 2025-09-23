package cmd

import (
	"log"

	"github.com/infraLinkit/mediaplatform-datasource/app"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Webserver CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.InitCfg()
		c := cfg.Initiate("api")

		// Migrate table
		c.DB.AutoMigrate(&entity.Campaign{}, &entity.CampaignDetail{}, &entity.MO{}, &entity.PixelStorage{}, &entity.Postback{}, &entity.SummaryCampaign{}, &entity.DataClicked{}, &entity.DataLanding{}, &entity.DataRedirect{}, &entity.DataTraffic{}, &entity.ApiPinReport{}, &entity.ApiPinPerformance{}, &entity.Menu{}, &entity.Country{}, &entity.Company{}, &entity.Domain{}, &entity.Operator{}, &entity.Partner{}, &entity.Service{}, &entity.AdnetList{}, &entity.SummaryMo{}, &entity.SummaryCr{}, &entity.SummaryCapping{}, &entity.SummaryRatio{}, &entity.Role{}, &entity.Permission{}, &entity.User{}, &entity.CcEmail{}, &entity.Email{}, &entity.DetailUser{}, &entity.UserAdnet{}, &entity.LpDesignType{}, &entity.Agency{}, &entity.Channel{}, &entity.MainstreamGroup{}, &entity.SummaryLanding{}, &entity.IPRangeCsvRow{}, &entity.IPRange{}, &entity.IncSummaryCampaign{}, &entity.IncSummaryCampaignHour{}, &entity.SummaryTraffic{}, &entity.UserCompany{})

		c.Rmqp.SetUpChannel("direct", true, cfg.RabbitMQPixelStorageExchangeName, true, cfg.RabbitMQPixelStorageQueueName)
		c.Rmqp.SetUpChannel("direct", true, cfg.RabbitMQRatioExchangeName, true, cfg.RabbitMQRatioQueueName)
		c.Rmqp.SetUpChannel("direct", true, cfg.RabbitMQCampaignManagementExchangeName, true, cfg.RabbitMQCampaignManagementQueueName)

		router := app.MapUrls(app.App3rdParty{
			Config: cfg,
			Logs:   c.Logs,
			DB:     c.DB,
			R:      c.R,
			RCP:    c.RCP,
			Rmqp:   c.Rmqp,
		})

		log.Fatal(router.Listen(":" + c.Config.AppApiPort))
	},
}
