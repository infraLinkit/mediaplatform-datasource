package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
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

				if px.Id < 0 {
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
							PixelUsedDate: px.PixelUsedDate,
						}})

					} else {

						px.PixelUsedDate = helper.GetFormatTime(h.Config.TZ, time.RFC3339)

						bodyReq, _ := json.Marshal(px)

						corId := "RTO" + helper.GetUniqId(h.Config.TZ)

						published := h.Rmqp.PublishMsg(rmqp.PublishItems{
							ExchangeName: h.Config.RabbitMQRatioExchangeName,
							QueueName:    h.Config.RabbitMQRatioQueueName,
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

func (h *IncomingHandler) Report(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	// Parse Traffic Data
	//m := c.Queries()

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) SetData(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	v := c.Params("v")

	if v == "set_target_daily_budget" {

		m := c.Queries()
		target_daily_budget := m["target_daily_budget"]
		country := strings.ToUpper(m["country"])
		operator := strings.ToUpper(m["operator"])

		if target_daily_budget == "" {
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{Code: fiber.StatusBadRequest, Message: "v is empty"})
		} else {

			gs, _ := h.DS.GetDataConfig("global_setting", "$")

			tdb := gs.TargetDailyBudget
			if target_daily_budget != tdb {
				tdb = target_daily_budget
			}
			var dc []entity.DataConfig
			dc = append(dc, entity.DataConfig{
				TargetDailyBudget: tdb,
				Country:           country,
				Operator:          operator,
			})

			bodyReq, _ := json.Marshal(entity.DataCampaignAction{
				Action: "UPDATE_CAMP_MONITOR_BUDGET", DataConfig: dc})

			corId := "CMN" + helper.GetUniqId(h.Config.TZ)

			published := h.Rmqp.PublishMsg(rmqp.PublishItems{
				ExchangeName: h.Config.RabbitMQCampaignManagementExchangeName,
				QueueName:    h.Config.RabbitMQCampaignManagementQueueName,
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

			/* gs, _ := h.DS.GetDataConfig("global_setting", "$")

			tdb := gs.TargetDailyBudget
			if target_daily_budget != tdb {
				tdb = target_daily_budget
			}

			gset, _ := json.Marshal(entity.GlobalSetting{
				CostPerConversion: gs.CPCR,
				AgencyFee:         gs.AgencyFee,
				TargetDailyBudget: tdb,
			})

			h.DS.SetData("global_setting", "$", string(gset))

			h.DS.UpdateCampaignMonitoringBudget(entity.DataConfig{
				TargetDailyBudget: tdb,
				Country:           country,
				Operator:          operator,
			})

			h.DS.UpdateReportSummaryCampaignMonitoringBudget(helper.GetFormatTime(h.Config.TZ, "2006-01-02"),
				entity.DataConfig{
					TargetDailyBudget: tdb,
					Country:           country,
					Operator:          operator,
				}) */

			return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{Code: fiber.StatusBadRequest, Message: "Request parameter unknown"})
	}
}
