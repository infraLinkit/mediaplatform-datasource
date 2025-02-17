package model

import (
	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

// CreateMenu inserts a new menu into the database
func (m *BaseModel) CreateMenu(menu *entity.Menu) error {
	return m.DB.Create(menu).Error
}

// GetAllMenus retrieves all menus
func (m *BaseModel) GetAllMenus() ([]entity.Menu, error) {
	var menus []entity.Menu
	err := m.DB.Find(&menus).Error
	return menus, err
}

// GetMenuByID retrieves a menu by ID
func (m *BaseModel) GetMenuByID(id uint) (*entity.Menu, error) {
	var menu entity.Menu
	err := m.DB.First(&menu, id).Error
	return &menu, err
}

// UpdateMenu updates an existing menu
func (m *BaseModel) UpdateMenu(menu *entity.Menu) error {
	return m.DB.Save(menu).Error
}

// DeleteMenu deletes a menu by ID
func (m *BaseModel) DeleteMenu(id uint) error {
	return m.DB.Delete(&entity.Menu{}, id).Error
}
