//
// Copyright (C) 2013 Space Monkey, Inc.
//

package flagfile

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// find the width of the int/uint types
const (
	_m = ^uint(0)
	_W = 32 * int(_m>>31&1+_m>>63&1)
)

func getParams(tags reflect.StructTag) (defaultStr, usage string) {
	defaultStr, usage = tags.Get("default"), tags.Get("usage")
	return
}

func Setup(prefix string, x interface{}) {
	rv := reflect.ValueOf(x).Elem()
	if rv.Kind() != reflect.Struct {
		panic("flagfile: setup must be passed pointer to struct")
	}
	typ := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if !field.CanAddr() {
			continue
		}
		ftyp := typ.Field(i)
		flagName := prefix + strings.ToLower(ftyp.Name)
		defaultStr, usage := getParams(ftyp.Tag)
		switch fvar := field.Addr().Interface().(type) {
		case *bool:
			defaultVal, err := strconv.ParseBool(defaultStr)
			if err != nil {
				panic(fmt.Sprintf("flagfile: %s", err))
			}
			flag.BoolVar(fvar, flagName, defaultVal, usage)

		case *time.Duration:
			defaultVal, err := time.ParseDuration(defaultStr)
			if err != nil {
				panic(fmt.Sprintf("flagfile: %s", err))
			}
			flag.DurationVar(fvar, flagName, defaultVal, usage)

		case *float64:
			defaultVal, err := strconv.ParseFloat(defaultStr, 64)
			if err != nil {
				panic(fmt.Sprintf("flagfile: %s", err))
			}
			flag.Float64Var(fvar, flagName, defaultVal, usage)

		case *int:
			defaultVal, err := strconv.ParseInt(defaultStr, 10, _W)
			if err != nil {
				panic(fmt.Sprintf("flagfile: %s", err))
			}
			flag.IntVar(fvar, flagName, int(defaultVal), usage)

		case *int64:
			defaultVal, err := strconv.ParseInt(defaultStr, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("flagfile: %s", err))
			}
			flag.Int64Var(fvar, flagName, defaultVal, usage)

		case *string:
			flag.StringVar(fvar, flagName, defaultStr, usage)

		case *uint:
			defaultVal, err := strconv.ParseUint(defaultStr, 10, _W)
			if err != nil {
				panic(fmt.Sprintf("flagfile: %s", err))
			}
			flag.UintVar(fvar, flagName, uint(defaultVal), usage)

		case *uint64:
			defaultVal, err := strconv.ParseUint(defaultStr, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("flagfile: %s", err))
			}
			flag.Uint64Var(fvar, flagName, defaultVal, usage)

		default:
			panic(fmt.Sprintf(
				"flagfile: can't set flag value for type %s", field.Type()))
		}
	}
}
