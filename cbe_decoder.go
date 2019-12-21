package cbe

import (
	"errors"
	"fmt"

	"github.com/kstenerud/go-cbe/rules"
	"github.com/kstenerud/go-compact-time"
)

// TODO: escape sequences

type CBEDecodeError struct {
	documentOffset int
	bufferOffset   int
	Err            error
}

func NewCBEDecodeError(documentOffset int, bufferOffset int, err error) *CBEDecodeError {
	newErr := new(CBEDecodeError)
	newErr.documentOffset = documentOffset
	newErr.bufferOffset = bufferOffset
	newErr.Err = err
	return newErr
}

func (this *CBEDecodeError) Error() string {
	return fmt.Sprintf("Offset %v (buffer offset %v): %v", this.documentOffset, this.bufferOffset, this.Err)
}

func (this *CBEDecodeError) Unwrap() error {
	return this.Err
}

// Callback functions that must be present in the receiver object.
type CBEDecoderCallbacks interface {
	OnNil() error
	OnBool(value bool) error
	OnPositiveInt(value uint64) error
	OnNegativeInt(value uint64) error
	OnFloat(value float64) error
	OnTime(time *compact_time.Time) error
	// Call order: list begin, item*, container end
	OnListBegin() error
	// Call order: map begin, (key, value)*, container end
	OnMapBegin() error
	// Call order: markup begin, name, (key, value)*, container end, item*, container end
	OnMarkupBegin() error
	// Call order: metadata begin, (key, value)*, container end
	OnMetadataBegin() error
	// Call order: comment begin, item*, container end
	OnCommentBegin() error
	OnContainerEnd() error
	// Call order: marker begin, id, item
	OnMarkerBegin() error
	// Call order: reference begin, id
	OnReferenceBegin() error
	// Call order: bytes begin, (array chunk begin, array data*)*, array end
	OnBytesBegin() error
	// Call order: string begin, (array chunk begin, array data*)*, array end
	OnStringBegin() error
	// Call order: URI begin, (array chunk begin, array data*)*, array end
	OnURIBegin() error
	OnArrayChunkBegin(byteCount uint64, isFinalChunk bool) error
	OnArrayData(bytes []byte) error
	OnDocumentEnd() error
}

// Biggest item is timestamp (10 bytes), longest tz is "America/Argentina/ComodRivadavia"
const cbeMaxPartialReadLength = 50

type decoderError struct {
	err error
}

type callbackError struct {
	err error
}

type CBEDecoder struct {
	rules                   rules.Rules
	streamOffset            uint64
	buffer                  *cbeDecodeBuffer
	mainBuffer              *cbeDecodeBuffer
	underflowBuffer         *cbeDecodeBuffer
	callbacks               CBEDecoderCallbacks
	inlineContainerType     InlineContainerType
	hasBegunInlineContainer bool
}

func (this *CBEDecoder) contextualizedError(err error) error {
	// TODO: also check if it's already contextualized
	// return NewDecodeError(documentOffset, bufferOffset, err)
	return err
}

func (this *CBEDecoder) initInlineContainer() (err error) {
	if this.inlineContainerType != InlineContainerTypeNone && !this.hasBegunInlineContainer {
		switch this.inlineContainerType {
		case InlineContainerTypeList:
			if err = this.rules.AddVersion(cbeCodecVersion); err != nil {
				return
			}
			if err = this.rules.BeginList(); err != nil {
				return
			}
			if err = this.callbacks.OnListBegin(); err != nil {
				return
			}
		case InlineContainerTypeMap:
			if err = this.rules.AddVersion(cbeCodecVersion); err != nil {
				return
			}
			if err = this.rules.BeginMap(); err != nil {
				return
			}
			if err = this.callbacks.OnMapBegin(); err != nil {
				return
			}
		}
		this.hasBegunInlineContainer = true
	}
	return
}

