package ipmi

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Data struct {
	updatedAt time.Time
	data      map[string]any
}

func (data *Data) Update(newData map[string]any) {
	data.updatedAt = time.Now()
	data.data = newData
}

func (data *Data) GetDiscovery() (any, error) {
	if err := checkOutDate(data.updatedAt); err != nil {
		return nil, err
	}
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
	if err := checkOutDate(data.updatedAt); err != nil {
		return nil, err
	}
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

func checkOutDate(compareDate time.Time) error {
	currentDate := time.Now().Add(-5 * time.Minute)
	if compareDate.Before(currentDate) {
		return errors.New("outdated data: maybe not working ipmitool?")
	}
	return nil
}
