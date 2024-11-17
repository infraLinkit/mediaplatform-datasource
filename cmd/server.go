package cmd

import (
	"log"

	"github.com/infraLinkit/mediaplatform-datasource/app"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Webserver CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.InitCfg()
		c := cfg.Initiate("api")

		c.Rmqp.SetUpChannel("direct", true, cfg.RabbitMQPixelStorageExchangeName, true, cfg.RabbitMQPixelStorageQueueName)

		router := app.MapUrls(app.App3rdParty{
			Config: cfg,
			Logs:   c.Logs,
			PS:     c.DB,
			R:      c.R,
			Rmqp:   c.Rmqp,
		})

		log.Fatal(router.Listen(":81"))
	},
}
