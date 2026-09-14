package ipmi

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Data struct {
	isRunning bool
	updatedAt time.Time
	data      map[string]any
}

func (data *Data) Update(newData map[string]any) {
	data.updatedAt = time.Now()
	data.data = newData
	data.isRunning = false
	time.Sleep(500 * time.Millisecond)
}

func (data *Data) SetRunning(status bool) {
	data.isRunning = status
}

func (data *Data) IsRunning() bool {
	return data.isRunning
}

func (data *Data) IsOutDated() bool {
	currentDate := time.Now().Add(-1 * time.Minute)
	if data.updatedAt.Before(currentDate) {
		return true
	}
	return false
}

func (data *Data) GetDiscovery() (any, error) {
	var discoveryData []map[string]string
	for moduleName := range data.data {
		tmpData := map[string]string{
			"{#MODULE}": moduleName,
		}
		discoveryData = append(discoveryData, tmpData)
	}
	jsonData, jErr := json.Marshal(discoveryData)
	if jErr != nil {
		return nil, errors.New("cannot convert discovery to json format")
	}
	return string(jsonData), nil
}

func (data *Data) GetFieldData(module string, field string) (any, error) {
	if moduleData, ok := data.data[module]; ok {
		if fieldData, ok := moduleData.(map[string]any)[field]; ok {
			return fieldData, nil
		} else {
			return nil, fmt.Errorf("not found field %s on power supply %s", field, module)
		}
	} else {
		return nil, fmt.Errorf("not found power supply: %s", module)
	}
}
