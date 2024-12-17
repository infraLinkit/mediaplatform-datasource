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

	if err = rueidis.DecodeSliceOfJSON(h.R.Conn().Do(ctx, h.R.Conn().B().JsonMget().Key(key).Path(path).Build()), &tempCfg); err != nil {

		h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
		return nil, err

	} else {

		if len(tempCfg) < 1 {
			err = errors.New("key is empty or not found")
			h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
			return dcfg, err
		} else {
			h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data config: %#v ...\n", key, tempCfg))
			dcfg = tempCfg[0][0]
			return dcfg, nil
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

	if err = rueidis.DecodeSliceOfJSON(h.R.Conn().Do(ctx, h.R.Conn().B().JsonMget().Key(key).Path(path).Build()), &tempCfg); err != nil {

		h.Logs.Warn(fmt.Sprintf("Cannot find data counter key (%s) or error: %#v ...\n", key, err))
		return nil, err

	} else {

		if len(tempCfg) < 1 {
			err = errors.New("key is empty or not found")
			h.Logs.Warn(fmt.Sprintf("Cannot find data config key (%s) or error: %#v ...\n", key, err))
			return dcfg, err
		} else {
			h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data config: %#v ...\n", key, dcfg))
			dcfg = tempCfg[0][0]
			return dcfg, nil
		}

	}

}

func (h *BaseModel) IncrCounterData(key string, path string, val float64) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonNumincrby().Key(key).Path(path).Value(val).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Increament error counter key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Increament success counter key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) AppendCounterData(key string, path string, o entity.DataCounterDetail) {

	// Get Config Data Landing
	ctx := context.Background()

	b, _ := json.Marshal(o)
	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonArrappend().Key(key).Path(path).Value(string(b)).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Append data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Append data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) SetCounterData(key string, path string, val string) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonSet().Key(key).Path(path).Value(val).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Set data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Set data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) IndexRedis(key string, field string) {

	h.R.Conn().B().FtCreate().Index(key).Prefix(1)

	//`$.campaign_id AS campaign_id TEXT $.pixel AS pixel TEXT $.user_agent AS user_agent TEXT $.os AS os TEXT $.browser AS browser TEXT $.ips AS ips TEXT $.user_is_rejected AS user_is_rejected TAG $.user_is_duplicated AS user_is_duplicated TAG $.refferal_url AS refferal_url TEXT $.handset_code AS handset_code TEXT $.handset_type AS handset_type TEXT $.pixel_is_used AS pixel_is_used TAG`

	//Create Indexing key
	result := h.R.Conn().Do(context.Background(), h.R.Conn().B().FtCreate().Index(key).OnJson().Schema().FieldName(field).Text().Build())
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

func (h *BaseModel) SetData(key string, path string, val string) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonSet().Key(key).Path(path).Value(val).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Set data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Set data success key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) SetExpireData(key string) {

	//Set key expire and delete automatically
	result := h.R.Conn().Do(context.Background(), h.R.Conn().B().Expire().Key(key).Seconds(h.Config.RedisKeyExpiration).Build())

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

	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonDel().Key(key).Path(path).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Del data error key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Del data success key (%s), path (%s) ...\n", key, path))
	}

}
