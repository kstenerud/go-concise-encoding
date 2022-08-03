// Copyright 2022 Karl Stenerud
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
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/event_parser"
	"github.com/kstenerud/go-concise-encoding/version"
)

func TestStuff(t *testing.T) {
	RunTests(t, "test.cte")
}

func RunTests(t *testing.T, sourceFile string) {
	suite, errors := loadTestSuite(sourceFile)
	if len(errors) > 0 {
		fmt.Printf("❌ %v\n", sourceFile)
		reportFailedTestLoad(errors)
		t.Fail()
		return
	}

	errors = suite.Run()
	if len(errors) > 0 {
		fmt.Printf("❌ %v\n", sourceFile)
		reportTestFailures(errors)
		t.Fail()
		return
	}
	fmt.Printf("✅ %v\n", sourceFile)
}

func reportFailedTestLoad(errors []error) {
	for _, err := range errors {
		fmt.Printf("  ❗ %v\n", err)
	}
}

func reportTestFailures(errors []error) {
	for _, err := range errors {
		fmt.Printf("  💣 %v\n", err)
	}
}

func loadTestSuite(testDescriptorFile string) (suite *TestSuite, errors []error) {
	defer func() { wrapPanic(recover(), "while loading test suite from %v", testDescriptorFile) }()

	file, err := os.Open(testDescriptorFile)
	if err != nil {
		panic(fmt.Errorf("unexpected error opening test suite file %v: %w", testDescriptorFile, err))
	}

	loadedTest, err := ce.UnmarshalCTE(file, suite, nil)
	if err != nil {
		panic(fmt.Errorf("malformed unit test: Unexpected CTE decode error in test suite file %v: %w", testDescriptorFile, err))
	}
	suite = loadedTest.(*TestSuite)
	errors = suite.PostDecodeInit(testDescriptorFile)

	return
}

const TestSuiteVersion = 1

type TestSuite struct {
	Type      map[string]interface{}
	CEVersion *int
	Tests     []UnitTest
	context   string
}

func (_this *TestSuite) PostDecodeInit(sourceFile string) (errors []error) {
	_this.context = sourceFile

	if _this.Type == nil {
		return []error{_this.errorf("missing type field")}
	}

	if v, ok := _this.Type["identifier"]; ok {
		if v != "ce test" {
			return []error{_this.errorf("%v: unrecognized type identifier", v)}
		}
	} else {
		return []error{_this.errorf("missing type identifier field")}
	}

	if v, ok := _this.Type["version"]; ok {
		if v != uint64(TestSuiteVersion) {
			if reflect.TypeOf(v).Kind() != reflect.Uint64 {
				return []error{_this.errorf("type version must be an unsigned integer")}
			}
			return []error{_this.errorf("ce test format version %v runner cannot load from version %v file", TestSuiteVersion, v)}
		}
	} else {
		return []error{_this.errorf("missing type version field")}
	}

	if _this.CEVersion == nil {
		return []error{_this.errorf("missing ceversion field")}
	}

	if *_this.CEVersion != version.ConciseEncodingVersion {
		return []error{_this.errorf("This codec is for Concise Encoding version %v and cannot run tests for version %v",
			version.ConciseEncodingVersion, _this.CEVersion)}
	}

	for index, test := range _this.Tests {
		if nextErrors := test.PostDecodeInit(*_this.CEVersion, _this.context, index); nextErrors != nil {
			errors = append(errors, nextErrors...)
		}
	}
	return
}

func (_this *TestSuite) Run() (errors []error) {
	for _, test := range _this.Tests {
		if nextErrors := test.Run(); nextErrors != nil {
			errors = append(errors, nextErrors...)
		}
	}
	return
}

func (_this *TestSuite) errorf(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return fmt.Errorf("%v: %v", _this.context, message)
}

type UnitTest struct {
	Name        string
	MustSucceed []*MustSucceedTest
	MustFail    []*MustFailTest
}

func (_this *UnitTest) PostDecodeInit(ceVersion int, context string, testIndex int) (errors []error) {
	context = fmt.Sprintf("%v, unit test #%v (%v)", context, testIndex+1, _this.Name)

	for index, test := range _this.MustSucceed {
		if err := test.PostDecodeInit(ceVersion, context, index); err != nil {
			errors = append(errors, err)
		}
	}

	for index, test := range _this.MustFail {
		if err := test.PostDecodeInit(ceVersion, context, index); err != nil {
			errors = append(errors, err)
		}
	}

	return
}