func (this *CBEDecoder) handleVersion() (err error) {
	var version uint64
	if version, err = this.buffer.DecodeUint(); err != nil {
		return
	}
	if err = this.rules.AddVersion(version); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleBool(value bool) (err error) {
	if err = this.rules.AddBool(); err != nil {
		return
	}
	if err = this.callbacks.OnBool(value); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleNil() (err error) {
	if err = this.rules.AddNil(); err != nil {
		return
	}
	if err = this.callbacks.OnNil(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleFloat(value float64) (err error) {
	if err = this.rules.AddFloat(value); err != nil {
		return
	}
	if err = this.callbacks.OnFloat(value); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) decodeFloat32() (err error) {
	var value float32
	if value, err = this.buffer.DecodeFloat32(); err != nil {
		return
	}
	return this.handleFloat(float64(value))
}

func (this *CBEDecoder) decodeFloat64() (err error) {
	var value float64
	if value, err = this.buffer.DecodeFloat64(); err != nil {
		return
	}
	return this.handleFloat(value)
}

func (this *CBEDecoder) decodeDecimal() (err error) {
	var value float64
	if value, err = this.buffer.DecodeFloat(); err != nil {
		return
	}
	return this.handleFloat(value)
}

func (this *CBEDecoder) handlePositiveInt(value uint64) (err error) {
	if err = this.rules.AddPositiveInt(value); err != nil {
		return
	}
	if err = this.callbacks.OnPositiveInt(value); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) decodePositiveInt8() (err error) {
	var value uint8
	if value, err = this.buffer.DecodeUint8(); err != nil {
		return
	}
	return this.handlePositiveInt(uint64(value))
}

func (this *CBEDecoder) decodePositiveInt16() (err error) {
	var value uint16
	if value, err = this.buffer.DecodeUint16(); err != nil {
		return
	}
	return this.handlePositiveInt(uint64(value))
}

func (this *CBEDecoder) decodePositiveInt32() (err error) {
	var value uint32
	if value, err = this.buffer.DecodeUint32(); err != nil {
		return
	}
	return this.handlePositiveInt(uint64(value))
}

func (this *CBEDecoder) decodePositiveInt64() (err error) {
	var value uint64
	if value, err = this.buffer.DecodeUint64(); err != nil {
		return
	}
	return this.handlePositiveInt(value)
}

func (this *CBEDecoder) decodePositiveInt() (err error) {
	var value uint64
	if value, err = this.buffer.DecodeUint(); err != nil {
		return
	}
	return this.handlePositiveInt(value)
}

func (this *CBEDecoder) handleNegativeInt(value uint64) (err error) {
	if err = this.rules.AddNegativeInt(); err != nil {
		return
	}
	if err = this.callbacks.OnNegativeInt(value); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) decodeNegativeInt8() (err error) {
	var value uint8
	if value, err = this.buffer.DecodeUint8(); err != nil {
		return
	}
	return this.handleNegativeInt(uint64(value))
}

func (this *CBEDecoder) decodeNegativeInt16() (err error) {
	var value uint16
	if value, err = this.buffer.DecodeUint16(); err != nil {
		return
	}
	return this.handleNegativeInt(uint64(value))
}

func (this *CBEDecoder) decodeNegativeInt32() (err error) {
	var value uint32
	if value, err = this.buffer.DecodeUint32(); err != nil {
		return
	}
	return this.handleNegativeInt(uint64(value))
}

func (this *CBEDecoder) decodeNegativeInt64() (err error) {
	var value uint64
	if value, err = this.buffer.DecodeUint64(); err != nil {
		return
	}
	return this.handleNegativeInt(value)
}

func (this *CBEDecoder) decodeNegativeInt() (err error) {
	var value uint64
	if value, err = this.buffer.DecodeUint(); err != nil {
		return
	}
	return this.handleNegativeInt(value)
}

func (this *CBEDecoder) handleTime(value *compact_time.Time) (err error) {
	if err = this.rules.AddTime(); err != nil {
		return
	}
	if err = this.callbacks.OnTime(value); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) decodeDate() (err error) {
	var value *compact_time.Time
	if value, err = this.buffer.DecodeDate(); err != nil {
		return
	}
	return this.handleTime(value)
}

func (this *CBEDecoder) decodeTime() (err error) {
	var value *compact_time.Time
	if value, err = this.buffer.DecodeTime(); err != nil {
		return
	}
	return this.handleTime(value)
}

func (this *CBEDecoder) decodeTimestamp() (err error) {
	var value *compact_time.Time
	if value, err = this.buffer.DecodeTimestamp(); err != nil {
		return
	}
	return this.handleTime(value)
}

func (this *CBEDecoder) handleList() (err error) {
	if err = this.rules.BeginList(); err != nil {
		return
	}
	if err = this.callbacks.OnListBegin(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleMap() (err error) {
	if err = this.rules.BeginMap(); err != nil {
		return
	}
	if err = this.callbacks.OnMapBegin(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleMarkup() (err error) {
	if err = this.rules.BeginMarkup(); err != nil {
		return
	}
	if err = this.callbacks.OnMarkupBegin(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleMetadata() (err error) {
	if err = this.rules.BeginMetadata(); err != nil {
		return
	}
	if err = this.callbacks.OnMetadataBegin(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleComment() (err error) {
	if err = this.rules.BeginComment(); err != nil {
		return
	}
	if err = this.callbacks.OnCommentBegin(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleEndContainer() (err error) {
	if err = this.rules.EndContainer(); err != nil {
		return
	}
	if err = this.callbacks.OnContainerEnd(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleBytes() (err error) {
	if err = this.rules.BeginBytes(); err != nil {
		return
	}
	if err = this.callbacks.OnBytesBegin(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleString() (err error) {
	if err = this.rules.BeginString(); err != nil {
		return
	}
	if err = this.callbacks.OnStringBegin(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleURI() (err error) {
	if err = this.rules.BeginURI(); err != nil {
		return
	}
	if err = this.callbacks.OnURIBegin(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleFixedLengthString(length int) (err error) {
	if err = this.rules.BeginString(); err != nil {
		return
	}
	if err = this.rules.BeginArrayChunk(uint64(length), true); err != nil {
		return
	}
	if err = this.callbacks.OnStringBegin(); err != nil {
		return
	}
	if err = this.callbacks.OnArrayChunkBegin(uint64(length), true); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleMarker() (err error) {
	if err = this.rules.BeginMarker(); err != nil {
		return
	}
	if err = this.callbacks.OnMarkerBegin(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handleReference() (err error) {
	if err = this.rules.BeginReference(); err != nil {
		return
	}
	if err = this.callbacks.OnReferenceBegin(); err != nil {
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) handlePadding() (err error) {
	// Nothing to do
	this.buffer.Commit()
	return
}

func (this *CBEDecoder) decodeObject() (err error) {
	var dataType typeField
	if dataType, err = this.buffer.DecodeType(); err != nil {
		return
	}

	asSmallInt := int8(dataType)
	if int64(asSmallInt) >= smallIntMin && int64(asSmallInt) <= smallIntMax {
		if asSmallInt < 0 {
			err = this.handleNegativeInt(uint64(-asSmallInt))
		} else {
			err = this.handlePositiveInt(uint64(asSmallInt))
		}
		this.buffer.Commit()
		return
	}

	switch dataType {
	case typeTrue:
		return this.handleBool(true)
	case typeFalse:
		return this.handleBool(false)
	case typeFloat32:
		return this.decodeFloat32()
	case typeFloat64:
		return this.decodeFloat64()
	case typeDecimal:
		return this.decodeDecimal()
	case typePosInt8:
		return this.decodePositiveInt8()
	case typePosInt16:
		return this.decodePositiveInt16()
	case typePosInt32:
		return this.decodePositiveInt32()
	case typePosInt64:
		return this.decodePositiveInt64()
	case typePosInt:
		return this.decodePositiveInt()
	case typeNegInt8:
		return this.decodeNegativeInt8()
	case typeNegInt16:
		return this.decodeNegativeInt16()
	case typeNegInt32:
		return this.decodeNegativeInt32()
	case typeNegInt64:
		return this.decodeNegativeInt64()
	case typeNegInt:
		return this.decodeNegativeInt()
	case typeDate:
		return this.decodeDate()
	case typeTime:
		return this.decodeTime()
	case typeTimestamp:
		return this.decodeTimestamp()
	case typeNil:
		return this.handleNil()
	case typePadding:
		return this.handlePadding()
	case typeList:
		return this.handleList()
	case typeMap:
		return this.handleMap()
	case typeMarkup:
		return this.handleMarkup()
	case typeMetadata:
		return this.handleMetadata()
	case typeComment:
		return this.handleComment()
	case typeEndContainer:
		return this.handleEndContainer()
	case typeBytes:
		return this.handleBytes()
	case typeURI:
		return this.handleURI()
	case typeString:
		return this.handleString()
	case typeString0:
		return this.handleFixedLengthString(0)
	case typeString1:
		return this.handleFixedLengthString(1)
	case typeString2:
		return this.handleFixedLengthString(2)
	case typeString3:
		return this.handleFixedLengthString(3)
	case typeString4:
		return this.handleFixedLengthString(4)
	case typeString5:
		return this.handleFixedLengthString(5)
	case typeString6:
		return this.handleFixedLengthString(6)
	case typeString7:
		return this.handleFixedLengthString(7)
	case typeString8:
		return this.handleFixedLengthString(8)
	case typeString9:
		return this.handleFixedLengthString(9)
	case typeString10:
		return this.handleFixedLengthString(10)
	case typeString11:
		return this.handleFixedLengthString(11)
	case typeString12:
		return this.handleFixedLengthString(12)
	case typeString13:
		return this.handleFixedLengthString(13)
	case typeString14:
		return this.handleFixedLengthString(14)
	case typeString15:
		return this.handleFixedLengthString(15)
	case typeMarker:
		return this.handleMarker()
	case typeReference:
		return this.handleReference()
	default:
		return fmt.Errorf("%02x: Unknown type code", dataType)
	}
}

func (this *CBEDecoder) handleChunkHeader() (err error) {
	var header uint64
	if header, err = this.buffer.DecodeUint(); err != nil {
		return
	}
	length := header >> 1
	isFinalChunk := header&1 == 0
	if err = this.rules.BeginArrayChunk(length, isFinalChunk); err != nil {
		return
	}
	if err = this.callbacks.OnArrayChunkBegin(length, isFinalChunk); err != nil {
		return
	}

	return
}

func (this *CBEDecoder) handleChunkData() (err error) {
	if this.rules.GetRemainingChunkByteCount() > 0 {
		decodeByteCount := this.buffer.GetUncommittedByteCount()
		if uint64(decodeByteCount) > this.rules.GetRemainingChunkByteCount() {
			decodeByteCount = int(this.rules.GetRemainingChunkByteCount())
		}
		bytes, err := this.buffer.DecodeBytes(decodeByteCount)
		if err != nil {
			return err
		}
		if err := this.rules.AddArrayData(bytes); err != nil {
			return err
		}
		if err := this.callbacks.OnArrayData(bytes); err != nil {
			return err
		}
		this.buffer.Commit()

		if this.rules.GetRemainingChunkByteCount() > 0 {
			return BufferExhaustedError
		}
	}
	return
}

func (this *CBEDecoder) handleDocumentComplete() (err error) {
	if this.hasBegunInlineContainer {
		if err = this.rules.EndContainer(); err != nil {
			return
		}
		if err = this.callbacks.OnContainerEnd(); err != nil {
			return
		}
	}

	if err = this.rules.EndDocument(); err != nil {
		return
	}
	if err = this.callbacks.OnDocumentEnd(); err != nil {
		return
	}

	return
}

func (this *CBEDecoder) feedOnceFromCurrentBuffer() (isComplete bool, err error) {
	switch this.rules.GetDecoderState() {
	case rules.DecoderStateAwaitingVersion:
		err = this.handleVersion()
	case rules.DecoderStateAwaitingArrayChunk:
		err = this.handleChunkHeader()
	case rules.DecoderStateAwaitingArrayData:
		err = this.handleChunkData()
	case rules.DecoderStateAwaitingEndDocument:
		isComplete = this.rules.IsDocumentComplete()
		if isComplete {
			err = this.handleDocumentComplete()
		}
	default:
		err = this.decodeObject()
	}
	return
}

func (this *CBEDecoder) feedFromUnderflow() (isComplete bool, err error) {
	this.buffer = this.underflowBuffer
	return this.feedOnceFromCurrentBuffer()
}

func (this *CBEDecoder) feedFromMain() (isComplete bool, err error) {
	this.buffer = this.mainBuffer
	for {
		isComplete, err = this.feedOnceFromCurrentBuffer()
		if err != nil || isComplete {
			return
		}
	}
}

func (this *CBEDecoder) fillUnderflowFromMain() int {
	return this.underflowBuffer.FillToByteCount(this.mainBuffer, cbeMaxPartialReadLength)

}

// ----------
// Public API
// ----------

func NewDecoder(inlineContainerType InlineContainerType, limits *rules.Limits, callbacks CBEDecoderCallbacks) *CBEDecoder {
	this := new(CBEDecoder)
	this.Init(inlineContainerType, limits, callbacks)
	return this
}

func (this *CBEDecoder) Init(inlineContainerType InlineContainerType, limits *rules.Limits, callbacks CBEDecoderCallbacks) {
	this.rules.Init(cbeCodecVersion, limits)
	this.underflowBuffer = NewDecodeBuffer(make([]byte, cbeMaxPartialReadLength))
	this.mainBuffer = NewDecodeBuffer(make([]byte, 0))
	this.callbacks = callbacks
	this.inlineContainerType = inlineContainerType
	this.Reset()
}

func (this *CBEDecoder) Reset() {
	this.rules.Reset()
	this.underflowBuffer.Clear()
	this.mainBuffer.Clear()

	this.streamOffset = 0 // TODO: Remove?
	this.buffer = this.mainBuffer
	this.hasBegunInlineContainer = false
}

// Feed bytes into the decoder to be decoded.
func (this *CBEDecoder) Feed(bytesToDecode []byte) (isComplete bool, err error) {
	if err = this.initInlineContainer(); err != nil {
		err = this.contextualizedError(err)
		return
	}

	this.mainBuffer.ReplaceBuffer(bytesToDecode)

	if this.buffer == this.underflowBuffer {
		this.underflowBuffer.Rollback()
		byteCountFilledFromMain := this.fillUnderflowFromMain()
		if isComplete, err = this.feedFromUnderflow(); err != nil {
			if !errors.Is(err, BufferExhaustedError) {
				err = this.contextualizedError(err)
			}
			err = nil
			return
		}

		unusedByteCount := this.underflowBuffer.GetUncommittedByteCount()
		newMainPosition := byteCountFilledFromMain - unusedByteCount
		this.mainBuffer.lastCommitPosition = newMainPosition
		this.mainBuffer.position = newMainPosition
		this.underflowBuffer.Clear()
		this.buffer = this.mainBuffer

		if isComplete {
			return
		}
	}

	if isComplete, err = this.feedFromMain(); err != nil {
		if !errors.Is(err, BufferExhaustedError) {
			err = this.contextualizedError(err)
			return
		}
		err = nil
		this.underflowBuffer.AddContents(this.mainBuffer.GetUncommittedBytes())
		this.buffer = this.underflowBuffer
		this.mainBuffer.Clear()
	}

	return
}

// Ensure the document is ended. Call this function if you've created a decoder
// with an inline container type. The method is idempotent, so it won't do any
// harm to call it multiple times or to call it when you don't have an inline
// container type.
func (this *CBEDecoder) EndDocument() error {
	if this.rules.GetDecoderState() == rules.DecoderStateEnded {
		return nil
	}
	return this.handleDocumentComplete()
}

// Convenience function to decode an entire document in a single call.
func (this *CBEDecoder) Decode(document []byte) (err error) {
	var isComplete bool
	if isComplete, err = this.Feed(document); err != nil {
		// Don't wrap errors from Feed()
		return
	}
	if !isComplete {
		return this.contextualizedError(fmt.Errorf("Incomplete document"))
	}
	return
}
