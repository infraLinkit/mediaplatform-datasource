package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"gorm.io/gorm"
)

var validate = validator.New()

type menuRequest struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Permission    string `json:"permission"`
	Url           string `json:"url"`
	Icon          string `json:"icon"`
	Parent        *int   `json:"parent"`
	Sort          string `json:"sort"`
	ShowOnSidebar bool   `json:"show_on_sidebar"`
}

// menuErrors applies MenuController's validation rules.
func (h *IncomingHandler) menuErrors(req *menuRequest, exceptID int) ([]string, error) {
	req.Code, req.Name, req.Url, req.Sort = strings.TrimSpace(req.Code), strings.TrimSpace(req.Name), strings.TrimSpace(req.Url), strings.TrimSpace(req.Sort)
	var errs []string
	unique := func(column, value string) error {
		if value == "" {
			errs = append(errs, "The "+column+" field is required.")
			return nil
		}
		taken, err := h.valueTaken("menus", column, value, exceptID)
		if taken {
			errs = append(errs, "The "+column+" has already been taken.")
		}
		return err
	}
	if err := unique("code", req.Code); err != nil {
		return nil, err
	}
	if req.Name == "" {
		errs = append(errs, "The name field is required.")
	}
	if req.Url == "" {
		errs = append(errs, "The url field is required.")
	}
	if err := unique("sort", req.Sort); err != nil {
		return nil, err
	}
	return errs, nil
}

func (h *IncomingHandler) CreateMenu(c *fiber.Ctx) error {
	var req menuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	errs, err := h.menuErrors(&req, 0)
	if err != nil {
		h.Logs.Errorf("CreateMenu: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create menu"})
	}
	if len(errs) > 0 {
		return validationFailed(c, errs)
	}

	menu := entity.Menu{Code: req.Code, Name: req.Name, Permission: req.Permission, Url: req.Url, Icon: req.Icon, Sort: req.Sort, ShowOnSidebar: req.ShowOnSidebar}
	query := h.DB
	// A top-level menu is NULL parent: Laravel's sidebar and parent list look for exactly that.
	if req.Parent == nil || *req.Parent == 0 {
		query = query.Omit("parent")
	} else {
		menu.Parent = *req.Parent
	}
	if req.Icon == "" {
		query = query.Omit("icon")
	}
	if req.Permission == "" {
		query = query.Omit("permission")
	}
	if err := query.Create(&menu).Error; err != nil {
		h.Logs.Errorf("CreateMenu: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create menu"})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

// GetAllMenus always reads the table: Blade reloads the list right after a save, so a cached copy would hide it.
func (h *IncomingHandler) GetAllMenus(c *fiber.Ctx) error {
	order := "id ASC"
	if c.Query("order") == "sort" {
		order = "sort ASC"
	}
	menus := []entity.Menu{}
	if err := h.DB.Order(order).Find(&menus).Error; err != nil {
		h.Logs.Errorf("GetAllMenus: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch menus",
		})
	}

	totalRecords := len(menus)

	if c.Query("all") == "true" || c.Query("limit") == "all" {
		return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
			Code:    fiber.StatusOK,
			Message: config.OK_DESC,
			Data:    menus,
		})
	}

	pagesize := PAGESIZE
	page, _ := strconv.Atoi(c.Query("page", "1"))
	start := (page - 1) * pagesize
	end := start + pagesize

	var displayMenus []entity.Menu
	if start < totalRecords {
		if end > totalRecords {
			end = totalRecords
		}
		displayMenus = menus[start:end]
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw:            page,
		Code:            fiber.StatusOK,
		Message:         config.OK_DESC,
		Data:            displayMenus,
		RecordsTotal:    totalRecords,
		RecordsFiltered: totalRecords,
	})
}

func (h *IncomingHandler) GetMenuByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponseWithData{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ID format",
			Data:    []entity.Menu{},
		})
	}

	menu, err := h.DS.GetMenuByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponseWithData{
			Code:    fiber.StatusNotFound,
			Message: "Menu not found",
			Data:    []entity.Menu{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
		Data:    menu,
	})
}