func (_this *UnitTest) Run() (errors []error) {
	for _, test := range _this.MustSucceed {
		if err := test.Run(); err != nil {
			errors = append(errors, err)
		}
	}

	for _, test := range _this.MustFail {
		if err := test.Run(); err != nil {
			errors = append(errors, err)
		}
	}

	return
}

type BaseTest struct {
	RawDocument bool
	Skip        bool
	ceVersion   int
	Cbe         []byte
	Cte         string
	context     string
}

func (_this *BaseTest) PostDecodeInit(ceVersion int, context string) error {
	if _this.Skip {
		return nil
	}

	_this.ceVersion = ceVersion
	_this.context = context
	_this.Cte = strings.TrimSpace(_this.Cte)
	if !_this.RawDocument {
		if len(_this.Cte) > 0 {
			_this.Cte = fmt.Sprintf("c%v\n%v", _this.ceVersion, _this.Cte)
		}
		if len(_this.Cbe) > 0 {
			b := make([]byte, len(_this.Cbe)+2)
			b[0] = 0x81
			b[1] = byte(_this.ceVersion)
			copy(b[2:], _this.Cbe)
			_this.Cbe = b
		}
	}

	return nil
}

func (_this *BaseTest) errorf(format string, args ...interface{}) error {
	if message := fmt.Sprintf(format, args...); len(message) > 0 {
		return fmt.Errorf("%v: %v", _this.context, message)
	} else {
		return fmt.Errorf("%v", _this.context)
	}
}

func (_this *BaseTest) wrapError(err error, format string, args ...interface{}) error {
	if message := fmt.Sprintf(format, args...); len(message) > 0 {
		return fmt.Errorf("%v: %v: %w", _this.context, message, err)
	} else {
		return fmt.Errorf("%v: %w", _this.context, err)
	}
}

type MustSucceedTest struct {
	BaseTest
	NoCteOutput bool
	Debug       bool
	Events      []string
	events      test.Events
}

func (_this *MustSucceedTest) PostDecodeInit(ceVersion int, context string, index int) error {
	if _this.Skip {
		return nil
	}
	context = fmt.Sprintf(`%v, "must succeed" test #%v`, context, index+1)
	if err := _this.BaseTest.PostDecodeInit(ceVersion, context); err != nil {
		return err
	}

	if len(_this.Cbe) == 0 && len(_this.Cte) == 0 {
		return _this.errorf("must have cbe and/or cte")
	}

	_this.events = event_parser.ParseEvents(_this.Events...)
	if !_this.RawDocument && len(_this.events) > 0 {
		newEvents := make(test.Events, 0, len(_this.events)+1)
		newEvents = append(newEvents, test.V(uint64(_this.ceVersion)))
		newEvents = append(newEvents, _this.events...)
		_this.events = newEvents
	}
	return nil
}

func (_this *MustSucceedTest) Run() error {
	if _this.Skip {
		return nil
	}

	if _this.Debug {
		debug.DebugOptions.PassThroughPanics = true
		defer func() {
			debug.DebugOptions.PassThroughPanics = false
		}()
	}

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("%v: %w", _this.context, v))
			default:
				panic(fmt.Errorf("%v: %v", _this.context, v))
			}
		}
	}()

	hasCte := len(_this.Cte) > 0
	hasCbe := len(_this.Cbe) > 0
	hasEvents := len(_this.events) > 0

	if hasCte {
		events, err := cteToEvents(_this.Cte)
		if err != nil {
			return _this.wrapError(err, "decoding CTE [%v]", _this.Cte)
		}
		if hasEvents {
			if !_this.events.AreEquivalentTo(events) {
				return _this.errorf("expected CTE [%v] to produce events [%v] but got [%v]",
					_this.Cte, _this.events, events)
			}
		}
		if !_this.NoCteOutput {
			document, err := eventsToCte(events)
			if err != nil {
				return _this.wrapError(err, "Encoding events [%v] to CTE", events)
			}
			if _this.Cte != document {
				return _this.errorf("re-encoding CTE [%v] produced unexpected CTE [%v]",
					_this.Cte, document)
			}
		}
		if hasCbe {
			document, err := eventsToCbe(events)
			if err != nil {
				return _this.wrapError(err, "Encoding events [%v] to CBE", events)
			}
			if !bytes.Equal(_this.Cbe, document) {
				return _this.errorf("expected CTE [%v] to encode to CBE [%v] but got [%v]",
					_this.Cte, asHex(_this.Cbe), asHex(document))
			}
		}
	}

	if hasCbe {
		events, err := cbeToEvents(_this.Cbe)
		if err != nil {
			return _this.wrapError(err, "decoding CBE [%v]", asHex(_this.Cbe))
		}
		if hasEvents {
			if !_this.events.AreEquivalentTo(events) {
				return _this.errorf("expected CBE [%v] to produce events [%v] but got [%v]",
					_this.Cbe, _this.events, events)
			}
		}
		document, err := eventsToCbe(events)
		if err != nil {
			return _this.wrapError(err, "Encoding events [%v] to CBE", events)
		}
		if !bytes.Equal(_this.Cbe, document) {
			return _this.errorf("re-encoding CBE [%v] produced unexpected CBE [%v]",
				asHex(_this.Cbe), asHex(document))
		}
		if hasCte && !_this.NoCteOutput {
			document, err := eventsToCte(events)
			if err != nil {
				return _this.wrapError(err, "Encoding events [%v] to CTE", events)
			}
			if _this.Cte != document {
				return _this.errorf("expected CBE [%v] to encode to CTE [%v] but got [%v]",
					asHex(_this.Cbe), _this.Cte, document)
			}
		}
	}
	return nil
}

