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
		c.DB.AutoMigrate(&entity.Campaign{}, &entity.CampaignDetail{}, &entity.MO{}, &entity.PixelStorage{}, &entity.Postback{}, &entity.SummaryCampaign{}, &entity.DataClicked{}, &entity.DataLanding{}, &entity.DataRedirect{}, &entity.DataTraffic{}, &entity.ApiPinReport{}, &entity.ApiPinPerformance{})

		c.Rmqp.SetUpChannel("direct", true, cfg.RabbitMQPixelStorageExchangeName, true, cfg.RabbitMQPixelStorageQueueName)
		c.Rmqp.SetUpChannel("direct", true, cfg.RabbitMQRatioExchangeName, true, cfg.RabbitMQRatioQueueName)
		c.Rmqp.SetUpChannel("direct", true, cfg.RabbitMQCampaignManagementExchangeName, true, cfg.RabbitMQCampaignManagementQueueName)

		router := app.MapUrls(app.App3rdParty{
			Config: cfg,
			Logs:   c.Logs,
			DB:     c.DB,
			R0:     c.R0,
			R1:     c.R1,
			Rmqp:   c.Rmqp,
		})

		log.Fatal(router.Listen(":" + c.Config.AppApiPort))
	},
}
