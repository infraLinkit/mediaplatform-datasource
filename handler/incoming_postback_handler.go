package handler

import (
	"encoding/json"
	"fmt"
	"strings"
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
				Expires:  time.Now().Add(3 * time.Second),
				HTTPOnly: true,
				SameSite: "lax",
			})

			if dc, err := h.DS.GetDataConfig(helper.Concat("-", p.URLServiceKey, "configIdx"), "$"); err == nil {

				pxData := entity.PixelStorage{
					URLServiceKey: p.URLServiceKey, Pixel: p.AffSub}

				var (
					px   entity.PixelStorage
					isPX bool
				)

				if dc.PostbackMethod == "ADNETCODE" {
					px, isPX = h.DS.GetByAdnetCode(pxData)
				} else if dc.PostbackMethod == "TOKEN" {
					px, isPX = h.DS.GetToken(pxData)
				} else if dc.PostbackMethod == "JSON-MSISDN" || dc.PostbackMethod == "XML-MSISDN" || dc.PostbackMethod == "HTML-MSISDN" {
					px, isPX = h.DS.GetPxByMsisdn(pxData)
				} else {
					px, isPX = h.DS.GetPx(pxData)
				}

				if !isPX {

					return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found or duplicate used, pixel : " + p.AffSub})

				} else {

					if px.IsUsed {

						return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusNotFound, Message: "NOK - Pixel already used", Data: entity.PixelStorageRsp{
							Adnet:         dc.Adnet,
							IsBillable:    dc.IsBillable,
							Pixel:         px.Pixel,
							Browser:       px.Browser,
							OS:            px.OS,
							Handset:       px.UserAgent,
							PubId:         px.PubId,
							PixelUsedDate: px.PixelUsedDate.Format(time.RFC3339),
						}})

					} else {

						px.PixelUsedDate = helper.GetCurrentTime(h.Config.TZ, time.RFC3339)

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
							Pixel:         px.Pixel,
							Browser:       px.Browser,
							OS:            px.OS,
							Handset:       px.UserAgent,
							PubId:         px.PubId,
							PixelUsedDate: helper.GetFormatTime(h.Config.TZ, time.RFC3339),
						}})
					}
				}

			} else {

				return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Campaign ID not found"})

			}
		}
	}
}

