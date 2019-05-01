package cbe

type typeField uint8

const (
	typeTrue         typeField = 0x6b
	typeFalse        typeField = 0x6c
	typeFloat32      typeField = 0x6d
	typeFloat64      typeField = 0x6e
	typeFloat128     typeField = 0x6f
	typePosInt8      typeField = 0x70
	typePosInt16     typeField = 0x71
	typePosInt32     typeField = 0x72
	typePosInt64     typeField = 0x73
	typePosInt128    typeField = 0x74
	typeDecimal32    typeField = 0x75
	typeDecimal64    typeField = 0x76
	typeDecimal128   typeField = 0x77
	typeNegInt8      typeField = 0x78
	typeNegInt16     typeField = 0x79
	typeNegInt32     typeField = 0x7a
	typeNegInt64     typeField = 0x7b
	typeNegInt128    typeField = 0x7c
	typeTime         typeField = 0x7d
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
	typeBinary       typeField = 0x91
	typeComment      typeField = 0x92
	typeList         typeField = 0x93
	typeMap          typeField = 0x94
	typeEndContainer typeField = 0x95
)

const (
	smallIntMin int64 = -106
	smallIntMax int64 = 106
)

const (
	length6Bit  int64 = 0
	length14Bit int64 = 1
	length30Bit int64 = 2
	length62Bit int64 = 3
)

type arrayType int

const (
	arrayTypeNone arrayType = iota
	arrayTypeBinary
	arrayTypeString
	arrayTypeComment
)

type containerType int

const (
	containerTypeNone containerType = iota
	containerTypeList
	containerTypeMap
)
