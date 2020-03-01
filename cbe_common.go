package concise_encoding

import (
	"fmt"
	"net/url"
	"reflect"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	timeType  = reflect.TypeOf(time.Time{})
	urlType   = reflect.TypeOf(url.URL{})
	pURLType  = reflect.TypeOf((*url.URL)(nil))
	bytesType = reflect.TypeOf([]uint8{})
)

func isFieldExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

const cbeCodecVersion = 1

type cbeTypeField uint8

const (
	cbeTypeDecimal      cbeTypeField = 0x65
	cbeTypePosInt       cbeTypeField = 0x66
	cbeTypeNegInt       cbeTypeField = 0x67
	cbeTypePosInt8      cbeTypeField = 0x68
	cbeTypeNegInt8      cbeTypeField = 0x69
	cbeTypePosInt16     cbeTypeField = 0x6a
	cbeTypeNegInt16     cbeTypeField = 0x6b
	cbeTypePosInt32     cbeTypeField = 0x6c
	cbeTypeNegInt32     cbeTypeField = 0x6d
	cbeTypePosInt64     cbeTypeField = 0x6e
	cbeTypeNegInt64     cbeTypeField = 0x6f
	cbeTypeFloat32      cbeTypeField = 0x70
	cbeTypeFloat64      cbeTypeField = 0x71
	cbeTypeUUID         cbeTypeField = 0x72
	cbeTypeReserved73   cbeTypeField = 0x73
	cbeTypeReserved74   cbeTypeField = 0x74
	cbeTypeReserved75   cbeTypeField = 0x75
	cbeTypeComment      cbeTypeField = 0x76
	cbeTypeMetadata     cbeTypeField = 0x77
	cbeTypeMarkup       cbeTypeField = 0x78
	cbeTypeMap          cbeTypeField = 0x79
	cbeTypeList         cbeTypeField = 0x7a
	cbeTypeEndContainer cbeTypeField = 0x7b
	cbeTypeFalse        cbeTypeField = 0x7c
	cbeTypeTrue         cbeTypeField = 0x7d
	cbeTypeNil          cbeTypeField = 0x7e
	cbeTypePadding      cbeTypeField = 0x7f
	cbeTypeString0      cbeTypeField = 0x80
	cbeTypeString1      cbeTypeField = 0x81
	cbeTypeString2      cbeTypeField = 0x82
	cbeTypeString3      cbeTypeField = 0x83
	cbeTypeString4      cbeTypeField = 0x84
	cbeTypeString5      cbeTypeField = 0x85
	cbeTypeString6      cbeTypeField = 0x86
	cbeTypeString7      cbeTypeField = 0x87
	cbeTypeString8      cbeTypeField = 0x88
	cbeTypeString9      cbeTypeField = 0x89
	cbeTypeString10     cbeTypeField = 0x8a
	cbeTypeString11     cbeTypeField = 0x8b
	cbeTypeString12     cbeTypeField = 0x8c
	cbeTypeString13     cbeTypeField = 0x8d
	cbeTypeString14     cbeTypeField = 0x8e
	cbeTypeString15     cbeTypeField = 0x8f
	cbeTypeString       cbeTypeField = 0x90
	cbeTypeBytes        cbeTypeField = 0x91
	cbeTypeURI          cbeTypeField = 0x92
	cbeTypeCustom       cbeTypeField = 0x93
	cbeTypeReserved94   cbeTypeField = 0x94
	cbeTypeReserved95   cbeTypeField = 0x95
	cbeTypeReserved96   cbeTypeField = 0x96
	cbeTypeMarker       cbeTypeField = 0x97
	cbeTypeReference    cbeTypeField = 0x98
	cbeTypeDate         cbeTypeField = 0x99
	cbeTypeTime         cbeTypeField = 0x9a
	cbeTypeTimestamp    cbeTypeField = 0x9b
)

const (
	cbeSmallIntMin int64 = -100
	cbeSmallIntMax int64 = 100
)

func MarshalCBE(object interface{}, useReferences bool) (document []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	encoder := NewCBEEncoder()
	iterator := NewRootObjectIterator(useReferences, encoder)
	iterator.Iterate(object)
	document = encoder.Document()
	return
}

func UnmarshalCBE(document []byte, template interface{}, shouldZeroCopy bool) (decoded interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	builder := NewBuilderFor(template)
	rules := NewRules(cbeCodecVersion, DefaultLimits(), builder)
	decoder := NewCBEDecoder(document, rules, shouldZeroCopy)
	decoder.Decode()
	decoded = builder.GetBuiltObject()
	return
}
