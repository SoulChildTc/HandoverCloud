package utils

import (
	"errors"
	"reflect"
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
