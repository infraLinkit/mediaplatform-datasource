package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"gorm.io/gorm"
)

type rolePermissionRequest struct {
	MenuID int  `json:"menu_id"`
	Status bool `json:"status"`
}

type roleRequest struct {
	Code        string                  `json:"code"`
	Name        string                  `json:"name"`
	Permissions []rolePermissionRequest `json:"permissions"`
}

// roleErrors applies RoleController's validation rules and messages.
func (h *IncomingHandler) roleErrors(req *roleRequest, exceptID int) ([]string, error) {
	req.Code, req.Name = strings.TrimSpace(req.Code), strings.TrimSpace(req.Name)
	var errs []string
	if req.Code == "" {
		errs = append(errs, "The code field is required.")
	} else {
		taken, err := h.valueTaken("roles", "code", req.Code, exceptID)
		if err != nil {
			return nil, err
		}
		if taken {
			errs = append(errs, "The code has already been taken.")
		}
	}
	if req.Name == "" {
		errs = append(errs, "role name is required")
	} else if len([]rune(req.Name)) > 50 {
		errs = append(errs, "role name max 50 characters")
	}
	return errs, nil
}

func syncRolePermissions(tx *gorm.DB, roleID int, permissions []rolePermissionRequest) error {
	rows := make([]assignmentRow, 0, len(permissions))
	for _, permission := range permissions {
		rows = append(rows, assignmentRow{Key: []int{permission.MenuID}, Status: permission.Status})
	}
	return syncAssignments(tx, "permissions", "role_id", roleID, []string{"menu_id"}, rows)
}

func (h *IncomingHandler) CreateRole(c *fiber.Ctx) error {
	var req roleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	errs, err := h.roleErrors(&req, 0)
	if err != nil {
		h.Logs.Errorf("CreateRole: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create role"})
	}
	if len(errs) > 0 {
		return validationFailed(c, errs)
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		role := entity.Role{Code: req.Code, Name: req.Name}
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		return syncRolePermissions(tx, role.ID, req.Permissions)
	})
	if err != nil {
		h.Logs.Errorf("CreateRole: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create role"})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

// GetRoleTable always reads the table: Blade reloads the list right after a save, so a cached copy would hide it.
func (h *IncomingHandler) GetRoleTable(c *fiber.Ctx) error {
	roles, err := h.DS.GetAllRolesWithPermission()
	if err != nil {
		h.Logs.Errorf("GetRoleTable: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch roles"})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: config.OK_DESC, Data: roles})
}

func (h *IncomingHandler) UpdateRole(c *fiber.Ctx) error {
	roleID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid role ID"})
	}
	var req roleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	if _, err := h.DS.GetRoleByID(roleID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Role not found"})
	}
	errs, err := h.roleErrors(&req, roleID)
	if err != nil {
		h.Logs.Errorf("UpdateRole: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update role"})
	}
	if len(errs) > 0 {
		return validationFailed(c, errs)
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		values := map[string]interface{}{"code": req.Code, "name": req.Name, "updated_at": time.Now()}
		if err := tx.Model(&entity.Role{}).Where("id = ?", roleID).Updates(values).Error; err != nil {
			return err
		}
		return syncRolePermissions(tx, roleID, req.Permissions)
	})
	if err != nil {
		h.Logs.Errorf("UpdateRole: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update role"})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

// DeleteRole follows RoleRepository::delete: the role's users lose their role and are switched off.
func (h *IncomingHandler) DeleteRole(c *fiber.Ctx) error {
	roleID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid role ID"})
	}
	if _, err := h.DS.GetRoleByID(roleID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Role not found"})
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&entity.Permission{}).Error; err != nil {
			return err
		}
		values := map[string]interface{}{"role_id": nil, "status": false, "updated_at": time.Now()}
		if err := tx.Model(&entity.User{}).Where("role_id = ?", roleID).Updates(values).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.Role{}, roleID).Error
	})
	if err != nil {
		h.Logs.Errorf("DeleteRole: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error when deleting data..."})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true, "message": "Success delete data"})
}

func (h *IncomingHandler) GetRoleByID(c *fiber.Ctx) error {
	roleID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid role ID"})
	}
	role, err := h.DS.GetRoleByID(roleID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Role not found"})
	}

	type grantedMenu struct {
		MenuID int  `json:"menu_id"`
		Status bool `json:"status"`
	}
	permissions := []grantedMenu{}
	if err := h.DB.Table("permissions").Select("menu_id, status").Where("role_id = ?", roleID).Scan(&permissions).Error; err != nil {
		h.Logs.Errorf("GetRoleByID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch role"})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
		Data: fiber.Map{
			"id":          role.ID,
			"code":        role.Code,
			"name":        role.Name,
			"permissions": permissions,
		},
	})
}
