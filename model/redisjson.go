package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/redis/rueidis"
)

func (h *BaseModel) GetDataConfig(key string, path string) (*entity.DataConfig, error) {

	// Get Config Data Landing
	ctx := context.Background()

	var (
		dcfg    *entity.DataConfig
		tempCfg [][]*entity.DataConfig // or []User is also scannable
		err     error
	)

	if err = rueidis.DecodeSliceOfJSON(h.R0.Conn().Do(ctx, h.R0.Conn().B().JsonMget().Key(key).Path(path).Build()), &tempCfg); err != nil {

		h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
		return nil, err

	} else {

		if len(tempCfg) > 0 && tempCfg != nil {
			h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data config: %#v ...\n", key, tempCfg))
			dcfg = tempCfg[0][0]
			return dcfg, nil
		} else {
			err = errors.New("key is empty or not found")
			h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
			return dcfg, err
		}

	}

}

func (h *BaseModel) GetDataConfigCounter(key string, path string) (*entity.DataCounter, error) {

	// Get Config Data Landing
	ctx := context.Background()

	var (
		dcfg    *entity.DataCounter
		tempCfg [][]*entity.DataCounter // or []User is also scannable
		err     error
	)

	if err = rueidis.DecodeSliceOfJSON(h.R0.Conn().Do(ctx, h.R0.Conn().B().JsonMget().Key(key).Path(path).Build()), &tempCfg); err != nil {

		h.Logs.Warn(fmt.Sprintf("Cannot find data counter key (%s) or error: %#v ...\n", key, err))
		return nil, err

	} else {

		if len(tempCfg) > 0 && tempCfg != nil {
			h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data config: %#v ...\n", key, tempCfg))
			dcfg = tempCfg[0][0]
			return dcfg, nil
		} else {
			err = errors.New("key is empty or not found")
			h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
			return dcfg, err
		}

	}

}

