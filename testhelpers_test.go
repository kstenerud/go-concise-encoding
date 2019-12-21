package cbe

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/kstenerud/go-cbe/rules"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

// ============================================================================
// Types
// ============================================================================

const (
	IDNil uint32 = 1000 + iota
	IDList
	IDEList
	IDMap
	IDEMap
	IDCmt
	IDECmt
	IDMeta
	IDEMeta
	IDMarkup
	IDEMattr
	IDEMarkup
	IDEnd
	IDMarker
	IDRef
	IDPad
	IDBytes
	IDStr
	IDURI
	IDEArr
)

type CBEToken interface {
	EncodeCBE(encoder *CBEEncoder) error
}

type NilToken int

func (this NilToken) String() string { return "Nil" }

func (this NilToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.Nil() }

func Nil() NilToken { return NilToken(IDNil) }

// ------------------------------------

type ListToken int

func (this ListToken) String() string { return "List" }

func (this ListToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.BeginList() }

func list() ListToken { return ListToken(IDList) }

// ------------------------------------

type EListToken int

func (this EListToken) String() string { return "EList" }

func (this EListToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.EndContainer() }

func elist() EListToken { return EListToken(IDEList) }

// ------------------------------------

type MapToken int

func (this MapToken) String() string { return "Map" }

func (this MapToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.BeginMap() }

func Map() MapToken { return MapToken(IDMap) }

// ------------------------------------

type EMapToken int

func (this EMapToken) String() string { return "EMap" }

func (this EMapToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.EndContainer() }

func emap() EMapToken { return EMapToken(IDEMap) }

// ------------------------------------

type CmtToken int

func (this CmtToken) String() string { return "Cmt" }

func (this CmtToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.BeginComment() }

func cmt() CmtToken { return CmtToken(IDCmt) }

// ------------------------------------

type ECmtToken int

func (this ECmtToken) String() string { return "ECmt" }

func (this ECmtToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.EndContainer() }

func ecmt() ECmtToken { return ECmtToken(IDECmt) }

// ------------------------------------

type MetaToken int

func (this MetaToken) String() string { return "Meta" }

func (this MetaToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.BeginMetadata() }

func meta() MetaToken { return MetaToken(IDMeta) }

// ------------------------------------

type EMetaToken int

func (this EMetaToken) String() string { return "EMeta" }

func (this EMetaToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.EndContainer() }

func emeta() EMetaToken { return EMetaToken(IDEMeta) }

// ------------------------------------

type MarkupToken int

func (this MarkupToken) String() string { return "Markup" }

func (this MarkupToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.BeginMarkup() }

func markup() MarkupToken { return MarkupToken(IDMarkup) }

// ------------------------------------

type EMattrToken int

func (this EMattrToken) String() string { return "EMattr" }

func (this EMattrToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.EndContainer() }

func emattr() EMattrToken { return EMattrToken(IDEMattr) }

// ------------------------------------

type EMarkupToken int

func (this EMarkupToken) String() string { return "EMarkup" }

func (this EMarkupToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.EndContainer() }

func emarkup() EMarkupToken { return EMarkupToken(IDEMarkup) }

// ------------------------------------

type EndToken int

func (this EndToken) String() string { return "End" }

func (this EndToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.EndContainer() }

func end() EndToken { return EndToken(IDEnd) }

// ------------------------------------

type MarkerToken int

func (this MarkerToken) String() string { return "Marker" }

func (this MarkerToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.BeginMarker() }

func marker() MarkerToken { return MarkerToken(IDMarker) }

// ------------------------------------

type RefToken int

func (this RefToken) String() string { return "Ref" }

func (this RefToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.BeginReference() }

func ref() RefToken { return RefToken(IDRef) }

// ------------------------------------

type PadToken int

func (this PadToken) String() string { return "Pad" }

func (this PadToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.Padding(1) }

func pad() PadToken { return PadToken(IDPad) }

// ------------------------------------

type BytesToken int

func (this BytesToken) String() string { return "Bytes" }

func (this BytesToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.BeginBytes() }

