package mapper

import "github.com/Dialosoft/src/pkg/errorsUtils"

// GetValueOrError gets the value of a pointer or returns an error
func GetValueOrError[T any](ptr *T, fieldName string) (T, error) {
    if ptr != nil {
        return *ptr, nil
    }
    var zero T
    return zero, errorsUtils.ErrParameterCannotBeNull
}