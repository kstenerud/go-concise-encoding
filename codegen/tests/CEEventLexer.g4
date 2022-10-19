lexer grammar CEEventLexer;

EVENT_AB:          'ab';                               // Array: Bits
EVENT_AB_ARGS:     'ab='     -> pushMode(MODE_BITS);   // Array: Bits
EVENT_ACL:         'acl='    -> pushMode(MODE_UINT);   // Array Chunk (last chunk)
EVENT_ACM:         'acm='    -> pushMode(MODE_UINT);   // Array Chunk (more chunks follow)
EVENT_ADB:         'adb='    -> pushMode(MODE_BITS);   // Array Data: Bits
EVENT_ADF16:       'adf16='  -> pushMode(MODE_FLOAT);  // Array Data: Float16
EVENT_ADF32:       'adf32='  -> pushMode(MODE_FLOAT);  // Array Data: Float32
EVENT_ADF64:       'adf64='  -> pushMode(MODE_FLOAT);  // Array Data: Float64
EVENT_ADI16:       'adi16='  -> pushMode(MODE_INT);    // Array Data: Int16
EVENT_ADI32:       'adi32='  -> pushMode(MODE_INT);    // Array Data: Int32
EVENT_ADI64:       'adi64='  -> pushMode(MODE_INT);    // Array Data: Int64
EVENT_ADI8:        'adi8='   -> pushMode(MODE_INT);    // Array Data: Int8
EVENT_ADT:         'adt='    -> pushMode(MODE_STRING); // Array Data: Text
EVENT_ADU16:       'adu16='  -> pushMode(MODE_UINT);   // Array Data: Uint16
EVENT_ADU16X:      'adu16x=' -> pushMode(MODE_UINTX);  // Array Data: Uint16 (hex)
EVENT_ADU32:       'adu32='  -> pushMode(MODE_UINT);   // Array Data: Uint32
EVENT_ADU32X:      'adu32x=' -> pushMode(MODE_UINTX);  // Array Data: Uint32 (hex)
EVENT_ADU64:       'adu64='  -> pushMode(MODE_UINT);   // Array Data: Uint64
EVENT_ADU64X:      'adu64x=' -> pushMode(MODE_UINTX);  // Array Data: Uint64 (hex)
EVENT_ADU8:        'adu8='   -> pushMode(MODE_UINT);   // Array Data: Uint8
EVENT_ADU8X:       'adu8x='  -> pushMode(MODE_UINTX);  // Array Data: Uint8 (hex)
EVENT_ADU:         'adu='    -> pushMode(MODE_UID);    // Array Data: UID
EVENT_AF16:        'af16';                             // Array: Float16
EVENT_AF16_ARGS:   'af16='   -> pushMode(MODE_FLOAT);  // Array: Float16
EVENT_AF32:        'af32';                             // Array: Float32
EVENT_AF32_ARGS:   'af32='   -> pushMode(MODE_FLOAT);  // Array: Float32
EVENT_AF64:        'af64';                             // Array: Float64
EVENT_AF64_ARGS:   'af64='   -> pushMode(MODE_FLOAT);  // Array: Float64
EVENT_AI16:        'ai16';                             // Array: Int16
EVENT_AI16_ARGS:   'ai16='   -> pushMode(MODE_INT);    // Array: Int16
EVENT_AI32:        'ai32';                             // Array: Int32
EVENT_AI32_ARGS:   'ai32='   -> pushMode(MODE_INT);    // Array: Int32
EVENT_AI64:        'ai64';                             // Array: Int64
EVENT_AI64_ARGS:   'ai64='   -> pushMode(MODE_INT);    // Array: Int64
EVENT_AI8:         'ai8';                              // Array: Int8
EVENT_AI8_ARGS:    'ai8='    -> pushMode(MODE_INT);    // Array: Int8
EVENT_AU16:        'au16';                             // Array: Uint16
EVENT_AU16_ARGS:   'au16='   -> pushMode(MODE_UINT);   // Array: Uint16
EVENT_AU16X:       'au16x';                            // Array: Uint16 (hex)
EVENT_AU16X_ARGS:  'au16x='  -> pushMode(MODE_UINTX);  // Array: Uint16 (hex)
EVENT_AU32:        'au32';                             // Array: Uint32
EVENT_AU32_ARGS:   'au32='   -> pushMode(MODE_UINT);   // Array: Uint32
EVENT_AU32X:       'au32x';                            // Array: Uint32 (hex)
EVENT_AU32X_ARGS:  'au32x='  -> pushMode(MODE_UINTX);  // Array: Uint32 (hex)
EVENT_AU64:        'au64';                             // Array: Uint64
EVENT_AU64_ARGS:   'au64='   -> pushMode(MODE_UINT);   // Array: Uint64
EVENT_AU64X:       'au64x';                            // Array: Uint64 (hex)
EVENT_AU64X_ARGS:  'au64x='  -> pushMode(MODE_UINTX);  // Array: Uint64 (hex)
EVENT_AU8:         'au8';                              // Array: Uint8
EVENT_AU8_ARGS:    'au8='    -> pushMode(MODE_UINT);   // Array: Uint8
EVENT_AU8X:        'au8x';                             // Array: Uint8 (hex)
EVENT_AU8X_ARGS:   'au8x='   -> pushMode(MODE_UINTX);  // Array: Uint8 (hex)
EVENT_AU:          'au';                               // Array: UID
EVENT_AU_ARGS:     'au='     -> pushMode(MODE_UID);    // Array: UID
EVENT_B:           'b=';                               // Boolean
EVENT_BAB:         'bab';                              // Begin Array: Bits
EVENT_BAF16:       'baf16';                            // Begin Array: Float16
EVENT_BAF32:       'baf32';                            // Begin Array: Float32
EVENT_BAF64:       'baf64';                            // Begin Array: Float64
EVENT_BAI16:       'bai16';                            // Begin Array: Int16
EVENT_BAI32:       'bai32';                            // Begin Array: Int32
EVENT_BAI64:       'bai64';                            // Begin Array: Int64
EVENT_BAI8:        'bai8';                             // Begin Array: Int8
EVENT_BAU16:       'bau16';                            // Begin Array: Uint16
EVENT_BAU32:       'bau32';                            // Begin Array: Uint32
EVENT_BAU64:       'bau64';                            // Begin Array: Uint64
EVENT_BAU8:        'bau8';                             // Begin Array: Uint8
EVENT_BAU:         'bau';                              // Begin Array: UID
EVENT_BCB:         'bcb='    -> pushMode(MODE_UINT);   // Begin Custom Binary
EVENT_BCT:         'bct='    -> pushMode(MODE_UINT);   // Begin Custom Text
EVENT_BMEDIA:      'bmedia=' -> pushMode(MODE_STRING); // Begin Media
EVENT_BREFR:       'brefr';                            // Begin Reference: Remote
EVENT_BRID:        'brid';                             // Begin Resource ID
EVENT_BS:          'bs';                               // Begin String
EVENT_CB:          'cb='     -> pushMode(MODE_CUSTOM_BINARY); // Custom Binary
EVENT_CM:          'cm';                               // Comment: Multiline
EVENT_CM_ARGS:     'cm='     -> pushMode(MODE_STRING); // Comment: Multiline
EVENT_CS:          'cs';                               // Comment: Single Line
EVENT_CS_ARGS:     'cs='     -> pushMode(MODE_STRING); // Comment: Single Line
EVENT_CT:          'ct='     -> pushMode(MODE_CUSTOM_TEXT); // Custom Text
EVENT_E:           'e';                                // End Container
EVENT_EDGE:        'edge';                             // Edge
EVENT_L:           'l';                                // List
EVENT_M:           'm';                                // Map
EVENT_MARK:        'mark='   -> pushMode(MODE_STRING); // Mark
EVENT_MEDIA:       'media='  -> pushMode(MODE_MEDIA);  // Media
EVENT_N:           'n=';                               // Number
EVENT_NODE:        'node';                             // Node
EVENT_NULL:        'null';                             // Null
EVENT_PAD:         'pad';                              // Padding
EVENT_REFL:        'refl='   -> pushMode(MODE_STRING); // Reference: Local
EVENT_REFR:        'refr';                             // Reference: Remote
EVENT_REFR_ARGS:   'refr='   -> pushMode(MODE_STRING); // Reference: Remote
EVENT_RID:         'rid';                              // Resource ID
EVENT_RID_ARGS:    'rid='    -> pushMode(MODE_STRING); // Resource ID
EVENT_SI:          'si='     -> pushMode(MODE_STRING); // Struct Instance
EVENT_ST:          'st='     -> pushMode(MODE_STRING); // Struct Template
EVENT_S:           's';                                // String
EVENT_S_ARGS:      's='      -> pushMode(MODE_STRING); // String
EVENT_T:           't='      -> pushMode(MODE_TIME);   // Time
EVENT_UID:         'uid=';                             // UID
EVENT_V:           'v='      -> pushMode(MODE_UINT);   // Version

