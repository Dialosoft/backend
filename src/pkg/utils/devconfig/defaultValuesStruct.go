package devconfig

import (
	"reflect"

	"github.com/Dialosoft/src/pkg/errorsUtils"
)

func SetDefaultValues(model interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errorsUtils.ErrInternalServer
	}
	v = v.Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
	}
	return nil
}