// UpdateMenu follows MenuRepository::update: empty fields keep their value, the sidebar switch is always written.
func (h *IncomingHandler) UpdateMenu(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid ID format"})
	}
	var req menuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	var menu entity.Menu
	if err := h.DB.First(&menu, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Menu not found"})
	}
	errs, err := h.menuErrors(&req, id)
	if err != nil {
		h.Logs.Errorf("UpdateMenu: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update menu"})
	}
	if len(errs) > 0 {
		return validationFailed(c, errs)
	}

	values := map[string]interface{}{
		"code":            req.Code,
		"name":            req.Name,
		"url":             req.Url,
		"sort":            req.Sort,
		"show_on_sidebar": req.ShowOnSidebar,
		"updated_at":      time.Now(),
	}
	if req.Permission != "" {
		values["permission"] = req.Permission
	}
	if req.Icon != "" {
		values["icon"] = req.Icon
	}
	if req.Parent != nil && *req.Parent != 0 {
		values["parent"] = *req.Parent
	}
	if err := h.DB.Model(&entity.Menu{}).Where("id = ?", id).Updates(values).Error; err != nil {
		h.Logs.Errorf("UpdateMenu: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update menu",
		})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

// DeleteMenu also removes the menu's role permissions, as MenuRepository::delete does.
func (h *IncomingHandler) DeleteMenu(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid ID format"})
	}
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_id = ?", id).Delete(&entity.Permission{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.Menu{}, id).Error
	})
	if err != nil {
		h.Logs.Errorf("DeleteMenu: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete menu",
		})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

type MenuItemResponse struct {
	ID       int                `json:"id"`
	Name     string             `json:"name"`
	Code     string             `json:"code"`
	Url      string             `json:"url"`
	Icon     string             `json:"icon"`
	Parent   int                `json:"parent"`
	Sort     string             `json:"sort"`
	HasChild []MenuItemResponse `json:"has_child"`
}

// GetSidebarMenu is MenuRepository::getSidebarMenu: granted menus of the caller's own role with show_on_sidebar on.
func (h *IncomingHandler) GetSidebarMenu(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int)

	var user entity.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": fiber.StatusOK, "message": "success", "data": []MenuItemResponse{}})
	}
	roleID := user.RoleID

	var parentMenus []entity.Menu
	query := h.DB.Table("menus").
		Select("menus.*").
		Joins("JOIN permissions rp ON rp.menu_id = menus.id").
		Where("(menus.parent IS NULL OR menus.parent = 0) AND rp.status = true AND menus.show_on_sidebar = true AND rp.role_id = ?", roleID)
	if err := query.Order("menus.sort ASC").Find(&parentMenus).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch menus"})
	}

	result := make([]MenuItemResponse, 0, len(parentMenus))
	for _, pm := range parentMenus {
		var children []entity.Menu
		childQuery := h.DB.Table("menus").
			Select("menus.*").
			Joins("JOIN permissions rp ON rp.menu_id = menus.id").
			Where("menus.parent = ? AND rp.status = true AND menus.show_on_sidebar = true AND rp.role_id = ?", pm.ID, roleID)
		_ = childQuery.Order("menus.sort ASC").Find(&children).Error

		childResponses := make([]MenuItemResponse, 0, len(children))
		for _, cm := range children {
			childResponses = append(childResponses, MenuItemResponse{
				ID:       cm.ID,
				Name:     cm.Name,
				Code:     cm.Code,
				Url:      cm.Url,
				Icon:     cm.Icon,
				Parent:   cm.Parent,
				Sort:     cm.Sort,
				HasChild: []MenuItemResponse{},
			})
		}

		result = append(result, MenuItemResponse{
			ID:       pm.ID,
			Name:     pm.Name,
			Code:     pm.Code,
			Url:      pm.Url,
			Icon:     pm.Icon,
			Parent:   pm.Parent,
			Sort:     pm.Sort,
			HasChild: childResponses,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "success",
		"data":    result,
	})
}