TRUE: 'true';
FALSE: 'false';

FLOAT_NAN: F_NAN;
FLOAT_SNAN: F_SNAN;
FLOAT_INF: F_INF;
FLOAT_DEC: F_FLOAT_DEC;
FLOAT_HEX: F_FLOAT_HEX;
INT_BIN: F_INT_BIN;
INT_OCT: F_INT_OCT;
INT_DEC: F_INT_DEC;
INT_HEX: F_INT_HEX;

UID: F_UID;

// ===========================================================================

mode MODE_UINT;
VALUE_UINT_BIN: F_UINT_BIN;
VALUE_UINT_OCT: F_UINT_OCT;
VALUE_UINT_DEC: F_UINT_DEC;
VALUE_UINT_HEX: F_UINT_HEX;
MODE_UINT_WS: F_WHITESPACE -> skip;

// ===========================================================================

mode MODE_UINTX;
VALUE_UINTX: F_DIGITS_HEX;
MODE_UINTX_WS: F_WHITESPACE -> skip;

// ===========================================================================

mode MODE_INT;
VALUE_INT_BIN: F_INT_BIN;
VALUE_INT_OCT: F_INT_OCT;
VALUE_INT_DEC: F_INT_DEC;
VALUE_INT_HEX: F_INT_HEX;
MODE_INT_WS: F_WHITESPACE -> skip;

