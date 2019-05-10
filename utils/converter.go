package utils

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ConvertStringToInterface(contentType string, contentValue string) (interface{}, error) {
	var regexpBooleanType = regexp.MustCompile(`bool`)
	var regexpIntType = regexp.MustCompile(`int`)
	var regexpFloatType = regexp.MustCompile(`float|double|decimal`)
	var regexpDateTimeType = regexp.MustCompile(`time|date`)
	var returnInterface interface{}
	var err error

	if regexpIntType.MatchString(contentType) {
		returnInterface, err = ConvertStringToInt(contentValue)
	} else if regexpFloatType.MatchString(contentType) {
		returnInterface, err = ConvertStringToFloat(contentValue)
	} else if regexpDateTimeType.MatchString(contentType) {
		returnInterface, err = time.Parse(time.RFC3339, contentValue)
	} else if regexpBooleanType.MatchString(contentType) {
		returnInterface, err = ConvertStringToBoolean(contentValue)
	} else {
		returnInterface = contentValue
		err = nil
	}

	return returnInterface, err
}

func ConvertStringToInt(valueContent string) (interface{}, error) {
	return strconv.ParseInt(valueContent, 10, 64)
}

func ConvertStringToFloat(valueContent string) (interface{}, error) {
	return strconv.ParseFloat(valueContent, 64)
}

func ConvertStringToBoolean(valueContent string) (interface{}, error) {
	var checkFalse = regexp.MustCompile(`\A(0|f|false|no)\z`)

	valueContent = strings.TrimSpace(valueContent)

	if valueContent != "" {
		if checkFalse.MatchString(valueContent) {
			return false, nil
		} else {
			return true, nil
		}
	} else {
		return false, nil
	}
}
