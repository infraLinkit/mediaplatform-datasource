package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/src/handler"
)

// Postback INTENTIONALLY PUBLIC: adnet/operator callbacks. Auth via signed params/IP allowlist.
func RegisterPostback(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Get("/postback/:urlservicekey/", h.Postback)
	grp.Get("/postback", h.PostbackV3)
	grp.Get("/postback_billed", h.PostbackBilled)
	grp.Get("/inquire/campid", h.InquiryCampID)
}
