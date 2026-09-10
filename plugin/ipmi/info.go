package ipmi

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"support"
)

var moduleNameRe = regexp.MustCompile(`^(\[(.+)\] )?\[(.+)\]$`)
var digitsDataRe = regexp.MustCompile(`^([\d]+\.?[\d]*) (RPM|A|W|V)+$`)
var tempDataRe = regexp.MustCompile(`^([\d]+)(C)\/[\d]+F$`)

func parseData(data string) (map[string]any, error) {
	pmData := make(map[string]any)

	pmModules := strings.Split(data, "\n\n")

	if len(pmModules) == 0 {
		return pmData, errors.New("modules list is empty")
	}

	pmError := errors.New("ipmitool incorrect input data")

	for _, moduleRawData := range pmModules {
		moduleData := strings.Split(moduleRawData, "\n")
		if len(moduleData) < 4 {
			return pmData, pmError
		}

		moduleTitle := strings.TrimSpace(moduleData[0])
		moduleNameResult := moduleNameRe.FindSubmatch([]byte(moduleTitle))

		if moduleNameResult == nil {
			return pmData, pmError
		}

		moduleName := string(moduleNameResult[3])
		for _, line := range moduleData[3:] {
			if strings.TrimSpace(line) == "" {
				break
			}
			lineData := strings.Split(line, "|")
			if len(lineData) != 2 {
				return pmData, pmError
			}

			fieldName := strings.TrimSpace(lineData[0])
			fieldValue := strings.TrimSpace(lineData[1])
			fieldInfo := []string{moduleName, fieldName}

			if fieldValue == "N/A" {
				support.AddDeepNested(pmData, fieldInfo, nil)
				continue
			}

			digitValue := digitsDataRe.FindSubmatch([]byte(fieldValue))
			if digitValue != nil {
				digit, err := support.GetDigit(digitValue[1])
				if err == nil {
					support.AddDeepNested(pmData, fieldInfo, digit)
					continue
				}
			}
			tempValue := tempDataRe.FindSubmatch([]byte(fieldValue))
			if tempValue != nil && string(tempValue[2]) == "C" {
				digit, err := support.GetDigit(tempValue[1])
				if err == nil {
					support.AddDeepNested(pmData, fieldInfo, digit)
					continue
				}
			}
			support.AddDeepNested(pmData, fieldInfo, fieldValue)
		}
	}
	return pmData, nil
}

func GetInfo(ctx context.Context, tool string) (map[string]any, error) {
	aData, err := ipmiData(ctx, tool)
	if err != nil {
		return nil, err
	}

	pmData, err := parseData(strings.ReplaceAll(aData, "\r\n", "\n"))
	if err != nil {
		return nil, err
	}

	return pmData, nil
}
