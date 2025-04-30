// Package devconfig provides functionality for handling default values in data structures.
// It uses Go's reflection mechanism to initialize nil pointer fields in structs,
// which helps prevent null pointer exceptions when accessing nested structures.
package devconfig

import (
	"reflect"

	"github.com/Dialosoft/src/pkg/errorsUtils"
)

// SetDefaultValues initializes nil pointer fields in a struct by creating new instances.
//
// This function uses Go's reflection capabilities to:
// 1. Check if the provided parameter is a pointer to a struct
// 2. Iterate through all fields of the struct
// 3. For each field that is a nil pointer, initialize it with a new zero value
//
// The reflect package is crucial for this functionality as it allows:
// - Runtime type inspection (using Kind(), Elem(), Type())
// - Value manipulation (using Set(), New())
// - Dynamic field iteration (using NumField(), Field())
//
// Example usage:
//
//	type Request struct {
//		Name string
//		Config *Configuration
//	}
//
//	type Configuration struct {
//		Setting1 string
//		Setting2 int
//	}
//
//	func ProcessRequest(c echo.Context) error {
//		var req Request
//		// At this point, req.Config is nil
//
//		// Initialize nil pointers with their zero values
//		err := devconfig.SetDefaultValues(&req)
//		if err != nil {
//			return response.ErrInternalServer(c, err, req, "layer")
//		}
//		// Now req.Config is a pointer to an empty Configuration struct
//		// This prevents nil pointer panics when accessing req.Config.Setting1
//
//		return nil
//	}
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