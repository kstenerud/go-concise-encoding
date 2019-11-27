package cbe

import "fmt"

type InlineContainerType int

const (
	InlineContainerTypeNone InlineContainerType = iota
	InlineContainerTypeList
	InlineContainerTypeMap
)

type typeField uint8

const (
	typeDecimal  typeField = 0x65
	typePosInt   typeField = 0x66
	typeNegInt   typeField = 0x67
	typePosInt8  typeField = 0x68
	typeNegInt8  typeField = 0x69
	typePosInt16 typeField = 0x6a
	typeNegInt16 typeField = 0x6b
	typePosInt32 typeField = 0x6c
	typeNegInt32 typeField = 0x6d
	typePosInt64 typeField = 0x6e
	typeNegInt64 typeField = 0x6f
	typeFloat32  typeField = 0x70
	typeFloat64  typeField = 0x71
	// RESERVED 0x72 - 0x74
	typeList         typeField = 0x75
	typeMap          typeField = 0x76
	typeMetadata     typeField = 0x77
	typeMarkup       typeField = 0x78
	typeEndContainer typeField = 0x79
	typeMarker       typeField = 0x7a
	typeReference    typeField = 0x7b
	typeFalse        typeField = 0x7c
	typeTrue         typeField = 0x7d
	typeNil          typeField = 0x7e
	typePadding      typeField = 0x7f
	typeString0      typeField = 0x80
	typeString1      typeField = 0x81
	typeString2      typeField = 0x82
	typeString3      typeField = 0x83
	typeString4      typeField = 0x84
	typeString5      typeField = 0x85
	typeString6      typeField = 0x86
	typeString7      typeField = 0x87
	typeString8      typeField = 0x88
	typeString9      typeField = 0x89
	typeString10     typeField = 0x8a
	typeString11     typeField = 0x8b
	typeString12     typeField = 0x8c
	typeString13     typeField = 0x8d
	typeString14     typeField = 0x8e
	typeString15     typeField = 0x8f
	typeString       typeField = 0x90
	typeBytes        typeField = 0x91
	typeURI          typeField = 0x92
	typeComment      typeField = 0x93
	// RESERVED 0x94 - 0x98
	typeDate      typeField = 0x99
	typeTime      typeField = 0x9a
	typeTimestamp typeField = 0x9b
)

const (
	smallIntMin int64 = -100
	smallIntMax int64 = 100
)

type arrayType int

const (
	arrayTypeNone arrayType = iota
	arrayTypeBytes
	arrayTypeString
	arrayTypeURI
	arrayTypeComment
)

type containerPhase int

const (
	containerPhaseNone containerPhase = iota
	containerPhaseList
	containerPhaseMap
	containerPhaseMetadata
	containerPhaseMarker
	containerPhaseMarkupAttributes
	containerPhaseMarkupContents
)

func ValidateCommentCharacter(ch int) error {
	if (ch >= 0x00 && ch <= 0x1f && ch != 0x09) ||
		(ch >= 0x7f && ch <= 0x9f) ||
		ch == 0x2028 || ch == 0x2029 {
		return fmt.Errorf("0x%04x: Invalid comment character", ch)
	}
	return nil
}

var digitsMax = [...]uint64{
	0,
	9,
	99,
	999,
	9999,
	99999,
	999999,
	9999999,
	99999999,
	999999999,
	9999999999,
	99999999999,
	999999999999,
	9999999999999,
	99999999999999,
	999999999999999,
	9999999999999999,
	99999999999999999,
	999999999999999999,
	9999999999999999999, // 19 digits
	// Max digits for uint64 is 20
}

func CountDigits(value uint64) int {
	// This is MUCH faster than the string method, and 4x faster than int(math.Log10(float64(value))) + 1
	// Subdividing any further yields no overall gains.
	if value <= digitsMax[10] {
		for i := 1; i < 10; i++ {
			if value <= digitsMax[i] {
				return i
			}
		}
		return 10
	}

	for i := 11; i < 20; i++ {
		if value <= digitsMax[i] {
			return i
		}
	}
	return 20
}
