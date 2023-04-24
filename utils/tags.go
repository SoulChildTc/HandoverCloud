package utils

import (
	"errors"
	"reflect"
	"strings"
)

func GetTagValue(obj any, fieldName string, tagName string) (string, error) {
	t := reflect.TypeOf(obj).Elem()
	f, exist := t.FieldByName(fieldName)
	if !exist {
		return "", errors.New("字段不存在")
	}
	tag, ok := f.Tag.Lookup(tagName)
	if !ok {
		return "", errors.New("tag不存在")
	}
	return tag, nil

}

func GetTagValueByNamespace(obj any, fieldNamespace string, tagName string) (string, error) {
	fields := strings.Split(fieldNamespace, ".")
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	for index, field := range fields {
		if index == 0 || index == len(fields)-1 {
			continue
		}

		value = value.FieldByName(field)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
	}

	f, exist := value.Type().FieldByName(fields[len(fields)-1])
	if !exist {
		return "", errors.New("字段不存在")
	}
	tag, ok := f.Tag.Lookup(tagName)
	if !ok {
		return "", errors.New("tag不存在")
	}
	return tag, nil
}