type MustFailTest struct {
	BaseTest
}

func (_this *MustFailTest) PostDecodeInit(ceVersion int, context string, index int) error {
	if _this.Skip {
		return nil
	}
	context = fmt.Sprintf(`%v, "must fail" test #%v`, context, index+1)
	if err := _this.BaseTest.PostDecodeInit(ceVersion, context); err != nil {
		return err
	}

	hasCte := len(_this.Cte) > 0
	hasCbe := len(_this.Cbe) > 0

	if (hasCte && hasCbe) || (!hasCte && !hasCbe) {
		return _this.errorf("must have either cbe or cte (not both)")
	}
	return nil
}

func (_this *MustFailTest) Run() error {
	if _this.Skip {
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("%v: %w", _this.context, v))
			default:
				panic(fmt.Errorf("%v: %v", _this.context, v))
			}
		}
	}()

	if len(_this.Cte) > 0 {
		if events, err := cteToEvents(_this.Cte); err == nil {
			return _this.errorf("expected CTE document [%v] to fail but it produced [%v]",
				_this.Cte, events)
		}
	} else {
		if events, err := cbeToEvents(_this.Cbe); err == nil {
			return _this.errorf("expected CBE document [%v] to fail but it produced [%v]",
				asHex(_this.Cbe), events)
		}
	}
	return nil
}

func asHex(data []byte) string {
	sb := strings.Builder{}
	for i, b := range data {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%02x", b))
	}
	return sb.String()
}

func cteToEvents(document string) (result test.Events, err error) {
	decoder := cte.NewDecoder(nil)
	buffer := bytes.NewBuffer([]byte(document))
	receiver, collection := test.NewEventCollector(nil)
	receiver = rules.NewRules(receiver, nil)
	if err = decoder.Decode(buffer, receiver); err != nil {
		return
	}
	result = collection.Events
	return
}

func cbeToEvents(document []byte) (result test.Events, err error) {
	decoder := cbe.NewDecoder(nil)
	buffer := bytes.NewBuffer(document)
	receiver, collection := test.NewEventCollector(nil)
	receiver = rules.NewRules(receiver, nil)
	if err = decoder.Decode(buffer, receiver); err != nil {
		return
	}
	result = collection.Events
	return
}

func eventsToCte(events test.Events) (result string, err error) {
	encoder := cte.NewEncoder(nil)
	outBuffer := &bytes.Buffer{}
	encoder.PrepareToEncode(outBuffer)
	receiver := rules.NewRules(encoder, nil)
	receiver.OnBeginDocument()
	for _, event := range events {
		event.Invoke(receiver)
	}
	receiver.OnEndDocument()
	result = outBuffer.String()
	return
}

func eventsToCbe(events test.Events) (result []byte, err error) {
	encoder := cbe.NewEncoder(nil)
	outBuffer := &bytes.Buffer{}
	encoder.PrepareToEncode(outBuffer)
	receiver := rules.NewRules(encoder, nil)
	receiver.OnBeginDocument()
	for _, event := range events {
		event.Invoke(receiver)
	}
	receiver.OnEndDocument()
	result = outBuffer.Bytes()
	return
}
