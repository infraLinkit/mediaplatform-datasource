package model

import (
	"context"
	"encoding/json"
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

		dcfg = tempCfg[0][0]
		h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data config: %#v ...\n", key, dcfg))
		return dcfg, nil

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

		dcfg = tempCfg[0][0]
		h.Logs.Debug(fmt.Sprintf("Found & Success parse json key (%s) data counter: %#v ...\n", key, dcfg))
		return dcfg, nil

	}

}

func (h *BaseModel) IncrCounterData(key string, path string, val float64) {

	// Get Config Data Landing
	ctx := context.Background()

	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonNumincrby().Key(key).Path(path).Value(val).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Append error counter key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Append success counter key (%s), path (%s) ...\n", key, path))
	}

}

func (h *BaseModel) AppendCounterData(key string, path string, o entity.DataCounterDetail) {

	// Get Config Data Landing
	ctx := context.Background()

	b, _ := json.Marshal(o)
	if err := h.R.Conn().Do(ctx, h.R.Conn().B().JsonArrappend().Key(key).Path(path).Value(string(b)).Build()).Error(); err != nil {
		h.Logs.Debug(fmt.Sprintf("Append error counter key (%s), path (%s), err : %#v ...\n", key, path, err))
	} else {
		h.Logs.Debug(fmt.Sprintf("Append success counter key (%s), path (%s) ...\n", key, path))
	}

}