func (h *BaseModel) IncrCounterData(key string, path string, val float64) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R0.Conn().Do(ctx, h.R0.Conn().B().JsonNumincrby().Key(key).Path(path).Value(val).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Increament error counter key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Increament success counter key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) AppendCounterData(key string, path string, o entity.DataCounterDetail) {

	// Get Config Data Landing
	ctx := context.Background()

	b, _ := json.Marshal(o)
	if err := h.R0.Conn().Do(ctx, h.R0.Conn().B().JsonArrappend().Key(key).Path(path).Value(string(b)).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Append data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Append data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) SetCounterData(key string, path string, val string) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R0.Conn().Do(ctx, h.R0.Conn().B().JsonSet().Key(key).Path(path).Value(val).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Set data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Set data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) IndexRedis(key string, field string) {

	h.R0.Conn().B().FtCreate().Index(key).Prefix(1)

	//`$.campaign_id AS campaign_id TEXT $.pixel AS pixel TEXT $.user_agent AS user_agent TEXT $.os AS os TEXT $.browser AS browser TEXT $.ips AS ips TEXT $.user_is_rejected AS user_is_rejected TAG $.user_is_duplicated AS user_is_duplicated TAG $.refferal_url AS refferal_url TEXT $.handset_code AS handset_code TEXT $.handset_type AS handset_type TEXT $.pixel_is_used AS pixel_is_used TAG`

	//Create Indexing key
	result := h.R0.Conn().Do(context.Background(), h.R0.Conn().B().FtCreate().Index(key).OnJson().Schema().FieldName(field).Text().Build())
	isIdxCreated, err := result.AsBool()

	if err != nil {
		h.Logs.Info(fmt.Sprintf("[v] Created FT idx key ( %s ) : %t, %#v...\n\r", key, isIdxCreated, result))
	} else {

		if !isIdxCreated {
			h.Logs.Info(fmt.Sprintf("[x] Failed created FT idx key ( %s ) : %#v...\n\r", key, isIdxCreated))
		} else {
			h.Logs.Info(fmt.Sprintf("[v] Success created FT idx key ( %s ) : %t, %#v...\n\r", key, isIdxCreated, result))
		}

	}

}

func (h *BaseModel) GetData(i []interface{}, key string, path string) []interface{} {

	// Get Config Data Landing
	ctx := context.Background()

	rueidis.DecodeSliceOfJSON(h.R0.Conn().Do(ctx, h.R0.Conn().B().JsonGet().Key(key).Path(path).Build()), &i)

	return i
}

func (h *BaseModel) SetData(key string, path string, val string) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R0.Conn().Do(ctx, h.R0.Conn().B().JsonSet().Key(key).Path(path).Value(val).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Set data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Set data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) SetExpireData(key string, expr int64) {

	//Set key expire and delete automatically
	result := h.R0.Conn().Do(context.Background(), h.R0.Conn().B().Expire().Key(key).Seconds(expr).Build())

	isExpiredCreated, err := result.AsBool()

	if err != nil {
		h.Logs.Info(fmt.Sprintf("[x] Error Set duration expire key ( %s ) : %t, %#v...\n\r", key, isExpiredCreated, result))
	} else {

		if !isExpiredCreated {
			h.Logs.Info(fmt.Sprintf("[x] Failed Set duration expire key ( %s ) : %#v...\n\r", key, isExpiredCreated))
		} else {
			h.Logs.Info(fmt.Sprintf("[v] Success Set duration expire key ( %s ) : %t, %#v...\n\r", key, isExpiredCreated, result))
		}

	}
}

func (h *BaseModel) DelData(key string, path string) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R0.Conn().Do(ctx, h.R0.Conn().B().JsonDel().Key(key).Path(path).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Del data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Del data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) GetAlertData(key string, path string) (*entity.AlertData, error) {

	// Get Config Data Landing
	ctx := context.Background()

	var (
		dcfg    *entity.AlertData
		tempCfg [][]*entity.AlertData // or []User is also scannable
		err     error
	)

	if err = rueidis.DecodeSliceOfJSON(h.R0.Conn().Do(ctx, h.R0.Conn().B().JsonMget().Key(key).Path(path).Build()), &tempCfg); err != nil {

		h.Logs.Warn(fmt.Sprintf("Cannot find data counter key (%s) or error: %#v ...\n", key, err))
		return nil, err

	} else {

		if len(tempCfg) > 0 && tempCfg != nil {
			h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data config: %#v ...\n", key, tempCfg))
			dcfg = tempCfg[0][0]
			return dcfg, nil
		} else {
			err = errors.New("key is empty or not found")
			h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
			return dcfg, err
		}

	}

}

func (h *BaseModel) GetDataSummary(key string, path string) (*entity.Summary, error) {

	// Get Config Data Landing
	ctx := context.Background()

	var (
		dcfg    *entity.Summary
		tempCfg [][]*entity.Summary // or []User is also scannable
		err     error
	)

	if err = rueidis.DecodeSliceOfJSON(h.R0.Conn().Do(ctx, h.R0.Conn().B().JsonMget().Key(key).Path(path).Build()), &tempCfg); err != nil {

		h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
		return nil, err

	} else {

		if len(tempCfg) > 0 && tempCfg != nil {
			h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data config: %#v ...\n", key, tempCfg))
			dcfg = tempCfg[0][0]
			return dcfg, nil
		} else {
			err = errors.New("key is empty or not found")
			h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
			return dcfg, err
		}

	}

}

func (h *BaseModel) RGetApiPinReport(key string, path string) ([]entity.ApiPinReport, bool) {

	var (
		isEmpty bool
		p       []entity.ApiPinReport
	)

	// Get Config Data Landing
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R1.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var pinreport [][]entity.ApiPinReport
		v.DecodeJSON(&pinreport)

		if len(pinreport) > 0 {
			isEmpty = false
			p = pinreport[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	return p, isEmpty
}

func (h *BaseModel) RGetApiPinPerformanceReport(key string, path string) ([]entity.ApiPinPerformance, bool) {

	var (
		isEmpty bool
		p       []entity.ApiPinPerformance
	)

	// Get Config Data Landing
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R1.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var pinperformancereport [][]entity.ApiPinPerformance
		v.DecodeJSON(&pinperformancereport)

		if len(pinperformancereport) > 0 {
			isEmpty = false
			p = pinperformancereport[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}

	return p, isEmpty
}

func (h *BaseModel) RGetDisplayCPAReport(key string, path string) ([]entity.SummaryCampaign, bool) {

	var (
		isEmpty bool
		p       []entity.SummaryCampaign
	)

	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R1.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var displaycpareport [][]entity.SummaryCampaign
		v.DecodeJSON(&displaycpareport)

		if len(displaycpareport) > 0 {
			isEmpty = false
			p = displaycpareport[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}
	return p, isEmpty
}

func (h *BaseModel) RGetDisplayCostReport(key string, path string) ([]entity.CostReport, bool) {
	var (
		isEmpty bool
		p       []entity.CostReport
	)
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R1.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var costreport [][]entity.CostReport
		v.DecodeJSON(&costreport)

		if len(costreport) > 0 {
			isEmpty = false
			p = costreport[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}
	return p, isEmpty
}
func (h *BaseModel) RGetDisplayCostReportDetail(key string, path string) ([]entity.CostReport, bool) {
	var (
		isEmpty bool
		p       []entity.CostReport
	)
	ctx := context.Background()

	data, _ := rueidis.JsonMGet(h.R1.Conn(), ctx, []string{key}, "$")

	for _, v := range data {
		var displaycostreport [][]entity.CostReport
		v.DecodeJSON(&displaycostreport)

		if len(displaycostreport) > 0 {
			isEmpty = false
			p = displaycostreport[0]
			h.Logs.Debug(fmt.Sprintf("Found & success parse json key (%s), total data : %d ...\n", key, len(p)))
		} else {
			isEmpty = true
			h.Logs.Debug(fmt.Sprintf("Data not found json key (%s) ...\n", key))
		}
	}
	return p, isEmpty
}
