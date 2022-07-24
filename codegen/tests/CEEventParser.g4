parser grammar CEEventParser;
options { tokenVocab=CEEventLexer; }

start: event EOF;

event
    : eventArrayBits
    | eventArrayChunkLast
    | eventArrayChunkMore
    | eventArrayDataBits
    | eventArrayDataFloat16
    | eventArrayDataFloat32
    | eventArrayDataFloat64
    | eventArrayDataInt16
    | eventArrayDataInt32
    | eventArrayDataInt64
    | eventArrayDataInt8
    | eventArrayDataText
    | eventArrayDataUID
    | eventArrayDataUint16
    | eventArrayDataUint16X
    | eventArrayDataUint32
    | eventArrayDataUint32X
    | eventArrayDataUint64
    | eventArrayDataUint64X
    | eventArrayDataUint8
    | eventArrayDataUint8X
    | eventArrayFloat16
    | eventArrayFloat32
    | eventArrayFloat64
    | eventArrayInt16
    | eventArrayInt32
    | eventArrayInt64
    | eventArrayInt8
    | eventArrayUID
    | eventArrayUint16
    | eventArrayUint16X
    | eventArrayUint32
    | eventArrayUint32X
    | eventArrayUint64
    | eventArrayUint64X
    | eventArrayUint8
    | eventArrayUint8X
    | eventBeginArrayBits
    | eventBeginArrayFloat16
    | eventBeginArrayFloat32
    | eventBeginArrayFloat64
    | eventBeginArrayInt16
    | eventBeginArrayInt32
    | eventBeginArrayInt64
    | eventBeginArrayInt8
    | eventBeginArrayUID
    | eventBeginArrayUint16
    | eventBeginArrayUint32
    | eventBeginArrayUint64
    | eventBeginArrayUint8
    | eventBeginCustomBinary
    | eventBeginCustomText
    | eventBeginMedia
    | eventBeginResourceId
    | eventBeginRemoteReference
    | eventBeginString
    | eventBoolean
    | eventCommentSingleLine
    | eventCommentMultiline
    | eventCustomBinary
    | eventCustomText
    | eventEdge
    | eventEndContainer
    | eventList
    | eventMap
    | eventMarker
    | eventMedia
    | eventNode
    | eventNull
    | eventNumber
    | eventPad
    | eventLocalReference
    | eventRemoteReference
    | eventResourceId
    | eventString
    | eventStructInstance
    | eventStructTemplate
    | eventTime
    | eventUID
    | eventVersion
    ;

