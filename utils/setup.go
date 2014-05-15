// Copyright (C) 2014 Space Monkey, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package utils provides a collection of nice flag/flagfile helpers without
requiring someone to actually use flagfile.Load() over flag.Parse()
*/
package utils

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

// Setup is meant to be called in init or in some other place prior to
// flagfile.Load() or flags.Parse(). Setup will take a struct and through
// reflect, create flags for all of the fields in the struct. The given prefix
// is added to all of the discovered flags. Setup doesn't return any value, but
// sets up the flag system to mutate the given struct when
// flag.Parse/flagfile.Load is called.
//
// If a prefix is non-zero length, it is joined with a separating period.
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
		flagName := strings.ToLower(ftyp.Name)
		if len(prefix) > 0 {
			flagName = fmt.Sprintf("%s.%s", prefix, flagName)
		}
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