func bin() BytesToken { return BytesToken(IDBytes) }

// ------------------------------------

type StrToken int

func (this StrToken) String() string { return "Str" }

func (this StrToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.BeginString() }

func str() StrToken { return StrToken(IDStr) }

// ------------------------------------

type URIToken int

func (this URIToken) String() string { return "URI" }

func (this URIToken) EncodeCBE(encoder *CBEEncoder) error { return encoder.BeginURI() }

func uri() URIToken { return URIToken(IDURI) }

// ------------------------------------

type NIntType uint64

func (this NIntType) String() string { return fmt.Sprintf("-%v", uint64(this)) }

func (this NIntType) EncodeCBE(encoder *CBEEncoder) error { return encoder.NegativeInt(uint64(this)) }

func neg(value uint64) NIntType { return NIntType(value) }

// ------------------------------------

type ChunkType uint64

var flagFinalChunk ChunkType = 0x1000000000000000

func (this ChunkType) String() string {
	isFinalChunk := this&flagFinalChunk != 0
	length := uint64(this & ^flagFinalChunk)
	return fmt.Sprintf("Chunk(%v, %v)", length, isFinalChunk)
}

func (this ChunkType) EncodeCBE(encoder *CBEEncoder) error {
	isFinalChunk := this&flagFinalChunk != 0
	length := uint64(this & ^flagFinalChunk)

	return encoder.BeginChunk(length, isFinalChunk)
}

func chunk(size uint64, isFinalChunk bool) ChunkType {
	result := ChunkType(size)
	if isFinalChunk {
		result |= flagFinalChunk
	}
	return result
}

// ------------------------------------

type EArrToken int

func (this EArrToken) String() string { return "EArr" }

func earr() EArrToken { return EArrToken(IDEArr) }

// ------------------------------------

func binBytes(values ...byte) []byte {
	return values
}

func strBytes(str string) []byte {
	return []byte(str)
}

func uriBytes(str string) []byte {
	_, err := url.Parse(str)
	if err != nil {
		_, err = url.Parse("http://parse.error")
	}
	return []byte(str)
}

func Date(year int, month int, day int) *compact_time.Time {
	return compact_time.NewDate(year, month, day)
}

func Time(hour int, minute int, second int, nanosecond int, areaLocation string) *compact_time.Time {
	return compact_time.NewTime(hour, minute, second, nanosecond, areaLocation)
}

func TimeLoc(hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths int) *compact_time.Time {
	return compact_time.NewTimeLatLong(hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths)
}

func TS(year, month, day, hour, minute, second, nanosecond int, areaLocation string) *compact_time.Time {
	return compact_time.NewTimestamp(year, month, day, hour, minute, second, nanosecond, areaLocation)
}

func TSLoc(year, month, day, hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths int) *compact_time.Time {
	return compact_time.NewTimestampLatLong(year, month, day, hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths)
}

// ------------------------------------

var stringGeneratorChars = [...]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
}

func genString(length int) string {
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteByte(stringGeneratorChars[i%len(stringGeneratorChars)])
	}
	return result.String()
}

func genBytes(length int) []byte {
	return []byte(genString(length))
}

func getPanicContents(function func()) (recovered interface{}) {
	defer func() {
		recovered = recover()
	}()
	function()
	return recovered
}

