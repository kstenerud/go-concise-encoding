/*
*/
package cbe

import "math"
import "time"
import "github.com/kstenerud/go-smalltime"

const
(
    maxValue6Bit int64 = 0x3f
    maxValue14Bit int64 = 0x3fff
    maxValue30Bit int64 = 0x3fffffff
)

func is6BitLength(value int64) bool {
	return value <= maxValue6Bit
}

func is14BitLength(value int64) bool {
	return value <= maxValue14Bit
}

func is30BitLength(value int64) bool {
	return value <= maxValue30Bit
}

func fitsInSmall(value int64) bool {
	return value >= smallIntMin && value <= smallIntMax
}

func uintFitsInSmall(value uint64) bool {
	return value <= uint64(smallIntMax)
}

func fitsInUInt8(value uint64) bool {
	return value <= math.MaxUint8
}

func fitsInUInt16(value uint64) bool {
	return value <= math.MaxUint16
}

func fitsInUInt32(value uint64) bool {
	return value <= math.MaxUint32
}

type Encoder struct {
	maxContainerDepth int
	// TODO: this is insufficient
	isInArray bool
	encoded []byte
}

func New(maxContainerDepth int) *Encoder {
	encoder := new(Encoder)
	encoder.maxContainerDepth = maxContainerDepth
	encoder.encoded = make([]byte, 0)
	return encoder
}

func (encoder *Encoder) addBytes(bytes []byte) *Encoder {
	encoder.encoded = append(encoder.encoded, bytes...)
	return encoder
}

func (encoder *Encoder) addPrimitive8(value byte) *Encoder {
	encoder.encoded = append(encoder.encoded, value)
	return encoder
}

func (encoder *Encoder) addPrimitive16(value uint16) *Encoder {
	return encoder.addBytes([]byte{byte(value), byte(value>>8)})
}

func (encoder *Encoder) addPrimitive32(value uint32) *Encoder {
	return encoder.addBytes([]byte{
			byte(value), byte(value>>8),
			byte(value>>16), byte(value>>24),
		})
}

func (encoder *Encoder) addPrimitive64(value uint64) *Encoder {
	return encoder.addBytes([]byte{
			byte(value), byte(value>>8), byte(value>>16),
			byte(value>>24), byte(value>>32), byte(value>>40),
			byte(value>>48), byte(value>>56),
		})
}

func (encoder *Encoder) addType(typeValue typeField) *Encoder {
	return encoder.addPrimitive8(byte(typeValue))
}

func (encoder *Encoder) addArrayLength(length int64) *Encoder {
	switch {
	case is6BitLength(length):
		return encoder.addPrimitive8(byte(length << 2 | length6Bit))
	case is14BitLength(length):
		return encoder.addPrimitive16(uint16(length << 2 | length14Bit))
	case is30BitLength(length):
		return encoder.addPrimitive32(uint32(length << 2 | length30Bit))
	default:
		return encoder.addPrimitive64(uint64(length << 2 | length62Bit))
	}
}

func (encoder *Encoder) enterArray() *Encoder {
	// TODO: sanity checks
	encoder.isInArray = true
	return encoder
}

func (encoder *Encoder) leaveArray() *Encoder {
	// TODO: sanity checks
	encoder.isInArray = false
	return encoder
}


func (encoder *Encoder) Padding(byteCount int) *Encoder {
	for i := 0; i < byteCount; i++ {
		encoder.addType(typePadding)
	}
	return encoder
}

func (encoder *Encoder) Nil() *Encoder {
	return encoder.addType(typeNil)
}

func (encoder *Encoder) UInt(value uint64) *Encoder {
	switch {
	case uintFitsInSmall(value):
		return encoder.addPrimitive8(byte(value))
	case fitsInUInt8(value):
		return encoder.addType(typePosInt8).addPrimitive8(uint8(value))
	case fitsInUInt16(value):
		return encoder.addType(typePosInt16).addPrimitive16(uint16(value))
	case fitsInUInt32(value):
		return encoder.addType(typePosInt32).addPrimitive32(uint32(value))
	default:
		return encoder.addType(typePosInt64).addPrimitive64(value)
	}
}

func (encoder *Encoder) Int(value int64) *Encoder {
	if value >= 0 {
		return encoder.UInt(uint64(value))
	}

	uvalue := uint64(-value);

	switch {
	case fitsInSmall(value):
		return encoder.addPrimitive8(byte(value))
	case fitsInUInt8(uvalue):
		return encoder.addType(typeNegInt8).addPrimitive8(uint8(uvalue))
	case fitsInUInt16(uvalue):
		return encoder.addType(typeNegInt16).addPrimitive16(uint16(uvalue))
	case fitsInUInt32(uvalue):
		return encoder.addType(typeNegInt32).addPrimitive32(uint32(uvalue))
	default:
		return encoder.addType(typeNegInt64).addPrimitive64(uvalue)
	}
}

func (encoder *Encoder) Float(value float64) *Encoder {
	asfloat32 := float32(value)
	if float64(asfloat32) == value {
		return encoder.addType(typeFloat32).addPrimitive32(math.Float32bits(asfloat32))
	}
	return encoder.addType(typeFloat64).addPrimitive64(math.Float64bits(value))
}

func (encoder *Encoder) Time(value time.Time) *Encoder {
	return encoder.addType(typeTime).addPrimitive64(uint64(smalltime.FromTime(value)))
}

func (encoder *Encoder) ListBegin() *Encoder {
	return encoder.addType(typeList)
}

func (encoder *Encoder) endContainer() *Encoder {
	return encoder.addType(typeEndContainer)
}

func (encoder *Encoder) ListEnd() *Encoder {
	return encoder.endContainer()
}

func (encoder *Encoder) MapBegin() *Encoder {
	return encoder.addType(typeMap)
}

func (encoder *Encoder) MapEnd() *Encoder {
	return encoder.endContainer()
}

func (encoder *Encoder) BinaryBegin(length int64) *Encoder {
	return encoder.enterArray().addType(typeBinary).addArrayLength(length)
}

func (encoder *Encoder) Binary(value []byte) *Encoder {
	wasInArray := encoder.isInArray
	if !wasInArray {
		encoder.BinaryBegin(int64(len(value)))
	}
	encoder.addBytes(value)
	if !wasInArray {
		return encoder.BinaryEnd()
	}
	return encoder
}

func (encoder *Encoder) BinaryEnd() *Encoder {
	return encoder.leaveArray()
}

func (encoder *Encoder) StringBegin(length int64) *Encoder {
	if(length <= 15) {
		return encoder.addType(typeString0 + typeField(length))
	}
	return encoder.enterArray().addType(typeString).addArrayLength(length)
}

func (encoder *Encoder) String(value string) *Encoder {
	// TODO: Differentiate the array types, sanity checks
	wasInArray := encoder.isInArray
	if !wasInArray {
		encoder.StringBegin(int64(len(value)))
	}
	encoder.addBytes([]byte(value))
	if !wasInArray {
		return encoder.StringEnd()
	}
	return encoder
}

func (encoder *Encoder) StringEnd() *Encoder {
	return encoder.leaveArray()
}

func (encoder *Encoder) Comment(value string) *Encoder {
	return encoder.addType(typeComment).addArrayLength(int64(len(value))).addBytes([]byte(value))
}

func (encoder *Encoder) Encoded() []byte {
	return encoder.encoded
}

