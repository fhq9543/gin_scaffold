package reflect2

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go_base/utils/errors"
)

const (
	ErrNonPrimitive = errors.Err("not primitive type")
)

// IsSlice check whether or not param is slice
func IsSlice(s interface{}) bool {
	return s != nil && reflect.TypeOf(s).Kind() == reflect.Slice
}

// Equaler is a interface that compare whether two object is equal
type Equaler interface {
	EqualTo(interface{}) bool
}

// IndirectType return real type of value without pointer
func IndirectType(v interface{}) reflect.Type {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		return typ.Elem()
	}

	return typ
}

func CanNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return true
	default:
		return false
	}
}

// UnmarshalPrimitive unmarshal bytes to primitive
func UnmarshalPrimitive(str string, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch k := v.Kind(); k {
	case reflect.Bool:
		v.SetBool(str[0] == 't')
	case reflect.String:
		v.SetString(str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(str, 10, 64)
		if err == nil {
			return err
		}

		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(str, v.Type().Bits())
		if err == nil {
			return err
		}

		v.SetFloat(n)
	default:
		return ErrNonPrimitive
	}

	return nil
}

func MarshalPrimitive(v reflect.Value) string {
	return fmt.Sprint(v.Interface())
}

func MarshalStruct(v interface{}, values map[string]string, tag string) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		vfield := value.Field(i)
		if !vfield.CanInterface() {
			continue
		}

		tfield := typ.Field(i)
		name := tfield.Name
		if n := tfield.Tag.Get(tag); n == "-" {
			continue
		} else if n != "" {
			name = n
		} else {
			name = strings.ToLower(name)
		}

		values[name] = MarshalPrimitive(vfield)
	}
}

type Values interface {
	Get(name string) string
}

type StringMap map[string]string

func (m StringMap) Get(name string) string {
	return m[name]
}

type StringSliceMap struct {
	Values    map[string][]string
	Seperator string
}

func (s StringSliceMap) Get(name string) string {
	return strings.Join(s.Values[name], s.Seperator)
}

func UnmarshalStruct(v interface{}, values Values, tag string) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	} else {
		panic("non-pointer type")
	}

	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		vfield := value.Field(i)
		if !vfield.CanSet() {
			continue
		}
		tfield := typ.Field(i)
		name := tfield.Name
		if n := tfield.Tag.Get(tag); n == "-" {
			continue
		} else if n != "" {
			name = n
		} else {
			name = strings.ToLower(name)
		}
		UnmarshalPrimitive(values.Get(name), vfield)
	}

	return
}

func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	switch val := reflect.ValueOf(v); val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return val.IsNil()
	case reflect.Invalid:
		return fmt.Sprint(v) == "<nil>"
	default:
		return false
	}
}

func TruncSliceCapToLen(vals ...interface{}) {
	for _, val := range vals {
		ptrVal := reflect.ValueOf(val)
		if kind := ptrVal.Kind(); kind != reflect.Ptr {
			panic(errors.Newf("expect pointer to slice, but got %s", kind.String()))
		}
		sliVal := ptrVal.Elem()
		if kind := sliVal.Kind(); kind != reflect.Slice {
			panic(errors.Newf("expect pointer to slice, but got pointer of %s", kind.String()))
		}
		len, cap := sliVal.Len(), sliVal.Cap()
		if len == cap {
			continue
		}
		newVal := reflect.MakeSlice(sliVal.Type(), len, len)
		reflect.Copy(newVal, sliVal)
		sliVal.Set(newVal)
	}
}
