package util

import "reflect"


func Contains[T any](list []T, target T) bool {
	for _, v := range list {
		if reflect.DeepEqual(v, target) {
			return true
		}
	}
	return false
}