// ===========================================================================

mode MODE_FLOAT;
VALUE_FLOAT_NAN: F_NAN;
VALUE_FLOAT_SNAN: F_SNAN;
VALUE_FLOAT_INF: F_INF;
VALUE_FLOAT_DEC: F_FLOAT_OR_INT_DEC;
VALUE_FLOAT_HEX: F_FLOAT_OR_INT_HEX;
MODE_FLOAT_WS: F_WHITESPACE -> skip;

// ===========================================================================

mode MODE_UID;
VALUE_UID: F_UID;
MODE_UID_WS: F_WHITESPACE -> skip;

// ===========================================================================

mode MODE_TIME;

TZ_PINT: F_DEC+;
TZ_NINT: F_NEG TZ_PINT;
TZ_INT: TZ_PINT | TZ_NINT;
TZ_COORD: TZ_INT ('.' TZ_PINT)?;
TZ_STRING: [A-Z] F_CHAR_NON_WS*;

TIME_ZONE
   : '/' TZ_STRING
   | '/' TZ_COORD '/' TZ_COORD
   | [+-] TZ_PINT
   ;

TIME: TZ_PINT ':' TZ_PINT ':' TZ_PINT ('.' TZ_PINT)? TIME_ZONE?;
DATE: TZ_INT '-' TZ_PINT '-' TZ_PINT;
DATETIME: TZ_INT '-' TZ_PINT '-' TZ_PINT '/' TZ_PINT ':' TZ_PINT ':' TZ_PINT ('.' TZ_PINT)? TIME_ZONE?;

MODE_TIME_WS: F_WHITESPACE -> skip;

// ===========================================================================

mode MODE_STRING;
STRING: F_STRING;

// ===========================================================================

mode MODE_BYTES;
MODE_BYTES_WS: F_WHITESPACE -> skip;
BYTE: F_BYTE;

