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

package test_runner

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func capturePanic(operation func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	operation()
	return
}

func runAndWrapPanic(operation func(), format string, args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			message := fmt.Sprintf(format, args...)
			switch e := r.(type) {
			case error:
				panic(fmt.Errorf("%v: %w", message, e))
			default:
				panic(fmt.Errorf("%v: %v", message, r))
			}
		}
	}()

	operation()
}

func wrapPanic(format string, args ...interface{}) {
	if r := recover(); r != nil {
		message := fmt.Sprintf(format, args...)

		switch e := r.(type) {
		case error:
			panic(fmt.Errorf("%v: %w", message, e))
		default:
			panic(fmt.Errorf("%v: %v", message, r))
		}
	}
}

func reportAnyError(r interface{}, trace bool, format string, args ...interface{}) (reportedError bool) {
	if r == nil {
		return
	}

	message := fmt.Sprintf(format, args...)
	fmt.Printf("%v: %v\n", message, r)
	if trace {
		fmt.Println(string(debug.Stack()))
	}
	return true
}

func desc(v interface{}) string {
	switch v := v.(type) {
	case string:
		return fmt.Sprintf("[%v]", v)
	case []byte:
		sb := strings.Builder{}
		sb.WriteByte('[')
		for i, b := range v {
			sb.WriteString(fmt.Sprintf("%02x", b))
			if i < len(v)-1 {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte(']')
		return sb.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toMembershipSet(valueArrays ...[]string) map[string]bool {
	set := make(map[string]bool)

	for _, array := range valueArrays {
		for _, value := range array {
			set[value] = true
		}
	}
	return set
}
