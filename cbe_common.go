package cbe

const cbeCodecVersion = 1

type typeField uint8

const (
	typeDecimal      typeField = 0x65
	typePosInt       typeField = 0x66
	typeNegInt       typeField = 0x67
	typePosInt8      typeField = 0x68
	typeNegInt8      typeField = 0x69
	typePosInt16     typeField = 0x6a
	typeNegInt16     typeField = 0x6b
	typePosInt32     typeField = 0x6c
	typeNegInt32     typeField = 0x6d
	typePosInt64     typeField = 0x6e
	typeNegInt64     typeField = 0x6f
	typeFloat32      typeField = 0x70
	typeFloat64      typeField = 0x71
	typeReserved72   typeField = 0x72
	typeReserved73   typeField = 0x73
	typeReserved74   typeField = 0x74
	typeReserved75   typeField = 0x75
	typeComment      typeField = 0x76
	typeMetadata     typeField = 0x77
	typeMarkup       typeField = 0x78
	typeMap          typeField = 0x79
	typeList         typeField = 0x7a
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
	typeReserved93   typeField = 0x93
	typeReserved94   typeField = 0x94
	typeReserved95   typeField = 0x95
	typeReserved96   typeField = 0x96
	typeMarker       typeField = 0x97
	typeReference    typeField = 0x98
	typeDate         typeField = 0x99
	typeTime         typeField = 0x9a
	typeTimestamp    typeField = 0x9b
)

const (
	smallIntMin int64 = -100
	smallIntMax int64 = 100
)