eventArrayBits:            EVENT_AB     VALUE_BIT*;
eventArrayChunkLast:       EVENT_ACL    VALUE_UINT_DEC*;
eventArrayChunkMore:       EVENT_ACM    VALUE_UINT_DEC*;
eventArrayDataBits:        EVENT_ADB    VALUE_BIT*;
eventArrayDataFloat16:     EVENT_ADF16  (VALUE_FLOAT_DEC | VALUE_FLOAT_HEX | VALUE_FLOAT_INF | VALUE_FLOAT_NAN | VALUE_FLOAT_SNAN)*;
eventArrayDataFloat32:     EVENT_ADF32  (VALUE_FLOAT_DEC | VALUE_FLOAT_HEX | VALUE_FLOAT_INF | VALUE_FLOAT_NAN | VALUE_FLOAT_SNAN)*;
eventArrayDataFloat64:     EVENT_ADF64  (VALUE_FLOAT_DEC | VALUE_FLOAT_HEX | VALUE_FLOAT_INF | VALUE_FLOAT_NAN | VALUE_FLOAT_SNAN)*;
eventArrayDataInt16:       EVENT_ADI16  (VALUE_INT_BIN | VALUE_INT_OCT | VALUE_INT_DEC | VALUE_INT_HEX)*;
eventArrayDataInt32:       EVENT_ADI32  (VALUE_INT_BIN | VALUE_INT_OCT | VALUE_INT_DEC | VALUE_INT_HEX)*;
eventArrayDataInt64:       EVENT_ADI64  (VALUE_INT_BIN | VALUE_INT_OCT | VALUE_INT_DEC | VALUE_INT_HEX)*;
eventArrayDataInt8:        EVENT_ADI8   (VALUE_INT_BIN | VALUE_INT_OCT | VALUE_INT_DEC | VALUE_INT_HEX)*;
eventArrayDataText:        EVENT_ADT    STRING?;
eventArrayDataUID:         EVENT_ADU    UID*;
eventArrayDataUint16:      EVENT_ADU16  (VALUE_UINT_BIN | VALUE_UINT_OCT | VALUE_UINT_DEC | VALUE_UINT_HEX)*;
eventArrayDataUint16X:     EVENT_ADU16X VALUE_UINTX*;
eventArrayDataUint32:      EVENT_ADU32  (VALUE_UINT_BIN | VALUE_UINT_OCT | VALUE_UINT_DEC | VALUE_UINT_HEX)*;
eventArrayDataUint32X:     EVENT_ADU32X VALUE_UINTX*;
eventArrayDataUint64:      EVENT_ADU64  (VALUE_UINT_BIN | VALUE_UINT_OCT | VALUE_UINT_DEC | VALUE_UINT_HEX)*;
eventArrayDataUint64X:     EVENT_ADU64X VALUE_UINTX*;
eventArrayDataUint8:       EVENT_ADU8   (VALUE_UINT_BIN | VALUE_UINT_OCT | VALUE_UINT_DEC | VALUE_UINT_HEX)*;
eventArrayDataUint8X:      EVENT_ADU8X  VALUE_UINTX*;
eventArrayFloat16:         EVENT_AF16   (VALUE_FLOAT_DEC | VALUE_FLOAT_HEX | VALUE_FLOAT_INF | VALUE_FLOAT_NAN | VALUE_FLOAT_SNAN)*;
eventArrayFloat32:         EVENT_AF32   (VALUE_FLOAT_DEC | VALUE_FLOAT_HEX | VALUE_FLOAT_INF | VALUE_FLOAT_NAN | VALUE_FLOAT_SNAN)*;
eventArrayFloat64:         EVENT_AF64   (VALUE_FLOAT_DEC | VALUE_FLOAT_HEX | VALUE_FLOAT_INF | VALUE_FLOAT_NAN | VALUE_FLOAT_SNAN)*;
eventArrayInt16:           EVENT_AI16   (VALUE_INT_BIN | VALUE_INT_OCT | VALUE_INT_DEC | VALUE_INT_HEX)*;
eventArrayInt32:           EVENT_AI32   (VALUE_INT_BIN | VALUE_INT_OCT | VALUE_INT_DEC | VALUE_INT_HEX)*;
eventArrayInt64:           EVENT_AI64   (VALUE_INT_BIN | VALUE_INT_OCT | VALUE_INT_DEC | VALUE_INT_HEX)*;
eventArrayInt8:            EVENT_AI8    (VALUE_INT_BIN | VALUE_INT_OCT | VALUE_INT_DEC | VALUE_INT_HEX)*;
eventArrayUID:             EVENT_AU     UID*;
eventArrayUint16:          EVENT_AU16   (VALUE_UINT_BIN | VALUE_UINT_OCT | VALUE_UINT_DEC | VALUE_UINT_HEX)*;
eventArrayUint16X:         EVENT_AU16X  VALUE_UINTX*;
eventArrayUint32:          EVENT_AU32   (VALUE_UINT_BIN | VALUE_UINT_OCT | VALUE_UINT_DEC | VALUE_UINT_HEX)*;
eventArrayUint32X:         EVENT_AU32X  VALUE_UINTX*;
eventArrayUint64:          EVENT_AU64   (VALUE_UINT_BIN | VALUE_UINT_OCT | VALUE_UINT_DEC | VALUE_UINT_HEX)*;
eventArrayUint64X:         EVENT_AU64X  VALUE_UINTX*;
eventArrayUint8:           EVENT_AU8    (VALUE_UINT_BIN | VALUE_UINT_OCT | VALUE_UINT_DEC | VALUE_UINT_HEX)*;
eventArrayUint8X:          EVENT_AU8X   VALUE_UINTX*;
eventBeginArrayBits:       EVENT_BAB    ;
eventBeginArrayFloat16:    EVENT_BAF16  ;
eventBeginArrayFloat32:    EVENT_BAF32  ;
eventBeginArrayFloat64:    EVENT_BAF64  ;
eventBeginArrayInt16:      EVENT_BAI16  ;
eventBeginArrayInt32:      EVENT_BAI32  ;
eventBeginArrayInt64:      EVENT_BAI64  ;
eventBeginArrayInt8:       EVENT_BAI8   ;
eventBeginArrayUID:        EVENT_BAU    ;
eventBeginArrayUint16:     EVENT_BAU16  ;
eventBeginArrayUint32:     EVENT_BAU32  ;
eventBeginArrayUint64:     EVENT_BAU64  ;
eventBeginArrayUint8:      EVENT_BAU8   ;
eventBeginCustomBinary:    EVENT_BCB    ;
eventBeginCustomText:      EVENT_BCT    ;
eventBeginMedia:           EVENT_BMEDIA ;
eventBeginRemoteReference: EVENT_BREFR  ;
eventBeginResourceId:      EVENT_BRID   ;
eventBeginString:          EVENT_BS     ;
eventBoolean:              EVENT_B      (TRUE | FALSE);
eventCommentMultiline:     EVENT_CM     ;
eventCommentSingleLine:    EVENT_CS     ;
eventCustomBinary:         EVENT_CB     VALUE_UINTX*;
eventCustomText:           EVENT_CT     STRING?;
eventEdge:                 EVENT_EDGE   ;
eventEndContainer:         EVENT_E      ;
eventList:                 EVENT_L      ;
eventMap:                  EVENT_M      ;
eventMarker:               EVENT_MARK   STRING;
eventMedia:                EVENT_MEDIA  MEDIA_TYPE MEDIA_DATA*;
eventNode:                 EVENT_NODE   ;
eventNull:                 EVENT_NULL   ;
eventNumber:               EVENT_N      (INT_BIN | INT_OCT | INT_DEC | INT_HEX | FLOAT_DEC | FLOAT_HEX | FLOAT_INF | FLOAT_NAN | FLOAT_SNAN);
eventPad:                  EVENT_PAD    ;
eventLocalReference:       EVENT_REFL   STRING;
eventRemoteReference:      EVENT_REFR   STRING?;
eventResourceId:           EVENT_RID    STRING?;
eventString:               EVENT_S      STRING?;
eventStructInstance:       EVENT_SI     STRING;
eventStructTemplate:       EVENT_ST     STRING;
eventTime:                 EVENT_T      (DATE | TIME | DATETIME);
eventUID:                  EVENT_UID    UID;
eventVersion:              EVENT_V      VALUE_UINT_DEC;
