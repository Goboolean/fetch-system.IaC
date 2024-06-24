package mapper

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Goboolean/fetch-system.IaC/pkg/influx"
)

func StructToPoint(in any) (influx.Point, error) {

	p := make(influx.Point)
	err := structToPoint(reflect.ValueOf(in), "", p)
	return p, err
}

func structToPoint(v reflect.Value, path string, p influx.Point) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := v.Type().Field(i)
		fv := v.Field(i)

		name := f.Tag.Get("name")
		if name == "" {
			name = f.Name
		}

		if f.Type.Kind() == reflect.Pointer {
			fv = fv.Elem()
		}

		switch fv.Type().Kind() {
		case reflect.Func, reflect.Chan, reflect.Interface, reflect.Invalid, reflect.UnsafePointer:
			return fmt.Errorf("unsupported type, field name: %v, type :%v", f.Name, f.Type.Name())
		case reflect.Struct:
			err := structToPoint(fv, joinString(path, name, "."), p)
			if err != nil {
				return err
			}
		case reflect.Array, reflect.Slice:
			err := arrayToPoint(fv, joinString(path, name, "."), p)
			if err != nil {
				return err
			}
		case reflect.Map:
			err := mapToPoint(fv, joinString(path, name, "."), p)
			if err != nil {
				return err
			}
		default:
			p[joinString(path, name)] = fv.Interface()
		}
	}
	return nil
}

func arrayToPoint(v reflect.Value, path string, p influx.Point) error {
	for i := 0; i < v.Len(); i++ {
		e := v.Index(i)
		idxString := joinString(".", strconv.FormatInt(int64(i), 10))

		if e.Kind() == reflect.Struct {
			err := structToPoint(e, joinString(path, idxString, "."), p)
			if err != nil {
				return err
			}
		}

		p[joinString(path, idxString)] = e.Interface()
	}
	return nil
}

func mapToPoint(v reflect.Value, path string, p influx.Point) error {
	iter := v.MapRange()

	for iter.Next() {

		keyString, err := encodeMapKey(iter.Key())
		if err != nil {
			return err
		}

		keyString = joinString(".", keyString)

		if iter.Value().Kind() == reflect.Struct {
			return structToPoint(v, joinString(path, keyString, "."), p)
		}

		p[joinString(path, keyString)] = iter.Value().Interface()

	}

	return nil
}

// 표준 라이브러리에서 그대로 가져옴
func encodeMapKey(key reflect.Value) (string, error) {

	if key.Kind() == reflect.String {
		return key.String(), nil
	}

	if tm, ok := key.Interface().(encoding.TextMarshaler); ok {
		if key.Kind() == reflect.Pointer && key.IsNil() {
			return "", nil
		}
		buf, err := tm.MarshalText()
		return string(buf), err
	}

	switch key.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(key.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(key.Uint(), 10), nil
	}
	return "", errors.New("encode map key unexpected map key type")
}

// strings.join()의 인터페이스가 여기에서 쓰기 장황해서 커스텀 구현
func joinString(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
