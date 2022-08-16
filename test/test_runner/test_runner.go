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
	"github.com/kstenerud/go-concise-encoding/options"
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

func wrapPanic(recovery interface{}, format string, args ...interface{}) {
	if recovery != nil {
		message := fmt.Sprintf(format, args...)

		switch e := recovery.(type) {
		case error:
			panic(fmt.Errorf("%v: %w", message, e))
		default:
			panic(fmt.Errorf("%v: %v", message, e))
		}
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
		if v != "ce-test" {
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
	Events      []string
	events      test.Events
	Debug       bool // When true, don't convert panics to errors.
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

	_this.events = event_parser.ParseEvents(_this.Events...)
	if !_this.RawDocument && len(_this.events) > 0 {
		newEvents := make(test.Events, 0, len(_this.events)+1)
		newEvents = append(newEvents, test.V(uint64(_this.ceVersion)))
		newEvents = append(newEvents, _this.events...)
		_this.events = newEvents
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
	LossyCTE    bool
	LossyCBE    bool
	LossyEvents bool
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

	return nil
}

func (_this *MustSucceedTest) runCte() error {
	hasEvents := len(_this.events) > 0

	if _this.Debug {
		fmt.Printf("%v: Convert CTE to events: [%v]", _this.context, _this.Cte)
	}
	events, err := _this.cteToEvents(_this.Cte)
	if err != nil {
		return _this.wrapError(err, "decoding CTE [%v]", _this.Cte)
	}
	if hasEvents && !_this.LossyCTE {
		if !_this.events.AreEquivalentTo(events) {
			return _this.errorf("expected CTE [%v] to produce events [%v] but got [%v]",
				_this.Cte, _this.events, events)
		}
	}
	if _this.Debug {
		fmt.Printf("%v: Convert events to CTE: [%v]", _this.context, _this.events)
	}
	if hasEvents {
		events = _this.events
	}
	document, err := _this.eventsToCte(events)
	if err != nil {
		return _this.wrapError(err, "Encoding events [%v] to CTE", events)
	}
	if !_this.LossyCTE {
		if _this.Cte != document {
			return _this.errorf("re-encoding events [%v] from CTE [%v] produced unexpected CTE [%v]",
				events, _this.Cte, document)
		}
	}

	return nil
}

func (_this *MustSucceedTest) runCbe() error {
	hasEvents := len(_this.events) > 0

	if _this.Debug {
		fmt.Printf("%v: Convert CBE to events: [%v]", _this.context, asHex(_this.Cbe))
	}
	events, err := _this.cbeToEvents(_this.Cbe)
	if err != nil {
		return _this.wrapError(err, "decoding CBE [%v]", asHex(_this.Cbe))
	}
	if hasEvents && !_this.LossyCBE {
		if !_this.events.AreEquivalentTo(events) {
			return _this.errorf("expected CBE [%v] to produce events [%v] but got [%v]",
				asHex(_this.Cbe), _this.events, events)
		}
	}
	if _this.Debug {
		fmt.Printf("%v: Convert events to CBE: [%v]", _this.context, _this.events)
	}
	if hasEvents {
		events = _this.events
	}
	document, err := _this.eventsToCbe(events)
	if err != nil {
		return _this.wrapError(err, "Encoding events [%v] to CBE", events)
	}
	if !_this.LossyCBE {
		if !bytes.Equal(_this.Cbe, document) {
			return _this.errorf("re-encoding events [%v] from CBE [%v] produced unexpected CBE [%v]",
				events, asHex(_this.Cbe), asHex(document))
		}
	}

	return nil
}

func (_this *MustSucceedTest) Run() error {
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
		if err := _this.runCte(); err != nil {
			return err
		}
	}

	if len(_this.Cbe) > 0 {
		if err := _this.runCbe(); err != nil {
			return err
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

	// A bit hacky, but it makes things easier to have Debug in the base class.
	_this.Debug = false

	context = fmt.Sprintf(`%v, "must fail" test #%v`, context, index+1)
	if err := _this.BaseTest.PostDecodeInit(ceVersion, context); err != nil {
		return err
	}

	total := 0
	if len(_this.Cte) > 0 {
		total++
	}
	if len(_this.Cbe) > 0 {
		total++
	}
	if len(_this.Events) > 0 {
		total++
	}

	if total != 1 {
		return _this.errorf("must have one and only one of: cbe, cte, events")
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
		if events, err := _this.cteToEvents(_this.Cte); err == nil {
			return _this.errorf("expected CTE document [%v] to fail but it produced [%v]",
				_this.Cte, events)
		}
	} else {
		if events, err := _this.cbeToEvents(_this.Cbe); err == nil {
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

func (_this *BaseTest) cteToEvents(document string) (result test.Events, err error) {
	opts := options.DefaultCEDecoderOptions()
	opts.DebugPanics = _this.Debug
	decoder := cte.NewDecoder(&opts)
	buffer := bytes.NewBuffer([]byte(document))
	receiver, collection := test.NewEventCollector(nil)
	receiver = rules.NewRules(receiver, nil)
	if err = decoder.Decode(buffer, receiver); err != nil {
		return
	}
	result = collection.Events
	return
}

func (_this *BaseTest) cbeToEvents(document []byte) (result test.Events, err error) {
	opts := options.DefaultCEDecoderOptions()
	opts.DebugPanics = _this.Debug
	decoder := cbe.NewDecoder(&opts)
	buffer := bytes.NewBuffer(document)
	receiver, collection := test.NewEventCollector(nil)
	receiver = rules.NewRules(receiver, nil)
	if err = decoder.Decode(buffer, receiver); err != nil {
		return
	}
	result = collection.Events
	return
}

func (_this *BaseTest) eventsToCte(events test.Events) (result string, err error) {
	if !_this.Debug {
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
	}

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

func (_this *BaseTest) eventsToCbe(events test.Events) (result []byte, err error) {
	if !_this.Debug {
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
	}

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
