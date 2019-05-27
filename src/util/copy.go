package util

import (
	"errors"
	"reflect"
	"time"
)

func Copy(src interface{}, dst interface{}) (err error) {
	dstValue := reflect.ValueOf(dst)
	dstType := dstValue.Type()

	set := false
	if dstType.Kind() != reflect.Ptr {
		return errors.New("Prune(non-pointer)")
	}

	if dstValue.IsNil() {
		return errors.New("Prune(nil " + dstType.String() + ")")
	}

	//only for struct
	if dstValue.Elem().Kind() != reflect.Struct {
		return errors.New("Prune(" + dstValue.String() + " is not a struct ) is" + dstValue.Elem().Kind().String())
	}

	srcValueElem := reflect.ValueOf(src).Elem()

	dstValueElem := dstValue.Elem()
	dstTypeElem := dstType.Elem()

	for i := 0; i < dstValueElem.NumField(); i++ {
		dstValue := dstValueElem.Field(i)

		srcValue := srcValueElem.FieldByName(dstTypeElem.Field(i).Name)
		if !srcValue.IsValid() {
			continue
		}

		if srcValue.Type() == dstValue.Type() {
			if dstValue.CanSet() {
				dstValue.Set(srcValue)
				set = true
			}
		} else {
			if value, ok := srcValue.Interface().(time.Time); ok {
				switch dstValue.Type().Kind() {
				case reflect.Int64, reflect.Int32:
					dstValue.SetInt(value.Unix())
					set = true
				case reflect.Uint64, reflect.Uint32:
					dstValue.SetUint(uint64(value.Unix()))
					set = true
				default:
					return errors.New("Prune " + dstValue.String() + ":" + dstValue.Type().Name() + " type is invalid")
				}
			}
		}
	}

	if !set {
		return errors.New("Nothing has been set on " + dstType.String() + "<-" + reflect.ValueOf(src).String())
	}
	return nil
}