func ShortCircuit(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

func Tokens(tokens ...interface{}) []interface{} {
	return tokens
}

// ============================================================================
// Decoder
// ============================================================================

type TokenDecoder struct {
	Tokens []interface{}
}

func NewTokenDecoder() *TokenDecoder {
	return new(TokenDecoder)
}

func (this *TokenDecoder) OnNil() error {
	this.Tokens = append(this.Tokens, Nil())
	return nil
}
func (this *TokenDecoder) OnBool(value bool) error {
	this.Tokens = append(this.Tokens, value)
	return nil
}
func (this *TokenDecoder) OnPositiveInt(value uint64) error {
	this.Tokens = append(this.Tokens, value)
	return nil
}
func (this *TokenDecoder) OnNegativeInt(value uint64) error {
	this.Tokens = append(this.Tokens, neg(value))
	return nil
}
func (this *TokenDecoder) OnFloat(value float64) error {
	this.Tokens = append(this.Tokens, value)
	return nil
}
func (this *TokenDecoder) OnTime(value *compact_time.Time) error {
	this.Tokens = append(this.Tokens, value)
	return nil
}
func (this *TokenDecoder) OnListBegin() error {
	this.Tokens = append(this.Tokens, list())
	return nil
}
func (this *TokenDecoder) OnMapBegin() error {
	this.Tokens = append(this.Tokens, Map())
	return nil
}
func (this *TokenDecoder) OnMarkupBegin() error {
	this.Tokens = append(this.Tokens, markup())
	return nil
}
func (this *TokenDecoder) OnMetadataBegin() error {
	this.Tokens = append(this.Tokens, meta())
	return nil
}
func (this *TokenDecoder) OnCommentBegin() error {
	this.Tokens = append(this.Tokens, cmt())
	return nil
}
func (this *TokenDecoder) OnContainerEnd() error {
	this.Tokens = append(this.Tokens, end())
	return nil
}
func (this *TokenDecoder) OnMarkerBegin() error {
	this.Tokens = append(this.Tokens, marker())
	return nil
}
func (this *TokenDecoder) OnReferenceBegin() error {
	this.Tokens = append(this.Tokens, ref())
	return nil
}
func (this *TokenDecoder) OnBytesBegin() error {
	this.Tokens = append(this.Tokens, bin())
	return nil
}
func (this *TokenDecoder) OnStringBegin() error {
	this.Tokens = append(this.Tokens, str())
	return nil
}
func (this *TokenDecoder) OnURIBegin() error {
	this.Tokens = append(this.Tokens, uri())
	return nil
}
func (this *TokenDecoder) OnArrayChunkBegin(byteCount uint64, isFinalChunk bool) error {
	this.Tokens = append(this.Tokens, chunk(byteCount, isFinalChunk))
	return nil
}
func (this *TokenDecoder) OnArrayData(bytes []byte) error {
	this.Tokens = append(this.Tokens, bytes)
	return nil
}
func (this *TokenDecoder) OnDocumentEnd() error {
	// Nothing to do
	return nil
}

// ===========================================

func decodeDocumentCommon(inlineContainerType InlineContainerType, autoEnd bool, limits *rules.Limits, encoded []byte) (tokens []interface{}, err error) {
	tokenDecoder := NewTokenDecoder()
	decoder := NewDecoder(inlineContainerType, limits, tokenDecoder)
	var isComplete bool
	if isComplete, err = decoder.Feed(encoded); err != nil {
		return
	}
	if !autoEnd && !isComplete {
		err = fmt.Errorf("Document incomplete")
		return
	}

	if err = decoder.EndDocument(); err != nil {
		return
	}

	tokens = tokenDecoder.Tokens
	return
}

func decodeDocumentWithAutoEnd(encoded []byte) (result []interface{}, err error) {
	return decodeDocumentCommon(InlineContainerTypeNone, true, rules.DefaultLimits(), encoded)
}

func decodeDocument(encoded []byte) (result []interface{}, err error) {
	return decodeDocumentCommon(InlineContainerTypeNone, false, rules.DefaultLimits(), encoded)
}

func decodeWithBufferSizeCommon(limits *rules.Limits, encoded []byte, bufferSize int) (tokens []interface{}, err error) {
	tokenDecoder := NewTokenDecoder()
	decoder := NewDecoder(InlineContainerTypeNone, limits, tokenDecoder)
	var isComplete bool
	for offset := 0; offset < len(encoded); offset += bufferSize {
		end := offset + bufferSize
		if end > len(encoded) {
			end = len(encoded)
		}
		if isComplete, err = decoder.Feed(encoded[offset:end]); err != nil {
			return
		}
		if isComplete && end != len(encoded) {
			err = fmt.Errorf("Unexpected end of document")
		}
	}
	if !isComplete {
		err = fmt.Errorf("Document incomplete")
	}
	// TODO: Check that all data was consumed
	tokens = tokenDecoder.Tokens
	return
}

func decodeWithBufferSize(encoded []byte, bufferSize int) (result []interface{}, err error) {
	return decodeWithBufferSizeCommon(rules.DefaultLimits(), encoded, bufferSize)
}

func tryDecode(encoded []byte) error {
	_, err := decodeDocument(encoded)
	return err
}

func encodeToken(encoder *CBEEncoder, value interface{}) error {
	switch v := value.(type) {
	case CBEToken:
		return v.EncodeCBE(encoder)
	case bool:
		return encoder.Bool(v)
	case uint:
		return encoder.PositiveInt(uint64(v))
	case uint8:
		return encoder.PositiveInt(uint64(v))
	case uint16:
		return encoder.PositiveInt(uint64(v))
	case uint32:
		return encoder.PositiveInt(uint64(v))
	case uint64:
		return encoder.PositiveInt(uint64(v))
	case int:
		if v < 0 {
			return encoder.NegativeInt(uint64(-v))
		}
		return encoder.PositiveInt(uint64(v))
	case int8:
		if v < 0 {
			return encoder.NegativeInt(uint64(-v))
		}
		return encoder.PositiveInt(uint64(v))
	case int16:
		if v < 0 {
			return encoder.NegativeInt(uint64(-v))
		}
		return encoder.PositiveInt(uint64(v))
	case int32:
		if v < 0 {
			return encoder.NegativeInt(uint64(-v))
		}
		return encoder.PositiveInt(uint64(v))
	case int64:
		if v < 0 {
			return encoder.NegativeInt(uint64(-v))
		}
		return encoder.PositiveInt(uint64(v))
	case float32:
		return encoder.Float(float64(v))
	case float64:
		return encoder.Float(float64(v))
	case *compact_time.Time:
		return encoder.CompactTime(v)
	case []byte:
		_, err := encoder.ArrayData(v)
		return err
	}
	return fmt.Errorf("%v has unhandled type", value)
}

func encode(src []interface{}) (encoded []byte, err error) {
	buffer := make([]byte, 10000)
	encoder := NewCBEEncoder(InlineContainerTypeNone, buffer, rules.DefaultLimits())
	for _, token := range src {
		if err = encodeToken(encoder, token); err != nil {
			return
		}
	}
	if err = encoder.End(); err != nil {
		return
	}
	encoded = encoder.EncodedBytes()
	return
}

// ============================================================================
// Assertions
// ============================================================================

func assertPanics(t *testing.T, function func()) {
	if getPanicContents(function) == nil {
		t.Errorf("Should have panicked but didn't")
	}
}

func assertDoesNotPanic(t *testing.T, function func()) {
	if result := getPanicContents(function); result != nil {
		t.Errorf("Should not have panicked, but did: %v", result)
	}
}

func assertSuccess(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func assertFailure(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Unexpected success")
	}
}

func assertDecodedEncoded(t *testing.T, expectedEncoded []byte, expectedDecoded []interface{}) {
	actualDecoded, err := decodeDocument(expectedEncoded)
	if err != nil {
		t.Fatalf("Decode Error: %v", err)
		return
	}
	if !equivalence.IsEquivalent(actualDecoded, expectedDecoded) {
		t.Fatalf("Expected decoded %v, actual %v", describe.Describe(expectedDecoded), describe.Describe(actualDecoded))
	}

	actualEncoded, err := encode(actualDecoded)
	if err != nil {
		t.Fatalf("Encode Error: %v", err)
		return
	}
	if !equivalence.IsEquivalent(actualEncoded, expectedEncoded) {
		t.Fatalf("Expected encoded %v, actual %v", describe.Describe(expectedEncoded), describe.Describe(actualEncoded))
	}
}

func assertDecoded(t *testing.T, encoded []byte, expected []interface{}) {
	actual, err := decodeDocument(encoded)
	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}
	if !equivalence.IsEquivalent(actual, expected) {
		t.Fatalf("Expected %v, actual %v", describe.Describe(expected), describe.Describe(actual))
	}
}