// ===========================================================================

mode MODE_BITS;
VALUE_BIT: F_BIT;
MODE_BITS_WS: F_WHITESPACE -> skip;

// ===========================================================================

mode MODE_CUSTOM_BINARY;
CUSTOM_BINARY_TYPE: F_DIGITS_DEC+ -> pushMode(MODE_BYTES);

// ===========================================================================

mode MODE_CUSTOM_TEXT;
CUSTOM_TEXT_TYPE: F_DIGITS_DEC+;
CUSTOM_TEXT_SEPARATOR: ' ' -> skip, pushMode(MODE_STRING);

// ===========================================================================

mode MODE_MEDIA;
MEDIA_TYPE: F_CHAR_NON_WS+ -> pushMode(MODE_BYTES);

// ===========================================================================

fragment F_CHAR_WS: [ \n\r\t];
fragment F_CHAR_NON_WS: ~[ \n\r\t];
fragment F_CHAR_NONNUL: ~[\u0000];

fragment F_WHITESPACE: F_CHAR_WS+;
fragment F_STRING: .*;
fragment F_BYTE: F_HEX F_HEX;

fragment F_PREFIX_BIN:    '0' [bB];
fragment F_PREFIX_OCT:    '0' [oO];
fragment F_PREFIX_HEX:    '0' [xX];
fragment F_NEG:           '-';
fragment F_BIT:           [0-1];
fragment F_OCT:           [0-7];
fragment F_DEC:           [0-9];
fragment F_HEX:           [0-9a-fA-F];
fragment F_DIGITS_BIN:    F_BIT+;
fragment F_DIGITS_OCT:    F_OCT+;
fragment F_DIGITS_DEC:    F_DEC+;
fragment F_DIGITS_HEX:    F_HEX+;
fragment F_FRACTION_DEC:  '.' F_DIGITS_DEC;
fragment F_FRACTION_HEX:  '.' F_DIGITS_HEX;
fragment F_EXPONENT_DEC:  [eE] [+\-]? F_DIGITS_DEC;
fragment F_EXPONENT_HEX:  [pP] [+\-]? F_DIGITS_DEC;
fragment F_NAN:           'nan';
fragment F_SNAN:          'snan';
fragment F_INF:           F_NEG? 'inf';

fragment F_UINT_BIN: F_PREFIX_BIN F_DIGITS_BIN;
fragment F_UINT_OCT: F_PREFIX_OCT F_DIGITS_OCT;
fragment F_UINT_DEC: F_DIGITS_DEC;
fragment F_UINT_HEX: F_PREFIX_HEX F_DIGITS_HEX;

fragment F_INT_BIN: F_NEG? F_UINT_BIN;
fragment F_INT_OCT: F_NEG? F_UINT_OCT;
fragment F_INT_DEC: F_NEG? F_UINT_DEC;
fragment F_INT_HEX: F_NEG? F_UINT_HEX;

fragment F_FLOAT_DEC: (F_INT_DEC F_FRACTION_DEC (F_EXPONENT_DEC)?) | (F_INT_DEC F_EXPONENT_DEC);
fragment F_FLOAT_HEX: (F_INT_HEX F_FRACTION_HEX (F_EXPONENT_HEX)?) | (F_INT_HEX F_EXPONENT_HEX);

fragment F_FLOAT_OR_INT_DEC: F_INT_DEC (F_FRACTION_DEC (F_EXPONENT_DEC)?)?;
fragment F_FLOAT_OR_INT_HEX: F_INT_HEX (F_FRACTION_HEX (F_EXPONENT_HEX)?)?;

fragment F_UID:
   F_HEX F_HEX F_HEX F_HEX F_HEX F_HEX F_HEX F_HEX '-'
   F_HEX F_HEX F_HEX F_HEX '-'
   F_HEX F_HEX F_HEX F_HEX '-'
   F_HEX F_HEX F_HEX F_HEX '-'
   F_HEX F_HEX F_HEX F_HEX F_HEX F_HEX F_HEX F_HEX F_HEX F_HEX F_HEX F_HEX
   ;
