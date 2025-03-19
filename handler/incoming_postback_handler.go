package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/wiliehidayat87/rmqp"
)

func (h *IncomingHandler) Postback(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	h.Logs.Debug(fmt.Sprintf("Receive request postback %#v ...\n", c.AllParams()))

	// Parse Postback Data
	//serv_id=[SERV_ID]&msisdn=[MSISDN]&px=[PIXEL]&trxid=[TRXID]
	p := entity.NewDataPostback(c)
	p.URLServiceKey = c.Params("urlservicekey")

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
				Expires:  time.Now().Add(1 * time.Minute),
				HTTPOnly: true,
				SameSite: "lax",
			})

			dc, _ := h.DS.GetDataConfig(helper.Concat("-", p.URLServiceKey, "configIdx"), "$")

			pxData := entity.PixelStorage{
				URLServiceKey: p.URLServiceKey, Pixel: p.AffSub}

			var (
				px  entity.PixelStorage
				err error
			)

			if dc.PostbackMethod == "ADNETCODE" {
				px, err = h.DS.GetByAdnetCode(pxData)
			} else if dc.PostbackMethod == "TOKEN" {
				px, err = h.DS.GetToken(pxData)
			} else {
				px, err = h.DS.GetPx(pxData)
			}

			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found"})

			} else {

				if px.ID < 0 {
					return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found"})
				} else {

					if px.IsUsed {

						return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusNotFound, Message: "NOK - Pixel already used", Data: entity.PixelStorageRsp{
							Adnet:         dc.Adnet,
							IsBillable:    dc.IsBillable,
							Pixel:         p.AffSub,
							Browser:       px.Browser,
							OS:            px.OS,
							Handset:       px.UserAgent,
							PubId:         px.PubId,
							PixelUsedDate: px.PixelUsedDate.Format(time.RFC3339),
						}})

					} else {

						px.PixelUsedDate = helper.GetCurrentTime(h.Config.TZ, time.RFC3339)

						bodyReq, _ := json.Marshal(px)

						corId := "POP" + helper.GetUniqId(h.Config.TZ)

						published := h.Rmqp.PublishMsg(rmqp.PublishItems{
							ExchangeName: h.Config.RabbitMQPopulatePostbackExchangeName,
							QueueName:    h.Config.RabbitMQPopulatePostbackQueueName,
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
							Adnet:         dc.Adnet,
							IsBillable:    dc.IsBillable,
							Pixel:         p.AffSub,
							Browser:       px.Browser,
							OS:            px.OS,
							Handset:       px.UserAgent,
							PubId:         px.PubId,
							PixelUsedDate: helper.GetFormatTime(h.Config.TZ, time.RFC3339),
						}})
					}
				}
			}
		}
	}
}