func assertDecodedWithAutoEnd(t *testing.T, encoded []byte, expected []interface{}) {
	actual, err := decodeDocumentWithAutoEnd(encoded)
	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}
	if !equivalence.IsEquivalent(actual, expected) {
		t.Fatalf("Expected %v, actual %v", describe.Describe(expected), describe.Describe(actual))
	}
}

func assertDecodedPiecemeal(t *testing.T, encoded []byte, minBufferSize int, maxBufferSize int, expected []interface{}) {
	for i := minBufferSize; i < maxBufferSize; i++ {
		actual, err := decodeWithBufferSize(encoded, i)
		if err != nil {
			t.Fatalf("Error: %v", err)
			return
		}
		if !equivalence.IsEquivalent(actual, expected) {
			t.Fatalf("Expected %v, actual %v", describe.Describe(expected), describe.Describe(actual))
		}
	}
}

// ============================================================================
// Encoder
// ============================================================================

func assertEncoded(t *testing.T, function func(*CBEEncoder) error, expected []byte) {
	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
	err := function(encoder)
	if err != nil {
		t.Fatal(err)
	}
	actual := encoder.EncodedBytes()
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected bytes %v, actual %v", describe.Describe(expected), describe.Describe(actual))
	}
}

// ============================================================================
// Marshal / Unmarshal
// ============================================================================

