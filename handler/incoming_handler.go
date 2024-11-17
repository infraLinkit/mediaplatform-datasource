package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/infraLinkit/mediaplatform-datasource/model"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
)

type IncomingHandler struct {
	Config *config.Cfg
	Logs   *logrus.Logger
	PS     *sql.DB
	Rmqp   rmqp.AMQP
	R      *rueidis.Storage
	DS     *model.BaseModel
}

func NewIncomingHandler(obj IncomingHandler) *IncomingHandler {

	b := model.NewBaseModel(model.BaseModel{
		Config:    obj.Config,
		Logs:      obj.Logs,
		DBPostgre: obj.PS,
		R:         obj.R,
	})

	return &IncomingHandler{
		Config: obj.Config,
		Logs:   obj.Logs,
		PS:     obj.PS,
		R:      obj.R,
		Rmqp:   obj.Rmqp,
		DS:     b,
	}
}

func (h *IncomingHandler) Postback(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	h.Logs.Debug(fmt.Sprintf("Receive request postback %#v ...\n", c.AllParams()))

	// Parse Postback Data
	//country=[COUNTRY]&operator=[OPERATOR]&partner=[PARTNER]&serv_id=[SERV_ID]&keyword=[KEYWORD]&subkeyword=[SUBKEYWORD]&msisdn=[MSISDN]&px=[PIXEL]&trxid=[TRXID]
	p := entity.NewDataPostback(c)

	// Validate Parameters
	if v := p.ValidateParams(h.Logs); v.Code == fiber.StatusBadRequest {

		return c.Status(v.Code).JSON(entity.GlobalResponse{Code: v.Code, Message: v.Message})

	} else {

		if c.Cookies(p.CookieKey) != "" {

			return c.Status(fiber.StatusForbidden).JSON(entity.GlobalResponse{Code: fiber.StatusForbidden, Message: "forbidden access"})

		} else {
			// Setup cookie if double requested within n-hour
			c.Cookie(&fiber.Cookie{
				Name:     p.CookieKey,
				Value:    "1",
				Expires:  time.Now().Add(1 * time.Hour),
				HTTPOnly: true,
				SameSite: "lax",
			})

			pxData := entity.PixelStorage{
				Country: p.Country, Operator: p.Operator, Partner: p.Partner, Service: p.ServiceId, Keyword: p.Keyword, IsBillable: p.IsBillable, Pixel: p.Px}

			if px, err := h.DS.GetPx(pxData); err != nil {

				return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found"})

			} else {

				pxData.Id = px.Id
				pxData.Msisdn = p.Msisdn
				pxData.TrxId = p.TrxId
				pxData.IsUsed = true
				pxData.PixelUsedDate = helper.GetFormatTime(h.Config.TZ, time.RFC3339)

				bodyReq, _ := json.Marshal(pxData)

				corId := "PBA" + helper.GetUniqId(h.Config.TZ)

				published := h.Rmqp.PublishMsg(rmqp.PublishItems{
					ExchangeName: h.Config.RabbitMQPostbackAdnetExchangeName,
					QueueName:    h.Config.RabbitMQPostbackAdnetQueueName,
					ContentType:  h.Config.RabbitMQDataType,
					CorId:        corId,
					Payload:      string(bodyReq),
					Priority:     0,
				})

				if !published {

					h.Logs.Debug(fmt.Sprintf("[x] Failed published: %s, Data: %s ...", corId, string(bodyReq)))

				} else {

					h.Logs.Debug(fmt.Sprintf("[v] Published: %s, Data: %s ...", corId, string(bodyReq)))
				}

				return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: "OK", Data: entity.PixelStorageRsp{
					Adnet:         px.Adnet,
					IsBillable:    px.IsBillable,
					Pixel:         px.Pixel,
					TrxId:         p.TrxId,
					Msisdn:        p.Msisdn,
					Browser:       px.Browser,
					OS:            px.OS,
					PubId:         px.PubId,
					Handset:       px.Handset,
					PixelUsedDate: pxData.PixelUsedDate,
				}})
			}
		}
	}
}

func (h *IncomingHandler) Report(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	// Parse Traffic Data
	//m := c.Queries()

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}
