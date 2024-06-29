package influx

import "errors"

var ErrFieldDoesNotExist = errors.New("extracting value: field does not exist")
var ErrInvalidFieldType = errors.New("extracting value: field type is not compatible")