// func assertMarshaled(t *testing.T, value interface{}, expected []byte) {
// 	encoder := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
// 	Marshal(encoder, InlineContainerTypeNone, value)
// 	actual := encoder.EncodedBytes()
// 	if !bytes.Equal(actual, expected) {
// 		t.Errorf("Expected %v, actual %v", hex.EncodeToString(expected), hex.EncodeToString(actual))
// 	}
// }

// func assertEncodesToExternalBuffer(t *testing.T, value interface{}, bufferSize int) {
// 	buffer := make([]byte, bufferSize)
// 	encoder := NewCBEEncoder(InlineContainerTypeNone, buffer, rules.DefaultLimits())
// 	if err := Marshal(encoder, InlineContainerTypeNone, value); err != nil {
// 		t.Errorf("Unexpected error while marshling: %v", err)
// 		return
// 	}

// 	encoder2 := NewCBEEncoder(InlineContainerTypeNone, nil, rules.DefaultLimits())
// 	Marshal(encoder2, InlineContainerTypeNone, value)
// 	expected := encoder2.EncodedBytes()
// 	if !bytes.Equal(buffer, expected) {
// 		t.Errorf("Expected %v, actual %v", expected, buffer)
// 	}
// }

// func assertFailsEncodingToExternalBuffer(t *testing.T, value interface{}, bufferSize int) {
// 	buffer := make([]byte, bufferSize)
// 	encoder := NewCBEEncoder(InlineContainerTypeNone, buffer, rules.DefaultLimits())
// 	assertPanics(t, func() {
// 		Marshal(encoder, InlineContainerTypeNone, value)
// 	})
// }

func assertMarshaledSize(t *testing.T, value interface{}, expectedSize int) {
	actualSize := CBEEncodedSize(InlineContainerTypeNone, value)
	if actualSize != expectedSize {
		t.Errorf("Expected size %v but got %v", expectedSize, actualSize)
	}
}

func assertMarshalUnmarshal(t *testing.T, expected interface{}, output interface{}) {
	document, err := MarshalCBE(expected)
	if err != nil {
		t.Error(err)
		return
	}
	err = UnmarshalCBE(document, output)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(output, expected) {
		t.Errorf("Expected %v, actual %v", describe.Describe(expected), describe.Describe(output))
	}
}

// func assertMarshalUnmarshalProduces(t *testing.T, input interface{}, expected interface{}) {
// 	document, err := MarshalCBE(input)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	var actual interface{}
// 	actual, err = UnmarshalCBE(document)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	fmt.Printf("### set %v from %v\n", actual, input)

// 	if !DeepEquivalence(actual, expected) {
// 		t.Errorf("Expected %t: <%v>, actual %t: <%v>", expected, expected, actual, actual)
// 	}
// }