func (h *IncomingHandler) Postback2(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	h.Logs.Debug(fmt.Sprintf("Receive request postback %#v ...\n", c.AllParams()))

	// Parse Postback Data
	p := entity.NewDataPostbackV2(c)
	//p.URLServiceKey = c.Params("urlservicekey")

	// Validate Parameters
	if v := p.ValidateParamsV2(h.Logs); v.Code == fiber.StatusBadRequest {

		return c.Status(v.Code).JSON(entity.GlobalResponse{Code: v.Code, Message: v.Message})

	} else {

		if c.Cookies(p.CookieKey) != "" {

			return c.Status(fiber.StatusForbidden).JSON(entity.GlobalResponse{Code: fiber.StatusForbidden, Message: "forbidden access"})

		} else {
			// Setup cookie if double requested within n-hour
			c.Cookie(&fiber.Cookie{
				Name:     p.CookieKey,
				Value:    "1",
				Expires:  time.Now().Add(3 * time.Second),
				HTTPOnly: true,
				SameSite: "lax",
			})

			if !strings.Contains(p.AffSub, "-") {

				return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Invalid pixel format, pixel : " + p.AffSub})

			} else {

				dataraw := strings.Split(p.AffSub, "-")
				p.URLServiceKey = helper.Concat("-", dataraw[0], dataraw[1])

				if dc, err := h.DS.GetDataConfig(helper.Concat("-", p.URLServiceKey, "configIdx"), "$"); err == nil {

					pxData := entity.PixelStorage{
						URLServiceKey:  p.URLServiceKey,
						Pxdate:         helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
						Pixel:          p.AffSub,
						PostbackMethod: p.Method,
					}

					var px entity.PixelStorage

					isPX := false

					switch p.Method {
					case "ADNETCODE":
						px, isPX = h.DS.GetByAdnetCode(pxData)
					case "TOKEN":
						px, isPX = h.DS.GetToken(pxData)
					case "JSON-MSISDN", "XML-MSISDN", "HTML-MSISDN":
						px, isPX = h.DS.GetPxByMsisdn(pxData)
					case "PIXEL":
						px, isPX = h.DS.GetPx(pxData)
					case "SPC-MVLS":

						campIdRemover := strings.NewReplacer(dc.URLServiceKey+"-", "")
						msisdn := campIdRemover.Replace(p.AffSub)

						isPX = true

						px = entity.PixelStorage{
							CampaignDetailId:  dc.Id,
							Pxdate:            helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
							URLServiceKey:     dc.URLServiceKey,
							CampaignId:        dc.CampaignId,
							Country:           dc.Country,
							Partner:           dc.Partner,
							Operator:          dc.Operator,
							Aggregator:        dc.Aggregator,
							Service:           dc.Service,
							ShortCode:         dc.ShortCode,
							Adnet:             dc.Adnet,
							Keyword:           dc.Keyword,
							Subkeyword:        dc.SubKeyword,
							IsBillable:        dc.IsBillable,
							Plan:              dc.Plan,
							URL:               dc.APIURL,
							URLType:           dc.URLType,
							Pixel:             "NA",
							Msisdn:            msisdn,
							TrxId:             "NA",
							Token:             "NA",
							IsUsed:            true,
							Browser:           "NA",
							OS:                "NA",
							Ip:                strings.Join(c.IPs(), ", "),
							ISP:               "NA",
							ReferralURL:       "NA",
							PubId:             dc.PubId,
							UserAgent:         "NA",
							TrafficSource:     false,
							TrafficSourceData: "NA",
							UserRejected:      false,
							UniqueClick:       false,
							UserDuplicated:    false,
							Handset:           "NA",
							HandsetCode:       "NA",
							HandsetType:       "NA",
							URLLanding:        dc.URLLanding,
							URLWarpLanding:    dc.URLWarpLanding,
							URLService:        dc.URLService,
							URLTFCORSmartlink: dc.URLTFCSmartlink,
							StatusCapping:     dc.StatusCapping,
							StatusRatio:       dc.StatusRatio,
							PO:                dc.PO,
							Cost:              dc.Cost,
							CampaignObjective: dc.Objective,
							Channel:           dc.Channel,
							Currency:          dc.Currency,
							PostbackMethod:    dc.PostbackMethod,
							LandingTime:       helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
							LandedTime:        float64(0),
							HttpStatus:        200,
							IsOperator:        false,
							CreatedAt:         helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
							UpdatedAt:         helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
						}

						h.DS.NewPixel(px)
					}

					if !isPX && p.Method == "" {

						return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found or duplicate used and parameter should have a method parameter, pixel : " + p.AffSub})

					} else {

						if px.IsUsed {

							return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusNotFound, Message: "NOK - Pixel already used", Data: entity.PixelStorageRsp{
								Adnet:         dc.Adnet,
								IsBillable:    dc.IsBillable,
								Pixel:         px.Pixel,
								Browser:       px.Browser,
								OS:            px.OS,
								Handset:       px.UserAgent,
								PubId:         px.PubId,
								PixelUsedDate: px.PixelUsedDate.Format(time.RFC3339),
							}})

						} else {

							px.PixelUsedDate = helper.GetCurrentTime(h.Config.TZ, time.RFC3339)

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
								Pixel:         px.Pixel,
								Browser:       px.Browser,
								OS:            px.OS,
								Handset:       px.UserAgent,
								PubId:         px.PubId,
								PixelUsedDate: helper.GetFormatTime(h.Config.TZ, time.RFC3339),
							}})
						}
					}

				} else {

					return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Campaign ID not found"})

				}
			}
		}
	}
}
