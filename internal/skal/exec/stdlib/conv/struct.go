package conv

import (
	"fmt"
	"log/slog"
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

func StructToLTable(value any, l *lua.LState) *lua.LTable {
	var (
		t = l.NewTable()
		v = reflect.ValueOf(value)
	)

	for i := range v.NumField() {
		f := v.Field(i)
		fn := v.Type().Field(i).Name

		if !v.Type().Field(i).IsExported() {
			slog.Debug(
				"Skipping unsettable field.",
				"field", fn,
			)
			continue
		}

		fv := f.Interface()
		var next lua.LValue
		switch f.Kind() {
		case reflect.String:
			next = lua.LString(fv.(string))

		case reflect.Int:
			next = lua.LNumber(fv.(int))
		case reflect.Int8:
			next = lua.LNumber(fv.(int8))
		case reflect.Int16:
			next = lua.LNumber(fv.(int16))
		case reflect.Int32:
			next = lua.LNumber(fv.(int32))
		case reflect.Int64:
			next = lua.LNumber(fv.(int64))

		case reflect.Bool:
			next = lua.LBool(fv.(bool))

		case reflect.Uint:
			next = lua.LNumber(fv.(uint))
		case reflect.Uint8:
			next = lua.LNumber(fv.(uint8))
		case reflect.Uint16:
			next = lua.LNumber(fv.(uint16))
		case reflect.Uint32:
			next = lua.LNumber(fv.(uint32))
		case reflect.Uint64:
			next = lua.LNumber(fv.(uint64))

		case reflect.Float32:
			next = lua.LNumber(fv.(float32))
		case reflect.Float64:
			next = lua.LNumber(fv.(float64))

		case reflect.Struct:
			next = StructToLTable(fv, l)

		default:
			slog.Error(
				"Failed to convert SToT type.",
				"type", fmt.Sprintf("T: %s\n", f.Kind()),
			)
			continue
		}

		l.SetField(t, fn, next)
	}

	return t
}
