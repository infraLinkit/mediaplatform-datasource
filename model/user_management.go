package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

// Create User Method
func (r *BaseModel) CreateUser(user *entity.User) (*entity.User, error) {
	var checkUser entity.User
	if err := r.DB.Where("email = ? OR username = ?", user.Email, user.Username).First(&checkUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	if err := r.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Create or Update Detail User
func (r *BaseModel) CreateOrUpdateDetailUser(detail *entity.DetailUser) (*entity.DetailUser, error) {
	var existing entity.DetailUser
	if err := r.DB.Where("user_id = ? AND country_id = ? AND operator_id = ? AND service_id = ?",
		detail.UserID, detail.CountryID, detail.OperatorID, detail.ServiceID).First(&existing).Error; err == nil {
		existing.Status = detail.Status
		if err := r.DB.Save(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}

	if err := r.DB.Create(detail).Error; err != nil {
		return nil, err
	}
	return detail, nil
}

// Create or Update User Adnet
func (r *BaseModel) CreateOrUpdateUserAdnet(adnet *entity.UserAdnet) (*entity.UserAdnet, error) {
	var existing entity.UserAdnet
	if err := r.DB.Where("user_id = ? AND adnet_id = ?", adnet.UserID, adnet.AdnetID).First(&existing).Error; err == nil {
		existing.Status = adnet.Status
		if err := r.DB.Save(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}

	if err := r.DB.Create(adnet).Error; err != nil {
		return nil, err
	}
	return adnet, nil
}

// Create or Update User Company
func (r *BaseModel) CreateOrUpdateUserCompany(company *entity.UserCompany) (*entity.UserCompany, error) {
	var existing entity.UserCompany
	if err := r.DB.Where("user_id = ? AND company_id = ?", company.UserID, company.CompanyID).First(&existing).Error; err == nil {
		existing.Status = company.Status
		if err := r.DB.Save(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}

	if err := r.DB.Create(company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

// Create or Update User Agency
func (r *BaseModel) CreateOrUpdateUserAgency(agency *entity.UserAgency) (*entity.UserAgency, error) {
	var existing entity.UserAgency
	if err := r.DB.Where("user_id = ? AND agency_id = ?", agency.UserID, agency.AgencyID).First(&existing).Error; err == nil {
		existing.Status = agency.Status
		if err := r.DB.Save(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}

	if err := r.DB.Create(agency).Error; err != nil {
		return nil, err
	}
	return agency, nil
}

// GetAllusers retrieves all users
func (r *BaseModel) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *BaseModel) GetAllUserWithRelation() ([]entity.UserManagementData, error) {
	var users []entity.User

	// Query without filtering
	query := r.DB.Where("is_verify = ?", true).
		Preload("Role").
		Preload("DetailUser", "status = ?", true).
		Preload("UserAdnet", "status = ?", true)

	// Execute Query
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	// Mapping Users
	var mappedUsers []entity.UserManagementData
	for _, user := range users {
		// Count unique countries in DetailUser
		countryMap := make(map[int]bool)
		for _, detail := range user.DetailUser {
			countryMap[detail.CountryID] = true
		}
		totalCountries := len(countryMap)

		// Count active UserAdnets
		totalAdnets := len(user.UserAdnet)

		mappedUsers = append(mappedUsers, entity.UserManagementData{
			ID:             user.ID,
			Username:       user.Username,
			Name:           user.Name,
			Email:          user.Email,
			Role:           user.Role.Name,
			Status:         user.Status,
			LastLogin:      user.LastLogin,
			IPAddress:      user.IPAddress,
			Handset:        user.Handset,
			TotalCountries: totalCountries,
			TotalAdnets:    totalAdnets,
		})
	}

	return mappedUsers, nil
}

// Get User Counts
func (r *BaseModel) GetUserCounts() (entity.UserCounts, error) {
	var total, active, nonActive int64

	// Get total users
	if err := r.DB.Model(&entity.User{}).Count(&total).Error; err != nil {
		return entity.UserCounts{}, err
	}

	// Get active users
	if err := r.DB.Model(&entity.User{}).Where("status = ?", true).Count(&active).Error; err != nil {
		return entity.UserCounts{}, err
	}

	// Get non-active users
	if err := r.DB.Model(&entity.User{}).Where("status = ?", false).Count(&nonActive).Error; err != nil {
		return entity.UserCounts{}, err
	}

	return entity.UserCounts{
		TotalUsers:    int(total),
		ActiveUsers:   int(active),
		NonActiveUsers: int(nonActive),
	}, nil
}

// GetUserByID retrieves a user by its ID
func (r *BaseModel) GetUserByID(userID int) (*entity.User, error) {
	var user entity.User
	if err := r.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update User Method
func (r *BaseModel) UpdateUser(userID int, userData *entity.User) (*entity.User, error) {
	var user entity.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Prevent updating email or username if they already exist for another user
	var checkUser entity.User
	if err := r.DB.Where("(email = ? OR username = ?) AND id != ?", userData.Email, userData.Username, userID).First(&checkUser).Error; err == nil {
		return nil, errors.New("email or username already in use")
	}

	if err := r.DB.Model(&user).Updates(userData).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *BaseModel) UpdateUserStatus(id int, status bool) (*entity.User, error) {
	user, err := r.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Only the status: saving the whole row would rewrite NULL columns as empty strings.
	if err := r.DB.Model(&entity.User{}).Where("id = ?", id).Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	user.Status = status

	return user, nil
}

// Delete User deletes a user by ID
func (r *BaseModel) DeleteUser(userID int) error {
	result := r.DB.Where("id = ?", userID).Delete(&entity.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *BaseModel) GetAllUserApprovalRequest() ([]entity.UserApprovalRequestData, error) {
	var users []entity.User
	err := r.DB.Where("is_verify = ?", false).Find(&users).Error
	if err != nil {
		return nil, err
	}

	// Mapping Users to Response Struct
	var userResponses []entity.UserApprovalRequestData
	for _, user := range users {
		userResponses = append(userResponses, entity.UserApprovalRequestData{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Email:    user.Email,
			TypeUser: user.TypeUser,
		})
	}

	return userResponses, nil
}

// Approve User
func (r *BaseModel) ApproveUser(id int, roleID int, verifyBy string) error {
	var user entity.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return errors.New("user not found")
	}

	now := time.Now()

	user.RoleID = roleID
	user.VerifyBy = verifyBy
	user.VerifyAt = &now
	user.IsVerify = true
	user.Status = true

	if err := r.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// GetUserDetail retrieves full user details with relations for editing / assign service
func (r *BaseModel) GetUserDetail(userID int) (*entity.UserDetailResponse, error) {
	var user entity.User
	if err := r.DB.Where("id = ?", userID).Preload("Role").First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var detailUser []entity.DetailUser
	r.DB.Where("user_id = ?", userID).Find(&detailUser)

	var userAdnet []entity.UserAdnet
	r.DB.Where("user_id = ?", userID).Find(&userAdnet)

	var userCompany []entity.UserCompany
	r.DB.Unscoped().Where("user_id = ?", userID).Find(&userCompany)

	var userAgency []entity.UserAgency
	r.DB.Unscoped().Where("user_id = ?", userID).Find(&userAgency)

	return &entity.UserDetailResponse{
		User:        user,
		DetailUser:  detailUser,
		UserAdnet:   userAdnet,
		UserCompany: userCompany,
		UserAgency:  userAgency,
	}, nil
}

// GetUserOptions retrieves master data for user management forms (create/edit/assign service)
func (r *BaseModel) GetUserOptions() (*entity.UserOptionsData, error) {
	var roles []entity.Role
	// RoleRepository::getAll without the admin role, which these forms never offer.
	r.DB.Where("id <> 1").Order("name ASC").Find(&roles)

	var companies []entity.Company
	r.DB.Order("id ASC").Find(&companies)

	var agencies []entity.Agency
	r.DB.Order("id ASC").Find(&agencies)

	var adnets []entity.AdnetList
	r.DB.Order("code ASC").Find(&adnets)

	type sRow struct {
		ID           int    `gorm:"column:id"`
		Service      string `gorm:"column:service"`
		Adn          string `gorm:"column:adn"`
		CountryID    int    `gorm:"column:country_id"`
		CountryCode  string `gorm:"column:country_code"`
		OperatorID   int    `gorm:"column:operator_id"`
		OperatorName string `gorm:"column:operator_name"`
	}

	var rows []sRow
	r.DB.Table("services as s").
		Select("s.id, s.service, s.adn, c.id as country_id, c.code as country_code, o.id as operator_id, o.name as operator_name").
		// One row per service even when a code or name repeats, as Eloquent's belongsTo gives one parent.
		Joins("JOIN countries c ON c.id = (SELECT MAX(id) FROM countries WHERE code = s.country)").
		Joins("JOIN operators o ON o.id = (SELECT MAX(id) FROM operators WHERE name = s.operator)").
		Order("s.service ASC").
		Scan(&rows)

	groupMap := make(map[string]*entity.ServiceGroupOption)
	var groupOrder []string

	for _, row := range rows {
		groupKey := fmt.Sprintf("[%s] %s", row.CountryCode, row.OperatorName)
		if _, exists := groupMap[groupKey]; !exists {
			groupMap[groupKey] = &entity.ServiceGroupOption{
				Group:      groupKey,
				CountryID:  row.CountryID,
				OperatorID: row.OperatorID,
				Services:   []entity.ServiceOption{},
			}
			groupOrder = append(groupOrder, groupKey)
		}
		groupMap[groupKey].Services = append(groupMap[groupKey].Services, entity.ServiceOption{
			ID:   row.ID,
			Code: row.Adn,
			Name: row.Service,
		})
	}

	var serviceGroups []entity.ServiceGroupOption
	for _, k := range groupOrder {
		serviceGroups = append(serviceGroups, *groupMap[k])
	}

	return &entity.UserOptionsData{
		Roles:     roles,
		Companies: companies,
		Agencies:  agencies,
		Adnets:    adnets,
		Services:  serviceGroups,
	}, nil
}

