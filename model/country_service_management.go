package model

import (
	"database/sql"
	"strings"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (m *BaseModel) CreateCountry(country *entity.Country) error {
	return m.DB.Create(country).Error
}

func (m *BaseModel) UpdateCountry(country *entity.Country) error {
	return m.DB.Updates(country).Error
}

func (r *BaseModel) GetCountry(o entity.GlobalRequestFromDataTable) ([]entity.Country, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Country{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%").Or("code ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Country
	for rows.Next() {
		var s entity.Country
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateCompany(company *entity.Company) error {
	return m.DB.Create(company).Error
}

func (m *BaseModel) UpdateCompany(company *entity.Company) error {
	return m.DB.Updates(company).Error
}

func (r *BaseModel) GetCompany(o entity.GlobalRequestFromDataTable) ([]entity.Company, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Company{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Company
	for rows.Next() {
		var s entity.Company
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateDomain(domain *entity.Domain) error {
	return m.DB.Create(domain).Error
}

func (m *BaseModel) UpdateDomain(domain *entity.Domain) error {
	return m.DB.Updates(domain).Error
}

func (r *BaseModel) GetDomain(o entity.GlobalRequestFromDataTable) ([]entity.Domain, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Domain{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("url ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("url").Rows()
	defer rows.Close()

	var ss []entity.Domain
	for rows.Next() {
		var s entity.Domain
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateOperator(operator *entity.Operator) error {
	return m.DB.Create(operator).Error
}

func (m *BaseModel) UpdateOperator(operator *entity.Operator) error {
	return m.DB.Updates(operator).Error
}

func (r *BaseModel) GetOperator(o entity.GlobalRequestFromDataTable) ([]entity.Operator, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Operator{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Operator
	for rows.Next() {
		var s entity.Operator
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreatePartner(partner *entity.Partner) error {
	return m.DB.Create(partner).Error
}

func (m *BaseModel) UpdatePartner(partner *entity.Partner) error {
	return m.DB.Updates(partner).Error
}

func (r *BaseModel) GetPartner(o entity.GlobalRequestFromDataTable) ([]entity.Partner, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Partner{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%").
			Or("operator ILIKE ?", "%"+search_value+"%").
			Or("aggregator ILIKE ?", "%"+search_value+"%").
			Or("country ILIKE ?", "%"+search_value+"%").
			Or("company ILIKE ?", "%"+search_value+"%").
			Or("client ILIKE ?", "%"+search_value+"%").
			Or("client_type ILIKE ?", "%"+search_value+"%").
			Or("url_postback ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Partner
	for rows.Next() {
		var s entity.Partner
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateService(service *entity.Service) error {
	return m.DB.Create(service).Error
}

func (m *BaseModel) UpdateService(service *entity.Service) error {
	return m.DB.Updates(service).Error
}

func (r *BaseModel) GetService(o entity.GlobalRequestFromDataTable) ([]entity.Service, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Service{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("service ILIKE ?", "%"+search_value+"%").
			Or("adn ILIKE ?", "%"+search_value+"%").
			Or("country ILIKE ?", "%"+search_value+"%").
			Or("operator ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("service").Rows()
	defer rows.Close()

	var ss []entity.Service
	for rows.Next() {
		var s entity.Service
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateAdnetList(adnet_list *entity.AdnetList) error {
	return m.DB.Create(adnet_list).Error
}

func (m *BaseModel) UpdateAdnetList(adnet_list *entity.AdnetList) error {
	return m.DB.Updates(adnet_list).Error
}

func (r *BaseModel) GetAdnetList(o entity.GlobalRequestFromDataTable) ([]entity.AdnetList, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.AdnetList{})
	query.Select("*, api_url as api_url_before")
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%").
			Or("code ILIKE ?", "%"+search_value+"%").
			Or("api_url ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()

	defer rows.Close()

	var ss []entity.AdnetList
	for rows.Next() {
		var s entity.AdnetList
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetAgency(o entity.GlobalRequestFromDataTable) ([]entity.Agency, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Agency{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%").Or("code ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Agency
	for rows.Next() {
		var s entity.Agency
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateAgency(agency *entity.Agency) error {
	return m.DB.Create(agency).Error
}

func (m *BaseModel) UpdateAgency(agency *entity.Agency) error {
	return m.DB.Updates(agency).Error
}

func (r *BaseModel) GetChannel(o entity.GlobalRequestFromDataTable) ([]entity.Channel, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.Channel{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.Channel
	for rows.Next() {
		var s entity.Channel
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateChannel(channel *entity.Channel) error {
	return m.DB.Create(channel).Error
}

func (m *BaseModel) UpdateChannel(channel *entity.Channel) error {
	return m.DB.Updates(channel).Error
}

func (r *BaseModel) GetMainstreamGroup(o entity.GlobalRequestFromDataTable) ([]entity.MainstreamGroup, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.MainstreamGroup{})
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("name ILIKE ?", "%"+search_value+"%").Or("agency ILIKE ?", "%"+search_value+"%").Or("channel ILIKE ?", "%"+search_value+"%").Or("service ILIKE ?", "%"+search_value+"%").Or("unique_domain ILIKE ?", "%"+search_value+"%")
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("name").Rows()
	defer rows.Close()

	var ss []entity.MainstreamGroup
	for rows.Next() {
		var s entity.MainstreamGroup
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (m *BaseModel) CreateMainstreamGroup(mainstreamGroup *entity.MainstreamGroup) error {
	return m.DB.Create(mainstreamGroup).Error
}

func (m *BaseModel) UpdateMainstreamGroup(mainstreamGroup *entity.MainstreamGroup) error {
	return m.DB.Updates(mainstreamGroup).Error
}
