// Copyright 2019 Karl Stenerud
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package common

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/kstenerud/go-concise-encoding/configuration"
)

func DecodeGoTags(field reflect.StructField) (tags GoTags) {
	requiresValue := func(kv []string, key string) {
		if len(kv) != 2 {
			panic(fmt.Errorf(`tag key "%s" requires a value`, key))
		}
	}

	tags.Name = field.Name
	tags.Order = math.MaxInt64

	tagString := strings.TrimSpace(field.Tag.Get("ce"))
	if len(tagString) == 0 {
		return
	}

	for _, entry := range strings.Split(tagString, ",") {
		kv := strings.Split(entry, "=")
		switch strings.TrimSpace(kv[0]) {
		/* TODO:
		 * - lowercase/origcase
		 * - omit specific value?
		 * - recurse/no_recurse?
		 * - type=f16, f10.x, i2, i8, i10, i16, string, vstring?
		 */
		case "omit":
			tags.OmitBehavior = configuration.OmitFieldAlways
		case "omit_empty":
			tags.OmitBehavior = configuration.OmitFieldEmpty
		case "omit_zero":
			tags.OmitBehavior = configuration.OmitFieldZero
		case "omit_never":
			tags.OmitBehavior = configuration.OmitFieldNever
		case "name":
			requiresValue(kv, "name")
			tags.Name = strings.TrimSpace(kv[1])
		case "order":
			order, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err != nil {
				panic(err)
			}
			tags.Order = order
		default:
			panic(fmt.Errorf("%v: Unknown Concise Encoding struct tag field decoding [%v] in field %v", entry, tagString, field.Name))
		}
	}

	return
}

type GoTags struct {
	Name         string
	OmitBehavior configuration.FieldOmitBehavior
	Order        int64
}
