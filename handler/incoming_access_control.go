package handler

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RequirePermission is the API counterpart of Laravel's check-role middleware: the caller's role must
// hold at least one of the given menu codes. It only runs once AuthGate has identified the caller, so
// with ENFORCE_API_AUTH off it lets the request through like every other management route.
func (h *IncomingHandler) RequirePermission(codes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, identified := c.Locals("user_id").(int)
		if !identified {
			return c.Next()
		}
		var granted int64
		err := h.DB.Table("permissions").
			Joins("JOIN menus ON menus.id = permissions.menu_id").
			Joins("JOIN users ON users.role_id = permissions.role_id").
			Where("users.id = ? AND users.deleted_at IS NULL AND permissions.status = true AND menus.code IN ?", userID, codes).
			Count(&granted).Error
		if err != nil {
			h.Logs.Errorf("RequirePermission: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to check permission"})
		}
		if granted == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "You do not have permission.", "status": "error"})
		}
		return c.Next()
	}
}

// laravelValidationMessage joins validation errors the way a Laravel ValidationException reports them.
func laravelValidationMessage(errs []string) string {
	switch len(errs) {
	case 0:
		return ""
	case 1:
		return errs[0]
	case 2:
		return errs[0] + " (and 1 more error)"
	}
	return fmt.Sprintf("%s (and %d more errors)", errs[0], len(errs)-1)
}

func validationFailed(c *fiber.Ctx, errs []string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": laravelValidationMessage(errs)})
}

// valueTaken mirrors Laravel's unique rule, which also counts soft-deleted rows.
func (h *IncomingHandler) valueTaken(table, column, value string, exceptID int) (bool, error) {
	query := h.DB.Table(table).Where(column+" = ?", value)
	if exceptID > 0 {
		query = query.Where("id <> ?", exceptID)
	}
	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

type assignmentRow struct {
	Key    []int
	Status bool
}

type storedAssignment struct {
	ID     int
	Key    []int
	Status bool
}

type assignmentPlan struct {
	Create     []assignmentRow
	Activate   []int
	Deactivate []int
}

func assignmentKey(key []int) string {
	return strings.Trim(strings.ReplaceAll(fmt.Sprint(key), " ", "-"), "[]")
}

// planAssignments works out the writes for one switch list. Every switch is posted on each save,
// switched off or not, and only rows whose status actually changed are touched.
func planAssignments(stored []storedAssignment, posted []assignmentRow) assignmentPlan {
	byKey := make(map[string][]storedAssignment, len(stored))
	for _, row := range stored {
		key := assignmentKey(row.Key)
		byKey[key] = append(byKey[key], row)
	}
	seen := make(map[string]bool, len(posted))
	var plan assignmentPlan
	for i := len(posted) - 1; i >= 0; i-- {
		row := posted[i]
		key := assignmentKey(row.Key)
		if seen[key] {
			continue
		}
		seen[key] = true
		existing := byKey[key]
		if len(existing) == 0 {
			plan.Create = append(plan.Create, row)
			continue
		}
		for _, current := range existing {
			if current.Status == row.Status {
				continue
			}
			if row.Status {
				plan.Activate = append(plan.Activate, current.ID)
			} else {
				plan.Deactivate = append(plan.Deactivate, current.ID)
			}
		}
	}
	return plan
}

// syncAssignments applies a posted switch list (role permissions, user services, adnets...) in one go.
func syncAssignments(tx *gorm.DB, table, ownerColumn string, ownerID int, keyColumns []string, posted []assignmentRow) error {
	if len(posted) == 0 {
		return nil
	}
	columns := append([]string{"id", "status"}, keyColumns...)
	rows, err := tx.Table(table).Select(columns).Where(ownerColumn+" = ?", ownerID).Rows()
	if err != nil {
		return err
	}
	var stored []storedAssignment
	for rows.Next() {
		var (
			id     int
			status sql.NullBool
		)
		keys := make([]sql.NullInt64, len(keyColumns))
		dest := []interface{}{&id, &status}
		for i := range keys {
			dest = append(dest, &keys[i])
		}
		if err := rows.Scan(dest...); err != nil {
			rows.Close()
			return err
		}
		row := storedAssignment{ID: id, Status: status.Bool, Key: make([]int, len(keys))}
		for i, key := range keys {
			row.Key[i] = int(key.Int64)
		}
		stored = append(stored, row)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	plan := planAssignments(stored, posted)
	now := time.Now()
	for status, ids := range map[bool][]int{true: plan.Activate, false: plan.Deactivate} {
		if len(ids) == 0 {
			continue
		}
		if err := tx.Table(table).Where("id IN ?", ids).Updates(map[string]interface{}{"status": status, "updated_at": now}).Error; err != nil {
			return err
		}
	}
	const batch = 200
	for start := 0; start < len(plan.Create); start += batch {
		end := start + batch
		if end > len(plan.Create) {
			end = len(plan.Create)
		}
		values := make([]map[string]interface{}, 0, end-start)
		for _, row := range plan.Create[start:end] {
			value := map[string]interface{}{ownerColumn: ownerID, "status": row.Status, "created_at": now, "updated_at": now}
			for i, column := range keyColumns {
				value[column] = row.Key[i]
			}
			values = append(values, value)
		}
		if err := tx.Table(table).Create(values).Error; err != nil {
			return err
		}
	}
	return nil
}
