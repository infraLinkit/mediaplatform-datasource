package model

import (
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) PinReport(o entity.ApiPinReport) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) PinPerformanceReport(o entity.ApiPinPerformance) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

