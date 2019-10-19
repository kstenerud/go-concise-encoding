package cbe

import "fmt"

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
	// RESERVED 0x72 - 0x76
	typeList         typeField = 0x77
	typeMapUnordered typeField = 0x78
	typeMapOrdered   typeField = 0x79
	typeMapMetadata  typeField = 0x7a
	typeEndContainer typeField = 0x7b
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

type ContainerType int

const (
	ContainerTypeNone ContainerType = iota
	ContainerTypeList
	ContainerTypeUnorderedMap
	ContainerTypeOrderedMap
	ContainerTypeMetadataMap
)

func ValidateCommentCharacter(ch int) error {
	if (ch >= 0x00 && ch <= 0x1f && ch != 0x09) ||
		(ch >= 0x7f && ch <= 0x9f) ||
		ch == 0x2028 || ch == 0x2029 {
		return fmt.Errorf("0x%04x: Invalid comment character", ch)
	}
	return nil
}
