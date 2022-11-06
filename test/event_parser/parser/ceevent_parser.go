// Code generated from /home/karl/Projects/go-concise-encoding/codegen/test/CEEventParser.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // CEEventParser

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type CEEventParser struct {
	*antlr.BaseParser
}

var ceeventparserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func ceeventparserParserInit() {
	staticData := &ceeventparserParserStaticData
	staticData.literalNames = []string{
		"", "'ab'", "'ab='", "'acl='", "'acm='", "'adb='", "'adf16='", "'adf32='",
		"'adf64='", "'adi16='", "'adi32='", "'adi64='", "'adi8='", "'adt='",
		"'adu16='", "'adu16x='", "'adu32='", "'adu32x='", "'adu64='", "'adu64x='",
		"'adu8='", "'adu8x='", "'adu='", "'af16'", "'af16='", "'af32'", "'af32='",
		"'af64'", "'af64='", "'ai16'", "'ai16='", "'ai32'", "'ai32='", "'ai64'",
		"'ai64='", "'ai8'", "'ai8='", "'au16'", "'au16='", "'au16x'", "'au16x='",
		"'au32'", "'au32='", "'au32x'", "'au32x='", "'au64'", "'au64='", "'au64x'",
		"'au64x='", "'au8'", "'au8='", "'au8x'", "'au8x='", "'au'", "'au='",
		"'b='", "'bab'", "'baf16'", "'baf32'", "'baf64'", "'bai16'", "'bai32'",
		"'bai64'", "'bai8'", "'bau16'", "'bau32'", "'bau64'", "'bau8'", "'bau'",
		"'bcb='", "'bct='", "'bmedia='", "'brefr'", "'brid'", "'bs'", "'cb='",
		"'cm'", "'cm='", "'cs'", "'cs='", "'ct='", "'e'", "'edge'", "'l'", "'m'",
		"'mark='", "'media='", "'n='", "'node'", "'null'", "'pad'", "'refl='",
		"'refr'", "'refr='", "'rid'", "'rid='", "'si='", "'st='", "'s'", "'s='",
		"'t='", "'uid='", "'v='", "'true'", "'false'", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "' '",
	}
	staticData.symbolicNames = []string{
		"", "EVENT_AB", "EVENT_AB_ARGS", "EVENT_ACL", "EVENT_ACM", "EVENT_ADB",
		"EVENT_ADF16", "EVENT_ADF32", "EVENT_ADF64", "EVENT_ADI16", "EVENT_ADI32",
		"EVENT_ADI64", "EVENT_ADI8", "EVENT_ADT", "EVENT_ADU16", "EVENT_ADU16X",
		"EVENT_ADU32", "EVENT_ADU32X", "EVENT_ADU64", "EVENT_ADU64X", "EVENT_ADU8",
		"EVENT_ADU8X", "EVENT_ADU", "EVENT_AF16", "EVENT_AF16_ARGS", "EVENT_AF32",
		"EVENT_AF32_ARGS", "EVENT_AF64", "EVENT_AF64_ARGS", "EVENT_AI16", "EVENT_AI16_ARGS",
		"EVENT_AI32", "EVENT_AI32_ARGS", "EVENT_AI64", "EVENT_AI64_ARGS", "EVENT_AI8",
		"EVENT_AI8_ARGS", "EVENT_AU16", "EVENT_AU16_ARGS", "EVENT_AU16X", "EVENT_AU16X_ARGS",
		"EVENT_AU32", "EVENT_AU32_ARGS", "EVENT_AU32X", "EVENT_AU32X_ARGS",
		"EVENT_AU64", "EVENT_AU64_ARGS", "EVENT_AU64X", "EVENT_AU64X_ARGS",
		"EVENT_AU8", "EVENT_AU8_ARGS", "EVENT_AU8X", "EVENT_AU8X_ARGS", "EVENT_AU",
		"EVENT_AU_ARGS", "EVENT_B", "EVENT_BAB", "EVENT_BAF16", "EVENT_BAF32",
		"EVENT_BAF64", "EVENT_BAI16", "EVENT_BAI32", "EVENT_BAI64", "EVENT_BAI8",
		"EVENT_BAU16", "EVENT_BAU32", "EVENT_BAU64", "EVENT_BAU8", "EVENT_BAU",
		"EVENT_BCB", "EVENT_BCT", "EVENT_BMEDIA", "EVENT_BREFR", "EVENT_BRID",
		"EVENT_BS", "EVENT_CB", "EVENT_CM", "EVENT_CM_ARGS", "EVENT_CS", "EVENT_CS_ARGS",
		"EVENT_CT", "EVENT_E", "EVENT_EDGE", "EVENT_L", "EVENT_M", "EVENT_MARK",
		"EVENT_MEDIA", "EVENT_N", "EVENT_NODE", "EVENT_NULL", "EVENT_PAD", "EVENT_REFL",
		"EVENT_REFR", "EVENT_REFR_ARGS", "EVENT_RID", "EVENT_RID_ARGS", "EVENT_SI",
		"EVENT_ST", "EVENT_S", "EVENT_S_ARGS", "EVENT_T", "EVENT_UID", "EVENT_V",
		"TRUE", "FALSE", "FLOAT_NAN", "FLOAT_SNAN", "FLOAT_INF", "FLOAT_DEC",
		"FLOAT_HEX", "INT_BIN", "INT_OCT", "INT_DEC", "INT_HEX", "UID", "VALUE_UINT_BIN",
		"VALUE_UINT_OCT", "VALUE_UINT_DEC", "VALUE_UINT_HEX", "MODE_UINT_WS",
		"VALUE_UINTX", "MODE_UINTX_WS", "VALUE_INT_BIN", "VALUE_INT_OCT", "VALUE_INT_DEC",
		"VALUE_INT_HEX", "MODE_INT_WS", "VALUE_FLOAT_NAN", "VALUE_FLOAT_SNAN",
		"VALUE_FLOAT_INF", "VALUE_FLOAT_DEC", "VALUE_FLOAT_HEX", "MODE_FLOAT_WS",
		"VALUE_UID", "MODE_UID_WS", "TZ_PINT", "TZ_NINT", "TZ_INT", "TZ_COORD",
		"TZ_STRING", "TIME_ZONE", "TIME", "DATE", "DATETIME", "MODE_TIME_WS",
		"STRING", "MODE_BYTES_WS", "BYTE", "VALUE_BIT", "MODE_BITS_WS", "CUSTOM_BINARY_TYPE",
		"CUSTOM_TEXT_TYPE", "CUSTOM_TEXT_SEPARATOR", "MEDIA_TYPE",
	}
	staticData.ruleNames = []string{
		"start", "event", "eventArrayBits", "eventArrayChunkLast", "eventArrayChunkMore",
		"eventArrayDataBits", "eventArrayDataFloat16", "eventArrayDataFloat32",
		"eventArrayDataFloat64", "eventArrayDataInt16", "eventArrayDataInt32",
		"eventArrayDataInt64", "eventArrayDataInt8", "eventArrayDataText", "eventArrayDataUID",
		"eventArrayDataUint16", "eventArrayDataUint16X", "eventArrayDataUint32",
		"eventArrayDataUint32X", "eventArrayDataUint64", "eventArrayDataUint64X",
		"eventArrayDataUint8", "eventArrayDataUint8X", "eventArrayFloat16",
		"eventArrayFloat32", "eventArrayFloat64", "eventArrayInt16", "eventArrayInt32",
		"eventArrayInt64", "eventArrayInt8", "eventArrayUID", "eventArrayUint16",
		"eventArrayUint16X", "eventArrayUint32", "eventArrayUint32X", "eventArrayUint64",
		"eventArrayUint64X", "eventArrayUint8", "eventArrayUint8X", "eventBeginArrayBits",
		"eventBeginArrayFloat16", "eventBeginArrayFloat32", "eventBeginArrayFloat64",
		"eventBeginArrayInt16", "eventBeginArrayInt32", "eventBeginArrayInt64",
		"eventBeginArrayInt8", "eventBeginArrayUID", "eventBeginArrayUint16",
		"eventBeginArrayUint32", "eventBeginArrayUint64", "eventBeginArrayUint8",
		"eventBeginCustomBinary", "eventBeginCustomText", "eventBeginMedia",
		"eventBeginRemoteReference", "eventBeginResourceId", "eventBeginString",
		"eventBoolean", "eventCommentMultiline", "eventCommentSingleLine", "eventCustomBinary",
		"eventCustomText", "eventEdge", "eventEndContainer", "eventList", "eventMap",
		"eventMarker", "eventMedia", "eventNode", "eventNull", "eventNumber",
		"eventPad", "eventLocalReference", "eventRemoteReference", "eventResourceId",
		"eventString", "eventStructInstance", "eventStructTemplate", "eventTime",
		"eventUID", "eventVersion",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 153, 695, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7,
		52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 2, 57, 7, 57,
		2, 58, 7, 58, 2, 59, 7, 59, 2, 60, 7, 60, 2, 61, 7, 61, 2, 62, 7, 62, 2,
		63, 7, 63, 2, 64, 7, 64, 2, 65, 7, 65, 2, 66, 7, 66, 2, 67, 7, 67, 2, 68,
		7, 68, 2, 69, 7, 69, 2, 70, 7, 70, 2, 71, 7, 71, 2, 72, 7, 72, 2, 73, 7,
		73, 2, 74, 7, 74, 2, 75, 7, 75, 2, 76, 7, 76, 2, 77, 7, 77, 2, 78, 7, 78,
		2, 79, 7, 79, 2, 80, 7, 80, 2, 81, 7, 81, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 248, 8, 1, 1, 2, 1, 2, 1, 2,
		5, 2, 253, 8, 2, 10, 2, 12, 2, 256, 9, 2, 3, 2, 258, 8, 2, 1, 3, 1, 3,
		5, 3, 262, 8, 3, 10, 3, 12, 3, 265, 9, 3, 1, 4, 1, 4, 5, 4, 269, 8, 4,
		10, 4, 12, 4, 272, 9, 4, 1, 5, 1, 5, 5, 5, 276, 8, 5, 10, 5, 12, 5, 279,
		9, 5, 1, 6, 1, 6, 5, 6, 283, 8, 6, 10, 6, 12, 6, 286, 9, 6, 1, 7, 1, 7,
		5, 7, 290, 8, 7, 10, 7, 12, 7, 293, 9, 7, 1, 8, 1, 8, 5, 8, 297, 8, 8,
		10, 8, 12, 8, 300, 9, 8, 1, 9, 1, 9, 5, 9, 304, 8, 9, 10, 9, 12, 9, 307,
		9, 9, 1, 10, 1, 10, 5, 10, 311, 8, 10, 10, 10, 12, 10, 314, 9, 10, 1, 11,
		1, 11, 5, 11, 318, 8, 11, 10, 11, 12, 11, 321, 9, 11, 1, 12, 1, 12, 5,
		12, 325, 8, 12, 10, 12, 12, 12, 328, 9, 12, 1, 13, 1, 13, 3, 13, 332, 8,
		13, 1, 14, 1, 14, 5, 14, 336, 8, 14, 10, 14, 12, 14, 339, 9, 14, 1, 15,
		1, 15, 5, 15, 343, 8, 15, 10, 15, 12, 15, 346, 9, 15, 1, 16, 1, 16, 5,
		16, 350, 8, 16, 10, 16, 12, 16, 353, 9, 16, 1, 17, 1, 17, 5, 17, 357, 8,
		17, 10, 17, 12, 17, 360, 9, 17, 1, 18, 1, 18, 5, 18, 364, 8, 18, 10, 18,
		12, 18, 367, 9, 18, 1, 19, 1, 19, 5, 19, 371, 8, 19, 10, 19, 12, 19, 374,
		9, 19, 1, 20, 1, 20, 5, 20, 378, 8, 20, 10, 20, 12, 20, 381, 9, 20, 1,
		21, 1, 21, 5, 21, 385, 8, 21, 10, 21, 12, 21, 388, 9, 21, 1, 22, 1, 22,
		5, 22, 392, 8, 22, 10, 22, 12, 22, 395, 9, 22, 1, 23, 1, 23, 1, 23, 5,
		23, 400, 8, 23, 10, 23, 12, 23, 403, 9, 23, 3, 23, 405, 8, 23, 1, 24, 1,
		24, 1, 24, 5, 24, 410, 8, 24, 10, 24, 12, 24, 413, 9, 24, 3, 24, 415, 8,
		24, 1, 25, 1, 25, 1, 25, 5, 25, 420, 8, 25, 10, 25, 12, 25, 423, 9, 25,
		3, 25, 425, 8, 25, 1, 26, 1, 26, 1, 26, 5, 26, 430, 8, 26, 10, 26, 12,
		26, 433, 9, 26, 3, 26, 435, 8, 26, 1, 27, 1, 27, 1, 27, 5, 27, 440, 8,
		27, 10, 27, 12, 27, 443, 9, 27, 3, 27, 445, 8, 27, 1, 28, 1, 28, 1, 28,
		5, 28, 450, 8, 28, 10, 28, 12, 28, 453, 9, 28, 3, 28, 455, 8, 28, 1, 29,
		1, 29, 1, 29, 5, 29, 460, 8, 29, 10, 29, 12, 29, 463, 9, 29, 3, 29, 465,
		8, 29, 1, 30, 1, 30, 1, 30, 5, 30, 470, 8, 30, 10, 30, 12, 30, 473, 9,
		30, 3, 30, 475, 8, 30, 1, 31, 1, 31, 1, 31, 5, 31, 480, 8, 31, 10, 31,
		12, 31, 483, 9, 31, 3, 31, 485, 8, 31, 1, 32, 1, 32, 1, 32, 5, 32, 490,
		8, 32, 10, 32, 12, 32, 493, 9, 32, 3, 32, 495, 8, 32, 1, 33, 1, 33, 1,
		33, 5, 33, 500, 8, 33, 10, 33, 12, 33, 503, 9, 33, 3, 33, 505, 8, 33, 1,
		34, 1, 34, 1, 34, 5, 34, 510, 8, 34, 10, 34, 12, 34, 513, 9, 34, 3, 34,
		515, 8, 34, 1, 35, 1, 35, 1, 35, 5, 35, 520, 8, 35, 10, 35, 12, 35, 523,
		9, 35, 3, 35, 525, 8, 35, 1, 36, 1, 36, 1, 36, 5, 36, 530, 8, 36, 10, 36,
		12, 36, 533, 9, 36, 3, 36, 535, 8, 36, 1, 37, 1, 37, 1, 37, 5, 37, 540,
		8, 37, 10, 37, 12, 37, 543, 9, 37, 3, 37, 545, 8, 37, 1, 38, 1, 38, 1,
		38, 5, 38, 550, 8, 38, 10, 38, 12, 38, 553, 9, 38, 3, 38, 555, 8, 38, 1,
		39, 1, 39, 1, 40, 1, 40, 1, 41, 1, 41, 1, 42, 1, 42, 1, 43, 1, 43, 1, 44,
		1, 44, 1, 45, 1, 45, 1, 46, 1, 46, 1, 47, 1, 47, 1, 48, 1, 48, 1, 49, 1,
		49, 1, 50, 1, 50, 1, 51, 1, 51, 1, 52, 1, 52, 1, 52, 1, 53, 1, 53, 1, 53,
		1, 54, 1, 54, 1, 54, 1, 55, 1, 55, 1, 56, 1, 56, 1, 57, 1, 57, 1, 58, 1,
		58, 1, 58, 1, 59, 1, 59, 1, 59, 3, 59, 604, 8, 59, 3, 59, 606, 8, 59, 1,
		60, 1, 60, 1, 60, 3, 60, 611, 8, 60, 3, 60, 613, 8, 60, 1, 61, 1, 61, 1,
		61, 5, 61, 618, 8, 61, 10, 61, 12, 61, 621, 9, 61, 1, 62, 1, 62, 1, 62,
		3, 62, 626, 8, 62, 1, 63, 1, 63, 1, 64, 1, 64, 1, 65, 1, 65, 1, 66, 1,
		66, 1, 67, 1, 67, 1, 67, 1, 68, 1, 68, 1, 68, 5, 68, 642, 8, 68, 10, 68,
		12, 68, 645, 9, 68, 1, 69, 1, 69, 1, 70, 1, 70, 1, 71, 1, 71, 1, 71, 1,
		72, 1, 72, 1, 73, 1, 73, 1, 73, 1, 74, 1, 74, 1, 74, 3, 74, 662, 8, 74,
		3, 74, 664, 8, 74, 1, 75, 1, 75, 1, 75, 3, 75, 669, 8, 75, 3, 75, 671,
		8, 75, 1, 76, 1, 76, 1, 76, 3, 76, 676, 8, 76, 3, 76, 678, 8, 76, 1, 77,
		1, 77, 1, 77, 1, 78, 1, 78, 1, 78, 1, 79, 1, 79, 1, 79, 1, 80, 1, 80, 1,
		80, 1, 81, 1, 81, 1, 81, 1, 81, 0, 0, 82, 0, 2, 4, 6, 8, 10, 12, 14, 16,
		18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52,
		54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88,
		90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 110, 112, 114, 116, 118, 120,
		122, 124, 126, 128, 130, 132, 134, 136, 138, 140, 142, 144, 146, 148, 150,
		152, 154, 156, 158, 160, 162, 0, 6, 1, 0, 127, 131, 1, 0, 122, 125, 1,
		0, 115, 118, 1, 0, 103, 104, 1, 0, 105, 113, 1, 0, 141, 143, 758, 0, 164,
		1, 0, 0, 0, 2, 247, 1, 0, 0, 0, 4, 257, 1, 0, 0, 0, 6, 259, 1, 0, 0, 0,
		8, 266, 1, 0, 0, 0, 10, 273, 1, 0, 0, 0, 12, 280, 1, 0, 0, 0, 14, 287,
		1, 0, 0, 0, 16, 294, 1, 0, 0, 0, 18, 301, 1, 0, 0, 0, 20, 308, 1, 0, 0,
		0, 22, 315, 1, 0, 0, 0, 24, 322, 1, 0, 0, 0, 26, 329, 1, 0, 0, 0, 28, 333,
		1, 0, 0, 0, 30, 340, 1, 0, 0, 0, 32, 347, 1, 0, 0, 0, 34, 354, 1, 0, 0,
		0, 36, 361, 1, 0, 0, 0, 38, 368, 1, 0, 0, 0, 40, 375, 1, 0, 0, 0, 42, 382,
		1, 0, 0, 0, 44, 389, 1, 0, 0, 0, 46, 404, 1, 0, 0, 0, 48, 414, 1, 0, 0,
		0, 50, 424, 1, 0, 0, 0, 52, 434, 1, 0, 0, 0, 54, 444, 1, 0, 0, 0, 56, 454,
		1, 0, 0, 0, 58, 464, 1, 0, 0, 0, 60, 474, 1, 0, 0, 0, 62, 484, 1, 0, 0,
		0, 64, 494, 1, 0, 0, 0, 66, 504, 1, 0, 0, 0, 68, 514, 1, 0, 0, 0, 70, 524,
		1, 0, 0, 0, 72, 534, 1, 0, 0, 0, 74, 544, 1, 0, 0, 0, 76, 554, 1, 0, 0,
		0, 78, 556, 1, 0, 0, 0, 80, 558, 1, 0, 0, 0, 82, 560, 1, 0, 0, 0, 84, 562,
		1, 0, 0, 0, 86, 564, 1, 0, 0, 0, 88, 566, 1, 0, 0, 0, 90, 568, 1, 0, 0,
		0, 92, 570, 1, 0, 0, 0, 94, 572, 1, 0, 0, 0, 96, 574, 1, 0, 0, 0, 98, 576,
		1, 0, 0, 0, 100, 578, 1, 0, 0, 0, 102, 580, 1, 0, 0, 0, 104, 582, 1, 0,
		0, 0, 106, 585, 1, 0, 0, 0, 108, 588, 1, 0, 0, 0, 110, 591, 1, 0, 0, 0,
		112, 593, 1, 0, 0, 0, 114, 595, 1, 0, 0, 0, 116, 597, 1, 0, 0, 0, 118,
		605, 1, 0, 0, 0, 120, 612, 1, 0, 0, 0, 122, 614, 1, 0, 0, 0, 124, 622,
		1, 0, 0, 0, 126, 627, 1, 0, 0, 0, 128, 629, 1, 0, 0, 0, 130, 631, 1, 0,
		0, 0, 132, 633, 1, 0, 0, 0, 134, 635, 1, 0, 0, 0, 136, 638, 1, 0, 0, 0,
		138, 646, 1, 0, 0, 0, 140, 648, 1, 0, 0, 0, 142, 650, 1, 0, 0, 0, 144,
		653, 1, 0, 0, 0, 146, 655, 1, 0, 0, 0, 148, 663, 1, 0, 0, 0, 150, 670,
		1, 0, 0, 0, 152, 677, 1, 0, 0, 0, 154, 679, 1, 0, 0, 0, 156, 682, 1, 0,
		0, 0, 158, 685, 1, 0, 0, 0, 160, 688, 1, 0, 0, 0, 162, 691, 1, 0, 0, 0,
		164, 165, 3, 2, 1, 0, 165, 166, 5, 0, 0, 1, 166, 1, 1, 0, 0, 0, 167, 248,
		3, 4, 2, 0, 168, 248, 3, 6, 3, 0, 169, 248, 3, 8, 4, 0, 170, 248, 3, 10,
		5, 0, 171, 248, 3, 12, 6, 0, 172, 248, 3, 14, 7, 0, 173, 248, 3, 16, 8,
		0, 174, 248, 3, 18, 9, 0, 175, 248, 3, 20, 10, 0, 176, 248, 3, 22, 11,
		0, 177, 248, 3, 24, 12, 0, 178, 248, 3, 26, 13, 0, 179, 248, 3, 28, 14,
		0, 180, 248, 3, 30, 15, 0, 181, 248, 3, 32, 16, 0, 182, 248, 3, 34, 17,
		0, 183, 248, 3, 36, 18, 0, 184, 248, 3, 38, 19, 0, 185, 248, 3, 40, 20,
		0, 186, 248, 3, 42, 21, 0, 187, 248, 3, 44, 22, 0, 188, 248, 3, 46, 23,
		0, 189, 248, 3, 48, 24, 0, 190, 248, 3, 50, 25, 0, 191, 248, 3, 52, 26,
		0, 192, 248, 3, 54, 27, 0, 193, 248, 3, 56, 28, 0, 194, 248, 3, 58, 29,
		0, 195, 248, 3, 60, 30, 0, 196, 248, 3, 62, 31, 0, 197, 248, 3, 64, 32,
		0, 198, 248, 3, 66, 33, 0, 199, 248, 3, 68, 34, 0, 200, 248, 3, 70, 35,
		0, 201, 248, 3, 72, 36, 0, 202, 248, 3, 74, 37, 0, 203, 248, 3, 76, 38,
		0, 204, 248, 3, 78, 39, 0, 205, 248, 3, 80, 40, 0, 206, 248, 3, 82, 41,
		0, 207, 248, 3, 84, 42, 0, 208, 248, 3, 86, 43, 0, 209, 248, 3, 88, 44,
		0, 210, 248, 3, 90, 45, 0, 211, 248, 3, 92, 46, 0, 212, 248, 3, 94, 47,
		0, 213, 248, 3, 96, 48, 0, 214, 248, 3, 98, 49, 0, 215, 248, 3, 100, 50,
		0, 216, 248, 3, 102, 51, 0, 217, 248, 3, 104, 52, 0, 218, 248, 3, 106,
		53, 0, 219, 248, 3, 108, 54, 0, 220, 248, 3, 112, 56, 0, 221, 248, 3, 110,
		55, 0, 222, 248, 3, 114, 57, 0, 223, 248, 3, 116, 58, 0, 224, 248, 3, 120,
		60, 0, 225, 248, 3, 118, 59, 0, 226, 248, 3, 122, 61, 0, 227, 248, 3, 124,
		62, 0, 228, 248, 3, 126, 63, 0, 229, 248, 3, 128, 64, 0, 230, 248, 3, 130,
		65, 0, 231, 248, 3, 132, 66, 0, 232, 248, 3, 134, 67, 0, 233, 248, 3, 136,
		68, 0, 234, 248, 3, 138, 69, 0, 235, 248, 3, 140, 70, 0, 236, 248, 3, 142,
		71, 0, 237, 248, 3, 144, 72, 0, 238, 248, 3, 146, 73, 0, 239, 248, 3, 148,
		74, 0, 240, 248, 3, 150, 75, 0, 241, 248, 3, 152, 76, 0, 242, 248, 3, 154,
		77, 0, 243, 248, 3, 156, 78, 0, 244, 248, 3, 158, 79, 0, 245, 248, 3, 160,
		80, 0, 246, 248, 3, 162, 81, 0, 247, 167, 1, 0, 0, 0, 247, 168, 1, 0, 0,
		0, 247, 169, 1, 0, 0, 0, 247, 170, 1, 0, 0, 0, 247, 171, 1, 0, 0, 0, 247,
		172, 1, 0, 0, 0, 247, 173, 1, 0, 0, 0, 247, 174, 1, 0, 0, 0, 247, 175,
		1, 0, 0, 0, 247, 176, 1, 0, 0, 0, 247, 177, 1, 0, 0, 0, 247, 178, 1, 0,
		0, 0, 247, 179, 1, 0, 0, 0, 247, 180, 1, 0, 0, 0, 247, 181, 1, 0, 0, 0,
		247, 182, 1, 0, 0, 0, 247, 183, 1, 0, 0, 0, 247, 184, 1, 0, 0, 0, 247,
		185, 1, 0, 0, 0, 247, 186, 1, 0, 0, 0, 247, 187, 1, 0, 0, 0, 247, 188,
		1, 0, 0, 0, 247, 189, 1, 0, 0, 0, 247, 190, 1, 0, 0, 0, 247, 191, 1, 0,
		0, 0, 247, 192, 1, 0, 0, 0, 247, 193, 1, 0, 0, 0, 247, 194, 1, 0, 0, 0,
		247, 195, 1, 0, 0, 0, 247, 196, 1, 0, 0, 0, 247, 197, 1, 0, 0, 0, 247,
		198, 1, 0, 0, 0, 247, 199, 1, 0, 0, 0, 247, 200, 1, 0, 0, 0, 247, 201,
		1, 0, 0, 0, 247, 202, 1, 0, 0, 0, 247, 203, 1, 0, 0, 0, 247, 204, 1, 0,
		0, 0, 247, 205, 1, 0, 0, 0, 247, 206, 1, 0, 0, 0, 247, 207, 1, 0, 0, 0,
		247, 208, 1, 0, 0, 0, 247, 209, 1, 0, 0, 0, 247, 210, 1, 0, 0, 0, 247,
		211, 1, 0, 0, 0, 247, 212, 1, 0, 0, 0, 247, 213, 1, 0, 0, 0, 247, 214,
		1, 0, 0, 0, 247, 215, 1, 0, 0, 0, 247, 216, 1, 0, 0, 0, 247, 217, 1, 0,
		0, 0, 247, 218, 1, 0, 0, 0, 247, 219, 1, 0, 0, 0, 247, 220, 1, 0, 0, 0,
		247, 221, 1, 0, 0, 0, 247, 222, 1, 0, 0, 0, 247, 223, 1, 0, 0, 0, 247,
		224, 1, 0, 0, 0, 247, 225, 1, 0, 0, 0, 247, 226, 1, 0, 0, 0, 247, 227,
		1, 0, 0, 0, 247, 228, 1, 0, 0, 0, 247, 229, 1, 0, 0, 0, 247, 230, 1, 0,
		0, 0, 247, 231, 1, 0, 0, 0, 247, 232, 1, 0, 0, 0, 247, 233, 1, 0, 0, 0,
		247, 234, 1, 0, 0, 0, 247, 235, 1, 0, 0, 0, 247, 236, 1, 0, 0, 0, 247,
		237, 1, 0, 0, 0, 247, 238, 1, 0, 0, 0, 247, 239, 1, 0, 0, 0, 247, 240,
		1, 0, 0, 0, 247, 241, 1, 0, 0, 0, 247, 242, 1, 0, 0, 0, 247, 243, 1, 0,
		0, 0, 247, 244, 1, 0, 0, 0, 247, 245, 1, 0, 0, 0, 247, 246, 1, 0, 0, 0,
		248, 3, 1, 0, 0, 0, 249, 258, 5, 1, 0, 0, 250, 254, 5, 2, 0, 0, 251, 253,
		5, 148, 0, 0, 252, 251, 1, 0, 0, 0, 253, 256, 1, 0, 0, 0, 254, 252, 1,
		0, 0, 0, 254, 255, 1, 0, 0, 0, 255, 258, 1, 0, 0, 0, 256, 254, 1, 0, 0,
		0, 257, 249, 1, 0, 0, 0, 257, 250, 1, 0, 0, 0, 258, 5, 1, 0, 0, 0, 259,
		263, 5, 3, 0, 0, 260, 262, 5, 117, 0, 0, 261, 260, 1, 0, 0, 0, 262, 265,
		1, 0, 0, 0, 263, 261, 1, 0, 0, 0, 263, 264, 1, 0, 0, 0, 264, 7, 1, 0, 0,
		0, 265, 263, 1, 0, 0, 0, 266, 270, 5, 4, 0, 0, 267, 269, 5, 117, 0, 0,
		268, 267, 1, 0, 0, 0, 269, 272, 1, 0, 0, 0, 270, 268, 1, 0, 0, 0, 270,
		271, 1, 0, 0, 0, 271, 9, 1, 0, 0, 0, 272, 270, 1, 0, 0, 0, 273, 277, 5,
		5, 0, 0, 274, 276, 5, 148, 0, 0, 275, 274, 1, 0, 0, 0, 276, 279, 1, 0,
		0, 0, 277, 275, 1, 0, 0, 0, 277, 278, 1, 0, 0, 0, 278, 11, 1, 0, 0, 0,
		279, 277, 1, 0, 0, 0, 280, 284, 5, 6, 0, 0, 281, 283, 7, 0, 0, 0, 282,
		281, 1, 0, 0, 0, 283, 286, 1, 0, 0, 0, 284, 282, 1, 0, 0, 0, 284, 285,
		1, 0, 0, 0, 285, 13, 1, 0, 0, 0, 286, 284, 1, 0, 0, 0, 287, 291, 5, 7,
		0, 0, 288, 290, 7, 0, 0, 0, 289, 288, 1, 0, 0, 0, 290, 293, 1, 0, 0, 0,
		291, 289, 1, 0, 0, 0, 291, 292, 1, 0, 0, 0, 292, 15, 1, 0, 0, 0, 293, 291,
		1, 0, 0, 0, 294, 298, 5, 8, 0, 0, 295, 297, 7, 0, 0, 0, 296, 295, 1, 0,
		0, 0, 297, 300, 1, 0, 0, 0, 298, 296, 1, 0, 0, 0, 298, 299, 1, 0, 0, 0,
		299, 17, 1, 0, 0, 0, 300, 298, 1, 0, 0, 0, 301, 305, 5, 9, 0, 0, 302, 304,
		7, 1, 0, 0, 303, 302, 1, 0, 0, 0, 304, 307, 1, 0, 0, 0, 305, 303, 1, 0,
		0, 0, 305, 306, 1, 0, 0, 0, 306, 19, 1, 0, 0, 0, 307, 305, 1, 0, 0, 0,
		308, 312, 5, 10, 0, 0, 309, 311, 7, 1, 0, 0, 310, 309, 1, 0, 0, 0, 311,
		314, 1, 0, 0, 0, 312, 310, 1, 0, 0, 0, 312, 313, 1, 0, 0, 0, 313, 21, 1,
		0, 0, 0, 314, 312, 1, 0, 0, 0, 315, 319, 5, 11, 0, 0, 316, 318, 7, 1, 0,
		0, 317, 316, 1, 0, 0, 0, 318, 321, 1, 0, 0, 0, 319, 317, 1, 0, 0, 0, 319,
		320, 1, 0, 0, 0, 320, 23, 1, 0, 0, 0, 321, 319, 1, 0, 0, 0, 322, 326, 5,
		12, 0, 0, 323, 325, 7, 1, 0, 0, 324, 323, 1, 0, 0, 0, 325, 328, 1, 0, 0,
		0, 326, 324, 1, 0, 0, 0, 326, 327, 1, 0, 0, 0, 327, 25, 1, 0, 0, 0, 328,
		326, 1, 0, 0, 0, 329, 331, 5, 13, 0, 0, 330, 332, 5, 145, 0, 0, 331, 330,
		1, 0, 0, 0, 331, 332, 1, 0, 0, 0, 332, 27, 1, 0, 0, 0, 333, 337, 5, 22,
		0, 0, 334, 336, 5, 133, 0, 0, 335, 334, 1, 0, 0, 0, 336, 339, 1, 0, 0,
		0, 337, 335, 1, 0, 0, 0, 337, 338, 1, 0, 0, 0, 338, 29, 1, 0, 0, 0, 339,
		337, 1, 0, 0, 0, 340, 344, 5, 14, 0, 0, 341, 343, 7, 2, 0, 0, 342, 341,
		1, 0, 0, 0, 343, 346, 1, 0, 0, 0, 344, 342, 1, 0, 0, 0, 344, 345, 1, 0,
		0, 0, 345, 31, 1, 0, 0, 0, 346, 344, 1, 0, 0, 0, 347, 351, 5, 15, 0, 0,
		348, 350, 5, 120, 0, 0, 349, 348, 1, 0, 0, 0, 350, 353, 1, 0, 0, 0, 351,
		349, 1, 0, 0, 0, 351, 352, 1, 0, 0, 0, 352, 33, 1, 0, 0, 0, 353, 351, 1,
		0, 0, 0, 354, 358, 5, 16, 0, 0, 355, 357, 7, 2, 0, 0, 356, 355, 1, 0, 0,
		0, 357, 360, 1, 0, 0, 0, 358, 356, 1, 0, 0, 0, 358, 359, 1, 0, 0, 0, 359,
		35, 1, 0, 0, 0, 360, 358, 1, 0, 0, 0, 361, 365, 5, 17, 0, 0, 362, 364,
		5, 120, 0, 0, 363, 362, 1, 0, 0, 0, 364, 367, 1, 0, 0, 0, 365, 363, 1,
		0, 0, 0, 365, 366, 1, 0, 0, 0, 366, 37, 1, 0, 0, 0, 367, 365, 1, 0, 0,
		0, 368, 372, 5, 18, 0, 0, 369, 371, 7, 2, 0, 0, 370, 369, 1, 0, 0, 0, 371,
		374, 1, 0, 0, 0, 372, 370, 1, 0, 0, 0, 372, 373, 1, 0, 0, 0, 373, 39, 1,
		0, 0, 0, 374, 372, 1, 0, 0, 0, 375, 379, 5, 19, 0, 0, 376, 378, 5, 120,
		0, 0, 377, 376, 1, 0, 0, 0, 378, 381, 1, 0, 0, 0, 379, 377, 1, 0, 0, 0,
		379, 380, 1, 0, 0, 0, 380, 41, 1, 0, 0, 0, 381, 379, 1, 0, 0, 0, 382, 386,
		5, 20, 0, 0, 383, 385, 7, 2, 0, 0, 384, 383, 1, 0, 0, 0, 385, 388, 1, 0,
		0, 0, 386, 384, 1, 0, 0, 0, 386, 387, 1, 0, 0, 0, 387, 43, 1, 0, 0, 0,
		388, 386, 1, 0, 0, 0, 389, 393, 5, 21, 0, 0, 390, 392, 5, 120, 0, 0, 391,
		390, 1, 0, 0, 0, 392, 395, 1, 0, 0, 0, 393, 391, 1, 0, 0, 0, 393, 394,
		1, 0, 0, 0, 394, 45, 1, 0, 0, 0, 395, 393, 1, 0, 0, 0, 396, 405, 5, 23,
		0, 0, 397, 401, 5, 24, 0, 0, 398, 400, 7, 0, 0, 0, 399, 398, 1, 0, 0, 0,
		400, 403, 1, 0, 0, 0, 401, 399, 1, 0, 0, 0, 401, 402, 1, 0, 0, 0, 402,
		405, 1, 0, 0, 0, 403, 401, 1, 0, 0, 0, 404, 396, 1, 0, 0, 0, 404, 397,
		1, 0, 0, 0, 405, 47, 1, 0, 0, 0, 406, 415, 5, 25, 0, 0, 407, 411, 5, 26,
		0, 0, 408, 410, 7, 0, 0, 0, 409, 408, 1, 0, 0, 0, 410, 413, 1, 0, 0, 0,
		411, 409, 1, 0, 0, 0, 411, 412, 1, 0, 0, 0, 412, 415, 1, 0, 0, 0, 413,
		411, 1, 0, 0, 0, 414, 406, 1, 0, 0, 0, 414, 407, 1, 0, 0, 0, 415, 49, 1,
		0, 0, 0, 416, 425, 5, 27, 0, 0, 417, 421, 5, 28, 0, 0, 418, 420, 7, 0,
		0, 0, 419, 418, 1, 0, 0, 0, 420, 423, 1, 0, 0, 0, 421, 419, 1, 0, 0, 0,
		421, 422, 1, 0, 0, 0, 422, 425, 1, 0, 0, 0, 423, 421, 1, 0, 0, 0, 424,
		416, 1, 0, 0, 0, 424, 417, 1, 0, 0, 0, 425, 51, 1, 0, 0, 0, 426, 435, 5,
		29, 0, 0, 427, 431, 5, 30, 0, 0, 428, 430, 7, 1, 0, 0, 429, 428, 1, 0,
		0, 0, 430, 433, 1, 0, 0, 0, 431, 429, 1, 0, 0, 0, 431, 432, 1, 0, 0, 0,
		432, 435, 1, 0, 0, 0, 433, 431, 1, 0, 0, 0, 434, 426, 1, 0, 0, 0, 434,
		427, 1, 0, 0, 0, 435, 53, 1, 0, 0, 0, 436, 445, 5, 31, 0, 0, 437, 441,
		5, 32, 0, 0, 438, 440, 7, 1, 0, 0, 439, 438, 1, 0, 0, 0, 440, 443, 1, 0,
		0, 0, 441, 439, 1, 0, 0, 0, 441, 442, 1, 0, 0, 0, 442, 445, 1, 0, 0, 0,
		443, 441, 1, 0, 0, 0, 444, 436, 1, 0, 0, 0, 444, 437, 1, 0, 0, 0, 445,
		55, 1, 0, 0, 0, 446, 455, 5, 33, 0, 0, 447, 451, 5, 34, 0, 0, 448, 450,
		7, 1, 0, 0, 449, 448, 1, 0, 0, 0, 450, 453, 1, 0, 0, 0, 451, 449, 1, 0,
		0, 0, 451, 452, 1, 0, 0, 0, 452, 455, 1, 0, 0, 0, 453, 451, 1, 0, 0, 0,
		454, 446, 1, 0, 0, 0, 454, 447, 1, 0, 0, 0, 455, 57, 1, 0, 0, 0, 456, 465,
		5, 35, 0, 0, 457, 461, 5, 36, 0, 0, 458, 460, 7, 1, 0, 0, 459, 458, 1,
		0, 0, 0, 460, 463, 1, 0, 0, 0, 461, 459, 1, 0, 0, 0, 461, 462, 1, 0, 0,
		0, 462, 465, 1, 0, 0, 0, 463, 461, 1, 0, 0, 0, 464, 456, 1, 0, 0, 0, 464,
		457, 1, 0, 0, 0, 465, 59, 1, 0, 0, 0, 466, 475, 5, 53, 0, 0, 467, 471,
		5, 54, 0, 0, 468, 470, 5, 133, 0, 0, 469, 468, 1, 0, 0, 0, 470, 473, 1,
		0, 0, 0, 471, 469, 1, 0, 0, 0, 471, 472, 1, 0, 0, 0, 472, 475, 1, 0, 0,
		0, 473, 471, 1, 0, 0, 0, 474, 466, 1, 0, 0, 0, 474, 467, 1, 0, 0, 0, 475,
		61, 1, 0, 0, 0, 476, 485, 5, 37, 0, 0, 477, 481, 5, 38, 0, 0, 478, 480,
		7, 2, 0, 0, 479, 478, 1, 0, 0, 0, 480, 483, 1, 0, 0, 0, 481, 479, 1, 0,
		0, 0, 481, 482, 1, 0, 0, 0, 482, 485, 1, 0, 0, 0, 483, 481, 1, 0, 0, 0,
		484, 476, 1, 0, 0, 0, 484, 477, 1, 0, 0, 0, 485, 63, 1, 0, 0, 0, 486, 495,
		5, 39, 0, 0, 487, 491, 5, 40, 0, 0, 488, 490, 5, 120, 0, 0, 489, 488, 1,
		0, 0, 0, 490, 493, 1, 0, 0, 0, 491, 489, 1, 0, 0, 0, 491, 492, 1, 0, 0,
		0, 492, 495, 1, 0, 0, 0, 493, 491, 1, 0, 0, 0, 494, 486, 1, 0, 0, 0, 494,
		487, 1, 0, 0, 0, 495, 65, 1, 0, 0, 0, 496, 505, 5, 41, 0, 0, 497, 501,
		5, 42, 0, 0, 498, 500, 7, 2, 0, 0, 499, 498, 1, 0, 0, 0, 500, 503, 1, 0,
		0, 0, 501, 499, 1, 0, 0, 0, 501, 502, 1, 0, 0, 0, 502, 505, 1, 0, 0, 0,
		503, 501, 1, 0, 0, 0, 504, 496, 1, 0, 0, 0, 504, 497, 1, 0, 0, 0, 505,
		67, 1, 0, 0, 0, 506, 515, 5, 43, 0, 0, 507, 511, 5, 44, 0, 0, 508, 510,
		5, 120, 0, 0, 509, 508, 1, 0, 0, 0, 510, 513, 1, 0, 0, 0, 511, 509, 1,
		0, 0, 0, 511, 512, 1, 0, 0, 0, 512, 515, 1, 0, 0, 0, 513, 511, 1, 0, 0,
		0, 514, 506, 1, 0, 0, 0, 514, 507, 1, 0, 0, 0, 515, 69, 1, 0, 0, 0, 516,
		525, 5, 45, 0, 0, 517, 521, 5, 46, 0, 0, 518, 520, 7, 2, 0, 0, 519, 518,
		1, 0, 0, 0, 520, 523, 1, 0, 0, 0, 521, 519, 1, 0, 0, 0, 521, 522, 1, 0,
		0, 0, 522, 525, 1, 0, 0, 0, 523, 521, 1, 0, 0, 0, 524, 516, 1, 0, 0, 0,
		524, 517, 1, 0, 0, 0, 525, 71, 1, 0, 0, 0, 526, 535, 5, 47, 0, 0, 527,
		531, 5, 48, 0, 0, 528, 530, 5, 120, 0, 0, 529, 528, 1, 0, 0, 0, 530, 533,
		1, 0, 0, 0, 531, 529, 1, 0, 0, 0, 531, 532, 1, 0, 0, 0, 532, 535, 1, 0,
		0, 0, 533, 531, 1, 0, 0, 0, 534, 526, 1, 0, 0, 0, 534, 527, 1, 0, 0, 0,
		535, 73, 1, 0, 0, 0, 536, 545, 5, 49, 0, 0, 537, 541, 5, 50, 0, 0, 538,
		540, 7, 2, 0, 0, 539, 538, 1, 0, 0, 0, 540, 543, 1, 0, 0, 0, 541, 539,
		1, 0, 0, 0, 541, 542, 1, 0, 0, 0, 542, 545, 1, 0, 0, 0, 543, 541, 1, 0,
		0, 0, 544, 536, 1, 0, 0, 0, 544, 537, 1, 0, 0, 0, 545, 75, 1, 0, 0, 0,
		546, 555, 5, 51, 0, 0, 547, 551, 5, 52, 0, 0, 548, 550, 5, 120, 0, 0, 549,
		548, 1, 0, 0, 0, 550, 553, 1, 0, 0, 0, 551, 549, 1, 0, 0, 0, 551, 552,
		1, 0, 0, 0, 552, 555, 1, 0, 0, 0, 553, 551, 1, 0, 0, 0, 554, 546, 1, 0,
		0, 0, 554, 547, 1, 0, 0, 0, 555, 77, 1, 0, 0, 0, 556, 557, 5, 56, 0, 0,
		557, 79, 1, 0, 0, 0, 558, 559, 5, 57, 0, 0, 559, 81, 1, 0, 0, 0, 560, 561,
		5, 58, 0, 0, 561, 83, 1, 0, 0, 0, 562, 563, 5, 59, 0, 0, 563, 85, 1, 0,
		0, 0, 564, 565, 5, 60, 0, 0, 565, 87, 1, 0, 0, 0, 566, 567, 5, 61, 0, 0,
		567, 89, 1, 0, 0, 0, 568, 569, 5, 62, 0, 0, 569, 91, 1, 0, 0, 0, 570, 571,
		5, 63, 0, 0, 571, 93, 1, 0, 0, 0, 572, 573, 5, 68, 0, 0, 573, 95, 1, 0,
		0, 0, 574, 575, 5, 64, 0, 0, 575, 97, 1, 0, 0, 0, 576, 577, 5, 65, 0, 0,
		577, 99, 1, 0, 0, 0, 578, 579, 5, 66, 0, 0, 579, 101, 1, 0, 0, 0, 580,
		581, 5, 67, 0, 0, 581, 103, 1, 0, 0, 0, 582, 583, 5, 69, 0, 0, 583, 584,
		5, 117, 0, 0, 584, 105, 1, 0, 0, 0, 585, 586, 5, 70, 0, 0, 586, 587, 5,
		117, 0, 0, 587, 107, 1, 0, 0, 0, 588, 589, 5, 71, 0, 0, 589, 590, 5, 145,
		0, 0, 590, 109, 1, 0, 0, 0, 591, 592, 5, 72, 0, 0, 592, 111, 1, 0, 0, 0,
		593, 594, 5, 73, 0, 0, 594, 113, 1, 0, 0, 0, 595, 596, 5, 74, 0, 0, 596,
		115, 1, 0, 0, 0, 597, 598, 5, 55, 0, 0, 598, 599, 7, 3, 0, 0, 599, 117,
		1, 0, 0, 0, 600, 606, 5, 76, 0, 0, 601, 603, 5, 77, 0, 0, 602, 604, 5,
		145, 0, 0, 603, 602, 1, 0, 0, 0, 603, 604, 1, 0, 0, 0, 604, 606, 1, 0,
		0, 0, 605, 600, 1, 0, 0, 0, 605, 601, 1, 0, 0, 0, 606, 119, 1, 0, 0, 0,
		607, 613, 5, 78, 0, 0, 608, 610, 5, 79, 0, 0, 609, 611, 5, 145, 0, 0, 610,
		609, 1, 0, 0, 0, 610, 611, 1, 0, 0, 0, 611, 613, 1, 0, 0, 0, 612, 607,
		1, 0, 0, 0, 612, 608, 1, 0, 0, 0, 613, 121, 1, 0, 0, 0, 614, 615, 5, 75,
		0, 0, 615, 619, 5, 150, 0, 0, 616, 618, 5, 147, 0, 0, 617, 616, 1, 0, 0,
		0, 618, 621, 1, 0, 0, 0, 619, 617, 1, 0, 0, 0, 619, 620, 1, 0, 0, 0, 620,
		123, 1, 0, 0, 0, 621, 619, 1, 0, 0, 0, 622, 623, 5, 80, 0, 0, 623, 625,
		5, 151, 0, 0, 624, 626, 5, 145, 0, 0, 625, 624, 1, 0, 0, 0, 625, 626, 1,
		0, 0, 0, 626, 125, 1, 0, 0, 0, 627, 628, 5, 82, 0, 0, 628, 127, 1, 0, 0,
		0, 629, 630, 5, 81, 0, 0, 630, 129, 1, 0, 0, 0, 631, 632, 5, 83, 0, 0,
		632, 131, 1, 0, 0, 0, 633, 634, 5, 84, 0, 0, 634, 133, 1, 0, 0, 0, 635,
		636, 5, 85, 0, 0, 636, 637, 5, 145, 0, 0, 637, 135, 1, 0, 0, 0, 638, 639,
		5, 86, 0, 0, 639, 643, 5, 153, 0, 0, 640, 642, 5, 147, 0, 0, 641, 640,
		1, 0, 0, 0, 642, 645, 1, 0, 0, 0, 643, 641, 1, 0, 0, 0, 643, 644, 1, 0,
		0, 0, 644, 137, 1, 0, 0, 0, 645, 643, 1, 0, 0, 0, 646, 647, 5, 88, 0, 0,
		647, 139, 1, 0, 0, 0, 648, 649, 5, 89, 0, 0, 649, 141, 1, 0, 0, 0, 650,
		651, 5, 87, 0, 0, 651, 652, 7, 4, 0, 0, 652, 143, 1, 0, 0, 0, 653, 654,
		5, 90, 0, 0, 654, 145, 1, 0, 0, 0, 655, 656, 5, 91, 0, 0, 656, 657, 5,
		145, 0, 0, 657, 147, 1, 0, 0, 0, 658, 664, 5, 92, 0, 0, 659, 661, 5, 93,
		0, 0, 660, 662, 5, 145, 0, 0, 661, 660, 1, 0, 0, 0, 661, 662, 1, 0, 0,
		0, 662, 664, 1, 0, 0, 0, 663, 658, 1, 0, 0, 0, 663, 659, 1, 0, 0, 0, 664,
		149, 1, 0, 0, 0, 665, 671, 5, 94, 0, 0, 666, 668, 5, 95, 0, 0, 667, 669,
		5, 145, 0, 0, 668, 667, 1, 0, 0, 0, 668, 669, 1, 0, 0, 0, 669, 671, 1,
		0, 0, 0, 670, 665, 1, 0, 0, 0, 670, 666, 1, 0, 0, 0, 671, 151, 1, 0, 0,
		0, 672, 678, 5, 98, 0, 0, 673, 675, 5, 99, 0, 0, 674, 676, 5, 145, 0, 0,
		675, 674, 1, 0, 0, 0, 675, 676, 1, 0, 0, 0, 676, 678, 1, 0, 0, 0, 677,
		672, 1, 0, 0, 0, 677, 673, 1, 0, 0, 0, 678, 153, 1, 0, 0, 0, 679, 680,
		5, 96, 0, 0, 680, 681, 5, 145, 0, 0, 681, 155, 1, 0, 0, 0, 682, 683, 5,
		97, 0, 0, 683, 684, 5, 145, 0, 0, 684, 157, 1, 0, 0, 0, 685, 686, 5, 100,
		0, 0, 686, 687, 7, 5, 0, 0, 687, 159, 1, 0, 0, 0, 688, 689, 5, 101, 0,
		0, 689, 690, 5, 114, 0, 0, 690, 161, 1, 0, 0, 0, 691, 692, 5, 102, 0, 0,
		692, 693, 5, 117, 0, 0, 693, 163, 1, 0, 0, 0, 68, 247, 254, 257, 263, 270,
		277, 284, 291, 298, 305, 312, 319, 326, 331, 337, 344, 351, 358, 365, 372,
		379, 386, 393, 401, 404, 411, 414, 421, 424, 431, 434, 441, 444, 451, 454,
		461, 464, 471, 474, 481, 484, 491, 494, 501, 504, 511, 514, 521, 524, 531,
		534, 541, 544, 551, 554, 603, 605, 610, 612, 619, 625, 643, 661, 663, 668,
		670, 675, 677,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// CEEventParserInit initializes any static state used to implement CEEventParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewCEEventParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func CEEventParserInit() {
	staticData := &ceeventparserParserStaticData
	staticData.once.Do(ceeventparserParserInit)
}

// NewCEEventParser produces a new parser instance for the optional input antlr.TokenStream.
func NewCEEventParser(input antlr.TokenStream) *CEEventParser {
	CEEventParserInit()
	this := new(CEEventParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &ceeventparserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "CEEventParser.g4"

	return this
}

// CEEventParser tokens.
const (
	CEEventParserEOF                   = antlr.TokenEOF
	CEEventParserEVENT_AB              = 1
	CEEventParserEVENT_AB_ARGS         = 2
	CEEventParserEVENT_ACL             = 3
	CEEventParserEVENT_ACM             = 4
	CEEventParserEVENT_ADB             = 5
	CEEventParserEVENT_ADF16           = 6
	CEEventParserEVENT_ADF32           = 7
	CEEventParserEVENT_ADF64           = 8
	CEEventParserEVENT_ADI16           = 9
	CEEventParserEVENT_ADI32           = 10
	CEEventParserEVENT_ADI64           = 11
	CEEventParserEVENT_ADI8            = 12
	CEEventParserEVENT_ADT             = 13
	CEEventParserEVENT_ADU16           = 14
	CEEventParserEVENT_ADU16X          = 15
	CEEventParserEVENT_ADU32           = 16
	CEEventParserEVENT_ADU32X          = 17
	CEEventParserEVENT_ADU64           = 18
	CEEventParserEVENT_ADU64X          = 19
	CEEventParserEVENT_ADU8            = 20
	CEEventParserEVENT_ADU8X           = 21
	CEEventParserEVENT_ADU             = 22
	CEEventParserEVENT_AF16            = 23
	CEEventParserEVENT_AF16_ARGS       = 24
	CEEventParserEVENT_AF32            = 25
	CEEventParserEVENT_AF32_ARGS       = 26
	CEEventParserEVENT_AF64            = 27
	CEEventParserEVENT_AF64_ARGS       = 28
	CEEventParserEVENT_AI16            = 29
	CEEventParserEVENT_AI16_ARGS       = 30
	CEEventParserEVENT_AI32            = 31
	CEEventParserEVENT_AI32_ARGS       = 32
	CEEventParserEVENT_AI64            = 33
	CEEventParserEVENT_AI64_ARGS       = 34
	CEEventParserEVENT_AI8             = 35
	CEEventParserEVENT_AI8_ARGS        = 36
	CEEventParserEVENT_AU16            = 37
	CEEventParserEVENT_AU16_ARGS       = 38
	CEEventParserEVENT_AU16X           = 39
	CEEventParserEVENT_AU16X_ARGS      = 40
	CEEventParserEVENT_AU32            = 41
	CEEventParserEVENT_AU32_ARGS       = 42
	CEEventParserEVENT_AU32X           = 43
	CEEventParserEVENT_AU32X_ARGS      = 44
	CEEventParserEVENT_AU64            = 45
	CEEventParserEVENT_AU64_ARGS       = 46
	CEEventParserEVENT_AU64X           = 47
	CEEventParserEVENT_AU64X_ARGS      = 48
	CEEventParserEVENT_AU8             = 49
	CEEventParserEVENT_AU8_ARGS        = 50
	CEEventParserEVENT_AU8X            = 51
	CEEventParserEVENT_AU8X_ARGS       = 52
	CEEventParserEVENT_AU              = 53
	CEEventParserEVENT_AU_ARGS         = 54
	CEEventParserEVENT_B               = 55
	CEEventParserEVENT_BAB             = 56
	CEEventParserEVENT_BAF16           = 57
	CEEventParserEVENT_BAF32           = 58
	CEEventParserEVENT_BAF64           = 59
	CEEventParserEVENT_BAI16           = 60
	CEEventParserEVENT_BAI32           = 61
	CEEventParserEVENT_BAI64           = 62
	CEEventParserEVENT_BAI8            = 63
	CEEventParserEVENT_BAU16           = 64
	CEEventParserEVENT_BAU32           = 65
	CEEventParserEVENT_BAU64           = 66
	CEEventParserEVENT_BAU8            = 67
	CEEventParserEVENT_BAU             = 68
	CEEventParserEVENT_BCB             = 69
	CEEventParserEVENT_BCT             = 70
	CEEventParserEVENT_BMEDIA          = 71
	CEEventParserEVENT_BREFR           = 72
	CEEventParserEVENT_BRID            = 73
	CEEventParserEVENT_BS              = 74
	CEEventParserEVENT_CB              = 75
	CEEventParserEVENT_CM              = 76
	CEEventParserEVENT_CM_ARGS         = 77
	CEEventParserEVENT_CS              = 78
	CEEventParserEVENT_CS_ARGS         = 79
	CEEventParserEVENT_CT              = 80
	CEEventParserEVENT_E               = 81
	CEEventParserEVENT_EDGE            = 82
	CEEventParserEVENT_L               = 83
	CEEventParserEVENT_M               = 84
	CEEventParserEVENT_MARK            = 85
	CEEventParserEVENT_MEDIA           = 86
	CEEventParserEVENT_N               = 87
	CEEventParserEVENT_NODE            = 88
	CEEventParserEVENT_NULL            = 89
	CEEventParserEVENT_PAD             = 90
	CEEventParserEVENT_REFL            = 91
	CEEventParserEVENT_REFR            = 92
	CEEventParserEVENT_REFR_ARGS       = 93
	CEEventParserEVENT_RID             = 94
	CEEventParserEVENT_RID_ARGS        = 95
	CEEventParserEVENT_SI              = 96
	CEEventParserEVENT_ST              = 97
	CEEventParserEVENT_S               = 98
	CEEventParserEVENT_S_ARGS          = 99
	CEEventParserEVENT_T               = 100
	CEEventParserEVENT_UID             = 101
	CEEventParserEVENT_V               = 102
	CEEventParserTRUE                  = 103
	CEEventParserFALSE                 = 104
	CEEventParserFLOAT_NAN             = 105
	CEEventParserFLOAT_SNAN            = 106
	CEEventParserFLOAT_INF             = 107
	CEEventParserFLOAT_DEC             = 108
	CEEventParserFLOAT_HEX             = 109
	CEEventParserINT_BIN               = 110
	CEEventParserINT_OCT               = 111
	CEEventParserINT_DEC               = 112
	CEEventParserINT_HEX               = 113
	CEEventParserUID                   = 114
	CEEventParserVALUE_UINT_BIN        = 115
	CEEventParserVALUE_UINT_OCT        = 116
	CEEventParserVALUE_UINT_DEC        = 117
	CEEventParserVALUE_UINT_HEX        = 118
	CEEventParserMODE_UINT_WS          = 119
	CEEventParserVALUE_UINTX           = 120
	CEEventParserMODE_UINTX_WS         = 121
	CEEventParserVALUE_INT_BIN         = 122
	CEEventParserVALUE_INT_OCT         = 123
	CEEventParserVALUE_INT_DEC         = 124
	CEEventParserVALUE_INT_HEX         = 125
	CEEventParserMODE_INT_WS           = 126
	CEEventParserVALUE_FLOAT_NAN       = 127
	CEEventParserVALUE_FLOAT_SNAN      = 128
	CEEventParserVALUE_FLOAT_INF       = 129
	CEEventParserVALUE_FLOAT_DEC       = 130
	CEEventParserVALUE_FLOAT_HEX       = 131
	CEEventParserMODE_FLOAT_WS         = 132
	CEEventParserVALUE_UID             = 133
	CEEventParserMODE_UID_WS           = 134
	CEEventParserTZ_PINT               = 135
	CEEventParserTZ_NINT               = 136
	CEEventParserTZ_INT                = 137
	CEEventParserTZ_COORD              = 138
	CEEventParserTZ_STRING             = 139
	CEEventParserTIME_ZONE             = 140
	CEEventParserTIME                  = 141
	CEEventParserDATE                  = 142
	CEEventParserDATETIME              = 143
	CEEventParserMODE_TIME_WS          = 144
	CEEventParserSTRING                = 145
	CEEventParserMODE_BYTES_WS         = 146
	CEEventParserBYTE                  = 147
	CEEventParserVALUE_BIT             = 148
	CEEventParserMODE_BITS_WS          = 149
	CEEventParserCUSTOM_BINARY_TYPE    = 150
	CEEventParserCUSTOM_TEXT_TYPE      = 151
	CEEventParserCUSTOM_TEXT_SEPARATOR = 152
	CEEventParserMEDIA_TYPE            = 153
)

// CEEventParser rules.
const (
	CEEventParserRULE_start                     = 0
	CEEventParserRULE_event                     = 1
	CEEventParserRULE_eventArrayBits            = 2
	CEEventParserRULE_eventArrayChunkLast       = 3
	CEEventParserRULE_eventArrayChunkMore       = 4
	CEEventParserRULE_eventArrayDataBits        = 5
	CEEventParserRULE_eventArrayDataFloat16     = 6
	CEEventParserRULE_eventArrayDataFloat32     = 7
	CEEventParserRULE_eventArrayDataFloat64     = 8
	CEEventParserRULE_eventArrayDataInt16       = 9
	CEEventParserRULE_eventArrayDataInt32       = 10
	CEEventParserRULE_eventArrayDataInt64       = 11
	CEEventParserRULE_eventArrayDataInt8        = 12
	CEEventParserRULE_eventArrayDataText        = 13
	CEEventParserRULE_eventArrayDataUID         = 14
	CEEventParserRULE_eventArrayDataUint16      = 15
	CEEventParserRULE_eventArrayDataUint16X     = 16
	CEEventParserRULE_eventArrayDataUint32      = 17
	CEEventParserRULE_eventArrayDataUint32X     = 18
	CEEventParserRULE_eventArrayDataUint64      = 19
	CEEventParserRULE_eventArrayDataUint64X     = 20
	CEEventParserRULE_eventArrayDataUint8       = 21
	CEEventParserRULE_eventArrayDataUint8X      = 22
	CEEventParserRULE_eventArrayFloat16         = 23
	CEEventParserRULE_eventArrayFloat32         = 24
	CEEventParserRULE_eventArrayFloat64         = 25
	CEEventParserRULE_eventArrayInt16           = 26
	CEEventParserRULE_eventArrayInt32           = 27
	CEEventParserRULE_eventArrayInt64           = 28
	CEEventParserRULE_eventArrayInt8            = 29
	CEEventParserRULE_eventArrayUID             = 30
	CEEventParserRULE_eventArrayUint16          = 31
	CEEventParserRULE_eventArrayUint16X         = 32
	CEEventParserRULE_eventArrayUint32          = 33
	CEEventParserRULE_eventArrayUint32X         = 34
	CEEventParserRULE_eventArrayUint64          = 35
	CEEventParserRULE_eventArrayUint64X         = 36
	CEEventParserRULE_eventArrayUint8           = 37
	CEEventParserRULE_eventArrayUint8X          = 38
	CEEventParserRULE_eventBeginArrayBits       = 39
	CEEventParserRULE_eventBeginArrayFloat16    = 40
	CEEventParserRULE_eventBeginArrayFloat32    = 41
	CEEventParserRULE_eventBeginArrayFloat64    = 42
	CEEventParserRULE_eventBeginArrayInt16      = 43
	CEEventParserRULE_eventBeginArrayInt32      = 44
	CEEventParserRULE_eventBeginArrayInt64      = 45
	CEEventParserRULE_eventBeginArrayInt8       = 46
	CEEventParserRULE_eventBeginArrayUID        = 47
	CEEventParserRULE_eventBeginArrayUint16     = 48
	CEEventParserRULE_eventBeginArrayUint32     = 49
	CEEventParserRULE_eventBeginArrayUint64     = 50
	CEEventParserRULE_eventBeginArrayUint8      = 51
	CEEventParserRULE_eventBeginCustomBinary    = 52
	CEEventParserRULE_eventBeginCustomText      = 53
	CEEventParserRULE_eventBeginMedia           = 54
	CEEventParserRULE_eventBeginRemoteReference = 55
	CEEventParserRULE_eventBeginResourceId      = 56
	CEEventParserRULE_eventBeginString          = 57
	CEEventParserRULE_eventBoolean              = 58
	CEEventParserRULE_eventCommentMultiline     = 59
	CEEventParserRULE_eventCommentSingleLine    = 60
	CEEventParserRULE_eventCustomBinary         = 61
	CEEventParserRULE_eventCustomText           = 62
	CEEventParserRULE_eventEdge                 = 63
	CEEventParserRULE_eventEndContainer         = 64
	CEEventParserRULE_eventList                 = 65
	CEEventParserRULE_eventMap                  = 66
	CEEventParserRULE_eventMarker               = 67
	CEEventParserRULE_eventMedia                = 68
	CEEventParserRULE_eventNode                 = 69
	CEEventParserRULE_eventNull                 = 70
	CEEventParserRULE_eventNumber               = 71
	CEEventParserRULE_eventPad                  = 72
	CEEventParserRULE_eventLocalReference       = 73
	CEEventParserRULE_eventRemoteReference      = 74
	CEEventParserRULE_eventResourceId           = 75
	CEEventParserRULE_eventString               = 76
	CEEventParserRULE_eventStructInstance       = 77
	CEEventParserRULE_eventStructTemplate       = 78
	CEEventParserRULE_eventTime                 = 79
	CEEventParserRULE_eventUID                  = 80
	CEEventParserRULE_eventVersion              = 81
)

// IStartContext is an interface to support dynamic dispatch.
type IStartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStartContext differentiates from other interfaces.
	IsStartContext()
}

type StartContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStartContext() *StartContext {
	var p = new(StartContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_start
	return p
}

func (*StartContext) IsStartContext() {}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	var p = new(StartContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_start

	return p
}

func (s *StartContext) GetParser() antlr.Parser { return s.parser }

func (s *StartContext) Event() IEventContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventContext)
}

func (s *StartContext) EOF() antlr.TerminalNode {
	return s.GetToken(CEEventParserEOF, 0)
}

func (s *StartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterStart(s)
	}
}

func (s *StartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitStart(s)
	}
}

func (p *CEEventParser) Start() (localctx IStartContext) {
	this := p
	_ = this

	localctx = NewStartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, CEEventParserRULE_start)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(164)
		p.Event()
	}
	{
		p.SetState(165)
		p.Match(CEEventParserEOF)
	}

	return localctx
}

// IEventContext is an interface to support dynamic dispatch.
type IEventContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventContext differentiates from other interfaces.
	IsEventContext()
}

type EventContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventContext() *EventContext {
	var p = new(EventContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_event
	return p
}

func (*EventContext) IsEventContext() {}

func NewEventContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventContext {
	var p = new(EventContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_event

	return p
}

func (s *EventContext) GetParser() antlr.Parser { return s.parser }

func (s *EventContext) EventArrayBits() IEventArrayBitsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayBitsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayBitsContext)
}

func (s *EventContext) EventArrayChunkLast() IEventArrayChunkLastContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayChunkLastContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayChunkLastContext)
}

func (s *EventContext) EventArrayChunkMore() IEventArrayChunkMoreContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayChunkMoreContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayChunkMoreContext)
}

func (s *EventContext) EventArrayDataBits() IEventArrayDataBitsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataBitsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataBitsContext)
}

func (s *EventContext) EventArrayDataFloat16() IEventArrayDataFloat16Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataFloat16Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataFloat16Context)
}

func (s *EventContext) EventArrayDataFloat32() IEventArrayDataFloat32Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataFloat32Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataFloat32Context)
}

func (s *EventContext) EventArrayDataFloat64() IEventArrayDataFloat64Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataFloat64Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataFloat64Context)
}

func (s *EventContext) EventArrayDataInt16() IEventArrayDataInt16Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataInt16Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataInt16Context)
}

func (s *EventContext) EventArrayDataInt32() IEventArrayDataInt32Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataInt32Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataInt32Context)
}

func (s *EventContext) EventArrayDataInt64() IEventArrayDataInt64Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataInt64Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataInt64Context)
}

func (s *EventContext) EventArrayDataInt8() IEventArrayDataInt8Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataInt8Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataInt8Context)
}

func (s *EventContext) EventArrayDataText() IEventArrayDataTextContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataTextContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataTextContext)
}

func (s *EventContext) EventArrayDataUID() IEventArrayDataUIDContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataUIDContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataUIDContext)
}

func (s *EventContext) EventArrayDataUint16() IEventArrayDataUint16Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataUint16Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataUint16Context)
}

func (s *EventContext) EventArrayDataUint16X() IEventArrayDataUint16XContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataUint16XContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataUint16XContext)
}

func (s *EventContext) EventArrayDataUint32() IEventArrayDataUint32Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataUint32Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataUint32Context)
}

func (s *EventContext) EventArrayDataUint32X() IEventArrayDataUint32XContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataUint32XContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataUint32XContext)
}

func (s *EventContext) EventArrayDataUint64() IEventArrayDataUint64Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataUint64Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataUint64Context)
}

func (s *EventContext) EventArrayDataUint64X() IEventArrayDataUint64XContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataUint64XContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataUint64XContext)
}

func (s *EventContext) EventArrayDataUint8() IEventArrayDataUint8Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataUint8Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataUint8Context)
}

func (s *EventContext) EventArrayDataUint8X() IEventArrayDataUint8XContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayDataUint8XContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayDataUint8XContext)
}

func (s *EventContext) EventArrayFloat16() IEventArrayFloat16Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayFloat16Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayFloat16Context)
}

func (s *EventContext) EventArrayFloat32() IEventArrayFloat32Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayFloat32Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayFloat32Context)
}

func (s *EventContext) EventArrayFloat64() IEventArrayFloat64Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayFloat64Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayFloat64Context)
}

func (s *EventContext) EventArrayInt16() IEventArrayInt16Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayInt16Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayInt16Context)
}

func (s *EventContext) EventArrayInt32() IEventArrayInt32Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayInt32Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayInt32Context)
}

func (s *EventContext) EventArrayInt64() IEventArrayInt64Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayInt64Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayInt64Context)
}

func (s *EventContext) EventArrayInt8() IEventArrayInt8Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayInt8Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayInt8Context)
}

func (s *EventContext) EventArrayUID() IEventArrayUIDContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayUIDContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayUIDContext)
}

func (s *EventContext) EventArrayUint16() IEventArrayUint16Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayUint16Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayUint16Context)
}

func (s *EventContext) EventArrayUint16X() IEventArrayUint16XContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayUint16XContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayUint16XContext)
}

func (s *EventContext) EventArrayUint32() IEventArrayUint32Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayUint32Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayUint32Context)
}

func (s *EventContext) EventArrayUint32X() IEventArrayUint32XContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayUint32XContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayUint32XContext)
}

func (s *EventContext) EventArrayUint64() IEventArrayUint64Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayUint64Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayUint64Context)
}

func (s *EventContext) EventArrayUint64X() IEventArrayUint64XContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayUint64XContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayUint64XContext)
}

func (s *EventContext) EventArrayUint8() IEventArrayUint8Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayUint8Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayUint8Context)
}

func (s *EventContext) EventArrayUint8X() IEventArrayUint8XContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventArrayUint8XContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventArrayUint8XContext)
}

func (s *EventContext) EventBeginArrayBits() IEventBeginArrayBitsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayBitsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayBitsContext)
}

func (s *EventContext) EventBeginArrayFloat16() IEventBeginArrayFloat16Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayFloat16Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayFloat16Context)
}

func (s *EventContext) EventBeginArrayFloat32() IEventBeginArrayFloat32Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayFloat32Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayFloat32Context)
}

func (s *EventContext) EventBeginArrayFloat64() IEventBeginArrayFloat64Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayFloat64Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayFloat64Context)
}

func (s *EventContext) EventBeginArrayInt16() IEventBeginArrayInt16Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayInt16Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayInt16Context)
}

func (s *EventContext) EventBeginArrayInt32() IEventBeginArrayInt32Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayInt32Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayInt32Context)
}

func (s *EventContext) EventBeginArrayInt64() IEventBeginArrayInt64Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayInt64Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayInt64Context)
}

func (s *EventContext) EventBeginArrayInt8() IEventBeginArrayInt8Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayInt8Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayInt8Context)
}

func (s *EventContext) EventBeginArrayUID() IEventBeginArrayUIDContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayUIDContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayUIDContext)
}

func (s *EventContext) EventBeginArrayUint16() IEventBeginArrayUint16Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayUint16Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayUint16Context)
}

func (s *EventContext) EventBeginArrayUint32() IEventBeginArrayUint32Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayUint32Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayUint32Context)
}

func (s *EventContext) EventBeginArrayUint64() IEventBeginArrayUint64Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayUint64Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayUint64Context)
}

func (s *EventContext) EventBeginArrayUint8() IEventBeginArrayUint8Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginArrayUint8Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginArrayUint8Context)
}

func (s *EventContext) EventBeginCustomBinary() IEventBeginCustomBinaryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginCustomBinaryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginCustomBinaryContext)
}

func (s *EventContext) EventBeginCustomText() IEventBeginCustomTextContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginCustomTextContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginCustomTextContext)
}

func (s *EventContext) EventBeginMedia() IEventBeginMediaContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginMediaContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginMediaContext)
}

func (s *EventContext) EventBeginResourceId() IEventBeginResourceIdContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginResourceIdContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginResourceIdContext)
}

func (s *EventContext) EventBeginRemoteReference() IEventBeginRemoteReferenceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginRemoteReferenceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginRemoteReferenceContext)
}

func (s *EventContext) EventBeginString() IEventBeginStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBeginStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBeginStringContext)
}

func (s *EventContext) EventBoolean() IEventBooleanContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventBooleanContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventBooleanContext)
}

func (s *EventContext) EventCommentSingleLine() IEventCommentSingleLineContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventCommentSingleLineContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventCommentSingleLineContext)
}

func (s *EventContext) EventCommentMultiline() IEventCommentMultilineContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventCommentMultilineContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventCommentMultilineContext)
}

func (s *EventContext) EventCustomBinary() IEventCustomBinaryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventCustomBinaryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventCustomBinaryContext)
}

func (s *EventContext) EventCustomText() IEventCustomTextContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventCustomTextContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventCustomTextContext)
}

func (s *EventContext) EventEdge() IEventEdgeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventEdgeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventEdgeContext)
}

func (s *EventContext) EventEndContainer() IEventEndContainerContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventEndContainerContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventEndContainerContext)
}

func (s *EventContext) EventList() IEventListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventListContext)
}

func (s *EventContext) EventMap() IEventMapContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventMapContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventMapContext)
}

func (s *EventContext) EventMarker() IEventMarkerContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventMarkerContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventMarkerContext)
}

func (s *EventContext) EventMedia() IEventMediaContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventMediaContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventMediaContext)
}

func (s *EventContext) EventNode() IEventNodeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventNodeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventNodeContext)
}

func (s *EventContext) EventNull() IEventNullContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventNullContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventNullContext)
}

func (s *EventContext) EventNumber() IEventNumberContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventNumberContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventNumberContext)
}

func (s *EventContext) EventPad() IEventPadContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventPadContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventPadContext)
}

func (s *EventContext) EventLocalReference() IEventLocalReferenceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventLocalReferenceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventLocalReferenceContext)
}

func (s *EventContext) EventRemoteReference() IEventRemoteReferenceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventRemoteReferenceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventRemoteReferenceContext)
}

func (s *EventContext) EventResourceId() IEventResourceIdContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventResourceIdContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventResourceIdContext)
}

func (s *EventContext) EventString() IEventStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventStringContext)
}

func (s *EventContext) EventStructInstance() IEventStructInstanceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventStructInstanceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventStructInstanceContext)
}

func (s *EventContext) EventStructTemplate() IEventStructTemplateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventStructTemplateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventStructTemplateContext)
}

func (s *EventContext) EventTime() IEventTimeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventTimeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventTimeContext)
}

func (s *EventContext) EventUID() IEventUIDContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventUIDContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventUIDContext)
}

func (s *EventContext) EventVersion() IEventVersionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventVersionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventVersionContext)
}

func (s *EventContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEvent(s)
	}
}

func (s *EventContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEvent(s)
	}
}

func (p *CEEventParser) Event() (localctx IEventContext) {
	this := p
	_ = this

	localctx = NewEventContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, CEEventParserRULE_event)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(247)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AB, CEEventParserEVENT_AB_ARGS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(167)
			p.EventArrayBits()
		}

	case CEEventParserEVENT_ACL:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(168)
			p.EventArrayChunkLast()
		}

	case CEEventParserEVENT_ACM:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(169)
			p.EventArrayChunkMore()
		}

	case CEEventParserEVENT_ADB:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(170)
			p.EventArrayDataBits()
		}

	case CEEventParserEVENT_ADF16:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(171)
			p.EventArrayDataFloat16()
		}

	case CEEventParserEVENT_ADF32:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(172)
			p.EventArrayDataFloat32()
		}

	case CEEventParserEVENT_ADF64:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(173)
			p.EventArrayDataFloat64()
		}

	case CEEventParserEVENT_ADI16:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(174)
			p.EventArrayDataInt16()
		}

	case CEEventParserEVENT_ADI32:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(175)
			p.EventArrayDataInt32()
		}

	case CEEventParserEVENT_ADI64:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(176)
			p.EventArrayDataInt64()
		}

	case CEEventParserEVENT_ADI8:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(177)
			p.EventArrayDataInt8()
		}

	case CEEventParserEVENT_ADT:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(178)
			p.EventArrayDataText()
		}

	case CEEventParserEVENT_ADU:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(179)
			p.EventArrayDataUID()
		}

	case CEEventParserEVENT_ADU16:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(180)
			p.EventArrayDataUint16()
		}

	case CEEventParserEVENT_ADU16X:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(181)
			p.EventArrayDataUint16X()
		}

	case CEEventParserEVENT_ADU32:
		p.EnterOuterAlt(localctx, 16)
		{
			p.SetState(182)
			p.EventArrayDataUint32()
		}

	case CEEventParserEVENT_ADU32X:
		p.EnterOuterAlt(localctx, 17)
		{
			p.SetState(183)
			p.EventArrayDataUint32X()
		}

	case CEEventParserEVENT_ADU64:
		p.EnterOuterAlt(localctx, 18)
		{
			p.SetState(184)
			p.EventArrayDataUint64()
		}

	case CEEventParserEVENT_ADU64X:
		p.EnterOuterAlt(localctx, 19)
		{
			p.SetState(185)
			p.EventArrayDataUint64X()
		}

	case CEEventParserEVENT_ADU8:
		p.EnterOuterAlt(localctx, 20)
		{
			p.SetState(186)
			p.EventArrayDataUint8()
		}

	case CEEventParserEVENT_ADU8X:
		p.EnterOuterAlt(localctx, 21)
		{
			p.SetState(187)
			p.EventArrayDataUint8X()
		}

	case CEEventParserEVENT_AF16, CEEventParserEVENT_AF16_ARGS:
		p.EnterOuterAlt(localctx, 22)
		{
			p.SetState(188)
			p.EventArrayFloat16()
		}

	case CEEventParserEVENT_AF32, CEEventParserEVENT_AF32_ARGS:
		p.EnterOuterAlt(localctx, 23)
		{
			p.SetState(189)
			p.EventArrayFloat32()
		}

	case CEEventParserEVENT_AF64, CEEventParserEVENT_AF64_ARGS:
		p.EnterOuterAlt(localctx, 24)
		{
			p.SetState(190)
			p.EventArrayFloat64()
		}

	case CEEventParserEVENT_AI16, CEEventParserEVENT_AI16_ARGS:
		p.EnterOuterAlt(localctx, 25)
		{
			p.SetState(191)
			p.EventArrayInt16()
		}

	case CEEventParserEVENT_AI32, CEEventParserEVENT_AI32_ARGS:
		p.EnterOuterAlt(localctx, 26)
		{
			p.SetState(192)
			p.EventArrayInt32()
		}

	case CEEventParserEVENT_AI64, CEEventParserEVENT_AI64_ARGS:
		p.EnterOuterAlt(localctx, 27)
		{
			p.SetState(193)
			p.EventArrayInt64()
		}

	case CEEventParserEVENT_AI8, CEEventParserEVENT_AI8_ARGS:
		p.EnterOuterAlt(localctx, 28)
		{
			p.SetState(194)
			p.EventArrayInt8()
		}

	case CEEventParserEVENT_AU, CEEventParserEVENT_AU_ARGS:
		p.EnterOuterAlt(localctx, 29)
		{
			p.SetState(195)
			p.EventArrayUID()
		}

	case CEEventParserEVENT_AU16, CEEventParserEVENT_AU16_ARGS:
		p.EnterOuterAlt(localctx, 30)
		{
			p.SetState(196)
			p.EventArrayUint16()
		}

	case CEEventParserEVENT_AU16X, CEEventParserEVENT_AU16X_ARGS:
		p.EnterOuterAlt(localctx, 31)
		{
			p.SetState(197)
			p.EventArrayUint16X()
		}

	case CEEventParserEVENT_AU32, CEEventParserEVENT_AU32_ARGS:
		p.EnterOuterAlt(localctx, 32)
		{
			p.SetState(198)
			p.EventArrayUint32()
		}

	case CEEventParserEVENT_AU32X, CEEventParserEVENT_AU32X_ARGS:
		p.EnterOuterAlt(localctx, 33)
		{
			p.SetState(199)
			p.EventArrayUint32X()
		}

	case CEEventParserEVENT_AU64, CEEventParserEVENT_AU64_ARGS:
		p.EnterOuterAlt(localctx, 34)
		{
			p.SetState(200)
			p.EventArrayUint64()
		}

	case CEEventParserEVENT_AU64X, CEEventParserEVENT_AU64X_ARGS:
		p.EnterOuterAlt(localctx, 35)
		{
			p.SetState(201)
			p.EventArrayUint64X()
		}

	case CEEventParserEVENT_AU8, CEEventParserEVENT_AU8_ARGS:
		p.EnterOuterAlt(localctx, 36)
		{
			p.SetState(202)
			p.EventArrayUint8()
		}

	case CEEventParserEVENT_AU8X, CEEventParserEVENT_AU8X_ARGS:
		p.EnterOuterAlt(localctx, 37)
		{
			p.SetState(203)
			p.EventArrayUint8X()
		}

	case CEEventParserEVENT_BAB:
		p.EnterOuterAlt(localctx, 38)
		{
			p.SetState(204)
			p.EventBeginArrayBits()
		}

	case CEEventParserEVENT_BAF16:
		p.EnterOuterAlt(localctx, 39)
		{
			p.SetState(205)
			p.EventBeginArrayFloat16()
		}

	case CEEventParserEVENT_BAF32:
		p.EnterOuterAlt(localctx, 40)
		{
			p.SetState(206)
			p.EventBeginArrayFloat32()
		}

	case CEEventParserEVENT_BAF64:
		p.EnterOuterAlt(localctx, 41)
		{
			p.SetState(207)
			p.EventBeginArrayFloat64()
		}

	case CEEventParserEVENT_BAI16:
		p.EnterOuterAlt(localctx, 42)
		{
			p.SetState(208)
			p.EventBeginArrayInt16()
		}

	case CEEventParserEVENT_BAI32:
		p.EnterOuterAlt(localctx, 43)
		{
			p.SetState(209)
			p.EventBeginArrayInt32()
		}

	case CEEventParserEVENT_BAI64:
		p.EnterOuterAlt(localctx, 44)
		{
			p.SetState(210)
			p.EventBeginArrayInt64()
		}

	case CEEventParserEVENT_BAI8:
		p.EnterOuterAlt(localctx, 45)
		{
			p.SetState(211)
			p.EventBeginArrayInt8()
		}

	case CEEventParserEVENT_BAU:
		p.EnterOuterAlt(localctx, 46)
		{
			p.SetState(212)
			p.EventBeginArrayUID()
		}

	case CEEventParserEVENT_BAU16:
		p.EnterOuterAlt(localctx, 47)
		{
			p.SetState(213)
			p.EventBeginArrayUint16()
		}

	case CEEventParserEVENT_BAU32:
		p.EnterOuterAlt(localctx, 48)
		{
			p.SetState(214)
			p.EventBeginArrayUint32()
		}

	case CEEventParserEVENT_BAU64:
		p.EnterOuterAlt(localctx, 49)
		{
			p.SetState(215)
			p.EventBeginArrayUint64()
		}

	case CEEventParserEVENT_BAU8:
		p.EnterOuterAlt(localctx, 50)
		{
			p.SetState(216)
			p.EventBeginArrayUint8()
		}

	case CEEventParserEVENT_BCB:
		p.EnterOuterAlt(localctx, 51)
		{
			p.SetState(217)
			p.EventBeginCustomBinary()
		}

	case CEEventParserEVENT_BCT:
		p.EnterOuterAlt(localctx, 52)
		{
			p.SetState(218)
			p.EventBeginCustomText()
		}

	case CEEventParserEVENT_BMEDIA:
		p.EnterOuterAlt(localctx, 53)
		{
			p.SetState(219)
			p.EventBeginMedia()
		}

	case CEEventParserEVENT_BRID:
		p.EnterOuterAlt(localctx, 54)
		{
			p.SetState(220)
			p.EventBeginResourceId()
		}

	case CEEventParserEVENT_BREFR:
		p.EnterOuterAlt(localctx, 55)
		{
			p.SetState(221)
			p.EventBeginRemoteReference()
		}

	case CEEventParserEVENT_BS:
		p.EnterOuterAlt(localctx, 56)
		{
			p.SetState(222)
			p.EventBeginString()
		}

	case CEEventParserEVENT_B:
		p.EnterOuterAlt(localctx, 57)
		{
			p.SetState(223)
			p.EventBoolean()
		}

	case CEEventParserEVENT_CS, CEEventParserEVENT_CS_ARGS:
		p.EnterOuterAlt(localctx, 58)
		{
			p.SetState(224)
			p.EventCommentSingleLine()
		}

	case CEEventParserEVENT_CM, CEEventParserEVENT_CM_ARGS:
		p.EnterOuterAlt(localctx, 59)
		{
			p.SetState(225)
			p.EventCommentMultiline()
		}

	case CEEventParserEVENT_CB:
		p.EnterOuterAlt(localctx, 60)
		{
			p.SetState(226)
			p.EventCustomBinary()
		}

	case CEEventParserEVENT_CT:
		p.EnterOuterAlt(localctx, 61)
		{
			p.SetState(227)
			p.EventCustomText()
		}

	case CEEventParserEVENT_EDGE:
		p.EnterOuterAlt(localctx, 62)
		{
			p.SetState(228)
			p.EventEdge()
		}

	case CEEventParserEVENT_E:
		p.EnterOuterAlt(localctx, 63)
		{
			p.SetState(229)
			p.EventEndContainer()
		}

	case CEEventParserEVENT_L:
		p.EnterOuterAlt(localctx, 64)
		{
			p.SetState(230)
			p.EventList()
		}

	case CEEventParserEVENT_M:
		p.EnterOuterAlt(localctx, 65)
		{
			p.SetState(231)
			p.EventMap()
		}

	case CEEventParserEVENT_MARK:
		p.EnterOuterAlt(localctx, 66)
		{
			p.SetState(232)
			p.EventMarker()
		}

	case CEEventParserEVENT_MEDIA:
		p.EnterOuterAlt(localctx, 67)
		{
			p.SetState(233)
			p.EventMedia()
		}

	case CEEventParserEVENT_NODE:
		p.EnterOuterAlt(localctx, 68)
		{
			p.SetState(234)
			p.EventNode()
		}

	case CEEventParserEVENT_NULL:
		p.EnterOuterAlt(localctx, 69)
		{
			p.SetState(235)
			p.EventNull()
		}

	case CEEventParserEVENT_N:
		p.EnterOuterAlt(localctx, 70)
		{
			p.SetState(236)
			p.EventNumber()
		}

	case CEEventParserEVENT_PAD:
		p.EnterOuterAlt(localctx, 71)
		{
			p.SetState(237)
			p.EventPad()
		}

	case CEEventParserEVENT_REFL:
		p.EnterOuterAlt(localctx, 72)
		{
			p.SetState(238)
			p.EventLocalReference()
		}

	case CEEventParserEVENT_REFR, CEEventParserEVENT_REFR_ARGS:
		p.EnterOuterAlt(localctx, 73)
		{
			p.SetState(239)
			p.EventRemoteReference()
		}

	case CEEventParserEVENT_RID, CEEventParserEVENT_RID_ARGS:
		p.EnterOuterAlt(localctx, 74)
		{
			p.SetState(240)
			p.EventResourceId()
		}

	case CEEventParserEVENT_S, CEEventParserEVENT_S_ARGS:
		p.EnterOuterAlt(localctx, 75)
		{
			p.SetState(241)
			p.EventString()
		}

	case CEEventParserEVENT_SI:
		p.EnterOuterAlt(localctx, 76)
		{
			p.SetState(242)
			p.EventStructInstance()
		}

	case CEEventParserEVENT_ST:
		p.EnterOuterAlt(localctx, 77)
		{
			p.SetState(243)
			p.EventStructTemplate()
		}

	case CEEventParserEVENT_T:
		p.EnterOuterAlt(localctx, 78)
		{
			p.SetState(244)
			p.EventTime()
		}

	case CEEventParserEVENT_UID:
		p.EnterOuterAlt(localctx, 79)
		{
			p.SetState(245)
			p.EventUID()
		}

	case CEEventParserEVENT_V:
		p.EnterOuterAlt(localctx, 80)
		{
			p.SetState(246)
			p.EventVersion()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayBitsContext is an interface to support dynamic dispatch.
type IEventArrayBitsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayBitsContext differentiates from other interfaces.
	IsEventArrayBitsContext()
}

type EventArrayBitsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayBitsContext() *EventArrayBitsContext {
	var p = new(EventArrayBitsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayBits
	return p
}

func (*EventArrayBitsContext) IsEventArrayBitsContext() {}

func NewEventArrayBitsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayBitsContext {
	var p = new(EventArrayBitsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayBits

	return p
}

func (s *EventArrayBitsContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayBitsContext) EVENT_AB() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AB, 0)
}

func (s *EventArrayBitsContext) EVENT_AB_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AB_ARGS, 0)
}

func (s *EventArrayBitsContext) AllVALUE_BIT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_BIT)
}

func (s *EventArrayBitsContext) VALUE_BIT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_BIT, i)
}

func (s *EventArrayBitsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayBitsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayBitsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayBits(s)
	}
}

func (s *EventArrayBitsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayBits(s)
	}
}

func (p *CEEventParser) EventArrayBits() (localctx IEventArrayBitsContext) {
	this := p
	_ = this

	localctx = NewEventArrayBitsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, CEEventParserRULE_eventArrayBits)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(257)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AB:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(249)
			p.Match(CEEventParserEVENT_AB)
		}

	case CEEventParserEVENT_AB_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(250)
			p.Match(CEEventParserEVENT_AB_ARGS)
		}
		p.SetState(254)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == CEEventParserVALUE_BIT {
			{
				p.SetState(251)
				p.Match(CEEventParserVALUE_BIT)
			}

			p.SetState(256)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayChunkLastContext is an interface to support dynamic dispatch.
type IEventArrayChunkLastContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayChunkLastContext differentiates from other interfaces.
	IsEventArrayChunkLastContext()
}

type EventArrayChunkLastContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayChunkLastContext() *EventArrayChunkLastContext {
	var p = new(EventArrayChunkLastContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayChunkLast
	return p
}

func (*EventArrayChunkLastContext) IsEventArrayChunkLastContext() {}

func NewEventArrayChunkLastContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayChunkLastContext {
	var p = new(EventArrayChunkLastContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayChunkLast

	return p
}

func (s *EventArrayChunkLastContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayChunkLastContext) EVENT_ACL() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ACL, 0)
}

func (s *EventArrayChunkLastContext) AllVALUE_UINT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_DEC)
}

func (s *EventArrayChunkLastContext) VALUE_UINT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, i)
}

func (s *EventArrayChunkLastContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayChunkLastContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayChunkLastContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayChunkLast(s)
	}
}

func (s *EventArrayChunkLastContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayChunkLast(s)
	}
}

func (p *CEEventParser) EventArrayChunkLast() (localctx IEventArrayChunkLastContext) {
	this := p
	_ = this

	localctx = NewEventArrayChunkLastContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, CEEventParserRULE_eventArrayChunkLast)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(259)
		p.Match(CEEventParserEVENT_ACL)
	}
	p.SetState(263)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CEEventParserVALUE_UINT_DEC {
		{
			p.SetState(260)
			p.Match(CEEventParserVALUE_UINT_DEC)
		}

		p.SetState(265)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayChunkMoreContext is an interface to support dynamic dispatch.
type IEventArrayChunkMoreContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayChunkMoreContext differentiates from other interfaces.
	IsEventArrayChunkMoreContext()
}

type EventArrayChunkMoreContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayChunkMoreContext() *EventArrayChunkMoreContext {
	var p = new(EventArrayChunkMoreContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayChunkMore
	return p
}

func (*EventArrayChunkMoreContext) IsEventArrayChunkMoreContext() {}

func NewEventArrayChunkMoreContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayChunkMoreContext {
	var p = new(EventArrayChunkMoreContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayChunkMore

	return p
}

func (s *EventArrayChunkMoreContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayChunkMoreContext) EVENT_ACM() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ACM, 0)
}

func (s *EventArrayChunkMoreContext) AllVALUE_UINT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_DEC)
}

func (s *EventArrayChunkMoreContext) VALUE_UINT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, i)
}

func (s *EventArrayChunkMoreContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayChunkMoreContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayChunkMoreContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayChunkMore(s)
	}
}

func (s *EventArrayChunkMoreContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayChunkMore(s)
	}
}

func (p *CEEventParser) EventArrayChunkMore() (localctx IEventArrayChunkMoreContext) {
	this := p
	_ = this

	localctx = NewEventArrayChunkMoreContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, CEEventParserRULE_eventArrayChunkMore)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(266)
		p.Match(CEEventParserEVENT_ACM)
	}
	p.SetState(270)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CEEventParserVALUE_UINT_DEC {
		{
			p.SetState(267)
			p.Match(CEEventParserVALUE_UINT_DEC)
		}

		p.SetState(272)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataBitsContext is an interface to support dynamic dispatch.
type IEventArrayDataBitsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataBitsContext differentiates from other interfaces.
	IsEventArrayDataBitsContext()
}

type EventArrayDataBitsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataBitsContext() *EventArrayDataBitsContext {
	var p = new(EventArrayDataBitsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataBits
	return p
}

func (*EventArrayDataBitsContext) IsEventArrayDataBitsContext() {}

func NewEventArrayDataBitsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataBitsContext {
	var p = new(EventArrayDataBitsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataBits

	return p
}

func (s *EventArrayDataBitsContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataBitsContext) EVENT_ADB() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADB, 0)
}

func (s *EventArrayDataBitsContext) AllVALUE_BIT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_BIT)
}

func (s *EventArrayDataBitsContext) VALUE_BIT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_BIT, i)
}

func (s *EventArrayDataBitsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataBitsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataBitsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataBits(s)
	}
}

func (s *EventArrayDataBitsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataBits(s)
	}
}

func (p *CEEventParser) EventArrayDataBits() (localctx IEventArrayDataBitsContext) {
	this := p
	_ = this

	localctx = NewEventArrayDataBitsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, CEEventParserRULE_eventArrayDataBits)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(273)
		p.Match(CEEventParserEVENT_ADB)
	}
	p.SetState(277)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CEEventParserVALUE_BIT {
		{
			p.SetState(274)
			p.Match(CEEventParserVALUE_BIT)
		}

		p.SetState(279)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataFloat16Context is an interface to support dynamic dispatch.
type IEventArrayDataFloat16Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataFloat16Context differentiates from other interfaces.
	IsEventArrayDataFloat16Context()
}

type EventArrayDataFloat16Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataFloat16Context() *EventArrayDataFloat16Context {
	var p = new(EventArrayDataFloat16Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataFloat16
	return p
}

func (*EventArrayDataFloat16Context) IsEventArrayDataFloat16Context() {}

func NewEventArrayDataFloat16Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataFloat16Context {
	var p = new(EventArrayDataFloat16Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataFloat16

	return p
}

func (s *EventArrayDataFloat16Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataFloat16Context) EVENT_ADF16() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADF16, 0)
}

func (s *EventArrayDataFloat16Context) AllVALUE_FLOAT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_DEC)
}

func (s *EventArrayDataFloat16Context) VALUE_FLOAT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_DEC, i)
}

func (s *EventArrayDataFloat16Context) AllVALUE_FLOAT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_HEX)
}

func (s *EventArrayDataFloat16Context) VALUE_FLOAT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_HEX, i)
}

func (s *EventArrayDataFloat16Context) AllVALUE_FLOAT_INF() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_INF)
}

func (s *EventArrayDataFloat16Context) VALUE_FLOAT_INF(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_INF, i)
}

func (s *EventArrayDataFloat16Context) AllVALUE_FLOAT_NAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_NAN)
}

func (s *EventArrayDataFloat16Context) VALUE_FLOAT_NAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_NAN, i)
}

func (s *EventArrayDataFloat16Context) AllVALUE_FLOAT_SNAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_SNAN)
}

func (s *EventArrayDataFloat16Context) VALUE_FLOAT_SNAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_SNAN, i)
}

func (s *EventArrayDataFloat16Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataFloat16Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataFloat16Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataFloat16(s)
	}
}

func (s *EventArrayDataFloat16Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataFloat16(s)
	}
}

func (p *CEEventParser) EventArrayDataFloat16() (localctx IEventArrayDataFloat16Context) {
	this := p
	_ = this

	localctx = NewEventArrayDataFloat16Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, CEEventParserRULE_eventArrayDataFloat16)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(280)
		p.Match(CEEventParserEVENT_ADF16)
	}
	p.SetState(284)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0 {
		{
			p.SetState(281)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(286)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataFloat32Context is an interface to support dynamic dispatch.
type IEventArrayDataFloat32Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataFloat32Context differentiates from other interfaces.
	IsEventArrayDataFloat32Context()
}

type EventArrayDataFloat32Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataFloat32Context() *EventArrayDataFloat32Context {
	var p = new(EventArrayDataFloat32Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataFloat32
	return p
}

func (*EventArrayDataFloat32Context) IsEventArrayDataFloat32Context() {}

func NewEventArrayDataFloat32Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataFloat32Context {
	var p = new(EventArrayDataFloat32Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataFloat32

	return p
}

func (s *EventArrayDataFloat32Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataFloat32Context) EVENT_ADF32() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADF32, 0)
}

func (s *EventArrayDataFloat32Context) AllVALUE_FLOAT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_DEC)
}

func (s *EventArrayDataFloat32Context) VALUE_FLOAT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_DEC, i)
}

func (s *EventArrayDataFloat32Context) AllVALUE_FLOAT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_HEX)
}

func (s *EventArrayDataFloat32Context) VALUE_FLOAT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_HEX, i)
}

func (s *EventArrayDataFloat32Context) AllVALUE_FLOAT_INF() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_INF)
}

func (s *EventArrayDataFloat32Context) VALUE_FLOAT_INF(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_INF, i)
}

func (s *EventArrayDataFloat32Context) AllVALUE_FLOAT_NAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_NAN)
}

func (s *EventArrayDataFloat32Context) VALUE_FLOAT_NAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_NAN, i)
}

func (s *EventArrayDataFloat32Context) AllVALUE_FLOAT_SNAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_SNAN)
}

func (s *EventArrayDataFloat32Context) VALUE_FLOAT_SNAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_SNAN, i)
}

func (s *EventArrayDataFloat32Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataFloat32Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataFloat32Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataFloat32(s)
	}
}

func (s *EventArrayDataFloat32Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataFloat32(s)
	}
}

func (p *CEEventParser) EventArrayDataFloat32() (localctx IEventArrayDataFloat32Context) {
	this := p
	_ = this

	localctx = NewEventArrayDataFloat32Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, CEEventParserRULE_eventArrayDataFloat32)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(287)
		p.Match(CEEventParserEVENT_ADF32)
	}
	p.SetState(291)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0 {
		{
			p.SetState(288)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(293)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataFloat64Context is an interface to support dynamic dispatch.
type IEventArrayDataFloat64Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataFloat64Context differentiates from other interfaces.
	IsEventArrayDataFloat64Context()
}

type EventArrayDataFloat64Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataFloat64Context() *EventArrayDataFloat64Context {
	var p = new(EventArrayDataFloat64Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataFloat64
	return p
}

func (*EventArrayDataFloat64Context) IsEventArrayDataFloat64Context() {}

func NewEventArrayDataFloat64Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataFloat64Context {
	var p = new(EventArrayDataFloat64Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataFloat64

	return p
}

func (s *EventArrayDataFloat64Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataFloat64Context) EVENT_ADF64() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADF64, 0)
}

func (s *EventArrayDataFloat64Context) AllVALUE_FLOAT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_DEC)
}

func (s *EventArrayDataFloat64Context) VALUE_FLOAT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_DEC, i)
}

func (s *EventArrayDataFloat64Context) AllVALUE_FLOAT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_HEX)
}

func (s *EventArrayDataFloat64Context) VALUE_FLOAT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_HEX, i)
}

func (s *EventArrayDataFloat64Context) AllVALUE_FLOAT_INF() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_INF)
}

func (s *EventArrayDataFloat64Context) VALUE_FLOAT_INF(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_INF, i)
}

func (s *EventArrayDataFloat64Context) AllVALUE_FLOAT_NAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_NAN)
}

func (s *EventArrayDataFloat64Context) VALUE_FLOAT_NAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_NAN, i)
}

func (s *EventArrayDataFloat64Context) AllVALUE_FLOAT_SNAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_SNAN)
}

func (s *EventArrayDataFloat64Context) VALUE_FLOAT_SNAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_SNAN, i)
}

func (s *EventArrayDataFloat64Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataFloat64Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataFloat64Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataFloat64(s)
	}
}

func (s *EventArrayDataFloat64Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataFloat64(s)
	}
}

func (p *CEEventParser) EventArrayDataFloat64() (localctx IEventArrayDataFloat64Context) {
	this := p
	_ = this

	localctx = NewEventArrayDataFloat64Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, CEEventParserRULE_eventArrayDataFloat64)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(294)
		p.Match(CEEventParserEVENT_ADF64)
	}
	p.SetState(298)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0 {
		{
			p.SetState(295)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(300)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataInt16Context is an interface to support dynamic dispatch.
type IEventArrayDataInt16Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataInt16Context differentiates from other interfaces.
	IsEventArrayDataInt16Context()
}

type EventArrayDataInt16Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataInt16Context() *EventArrayDataInt16Context {
	var p = new(EventArrayDataInt16Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataInt16
	return p
}

func (*EventArrayDataInt16Context) IsEventArrayDataInt16Context() {}

func NewEventArrayDataInt16Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataInt16Context {
	var p = new(EventArrayDataInt16Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataInt16

	return p
}

func (s *EventArrayDataInt16Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataInt16Context) EVENT_ADI16() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADI16, 0)
}

func (s *EventArrayDataInt16Context) AllVALUE_INT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_BIN)
}

func (s *EventArrayDataInt16Context) VALUE_INT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_BIN, i)
}

func (s *EventArrayDataInt16Context) AllVALUE_INT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_OCT)
}

func (s *EventArrayDataInt16Context) VALUE_INT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_OCT, i)
}

func (s *EventArrayDataInt16Context) AllVALUE_INT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_DEC)
}

func (s *EventArrayDataInt16Context) VALUE_INT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_DEC, i)
}

func (s *EventArrayDataInt16Context) AllVALUE_INT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_HEX)
}

func (s *EventArrayDataInt16Context) VALUE_INT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_HEX, i)
}

func (s *EventArrayDataInt16Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataInt16Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataInt16Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataInt16(s)
	}
}

func (s *EventArrayDataInt16Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataInt16(s)
	}
}

func (p *CEEventParser) EventArrayDataInt16() (localctx IEventArrayDataInt16Context) {
	this := p
	_ = this

	localctx = NewEventArrayDataInt16Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, CEEventParserRULE_eventArrayDataInt16)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(301)
		p.Match(CEEventParserEVENT_ADI16)
	}
	p.SetState(305)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0 {
		{
			p.SetState(302)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(307)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataInt32Context is an interface to support dynamic dispatch.
type IEventArrayDataInt32Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataInt32Context differentiates from other interfaces.
	IsEventArrayDataInt32Context()
}

type EventArrayDataInt32Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataInt32Context() *EventArrayDataInt32Context {
	var p = new(EventArrayDataInt32Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataInt32
	return p
}

func (*EventArrayDataInt32Context) IsEventArrayDataInt32Context() {}

func NewEventArrayDataInt32Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataInt32Context {
	var p = new(EventArrayDataInt32Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataInt32

	return p
}

func (s *EventArrayDataInt32Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataInt32Context) EVENT_ADI32() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADI32, 0)
}

func (s *EventArrayDataInt32Context) AllVALUE_INT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_BIN)
}

func (s *EventArrayDataInt32Context) VALUE_INT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_BIN, i)
}

func (s *EventArrayDataInt32Context) AllVALUE_INT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_OCT)
}

func (s *EventArrayDataInt32Context) VALUE_INT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_OCT, i)
}

func (s *EventArrayDataInt32Context) AllVALUE_INT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_DEC)
}

func (s *EventArrayDataInt32Context) VALUE_INT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_DEC, i)
}

func (s *EventArrayDataInt32Context) AllVALUE_INT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_HEX)
}

func (s *EventArrayDataInt32Context) VALUE_INT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_HEX, i)
}

func (s *EventArrayDataInt32Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataInt32Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataInt32Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataInt32(s)
	}
}

func (s *EventArrayDataInt32Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataInt32(s)
	}
}

func (p *CEEventParser) EventArrayDataInt32() (localctx IEventArrayDataInt32Context) {
	this := p
	_ = this

	localctx = NewEventArrayDataInt32Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, CEEventParserRULE_eventArrayDataInt32)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(308)
		p.Match(CEEventParserEVENT_ADI32)
	}
	p.SetState(312)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0 {
		{
			p.SetState(309)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(314)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataInt64Context is an interface to support dynamic dispatch.
type IEventArrayDataInt64Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataInt64Context differentiates from other interfaces.
	IsEventArrayDataInt64Context()
}

type EventArrayDataInt64Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataInt64Context() *EventArrayDataInt64Context {
	var p = new(EventArrayDataInt64Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataInt64
	return p
}

func (*EventArrayDataInt64Context) IsEventArrayDataInt64Context() {}

func NewEventArrayDataInt64Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataInt64Context {
	var p = new(EventArrayDataInt64Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataInt64

	return p
}

func (s *EventArrayDataInt64Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataInt64Context) EVENT_ADI64() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADI64, 0)
}

func (s *EventArrayDataInt64Context) AllVALUE_INT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_BIN)
}

func (s *EventArrayDataInt64Context) VALUE_INT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_BIN, i)
}

func (s *EventArrayDataInt64Context) AllVALUE_INT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_OCT)
}

func (s *EventArrayDataInt64Context) VALUE_INT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_OCT, i)
}

func (s *EventArrayDataInt64Context) AllVALUE_INT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_DEC)
}

func (s *EventArrayDataInt64Context) VALUE_INT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_DEC, i)
}

func (s *EventArrayDataInt64Context) AllVALUE_INT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_HEX)
}

func (s *EventArrayDataInt64Context) VALUE_INT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_HEX, i)
}

func (s *EventArrayDataInt64Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataInt64Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataInt64Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataInt64(s)
	}
}

func (s *EventArrayDataInt64Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataInt64(s)
	}
}

func (p *CEEventParser) EventArrayDataInt64() (localctx IEventArrayDataInt64Context) {
	this := p
	_ = this

	localctx = NewEventArrayDataInt64Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, CEEventParserRULE_eventArrayDataInt64)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(315)
		p.Match(CEEventParserEVENT_ADI64)
	}
	p.SetState(319)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0 {
		{
			p.SetState(316)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(321)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataInt8Context is an interface to support dynamic dispatch.
type IEventArrayDataInt8Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataInt8Context differentiates from other interfaces.
	IsEventArrayDataInt8Context()
}

type EventArrayDataInt8Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataInt8Context() *EventArrayDataInt8Context {
	var p = new(EventArrayDataInt8Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataInt8
	return p
}

func (*EventArrayDataInt8Context) IsEventArrayDataInt8Context() {}

func NewEventArrayDataInt8Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataInt8Context {
	var p = new(EventArrayDataInt8Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataInt8

	return p
}

func (s *EventArrayDataInt8Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataInt8Context) EVENT_ADI8() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADI8, 0)
}

func (s *EventArrayDataInt8Context) AllVALUE_INT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_BIN)
}

func (s *EventArrayDataInt8Context) VALUE_INT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_BIN, i)
}

func (s *EventArrayDataInt8Context) AllVALUE_INT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_OCT)
}

func (s *EventArrayDataInt8Context) VALUE_INT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_OCT, i)
}

func (s *EventArrayDataInt8Context) AllVALUE_INT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_DEC)
}

func (s *EventArrayDataInt8Context) VALUE_INT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_DEC, i)
}

func (s *EventArrayDataInt8Context) AllVALUE_INT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_HEX)
}

func (s *EventArrayDataInt8Context) VALUE_INT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_HEX, i)
}

func (s *EventArrayDataInt8Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataInt8Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataInt8Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataInt8(s)
	}
}

func (s *EventArrayDataInt8Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataInt8(s)
	}
}

func (p *CEEventParser) EventArrayDataInt8() (localctx IEventArrayDataInt8Context) {
	this := p
	_ = this

	localctx = NewEventArrayDataInt8Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, CEEventParserRULE_eventArrayDataInt8)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(322)
		p.Match(CEEventParserEVENT_ADI8)
	}
	p.SetState(326)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0 {
		{
			p.SetState(323)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(328)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataTextContext is an interface to support dynamic dispatch.
type IEventArrayDataTextContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataTextContext differentiates from other interfaces.
	IsEventArrayDataTextContext()
}

type EventArrayDataTextContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataTextContext() *EventArrayDataTextContext {
	var p = new(EventArrayDataTextContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataText
	return p
}

func (*EventArrayDataTextContext) IsEventArrayDataTextContext() {}

func NewEventArrayDataTextContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataTextContext {
	var p = new(EventArrayDataTextContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataText

	return p
}

func (s *EventArrayDataTextContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataTextContext) EVENT_ADT() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADT, 0)
}

func (s *EventArrayDataTextContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventArrayDataTextContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataTextContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataTextContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataText(s)
	}
}

func (s *EventArrayDataTextContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataText(s)
	}
}

func (p *CEEventParser) EventArrayDataText() (localctx IEventArrayDataTextContext) {
	this := p
	_ = this

	localctx = NewEventArrayDataTextContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, CEEventParserRULE_eventArrayDataText)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(329)
		p.Match(CEEventParserEVENT_ADT)
	}
	p.SetState(331)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CEEventParserSTRING {
		{
			p.SetState(330)
			p.Match(CEEventParserSTRING)
		}

	}

	return localctx
}

// IEventArrayDataUIDContext is an interface to support dynamic dispatch.
type IEventArrayDataUIDContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataUIDContext differentiates from other interfaces.
	IsEventArrayDataUIDContext()
}

type EventArrayDataUIDContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataUIDContext() *EventArrayDataUIDContext {
	var p = new(EventArrayDataUIDContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataUID
	return p
}

func (*EventArrayDataUIDContext) IsEventArrayDataUIDContext() {}

func NewEventArrayDataUIDContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataUIDContext {
	var p = new(EventArrayDataUIDContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataUID

	return p
}

func (s *EventArrayDataUIDContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataUIDContext) EVENT_ADU() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADU, 0)
}

func (s *EventArrayDataUIDContext) AllVALUE_UID() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UID)
}

func (s *EventArrayDataUIDContext) VALUE_UID(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UID, i)
}

func (s *EventArrayDataUIDContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataUIDContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataUIDContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataUID(s)
	}
}

func (s *EventArrayDataUIDContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataUID(s)
	}
}

func (p *CEEventParser) EventArrayDataUID() (localctx IEventArrayDataUIDContext) {
	this := p
	_ = this

	localctx = NewEventArrayDataUIDContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, CEEventParserRULE_eventArrayDataUID)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(333)
		p.Match(CEEventParserEVENT_ADU)
	}
	p.SetState(337)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CEEventParserVALUE_UID {
		{
			p.SetState(334)
			p.Match(CEEventParserVALUE_UID)
		}

		p.SetState(339)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataUint16Context is an interface to support dynamic dispatch.
type IEventArrayDataUint16Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataUint16Context differentiates from other interfaces.
	IsEventArrayDataUint16Context()
}

type EventArrayDataUint16Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataUint16Context() *EventArrayDataUint16Context {
	var p = new(EventArrayDataUint16Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint16
	return p
}

func (*EventArrayDataUint16Context) IsEventArrayDataUint16Context() {}

func NewEventArrayDataUint16Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataUint16Context {
	var p = new(EventArrayDataUint16Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint16

	return p
}

func (s *EventArrayDataUint16Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataUint16Context) EVENT_ADU16() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADU16, 0)
}

func (s *EventArrayDataUint16Context) AllVALUE_UINT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_BIN)
}

func (s *EventArrayDataUint16Context) VALUE_UINT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_BIN, i)
}

func (s *EventArrayDataUint16Context) AllVALUE_UINT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_OCT)
}

func (s *EventArrayDataUint16Context) VALUE_UINT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_OCT, i)
}

func (s *EventArrayDataUint16Context) AllVALUE_UINT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_DEC)
}

func (s *EventArrayDataUint16Context) VALUE_UINT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, i)
}

func (s *EventArrayDataUint16Context) AllVALUE_UINT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_HEX)
}

func (s *EventArrayDataUint16Context) VALUE_UINT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_HEX, i)
}

func (s *EventArrayDataUint16Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataUint16Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataUint16Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataUint16(s)
	}
}

func (s *EventArrayDataUint16Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataUint16(s)
	}
}

func (p *CEEventParser) EventArrayDataUint16() (localctx IEventArrayDataUint16Context) {
	this := p
	_ = this

	localctx = NewEventArrayDataUint16Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, CEEventParserRULE_eventArrayDataUint16)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(340)
		p.Match(CEEventParserEVENT_ADU16)
	}
	p.SetState(344)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0 {
		{
			p.SetState(341)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(346)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataUint16XContext is an interface to support dynamic dispatch.
type IEventArrayDataUint16XContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataUint16XContext differentiates from other interfaces.
	IsEventArrayDataUint16XContext()
}

type EventArrayDataUint16XContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataUint16XContext() *EventArrayDataUint16XContext {
	var p = new(EventArrayDataUint16XContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint16X
	return p
}

func (*EventArrayDataUint16XContext) IsEventArrayDataUint16XContext() {}

func NewEventArrayDataUint16XContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataUint16XContext {
	var p = new(EventArrayDataUint16XContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint16X

	return p
}

func (s *EventArrayDataUint16XContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataUint16XContext) EVENT_ADU16X() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADU16X, 0)
}

func (s *EventArrayDataUint16XContext) AllVALUE_UINTX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINTX)
}

func (s *EventArrayDataUint16XContext) VALUE_UINTX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINTX, i)
}

func (s *EventArrayDataUint16XContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataUint16XContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataUint16XContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataUint16X(s)
	}
}

func (s *EventArrayDataUint16XContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataUint16X(s)
	}
}

func (p *CEEventParser) EventArrayDataUint16X() (localctx IEventArrayDataUint16XContext) {
	this := p
	_ = this

	localctx = NewEventArrayDataUint16XContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, CEEventParserRULE_eventArrayDataUint16X)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(347)
		p.Match(CEEventParserEVENT_ADU16X)
	}
	p.SetState(351)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CEEventParserVALUE_UINTX {
		{
			p.SetState(348)
			p.Match(CEEventParserVALUE_UINTX)
		}

		p.SetState(353)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataUint32Context is an interface to support dynamic dispatch.
type IEventArrayDataUint32Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataUint32Context differentiates from other interfaces.
	IsEventArrayDataUint32Context()
}

type EventArrayDataUint32Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataUint32Context() *EventArrayDataUint32Context {
	var p = new(EventArrayDataUint32Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint32
	return p
}

func (*EventArrayDataUint32Context) IsEventArrayDataUint32Context() {}

func NewEventArrayDataUint32Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataUint32Context {
	var p = new(EventArrayDataUint32Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint32

	return p
}

func (s *EventArrayDataUint32Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataUint32Context) EVENT_ADU32() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADU32, 0)
}

func (s *EventArrayDataUint32Context) AllVALUE_UINT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_BIN)
}

func (s *EventArrayDataUint32Context) VALUE_UINT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_BIN, i)
}

func (s *EventArrayDataUint32Context) AllVALUE_UINT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_OCT)
}

func (s *EventArrayDataUint32Context) VALUE_UINT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_OCT, i)
}

func (s *EventArrayDataUint32Context) AllVALUE_UINT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_DEC)
}

func (s *EventArrayDataUint32Context) VALUE_UINT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, i)
}

func (s *EventArrayDataUint32Context) AllVALUE_UINT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_HEX)
}

func (s *EventArrayDataUint32Context) VALUE_UINT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_HEX, i)
}

func (s *EventArrayDataUint32Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataUint32Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataUint32Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataUint32(s)
	}
}

func (s *EventArrayDataUint32Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataUint32(s)
	}
}

func (p *CEEventParser) EventArrayDataUint32() (localctx IEventArrayDataUint32Context) {
	this := p
	_ = this

	localctx = NewEventArrayDataUint32Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, CEEventParserRULE_eventArrayDataUint32)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(354)
		p.Match(CEEventParserEVENT_ADU32)
	}
	p.SetState(358)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0 {
		{
			p.SetState(355)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(360)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataUint32XContext is an interface to support dynamic dispatch.
type IEventArrayDataUint32XContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataUint32XContext differentiates from other interfaces.
	IsEventArrayDataUint32XContext()
}

type EventArrayDataUint32XContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataUint32XContext() *EventArrayDataUint32XContext {
	var p = new(EventArrayDataUint32XContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint32X
	return p
}

func (*EventArrayDataUint32XContext) IsEventArrayDataUint32XContext() {}

func NewEventArrayDataUint32XContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataUint32XContext {
	var p = new(EventArrayDataUint32XContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint32X

	return p
}

func (s *EventArrayDataUint32XContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataUint32XContext) EVENT_ADU32X() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADU32X, 0)
}

func (s *EventArrayDataUint32XContext) AllVALUE_UINTX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINTX)
}

func (s *EventArrayDataUint32XContext) VALUE_UINTX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINTX, i)
}

func (s *EventArrayDataUint32XContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataUint32XContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataUint32XContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataUint32X(s)
	}
}

func (s *EventArrayDataUint32XContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataUint32X(s)
	}
}

func (p *CEEventParser) EventArrayDataUint32X() (localctx IEventArrayDataUint32XContext) {
	this := p
	_ = this

	localctx = NewEventArrayDataUint32XContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, CEEventParserRULE_eventArrayDataUint32X)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(361)
		p.Match(CEEventParserEVENT_ADU32X)
	}
	p.SetState(365)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CEEventParserVALUE_UINTX {
		{
			p.SetState(362)
			p.Match(CEEventParserVALUE_UINTX)
		}

		p.SetState(367)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataUint64Context is an interface to support dynamic dispatch.
type IEventArrayDataUint64Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataUint64Context differentiates from other interfaces.
	IsEventArrayDataUint64Context()
}

type EventArrayDataUint64Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataUint64Context() *EventArrayDataUint64Context {
	var p = new(EventArrayDataUint64Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint64
	return p
}

func (*EventArrayDataUint64Context) IsEventArrayDataUint64Context() {}

func NewEventArrayDataUint64Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataUint64Context {
	var p = new(EventArrayDataUint64Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint64

	return p
}

func (s *EventArrayDataUint64Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataUint64Context) EVENT_ADU64() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADU64, 0)
}

func (s *EventArrayDataUint64Context) AllVALUE_UINT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_BIN)
}

func (s *EventArrayDataUint64Context) VALUE_UINT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_BIN, i)
}

func (s *EventArrayDataUint64Context) AllVALUE_UINT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_OCT)
}

func (s *EventArrayDataUint64Context) VALUE_UINT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_OCT, i)
}

func (s *EventArrayDataUint64Context) AllVALUE_UINT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_DEC)
}

func (s *EventArrayDataUint64Context) VALUE_UINT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, i)
}

func (s *EventArrayDataUint64Context) AllVALUE_UINT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_HEX)
}

func (s *EventArrayDataUint64Context) VALUE_UINT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_HEX, i)
}

func (s *EventArrayDataUint64Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataUint64Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataUint64Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataUint64(s)
	}
}

func (s *EventArrayDataUint64Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataUint64(s)
	}
}

func (p *CEEventParser) EventArrayDataUint64() (localctx IEventArrayDataUint64Context) {
	this := p
	_ = this

	localctx = NewEventArrayDataUint64Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, CEEventParserRULE_eventArrayDataUint64)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(368)
		p.Match(CEEventParserEVENT_ADU64)
	}
	p.SetState(372)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0 {
		{
			p.SetState(369)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(374)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataUint64XContext is an interface to support dynamic dispatch.
type IEventArrayDataUint64XContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataUint64XContext differentiates from other interfaces.
	IsEventArrayDataUint64XContext()
}

type EventArrayDataUint64XContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataUint64XContext() *EventArrayDataUint64XContext {
	var p = new(EventArrayDataUint64XContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint64X
	return p
}

func (*EventArrayDataUint64XContext) IsEventArrayDataUint64XContext() {}

func NewEventArrayDataUint64XContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataUint64XContext {
	var p = new(EventArrayDataUint64XContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint64X

	return p
}

func (s *EventArrayDataUint64XContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataUint64XContext) EVENT_ADU64X() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADU64X, 0)
}

func (s *EventArrayDataUint64XContext) AllVALUE_UINTX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINTX)
}

func (s *EventArrayDataUint64XContext) VALUE_UINTX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINTX, i)
}

func (s *EventArrayDataUint64XContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataUint64XContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataUint64XContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataUint64X(s)
	}
}

func (s *EventArrayDataUint64XContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataUint64X(s)
	}
}

func (p *CEEventParser) EventArrayDataUint64X() (localctx IEventArrayDataUint64XContext) {
	this := p
	_ = this

	localctx = NewEventArrayDataUint64XContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, CEEventParserRULE_eventArrayDataUint64X)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(375)
		p.Match(CEEventParserEVENT_ADU64X)
	}
	p.SetState(379)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CEEventParserVALUE_UINTX {
		{
			p.SetState(376)
			p.Match(CEEventParserVALUE_UINTX)
		}

		p.SetState(381)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataUint8Context is an interface to support dynamic dispatch.
type IEventArrayDataUint8Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataUint8Context differentiates from other interfaces.
	IsEventArrayDataUint8Context()
}

type EventArrayDataUint8Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataUint8Context() *EventArrayDataUint8Context {
	var p = new(EventArrayDataUint8Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint8
	return p
}

func (*EventArrayDataUint8Context) IsEventArrayDataUint8Context() {}

func NewEventArrayDataUint8Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataUint8Context {
	var p = new(EventArrayDataUint8Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint8

	return p
}

func (s *EventArrayDataUint8Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataUint8Context) EVENT_ADU8() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADU8, 0)
}

func (s *EventArrayDataUint8Context) AllVALUE_UINT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_BIN)
}

func (s *EventArrayDataUint8Context) VALUE_UINT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_BIN, i)
}

func (s *EventArrayDataUint8Context) AllVALUE_UINT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_OCT)
}

func (s *EventArrayDataUint8Context) VALUE_UINT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_OCT, i)
}

func (s *EventArrayDataUint8Context) AllVALUE_UINT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_DEC)
}

func (s *EventArrayDataUint8Context) VALUE_UINT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, i)
}

func (s *EventArrayDataUint8Context) AllVALUE_UINT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_HEX)
}

func (s *EventArrayDataUint8Context) VALUE_UINT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_HEX, i)
}

func (s *EventArrayDataUint8Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataUint8Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataUint8Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataUint8(s)
	}
}

func (s *EventArrayDataUint8Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataUint8(s)
	}
}

func (p *CEEventParser) EventArrayDataUint8() (localctx IEventArrayDataUint8Context) {
	this := p
	_ = this

	localctx = NewEventArrayDataUint8Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, CEEventParserRULE_eventArrayDataUint8)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(382)
		p.Match(CEEventParserEVENT_ADU8)
	}
	p.SetState(386)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0 {
		{
			p.SetState(383)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(388)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayDataUint8XContext is an interface to support dynamic dispatch.
type IEventArrayDataUint8XContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayDataUint8XContext differentiates from other interfaces.
	IsEventArrayDataUint8XContext()
}

type EventArrayDataUint8XContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayDataUint8XContext() *EventArrayDataUint8XContext {
	var p = new(EventArrayDataUint8XContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint8X
	return p
}

func (*EventArrayDataUint8XContext) IsEventArrayDataUint8XContext() {}

func NewEventArrayDataUint8XContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayDataUint8XContext {
	var p = new(EventArrayDataUint8XContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayDataUint8X

	return p
}

func (s *EventArrayDataUint8XContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayDataUint8XContext) EVENT_ADU8X() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ADU8X, 0)
}

func (s *EventArrayDataUint8XContext) AllVALUE_UINTX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINTX)
}

func (s *EventArrayDataUint8XContext) VALUE_UINTX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINTX, i)
}

func (s *EventArrayDataUint8XContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayDataUint8XContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayDataUint8XContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayDataUint8X(s)
	}
}

func (s *EventArrayDataUint8XContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayDataUint8X(s)
	}
}

func (p *CEEventParser) EventArrayDataUint8X() (localctx IEventArrayDataUint8XContext) {
	this := p
	_ = this

	localctx = NewEventArrayDataUint8XContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, CEEventParserRULE_eventArrayDataUint8X)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(389)
		p.Match(CEEventParserEVENT_ADU8X)
	}
	p.SetState(393)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CEEventParserVALUE_UINTX {
		{
			p.SetState(390)
			p.Match(CEEventParserVALUE_UINTX)
		}

		p.SetState(395)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventArrayFloat16Context is an interface to support dynamic dispatch.
type IEventArrayFloat16Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayFloat16Context differentiates from other interfaces.
	IsEventArrayFloat16Context()
}

type EventArrayFloat16Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayFloat16Context() *EventArrayFloat16Context {
	var p = new(EventArrayFloat16Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayFloat16
	return p
}

func (*EventArrayFloat16Context) IsEventArrayFloat16Context() {}

func NewEventArrayFloat16Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayFloat16Context {
	var p = new(EventArrayFloat16Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayFloat16

	return p
}

func (s *EventArrayFloat16Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayFloat16Context) EVENT_AF16() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AF16, 0)
}

func (s *EventArrayFloat16Context) EVENT_AF16_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AF16_ARGS, 0)
}

func (s *EventArrayFloat16Context) AllVALUE_FLOAT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_DEC)
}

func (s *EventArrayFloat16Context) VALUE_FLOAT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_DEC, i)
}

func (s *EventArrayFloat16Context) AllVALUE_FLOAT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_HEX)
}

func (s *EventArrayFloat16Context) VALUE_FLOAT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_HEX, i)
}

func (s *EventArrayFloat16Context) AllVALUE_FLOAT_INF() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_INF)
}

func (s *EventArrayFloat16Context) VALUE_FLOAT_INF(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_INF, i)
}

func (s *EventArrayFloat16Context) AllVALUE_FLOAT_NAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_NAN)
}

func (s *EventArrayFloat16Context) VALUE_FLOAT_NAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_NAN, i)
}

func (s *EventArrayFloat16Context) AllVALUE_FLOAT_SNAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_SNAN)
}

func (s *EventArrayFloat16Context) VALUE_FLOAT_SNAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_SNAN, i)
}

func (s *EventArrayFloat16Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayFloat16Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayFloat16Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayFloat16(s)
	}
}

func (s *EventArrayFloat16Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayFloat16(s)
	}
}

func (p *CEEventParser) EventArrayFloat16() (localctx IEventArrayFloat16Context) {
	this := p
	_ = this

	localctx = NewEventArrayFloat16Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, CEEventParserRULE_eventArrayFloat16)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(404)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AF16:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(396)
			p.Match(CEEventParserEVENT_AF16)
		}

	case CEEventParserEVENT_AF16_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(397)
			p.Match(CEEventParserEVENT_AF16_ARGS)
		}
		p.SetState(401)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0 {
			{
				p.SetState(398)
				_la = p.GetTokenStream().LA(1)

				if !(((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			p.SetState(403)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayFloat32Context is an interface to support dynamic dispatch.
type IEventArrayFloat32Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayFloat32Context differentiates from other interfaces.
	IsEventArrayFloat32Context()
}

type EventArrayFloat32Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayFloat32Context() *EventArrayFloat32Context {
	var p = new(EventArrayFloat32Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayFloat32
	return p
}

func (*EventArrayFloat32Context) IsEventArrayFloat32Context() {}

func NewEventArrayFloat32Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayFloat32Context {
	var p = new(EventArrayFloat32Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayFloat32

	return p
}

func (s *EventArrayFloat32Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayFloat32Context) EVENT_AF32() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AF32, 0)
}

func (s *EventArrayFloat32Context) EVENT_AF32_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AF32_ARGS, 0)
}

func (s *EventArrayFloat32Context) AllVALUE_FLOAT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_DEC)
}

func (s *EventArrayFloat32Context) VALUE_FLOAT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_DEC, i)
}

func (s *EventArrayFloat32Context) AllVALUE_FLOAT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_HEX)
}

func (s *EventArrayFloat32Context) VALUE_FLOAT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_HEX, i)
}

func (s *EventArrayFloat32Context) AllVALUE_FLOAT_INF() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_INF)
}

func (s *EventArrayFloat32Context) VALUE_FLOAT_INF(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_INF, i)
}

func (s *EventArrayFloat32Context) AllVALUE_FLOAT_NAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_NAN)
}

func (s *EventArrayFloat32Context) VALUE_FLOAT_NAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_NAN, i)
}

func (s *EventArrayFloat32Context) AllVALUE_FLOAT_SNAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_SNAN)
}

func (s *EventArrayFloat32Context) VALUE_FLOAT_SNAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_SNAN, i)
}

func (s *EventArrayFloat32Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayFloat32Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayFloat32Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayFloat32(s)
	}
}

func (s *EventArrayFloat32Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayFloat32(s)
	}
}

func (p *CEEventParser) EventArrayFloat32() (localctx IEventArrayFloat32Context) {
	this := p
	_ = this

	localctx = NewEventArrayFloat32Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, CEEventParserRULE_eventArrayFloat32)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(414)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AF32:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(406)
			p.Match(CEEventParserEVENT_AF32)
		}

	case CEEventParserEVENT_AF32_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(407)
			p.Match(CEEventParserEVENT_AF32_ARGS)
		}
		p.SetState(411)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0 {
			{
				p.SetState(408)
				_la = p.GetTokenStream().LA(1)

				if !(((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			p.SetState(413)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayFloat64Context is an interface to support dynamic dispatch.
type IEventArrayFloat64Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayFloat64Context differentiates from other interfaces.
	IsEventArrayFloat64Context()
}

type EventArrayFloat64Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayFloat64Context() *EventArrayFloat64Context {
	var p = new(EventArrayFloat64Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayFloat64
	return p
}

func (*EventArrayFloat64Context) IsEventArrayFloat64Context() {}

func NewEventArrayFloat64Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayFloat64Context {
	var p = new(EventArrayFloat64Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayFloat64

	return p
}

func (s *EventArrayFloat64Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayFloat64Context) EVENT_AF64() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AF64, 0)
}

func (s *EventArrayFloat64Context) EVENT_AF64_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AF64_ARGS, 0)
}

func (s *EventArrayFloat64Context) AllVALUE_FLOAT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_DEC)
}

func (s *EventArrayFloat64Context) VALUE_FLOAT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_DEC, i)
}

func (s *EventArrayFloat64Context) AllVALUE_FLOAT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_HEX)
}

func (s *EventArrayFloat64Context) VALUE_FLOAT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_HEX, i)
}

func (s *EventArrayFloat64Context) AllVALUE_FLOAT_INF() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_INF)
}

func (s *EventArrayFloat64Context) VALUE_FLOAT_INF(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_INF, i)
}

func (s *EventArrayFloat64Context) AllVALUE_FLOAT_NAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_NAN)
}

func (s *EventArrayFloat64Context) VALUE_FLOAT_NAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_NAN, i)
}

func (s *EventArrayFloat64Context) AllVALUE_FLOAT_SNAN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_FLOAT_SNAN)
}

func (s *EventArrayFloat64Context) VALUE_FLOAT_SNAN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_FLOAT_SNAN, i)
}

func (s *EventArrayFloat64Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayFloat64Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayFloat64Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayFloat64(s)
	}
}

func (s *EventArrayFloat64Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayFloat64(s)
	}
}

func (p *CEEventParser) EventArrayFloat64() (localctx IEventArrayFloat64Context) {
	this := p
	_ = this

	localctx = NewEventArrayFloat64Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, CEEventParserRULE_eventArrayFloat64)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(424)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AF64:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(416)
			p.Match(CEEventParserEVENT_AF64)
		}

	case CEEventParserEVENT_AF64_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(417)
			p.Match(CEEventParserEVENT_AF64_ARGS)
		}
		p.SetState(421)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0 {
			{
				p.SetState(418)
				_la = p.GetTokenStream().LA(1)

				if !(((_la-127)&-(0x1f+1)) == 0 && ((1<<uint((_la-127)))&((1<<(CEEventParserVALUE_FLOAT_NAN-127))|(1<<(CEEventParserVALUE_FLOAT_SNAN-127))|(1<<(CEEventParserVALUE_FLOAT_INF-127))|(1<<(CEEventParserVALUE_FLOAT_DEC-127))|(1<<(CEEventParserVALUE_FLOAT_HEX-127)))) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			p.SetState(423)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayInt16Context is an interface to support dynamic dispatch.
type IEventArrayInt16Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayInt16Context differentiates from other interfaces.
	IsEventArrayInt16Context()
}

type EventArrayInt16Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayInt16Context() *EventArrayInt16Context {
	var p = new(EventArrayInt16Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayInt16
	return p
}

func (*EventArrayInt16Context) IsEventArrayInt16Context() {}

func NewEventArrayInt16Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayInt16Context {
	var p = new(EventArrayInt16Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayInt16

	return p
}

func (s *EventArrayInt16Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayInt16Context) EVENT_AI16() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AI16, 0)
}

func (s *EventArrayInt16Context) EVENT_AI16_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AI16_ARGS, 0)
}

func (s *EventArrayInt16Context) AllVALUE_INT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_BIN)
}

func (s *EventArrayInt16Context) VALUE_INT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_BIN, i)
}

func (s *EventArrayInt16Context) AllVALUE_INT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_OCT)
}

func (s *EventArrayInt16Context) VALUE_INT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_OCT, i)
}

func (s *EventArrayInt16Context) AllVALUE_INT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_DEC)
}

func (s *EventArrayInt16Context) VALUE_INT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_DEC, i)
}

func (s *EventArrayInt16Context) AllVALUE_INT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_HEX)
}

func (s *EventArrayInt16Context) VALUE_INT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_HEX, i)
}

func (s *EventArrayInt16Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayInt16Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayInt16Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayInt16(s)
	}
}

func (s *EventArrayInt16Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayInt16(s)
	}
}

func (p *CEEventParser) EventArrayInt16() (localctx IEventArrayInt16Context) {
	this := p
	_ = this

	localctx = NewEventArrayInt16Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, CEEventParserRULE_eventArrayInt16)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(434)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AI16:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(426)
			p.Match(CEEventParserEVENT_AI16)
		}

	case CEEventParserEVENT_AI16_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(427)
			p.Match(CEEventParserEVENT_AI16_ARGS)
		}
		p.SetState(431)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0 {
			{
				p.SetState(428)
				_la = p.GetTokenStream().LA(1)

				if !(((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			p.SetState(433)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayInt32Context is an interface to support dynamic dispatch.
type IEventArrayInt32Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayInt32Context differentiates from other interfaces.
	IsEventArrayInt32Context()
}

type EventArrayInt32Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayInt32Context() *EventArrayInt32Context {
	var p = new(EventArrayInt32Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayInt32
	return p
}

func (*EventArrayInt32Context) IsEventArrayInt32Context() {}

func NewEventArrayInt32Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayInt32Context {
	var p = new(EventArrayInt32Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayInt32

	return p
}

func (s *EventArrayInt32Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayInt32Context) EVENT_AI32() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AI32, 0)
}

func (s *EventArrayInt32Context) EVENT_AI32_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AI32_ARGS, 0)
}

func (s *EventArrayInt32Context) AllVALUE_INT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_BIN)
}

func (s *EventArrayInt32Context) VALUE_INT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_BIN, i)
}

func (s *EventArrayInt32Context) AllVALUE_INT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_OCT)
}

func (s *EventArrayInt32Context) VALUE_INT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_OCT, i)
}

func (s *EventArrayInt32Context) AllVALUE_INT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_DEC)
}

func (s *EventArrayInt32Context) VALUE_INT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_DEC, i)
}

func (s *EventArrayInt32Context) AllVALUE_INT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_HEX)
}

func (s *EventArrayInt32Context) VALUE_INT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_HEX, i)
}

func (s *EventArrayInt32Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayInt32Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayInt32Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayInt32(s)
	}
}

func (s *EventArrayInt32Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayInt32(s)
	}
}

func (p *CEEventParser) EventArrayInt32() (localctx IEventArrayInt32Context) {
	this := p
	_ = this

	localctx = NewEventArrayInt32Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, CEEventParserRULE_eventArrayInt32)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(444)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AI32:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(436)
			p.Match(CEEventParserEVENT_AI32)
		}

	case CEEventParserEVENT_AI32_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(437)
			p.Match(CEEventParserEVENT_AI32_ARGS)
		}
		p.SetState(441)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0 {
			{
				p.SetState(438)
				_la = p.GetTokenStream().LA(1)

				if !(((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			p.SetState(443)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayInt64Context is an interface to support dynamic dispatch.
type IEventArrayInt64Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayInt64Context differentiates from other interfaces.
	IsEventArrayInt64Context()
}

type EventArrayInt64Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayInt64Context() *EventArrayInt64Context {
	var p = new(EventArrayInt64Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayInt64
	return p
}

func (*EventArrayInt64Context) IsEventArrayInt64Context() {}

func NewEventArrayInt64Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayInt64Context {
	var p = new(EventArrayInt64Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayInt64

	return p
}

func (s *EventArrayInt64Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayInt64Context) EVENT_AI64() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AI64, 0)
}

func (s *EventArrayInt64Context) EVENT_AI64_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AI64_ARGS, 0)
}

func (s *EventArrayInt64Context) AllVALUE_INT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_BIN)
}

func (s *EventArrayInt64Context) VALUE_INT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_BIN, i)
}

func (s *EventArrayInt64Context) AllVALUE_INT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_OCT)
}

func (s *EventArrayInt64Context) VALUE_INT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_OCT, i)
}

func (s *EventArrayInt64Context) AllVALUE_INT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_DEC)
}

func (s *EventArrayInt64Context) VALUE_INT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_DEC, i)
}

func (s *EventArrayInt64Context) AllVALUE_INT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_HEX)
}

func (s *EventArrayInt64Context) VALUE_INT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_HEX, i)
}

func (s *EventArrayInt64Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayInt64Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayInt64Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayInt64(s)
	}
}

func (s *EventArrayInt64Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayInt64(s)
	}
}

func (p *CEEventParser) EventArrayInt64() (localctx IEventArrayInt64Context) {
	this := p
	_ = this

	localctx = NewEventArrayInt64Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, CEEventParserRULE_eventArrayInt64)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(454)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AI64:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(446)
			p.Match(CEEventParserEVENT_AI64)
		}

	case CEEventParserEVENT_AI64_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(447)
			p.Match(CEEventParserEVENT_AI64_ARGS)
		}
		p.SetState(451)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0 {
			{
				p.SetState(448)
				_la = p.GetTokenStream().LA(1)

				if !(((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			p.SetState(453)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayInt8Context is an interface to support dynamic dispatch.
type IEventArrayInt8Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayInt8Context differentiates from other interfaces.
	IsEventArrayInt8Context()
}

type EventArrayInt8Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayInt8Context() *EventArrayInt8Context {
	var p = new(EventArrayInt8Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayInt8
	return p
}

func (*EventArrayInt8Context) IsEventArrayInt8Context() {}

func NewEventArrayInt8Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayInt8Context {
	var p = new(EventArrayInt8Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayInt8

	return p
}

func (s *EventArrayInt8Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayInt8Context) EVENT_AI8() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AI8, 0)
}

func (s *EventArrayInt8Context) EVENT_AI8_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AI8_ARGS, 0)
}

func (s *EventArrayInt8Context) AllVALUE_INT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_BIN)
}

func (s *EventArrayInt8Context) VALUE_INT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_BIN, i)
}

func (s *EventArrayInt8Context) AllVALUE_INT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_OCT)
}

func (s *EventArrayInt8Context) VALUE_INT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_OCT, i)
}

func (s *EventArrayInt8Context) AllVALUE_INT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_DEC)
}

func (s *EventArrayInt8Context) VALUE_INT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_DEC, i)
}

func (s *EventArrayInt8Context) AllVALUE_INT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_INT_HEX)
}

func (s *EventArrayInt8Context) VALUE_INT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_INT_HEX, i)
}

func (s *EventArrayInt8Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayInt8Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayInt8Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayInt8(s)
	}
}

func (s *EventArrayInt8Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayInt8(s)
	}
}

func (p *CEEventParser) EventArrayInt8() (localctx IEventArrayInt8Context) {
	this := p
	_ = this

	localctx = NewEventArrayInt8Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, CEEventParserRULE_eventArrayInt8)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(464)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AI8:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(456)
			p.Match(CEEventParserEVENT_AI8)
		}

	case CEEventParserEVENT_AI8_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(457)
			p.Match(CEEventParserEVENT_AI8_ARGS)
		}
		p.SetState(461)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0 {
			{
				p.SetState(458)
				_la = p.GetTokenStream().LA(1)

				if !(((_la-122)&-(0x1f+1)) == 0 && ((1<<uint((_la-122)))&((1<<(CEEventParserVALUE_INT_BIN-122))|(1<<(CEEventParserVALUE_INT_OCT-122))|(1<<(CEEventParserVALUE_INT_DEC-122))|(1<<(CEEventParserVALUE_INT_HEX-122)))) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			p.SetState(463)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayUIDContext is an interface to support dynamic dispatch.
type IEventArrayUIDContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayUIDContext differentiates from other interfaces.
	IsEventArrayUIDContext()
}

type EventArrayUIDContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayUIDContext() *EventArrayUIDContext {
	var p = new(EventArrayUIDContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayUID
	return p
}

func (*EventArrayUIDContext) IsEventArrayUIDContext() {}

func NewEventArrayUIDContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayUIDContext {
	var p = new(EventArrayUIDContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayUID

	return p
}

func (s *EventArrayUIDContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayUIDContext) EVENT_AU() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU, 0)
}

func (s *EventArrayUIDContext) EVENT_AU_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU_ARGS, 0)
}

func (s *EventArrayUIDContext) AllVALUE_UID() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UID)
}

func (s *EventArrayUIDContext) VALUE_UID(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UID, i)
}

func (s *EventArrayUIDContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayUIDContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayUIDContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayUID(s)
	}
}

func (s *EventArrayUIDContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayUID(s)
	}
}

func (p *CEEventParser) EventArrayUID() (localctx IEventArrayUIDContext) {
	this := p
	_ = this

	localctx = NewEventArrayUIDContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, CEEventParserRULE_eventArrayUID)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(474)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AU:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(466)
			p.Match(CEEventParserEVENT_AU)
		}

	case CEEventParserEVENT_AU_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(467)
			p.Match(CEEventParserEVENT_AU_ARGS)
		}
		p.SetState(471)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == CEEventParserVALUE_UID {
			{
				p.SetState(468)
				p.Match(CEEventParserVALUE_UID)
			}

			p.SetState(473)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayUint16Context is an interface to support dynamic dispatch.
type IEventArrayUint16Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayUint16Context differentiates from other interfaces.
	IsEventArrayUint16Context()
}

type EventArrayUint16Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayUint16Context() *EventArrayUint16Context {
	var p = new(EventArrayUint16Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayUint16
	return p
}

func (*EventArrayUint16Context) IsEventArrayUint16Context() {}

func NewEventArrayUint16Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayUint16Context {
	var p = new(EventArrayUint16Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayUint16

	return p
}

func (s *EventArrayUint16Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayUint16Context) EVENT_AU16() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU16, 0)
}

func (s *EventArrayUint16Context) EVENT_AU16_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU16_ARGS, 0)
}

func (s *EventArrayUint16Context) AllVALUE_UINT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_BIN)
}

func (s *EventArrayUint16Context) VALUE_UINT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_BIN, i)
}

func (s *EventArrayUint16Context) AllVALUE_UINT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_OCT)
}

func (s *EventArrayUint16Context) VALUE_UINT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_OCT, i)
}

func (s *EventArrayUint16Context) AllVALUE_UINT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_DEC)
}

func (s *EventArrayUint16Context) VALUE_UINT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, i)
}

func (s *EventArrayUint16Context) AllVALUE_UINT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_HEX)
}

func (s *EventArrayUint16Context) VALUE_UINT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_HEX, i)
}

func (s *EventArrayUint16Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayUint16Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayUint16Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayUint16(s)
	}
}

func (s *EventArrayUint16Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayUint16(s)
	}
}

func (p *CEEventParser) EventArrayUint16() (localctx IEventArrayUint16Context) {
	this := p
	_ = this

	localctx = NewEventArrayUint16Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, CEEventParserRULE_eventArrayUint16)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(484)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AU16:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(476)
			p.Match(CEEventParserEVENT_AU16)
		}

	case CEEventParserEVENT_AU16_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(477)
			p.Match(CEEventParserEVENT_AU16_ARGS)
		}
		p.SetState(481)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0 {
			{
				p.SetState(478)
				_la = p.GetTokenStream().LA(1)

				if !(((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			p.SetState(483)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayUint16XContext is an interface to support dynamic dispatch.
type IEventArrayUint16XContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayUint16XContext differentiates from other interfaces.
	IsEventArrayUint16XContext()
}

type EventArrayUint16XContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayUint16XContext() *EventArrayUint16XContext {
	var p = new(EventArrayUint16XContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayUint16X
	return p
}

func (*EventArrayUint16XContext) IsEventArrayUint16XContext() {}

func NewEventArrayUint16XContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayUint16XContext {
	var p = new(EventArrayUint16XContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayUint16X

	return p
}

func (s *EventArrayUint16XContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayUint16XContext) EVENT_AU16X() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU16X, 0)
}

func (s *EventArrayUint16XContext) EVENT_AU16X_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU16X_ARGS, 0)
}

func (s *EventArrayUint16XContext) AllVALUE_UINTX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINTX)
}

func (s *EventArrayUint16XContext) VALUE_UINTX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINTX, i)
}

func (s *EventArrayUint16XContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayUint16XContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayUint16XContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayUint16X(s)
	}
}

func (s *EventArrayUint16XContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayUint16X(s)
	}
}

func (p *CEEventParser) EventArrayUint16X() (localctx IEventArrayUint16XContext) {
	this := p
	_ = this

	localctx = NewEventArrayUint16XContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, CEEventParserRULE_eventArrayUint16X)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(494)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AU16X:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(486)
			p.Match(CEEventParserEVENT_AU16X)
		}

	case CEEventParserEVENT_AU16X_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(487)
			p.Match(CEEventParserEVENT_AU16X_ARGS)
		}
		p.SetState(491)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == CEEventParserVALUE_UINTX {
			{
				p.SetState(488)
				p.Match(CEEventParserVALUE_UINTX)
			}

			p.SetState(493)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayUint32Context is an interface to support dynamic dispatch.
type IEventArrayUint32Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayUint32Context differentiates from other interfaces.
	IsEventArrayUint32Context()
}

type EventArrayUint32Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayUint32Context() *EventArrayUint32Context {
	var p = new(EventArrayUint32Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayUint32
	return p
}

func (*EventArrayUint32Context) IsEventArrayUint32Context() {}

func NewEventArrayUint32Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayUint32Context {
	var p = new(EventArrayUint32Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayUint32

	return p
}

func (s *EventArrayUint32Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayUint32Context) EVENT_AU32() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU32, 0)
}

func (s *EventArrayUint32Context) EVENT_AU32_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU32_ARGS, 0)
}

func (s *EventArrayUint32Context) AllVALUE_UINT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_BIN)
}

func (s *EventArrayUint32Context) VALUE_UINT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_BIN, i)
}

func (s *EventArrayUint32Context) AllVALUE_UINT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_OCT)
}

func (s *EventArrayUint32Context) VALUE_UINT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_OCT, i)
}

func (s *EventArrayUint32Context) AllVALUE_UINT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_DEC)
}

func (s *EventArrayUint32Context) VALUE_UINT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, i)
}

func (s *EventArrayUint32Context) AllVALUE_UINT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_HEX)
}

func (s *EventArrayUint32Context) VALUE_UINT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_HEX, i)
}

func (s *EventArrayUint32Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayUint32Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayUint32Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayUint32(s)
	}
}

func (s *EventArrayUint32Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayUint32(s)
	}
}

func (p *CEEventParser) EventArrayUint32() (localctx IEventArrayUint32Context) {
	this := p
	_ = this

	localctx = NewEventArrayUint32Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, CEEventParserRULE_eventArrayUint32)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(504)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AU32:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(496)
			p.Match(CEEventParserEVENT_AU32)
		}

	case CEEventParserEVENT_AU32_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(497)
			p.Match(CEEventParserEVENT_AU32_ARGS)
		}
		p.SetState(501)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0 {
			{
				p.SetState(498)
				_la = p.GetTokenStream().LA(1)

				if !(((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			p.SetState(503)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayUint32XContext is an interface to support dynamic dispatch.
type IEventArrayUint32XContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayUint32XContext differentiates from other interfaces.
	IsEventArrayUint32XContext()
}

type EventArrayUint32XContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayUint32XContext() *EventArrayUint32XContext {
	var p = new(EventArrayUint32XContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayUint32X
	return p
}

func (*EventArrayUint32XContext) IsEventArrayUint32XContext() {}

func NewEventArrayUint32XContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayUint32XContext {
	var p = new(EventArrayUint32XContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayUint32X

	return p
}

func (s *EventArrayUint32XContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayUint32XContext) EVENT_AU32X() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU32X, 0)
}

func (s *EventArrayUint32XContext) EVENT_AU32X_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU32X_ARGS, 0)
}

func (s *EventArrayUint32XContext) AllVALUE_UINTX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINTX)
}

func (s *EventArrayUint32XContext) VALUE_UINTX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINTX, i)
}

func (s *EventArrayUint32XContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayUint32XContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayUint32XContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayUint32X(s)
	}
}

func (s *EventArrayUint32XContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayUint32X(s)
	}
}

func (p *CEEventParser) EventArrayUint32X() (localctx IEventArrayUint32XContext) {
	this := p
	_ = this

	localctx = NewEventArrayUint32XContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, CEEventParserRULE_eventArrayUint32X)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(514)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AU32X:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(506)
			p.Match(CEEventParserEVENT_AU32X)
		}

	case CEEventParserEVENT_AU32X_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(507)
			p.Match(CEEventParserEVENT_AU32X_ARGS)
		}
		p.SetState(511)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == CEEventParserVALUE_UINTX {
			{
				p.SetState(508)
				p.Match(CEEventParserVALUE_UINTX)
			}

			p.SetState(513)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayUint64Context is an interface to support dynamic dispatch.
type IEventArrayUint64Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayUint64Context differentiates from other interfaces.
	IsEventArrayUint64Context()
}

type EventArrayUint64Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayUint64Context() *EventArrayUint64Context {
	var p = new(EventArrayUint64Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayUint64
	return p
}

func (*EventArrayUint64Context) IsEventArrayUint64Context() {}

func NewEventArrayUint64Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayUint64Context {
	var p = new(EventArrayUint64Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayUint64

	return p
}

func (s *EventArrayUint64Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayUint64Context) EVENT_AU64() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU64, 0)
}

func (s *EventArrayUint64Context) EVENT_AU64_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU64_ARGS, 0)
}

func (s *EventArrayUint64Context) AllVALUE_UINT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_BIN)
}

func (s *EventArrayUint64Context) VALUE_UINT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_BIN, i)
}

func (s *EventArrayUint64Context) AllVALUE_UINT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_OCT)
}

func (s *EventArrayUint64Context) VALUE_UINT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_OCT, i)
}

func (s *EventArrayUint64Context) AllVALUE_UINT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_DEC)
}

func (s *EventArrayUint64Context) VALUE_UINT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, i)
}

func (s *EventArrayUint64Context) AllVALUE_UINT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_HEX)
}

func (s *EventArrayUint64Context) VALUE_UINT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_HEX, i)
}

func (s *EventArrayUint64Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayUint64Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayUint64Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayUint64(s)
	}
}

func (s *EventArrayUint64Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayUint64(s)
	}
}

func (p *CEEventParser) EventArrayUint64() (localctx IEventArrayUint64Context) {
	this := p
	_ = this

	localctx = NewEventArrayUint64Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, CEEventParserRULE_eventArrayUint64)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(524)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AU64:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(516)
			p.Match(CEEventParserEVENT_AU64)
		}

	case CEEventParserEVENT_AU64_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(517)
			p.Match(CEEventParserEVENT_AU64_ARGS)
		}
		p.SetState(521)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0 {
			{
				p.SetState(518)
				_la = p.GetTokenStream().LA(1)

				if !(((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			p.SetState(523)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayUint64XContext is an interface to support dynamic dispatch.
type IEventArrayUint64XContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayUint64XContext differentiates from other interfaces.
	IsEventArrayUint64XContext()
}

type EventArrayUint64XContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayUint64XContext() *EventArrayUint64XContext {
	var p = new(EventArrayUint64XContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayUint64X
	return p
}

func (*EventArrayUint64XContext) IsEventArrayUint64XContext() {}

func NewEventArrayUint64XContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayUint64XContext {
	var p = new(EventArrayUint64XContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayUint64X

	return p
}

func (s *EventArrayUint64XContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayUint64XContext) EVENT_AU64X() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU64X, 0)
}

func (s *EventArrayUint64XContext) EVENT_AU64X_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU64X_ARGS, 0)
}

func (s *EventArrayUint64XContext) AllVALUE_UINTX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINTX)
}

func (s *EventArrayUint64XContext) VALUE_UINTX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINTX, i)
}

func (s *EventArrayUint64XContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayUint64XContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayUint64XContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayUint64X(s)
	}
}

func (s *EventArrayUint64XContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayUint64X(s)
	}
}

func (p *CEEventParser) EventArrayUint64X() (localctx IEventArrayUint64XContext) {
	this := p
	_ = this

	localctx = NewEventArrayUint64XContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, CEEventParserRULE_eventArrayUint64X)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(534)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AU64X:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(526)
			p.Match(CEEventParserEVENT_AU64X)
		}

	case CEEventParserEVENT_AU64X_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(527)
			p.Match(CEEventParserEVENT_AU64X_ARGS)
		}
		p.SetState(531)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == CEEventParserVALUE_UINTX {
			{
				p.SetState(528)
				p.Match(CEEventParserVALUE_UINTX)
			}

			p.SetState(533)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayUint8Context is an interface to support dynamic dispatch.
type IEventArrayUint8Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayUint8Context differentiates from other interfaces.
	IsEventArrayUint8Context()
}

type EventArrayUint8Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayUint8Context() *EventArrayUint8Context {
	var p = new(EventArrayUint8Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayUint8
	return p
}

func (*EventArrayUint8Context) IsEventArrayUint8Context() {}

func NewEventArrayUint8Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayUint8Context {
	var p = new(EventArrayUint8Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayUint8

	return p
}

func (s *EventArrayUint8Context) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayUint8Context) EVENT_AU8() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU8, 0)
}

func (s *EventArrayUint8Context) EVENT_AU8_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU8_ARGS, 0)
}

func (s *EventArrayUint8Context) AllVALUE_UINT_BIN() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_BIN)
}

func (s *EventArrayUint8Context) VALUE_UINT_BIN(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_BIN, i)
}

func (s *EventArrayUint8Context) AllVALUE_UINT_OCT() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_OCT)
}

func (s *EventArrayUint8Context) VALUE_UINT_OCT(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_OCT, i)
}

func (s *EventArrayUint8Context) AllVALUE_UINT_DEC() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_DEC)
}

func (s *EventArrayUint8Context) VALUE_UINT_DEC(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, i)
}

func (s *EventArrayUint8Context) AllVALUE_UINT_HEX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINT_HEX)
}

func (s *EventArrayUint8Context) VALUE_UINT_HEX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_HEX, i)
}

func (s *EventArrayUint8Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayUint8Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayUint8Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayUint8(s)
	}
}

func (s *EventArrayUint8Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayUint8(s)
	}
}

func (p *CEEventParser) EventArrayUint8() (localctx IEventArrayUint8Context) {
	this := p
	_ = this

	localctx = NewEventArrayUint8Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, CEEventParserRULE_eventArrayUint8)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(544)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AU8:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(536)
			p.Match(CEEventParserEVENT_AU8)
		}

	case CEEventParserEVENT_AU8_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(537)
			p.Match(CEEventParserEVENT_AU8_ARGS)
		}
		p.SetState(541)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0 {
			{
				p.SetState(538)
				_la = p.GetTokenStream().LA(1)

				if !(((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CEEventParserVALUE_UINT_BIN-115))|(1<<(CEEventParserVALUE_UINT_OCT-115))|(1<<(CEEventParserVALUE_UINT_DEC-115))|(1<<(CEEventParserVALUE_UINT_HEX-115)))) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

			p.SetState(543)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventArrayUint8XContext is an interface to support dynamic dispatch.
type IEventArrayUint8XContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventArrayUint8XContext differentiates from other interfaces.
	IsEventArrayUint8XContext()
}

type EventArrayUint8XContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventArrayUint8XContext() *EventArrayUint8XContext {
	var p = new(EventArrayUint8XContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventArrayUint8X
	return p
}

func (*EventArrayUint8XContext) IsEventArrayUint8XContext() {}

func NewEventArrayUint8XContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventArrayUint8XContext {
	var p = new(EventArrayUint8XContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventArrayUint8X

	return p
}

func (s *EventArrayUint8XContext) GetParser() antlr.Parser { return s.parser }

func (s *EventArrayUint8XContext) EVENT_AU8X() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU8X, 0)
}

func (s *EventArrayUint8XContext) EVENT_AU8X_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_AU8X_ARGS, 0)
}

func (s *EventArrayUint8XContext) AllVALUE_UINTX() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserVALUE_UINTX)
}

func (s *EventArrayUint8XContext) VALUE_UINTX(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINTX, i)
}

func (s *EventArrayUint8XContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventArrayUint8XContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventArrayUint8XContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventArrayUint8X(s)
	}
}

func (s *EventArrayUint8XContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventArrayUint8X(s)
	}
}

func (p *CEEventParser) EventArrayUint8X() (localctx IEventArrayUint8XContext) {
	this := p
	_ = this

	localctx = NewEventArrayUint8XContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, CEEventParserRULE_eventArrayUint8X)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(554)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_AU8X:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(546)
			p.Match(CEEventParserEVENT_AU8X)
		}

	case CEEventParserEVENT_AU8X_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(547)
			p.Match(CEEventParserEVENT_AU8X_ARGS)
		}
		p.SetState(551)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == CEEventParserVALUE_UINTX {
			{
				p.SetState(548)
				p.Match(CEEventParserVALUE_UINTX)
			}

			p.SetState(553)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventBeginArrayBitsContext is an interface to support dynamic dispatch.
type IEventBeginArrayBitsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayBitsContext differentiates from other interfaces.
	IsEventBeginArrayBitsContext()
}

type EventBeginArrayBitsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayBitsContext() *EventBeginArrayBitsContext {
	var p = new(EventBeginArrayBitsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayBits
	return p
}

func (*EventBeginArrayBitsContext) IsEventBeginArrayBitsContext() {}

func NewEventBeginArrayBitsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayBitsContext {
	var p = new(EventBeginArrayBitsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayBits

	return p
}

func (s *EventBeginArrayBitsContext) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayBitsContext) EVENT_BAB() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAB, 0)
}

func (s *EventBeginArrayBitsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayBitsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayBitsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayBits(s)
	}
}

func (s *EventBeginArrayBitsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayBits(s)
	}
}

func (p *CEEventParser) EventBeginArrayBits() (localctx IEventBeginArrayBitsContext) {
	this := p
	_ = this

	localctx = NewEventBeginArrayBitsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, CEEventParserRULE_eventBeginArrayBits)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(556)
		p.Match(CEEventParserEVENT_BAB)
	}

	return localctx
}

// IEventBeginArrayFloat16Context is an interface to support dynamic dispatch.
type IEventBeginArrayFloat16Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayFloat16Context differentiates from other interfaces.
	IsEventBeginArrayFloat16Context()
}

type EventBeginArrayFloat16Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayFloat16Context() *EventBeginArrayFloat16Context {
	var p = new(EventBeginArrayFloat16Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayFloat16
	return p
}

func (*EventBeginArrayFloat16Context) IsEventBeginArrayFloat16Context() {}

func NewEventBeginArrayFloat16Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayFloat16Context {
	var p = new(EventBeginArrayFloat16Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayFloat16

	return p
}

func (s *EventBeginArrayFloat16Context) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayFloat16Context) EVENT_BAF16() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAF16, 0)
}

func (s *EventBeginArrayFloat16Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayFloat16Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayFloat16Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayFloat16(s)
	}
}

func (s *EventBeginArrayFloat16Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayFloat16(s)
	}
}

func (p *CEEventParser) EventBeginArrayFloat16() (localctx IEventBeginArrayFloat16Context) {
	this := p
	_ = this

	localctx = NewEventBeginArrayFloat16Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, CEEventParserRULE_eventBeginArrayFloat16)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(558)
		p.Match(CEEventParserEVENT_BAF16)
	}

	return localctx
}

// IEventBeginArrayFloat32Context is an interface to support dynamic dispatch.
type IEventBeginArrayFloat32Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayFloat32Context differentiates from other interfaces.
	IsEventBeginArrayFloat32Context()
}

type EventBeginArrayFloat32Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayFloat32Context() *EventBeginArrayFloat32Context {
	var p = new(EventBeginArrayFloat32Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayFloat32
	return p
}

func (*EventBeginArrayFloat32Context) IsEventBeginArrayFloat32Context() {}

func NewEventBeginArrayFloat32Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayFloat32Context {
	var p = new(EventBeginArrayFloat32Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayFloat32

	return p
}

func (s *EventBeginArrayFloat32Context) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayFloat32Context) EVENT_BAF32() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAF32, 0)
}

func (s *EventBeginArrayFloat32Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayFloat32Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayFloat32Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayFloat32(s)
	}
}

func (s *EventBeginArrayFloat32Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayFloat32(s)
	}
}

func (p *CEEventParser) EventBeginArrayFloat32() (localctx IEventBeginArrayFloat32Context) {
	this := p
	_ = this

	localctx = NewEventBeginArrayFloat32Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, CEEventParserRULE_eventBeginArrayFloat32)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(560)
		p.Match(CEEventParserEVENT_BAF32)
	}

	return localctx
}

// IEventBeginArrayFloat64Context is an interface to support dynamic dispatch.
type IEventBeginArrayFloat64Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayFloat64Context differentiates from other interfaces.
	IsEventBeginArrayFloat64Context()
}

type EventBeginArrayFloat64Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayFloat64Context() *EventBeginArrayFloat64Context {
	var p = new(EventBeginArrayFloat64Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayFloat64
	return p
}

func (*EventBeginArrayFloat64Context) IsEventBeginArrayFloat64Context() {}

func NewEventBeginArrayFloat64Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayFloat64Context {
	var p = new(EventBeginArrayFloat64Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayFloat64

	return p
}

func (s *EventBeginArrayFloat64Context) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayFloat64Context) EVENT_BAF64() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAF64, 0)
}

func (s *EventBeginArrayFloat64Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayFloat64Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayFloat64Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayFloat64(s)
	}
}

func (s *EventBeginArrayFloat64Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayFloat64(s)
	}
}

func (p *CEEventParser) EventBeginArrayFloat64() (localctx IEventBeginArrayFloat64Context) {
	this := p
	_ = this

	localctx = NewEventBeginArrayFloat64Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, CEEventParserRULE_eventBeginArrayFloat64)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(562)
		p.Match(CEEventParserEVENT_BAF64)
	}

	return localctx
}

// IEventBeginArrayInt16Context is an interface to support dynamic dispatch.
type IEventBeginArrayInt16Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayInt16Context differentiates from other interfaces.
	IsEventBeginArrayInt16Context()
}

type EventBeginArrayInt16Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayInt16Context() *EventBeginArrayInt16Context {
	var p = new(EventBeginArrayInt16Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayInt16
	return p
}

func (*EventBeginArrayInt16Context) IsEventBeginArrayInt16Context() {}

func NewEventBeginArrayInt16Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayInt16Context {
	var p = new(EventBeginArrayInt16Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayInt16

	return p
}

func (s *EventBeginArrayInt16Context) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayInt16Context) EVENT_BAI16() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAI16, 0)
}

func (s *EventBeginArrayInt16Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayInt16Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayInt16Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayInt16(s)
	}
}

func (s *EventBeginArrayInt16Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayInt16(s)
	}
}

func (p *CEEventParser) EventBeginArrayInt16() (localctx IEventBeginArrayInt16Context) {
	this := p
	_ = this

	localctx = NewEventBeginArrayInt16Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, CEEventParserRULE_eventBeginArrayInt16)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(564)
		p.Match(CEEventParserEVENT_BAI16)
	}

	return localctx
}

// IEventBeginArrayInt32Context is an interface to support dynamic dispatch.
type IEventBeginArrayInt32Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayInt32Context differentiates from other interfaces.
	IsEventBeginArrayInt32Context()
}

type EventBeginArrayInt32Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayInt32Context() *EventBeginArrayInt32Context {
	var p = new(EventBeginArrayInt32Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayInt32
	return p
}

func (*EventBeginArrayInt32Context) IsEventBeginArrayInt32Context() {}

func NewEventBeginArrayInt32Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayInt32Context {
	var p = new(EventBeginArrayInt32Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayInt32

	return p
}

func (s *EventBeginArrayInt32Context) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayInt32Context) EVENT_BAI32() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAI32, 0)
}

func (s *EventBeginArrayInt32Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayInt32Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayInt32Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayInt32(s)
	}
}

func (s *EventBeginArrayInt32Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayInt32(s)
	}
}

func (p *CEEventParser) EventBeginArrayInt32() (localctx IEventBeginArrayInt32Context) {
	this := p
	_ = this

	localctx = NewEventBeginArrayInt32Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, CEEventParserRULE_eventBeginArrayInt32)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(566)
		p.Match(CEEventParserEVENT_BAI32)
	}

	return localctx
}

// IEventBeginArrayInt64Context is an interface to support dynamic dispatch.
type IEventBeginArrayInt64Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayInt64Context differentiates from other interfaces.
	IsEventBeginArrayInt64Context()
}

type EventBeginArrayInt64Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayInt64Context() *EventBeginArrayInt64Context {
	var p = new(EventBeginArrayInt64Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayInt64
	return p
}

func (*EventBeginArrayInt64Context) IsEventBeginArrayInt64Context() {}

func NewEventBeginArrayInt64Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayInt64Context {
	var p = new(EventBeginArrayInt64Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayInt64

	return p
}

func (s *EventBeginArrayInt64Context) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayInt64Context) EVENT_BAI64() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAI64, 0)
}

func (s *EventBeginArrayInt64Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayInt64Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayInt64Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayInt64(s)
	}
}

func (s *EventBeginArrayInt64Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayInt64(s)
	}
}

func (p *CEEventParser) EventBeginArrayInt64() (localctx IEventBeginArrayInt64Context) {
	this := p
	_ = this

	localctx = NewEventBeginArrayInt64Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, CEEventParserRULE_eventBeginArrayInt64)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(568)
		p.Match(CEEventParserEVENT_BAI64)
	}

	return localctx
}

// IEventBeginArrayInt8Context is an interface to support dynamic dispatch.
type IEventBeginArrayInt8Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayInt8Context differentiates from other interfaces.
	IsEventBeginArrayInt8Context()
}

type EventBeginArrayInt8Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayInt8Context() *EventBeginArrayInt8Context {
	var p = new(EventBeginArrayInt8Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayInt8
	return p
}

func (*EventBeginArrayInt8Context) IsEventBeginArrayInt8Context() {}

func NewEventBeginArrayInt8Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayInt8Context {
	var p = new(EventBeginArrayInt8Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayInt8

	return p
}

func (s *EventBeginArrayInt8Context) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayInt8Context) EVENT_BAI8() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAI8, 0)
}

func (s *EventBeginArrayInt8Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayInt8Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayInt8Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayInt8(s)
	}
}

func (s *EventBeginArrayInt8Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayInt8(s)
	}
}

func (p *CEEventParser) EventBeginArrayInt8() (localctx IEventBeginArrayInt8Context) {
	this := p
	_ = this

	localctx = NewEventBeginArrayInt8Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, CEEventParserRULE_eventBeginArrayInt8)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(570)
		p.Match(CEEventParserEVENT_BAI8)
	}

	return localctx
}

// IEventBeginArrayUIDContext is an interface to support dynamic dispatch.
type IEventBeginArrayUIDContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayUIDContext differentiates from other interfaces.
	IsEventBeginArrayUIDContext()
}

type EventBeginArrayUIDContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayUIDContext() *EventBeginArrayUIDContext {
	var p = new(EventBeginArrayUIDContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayUID
	return p
}

func (*EventBeginArrayUIDContext) IsEventBeginArrayUIDContext() {}

func NewEventBeginArrayUIDContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayUIDContext {
	var p = new(EventBeginArrayUIDContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayUID

	return p
}

func (s *EventBeginArrayUIDContext) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayUIDContext) EVENT_BAU() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAU, 0)
}

func (s *EventBeginArrayUIDContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayUIDContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayUIDContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayUID(s)
	}
}

func (s *EventBeginArrayUIDContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayUID(s)
	}
}

func (p *CEEventParser) EventBeginArrayUID() (localctx IEventBeginArrayUIDContext) {
	this := p
	_ = this

	localctx = NewEventBeginArrayUIDContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, CEEventParserRULE_eventBeginArrayUID)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(572)
		p.Match(CEEventParserEVENT_BAU)
	}

	return localctx
}

// IEventBeginArrayUint16Context is an interface to support dynamic dispatch.
type IEventBeginArrayUint16Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayUint16Context differentiates from other interfaces.
	IsEventBeginArrayUint16Context()
}

type EventBeginArrayUint16Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayUint16Context() *EventBeginArrayUint16Context {
	var p = new(EventBeginArrayUint16Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayUint16
	return p
}

func (*EventBeginArrayUint16Context) IsEventBeginArrayUint16Context() {}

func NewEventBeginArrayUint16Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayUint16Context {
	var p = new(EventBeginArrayUint16Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayUint16

	return p
}

func (s *EventBeginArrayUint16Context) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayUint16Context) EVENT_BAU16() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAU16, 0)
}

func (s *EventBeginArrayUint16Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayUint16Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayUint16Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayUint16(s)
	}
}

func (s *EventBeginArrayUint16Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayUint16(s)
	}
}

func (p *CEEventParser) EventBeginArrayUint16() (localctx IEventBeginArrayUint16Context) {
	this := p
	_ = this

	localctx = NewEventBeginArrayUint16Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, CEEventParserRULE_eventBeginArrayUint16)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(574)
		p.Match(CEEventParserEVENT_BAU16)
	}

	return localctx
}

// IEventBeginArrayUint32Context is an interface to support dynamic dispatch.
type IEventBeginArrayUint32Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayUint32Context differentiates from other interfaces.
	IsEventBeginArrayUint32Context()
}

type EventBeginArrayUint32Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayUint32Context() *EventBeginArrayUint32Context {
	var p = new(EventBeginArrayUint32Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayUint32
	return p
}

func (*EventBeginArrayUint32Context) IsEventBeginArrayUint32Context() {}

func NewEventBeginArrayUint32Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayUint32Context {
	var p = new(EventBeginArrayUint32Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayUint32

	return p
}

func (s *EventBeginArrayUint32Context) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayUint32Context) EVENT_BAU32() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAU32, 0)
}

func (s *EventBeginArrayUint32Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayUint32Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayUint32Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayUint32(s)
	}
}

func (s *EventBeginArrayUint32Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayUint32(s)
	}
}

func (p *CEEventParser) EventBeginArrayUint32() (localctx IEventBeginArrayUint32Context) {
	this := p
	_ = this

	localctx = NewEventBeginArrayUint32Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, CEEventParserRULE_eventBeginArrayUint32)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(576)
		p.Match(CEEventParserEVENT_BAU32)
	}

	return localctx
}

// IEventBeginArrayUint64Context is an interface to support dynamic dispatch.
type IEventBeginArrayUint64Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayUint64Context differentiates from other interfaces.
	IsEventBeginArrayUint64Context()
}

type EventBeginArrayUint64Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayUint64Context() *EventBeginArrayUint64Context {
	var p = new(EventBeginArrayUint64Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayUint64
	return p
}

func (*EventBeginArrayUint64Context) IsEventBeginArrayUint64Context() {}

func NewEventBeginArrayUint64Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayUint64Context {
	var p = new(EventBeginArrayUint64Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayUint64

	return p
}

func (s *EventBeginArrayUint64Context) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayUint64Context) EVENT_BAU64() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAU64, 0)
}

func (s *EventBeginArrayUint64Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayUint64Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayUint64Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayUint64(s)
	}
}

func (s *EventBeginArrayUint64Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayUint64(s)
	}
}

func (p *CEEventParser) EventBeginArrayUint64() (localctx IEventBeginArrayUint64Context) {
	this := p
	_ = this

	localctx = NewEventBeginArrayUint64Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, CEEventParserRULE_eventBeginArrayUint64)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(578)
		p.Match(CEEventParserEVENT_BAU64)
	}

	return localctx
}

// IEventBeginArrayUint8Context is an interface to support dynamic dispatch.
type IEventBeginArrayUint8Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginArrayUint8Context differentiates from other interfaces.
	IsEventBeginArrayUint8Context()
}

type EventBeginArrayUint8Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginArrayUint8Context() *EventBeginArrayUint8Context {
	var p = new(EventBeginArrayUint8Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginArrayUint8
	return p
}

func (*EventBeginArrayUint8Context) IsEventBeginArrayUint8Context() {}

func NewEventBeginArrayUint8Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginArrayUint8Context {
	var p = new(EventBeginArrayUint8Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginArrayUint8

	return p
}

func (s *EventBeginArrayUint8Context) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginArrayUint8Context) EVENT_BAU8() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BAU8, 0)
}

func (s *EventBeginArrayUint8Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginArrayUint8Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginArrayUint8Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginArrayUint8(s)
	}
}

func (s *EventBeginArrayUint8Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginArrayUint8(s)
	}
}

func (p *CEEventParser) EventBeginArrayUint8() (localctx IEventBeginArrayUint8Context) {
	this := p
	_ = this

	localctx = NewEventBeginArrayUint8Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, CEEventParserRULE_eventBeginArrayUint8)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(580)
		p.Match(CEEventParserEVENT_BAU8)
	}

	return localctx
}

// IEventBeginCustomBinaryContext is an interface to support dynamic dispatch.
type IEventBeginCustomBinaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginCustomBinaryContext differentiates from other interfaces.
	IsEventBeginCustomBinaryContext()
}

type EventBeginCustomBinaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginCustomBinaryContext() *EventBeginCustomBinaryContext {
	var p = new(EventBeginCustomBinaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginCustomBinary
	return p
}

func (*EventBeginCustomBinaryContext) IsEventBeginCustomBinaryContext() {}

func NewEventBeginCustomBinaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginCustomBinaryContext {
	var p = new(EventBeginCustomBinaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginCustomBinary

	return p
}

func (s *EventBeginCustomBinaryContext) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginCustomBinaryContext) EVENT_BCB() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BCB, 0)
}

func (s *EventBeginCustomBinaryContext) VALUE_UINT_DEC() antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, 0)
}

func (s *EventBeginCustomBinaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginCustomBinaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginCustomBinaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginCustomBinary(s)
	}
}

func (s *EventBeginCustomBinaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginCustomBinary(s)
	}
}

func (p *CEEventParser) EventBeginCustomBinary() (localctx IEventBeginCustomBinaryContext) {
	this := p
	_ = this

	localctx = NewEventBeginCustomBinaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, CEEventParserRULE_eventBeginCustomBinary)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(582)
		p.Match(CEEventParserEVENT_BCB)
	}
	{
		p.SetState(583)
		p.Match(CEEventParserVALUE_UINT_DEC)
	}

	return localctx
}

// IEventBeginCustomTextContext is an interface to support dynamic dispatch.
type IEventBeginCustomTextContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginCustomTextContext differentiates from other interfaces.
	IsEventBeginCustomTextContext()
}

type EventBeginCustomTextContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginCustomTextContext() *EventBeginCustomTextContext {
	var p = new(EventBeginCustomTextContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginCustomText
	return p
}

func (*EventBeginCustomTextContext) IsEventBeginCustomTextContext() {}

func NewEventBeginCustomTextContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginCustomTextContext {
	var p = new(EventBeginCustomTextContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginCustomText

	return p
}

func (s *EventBeginCustomTextContext) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginCustomTextContext) EVENT_BCT() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BCT, 0)
}

func (s *EventBeginCustomTextContext) VALUE_UINT_DEC() antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, 0)
}

func (s *EventBeginCustomTextContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginCustomTextContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginCustomTextContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginCustomText(s)
	}
}

func (s *EventBeginCustomTextContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginCustomText(s)
	}
}

func (p *CEEventParser) EventBeginCustomText() (localctx IEventBeginCustomTextContext) {
	this := p
	_ = this

	localctx = NewEventBeginCustomTextContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, CEEventParserRULE_eventBeginCustomText)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(585)
		p.Match(CEEventParserEVENT_BCT)
	}
	{
		p.SetState(586)
		p.Match(CEEventParserVALUE_UINT_DEC)
	}

	return localctx
}

// IEventBeginMediaContext is an interface to support dynamic dispatch.
type IEventBeginMediaContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginMediaContext differentiates from other interfaces.
	IsEventBeginMediaContext()
}

type EventBeginMediaContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginMediaContext() *EventBeginMediaContext {
	var p = new(EventBeginMediaContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginMedia
	return p
}

func (*EventBeginMediaContext) IsEventBeginMediaContext() {}

func NewEventBeginMediaContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginMediaContext {
	var p = new(EventBeginMediaContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginMedia

	return p
}

func (s *EventBeginMediaContext) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginMediaContext) EVENT_BMEDIA() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BMEDIA, 0)
}

func (s *EventBeginMediaContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventBeginMediaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginMediaContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginMediaContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginMedia(s)
	}
}

func (s *EventBeginMediaContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginMedia(s)
	}
}

func (p *CEEventParser) EventBeginMedia() (localctx IEventBeginMediaContext) {
	this := p
	_ = this

	localctx = NewEventBeginMediaContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, CEEventParserRULE_eventBeginMedia)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(588)
		p.Match(CEEventParserEVENT_BMEDIA)
	}
	{
		p.SetState(589)
		p.Match(CEEventParserSTRING)
	}

	return localctx
}

// IEventBeginRemoteReferenceContext is an interface to support dynamic dispatch.
type IEventBeginRemoteReferenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginRemoteReferenceContext differentiates from other interfaces.
	IsEventBeginRemoteReferenceContext()
}

type EventBeginRemoteReferenceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginRemoteReferenceContext() *EventBeginRemoteReferenceContext {
	var p = new(EventBeginRemoteReferenceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginRemoteReference
	return p
}

func (*EventBeginRemoteReferenceContext) IsEventBeginRemoteReferenceContext() {}

func NewEventBeginRemoteReferenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginRemoteReferenceContext {
	var p = new(EventBeginRemoteReferenceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginRemoteReference

	return p
}

func (s *EventBeginRemoteReferenceContext) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginRemoteReferenceContext) EVENT_BREFR() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BREFR, 0)
}

func (s *EventBeginRemoteReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginRemoteReferenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginRemoteReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginRemoteReference(s)
	}
}

func (s *EventBeginRemoteReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginRemoteReference(s)
	}
}

func (p *CEEventParser) EventBeginRemoteReference() (localctx IEventBeginRemoteReferenceContext) {
	this := p
	_ = this

	localctx = NewEventBeginRemoteReferenceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, CEEventParserRULE_eventBeginRemoteReference)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(591)
		p.Match(CEEventParserEVENT_BREFR)
	}

	return localctx
}

// IEventBeginResourceIdContext is an interface to support dynamic dispatch.
type IEventBeginResourceIdContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginResourceIdContext differentiates from other interfaces.
	IsEventBeginResourceIdContext()
}

type EventBeginResourceIdContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginResourceIdContext() *EventBeginResourceIdContext {
	var p = new(EventBeginResourceIdContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginResourceId
	return p
}

func (*EventBeginResourceIdContext) IsEventBeginResourceIdContext() {}

func NewEventBeginResourceIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginResourceIdContext {
	var p = new(EventBeginResourceIdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginResourceId

	return p
}

func (s *EventBeginResourceIdContext) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginResourceIdContext) EVENT_BRID() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BRID, 0)
}

func (s *EventBeginResourceIdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginResourceIdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginResourceIdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginResourceId(s)
	}
}

func (s *EventBeginResourceIdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginResourceId(s)
	}
}

func (p *CEEventParser) EventBeginResourceId() (localctx IEventBeginResourceIdContext) {
	this := p
	_ = this

	localctx = NewEventBeginResourceIdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, CEEventParserRULE_eventBeginResourceId)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(593)
		p.Match(CEEventParserEVENT_BRID)
	}

	return localctx
}

// IEventBeginStringContext is an interface to support dynamic dispatch.
type IEventBeginStringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBeginStringContext differentiates from other interfaces.
	IsEventBeginStringContext()
}

type EventBeginStringContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBeginStringContext() *EventBeginStringContext {
	var p = new(EventBeginStringContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBeginString
	return p
}

func (*EventBeginStringContext) IsEventBeginStringContext() {}

func NewEventBeginStringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBeginStringContext {
	var p = new(EventBeginStringContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBeginString

	return p
}

func (s *EventBeginStringContext) GetParser() antlr.Parser { return s.parser }

func (s *EventBeginStringContext) EVENT_BS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_BS, 0)
}

func (s *EventBeginStringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBeginStringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBeginStringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBeginString(s)
	}
}

func (s *EventBeginStringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBeginString(s)
	}
}

func (p *CEEventParser) EventBeginString() (localctx IEventBeginStringContext) {
	this := p
	_ = this

	localctx = NewEventBeginStringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, CEEventParserRULE_eventBeginString)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(595)
		p.Match(CEEventParserEVENT_BS)
	}

	return localctx
}

// IEventBooleanContext is an interface to support dynamic dispatch.
type IEventBooleanContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventBooleanContext differentiates from other interfaces.
	IsEventBooleanContext()
}

type EventBooleanContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventBooleanContext() *EventBooleanContext {
	var p = new(EventBooleanContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventBoolean
	return p
}

func (*EventBooleanContext) IsEventBooleanContext() {}

func NewEventBooleanContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventBooleanContext {
	var p = new(EventBooleanContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventBoolean

	return p
}

func (s *EventBooleanContext) GetParser() antlr.Parser { return s.parser }

func (s *EventBooleanContext) EVENT_B() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_B, 0)
}

func (s *EventBooleanContext) TRUE() antlr.TerminalNode {
	return s.GetToken(CEEventParserTRUE, 0)
}

func (s *EventBooleanContext) FALSE() antlr.TerminalNode {
	return s.GetToken(CEEventParserFALSE, 0)
}

func (s *EventBooleanContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventBooleanContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventBooleanContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventBoolean(s)
	}
}

func (s *EventBooleanContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventBoolean(s)
	}
}

func (p *CEEventParser) EventBoolean() (localctx IEventBooleanContext) {
	this := p
	_ = this

	localctx = NewEventBooleanContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, CEEventParserRULE_eventBoolean)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(597)
		p.Match(CEEventParserEVENT_B)
	}
	{
		p.SetState(598)
		_la = p.GetTokenStream().LA(1)

		if !(_la == CEEventParserTRUE || _la == CEEventParserFALSE) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IEventCommentMultilineContext is an interface to support dynamic dispatch.
type IEventCommentMultilineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventCommentMultilineContext differentiates from other interfaces.
	IsEventCommentMultilineContext()
}

type EventCommentMultilineContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventCommentMultilineContext() *EventCommentMultilineContext {
	var p = new(EventCommentMultilineContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventCommentMultiline
	return p
}

func (*EventCommentMultilineContext) IsEventCommentMultilineContext() {}

func NewEventCommentMultilineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventCommentMultilineContext {
	var p = new(EventCommentMultilineContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventCommentMultiline

	return p
}

func (s *EventCommentMultilineContext) GetParser() antlr.Parser { return s.parser }

func (s *EventCommentMultilineContext) EVENT_CM() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_CM, 0)
}

func (s *EventCommentMultilineContext) EVENT_CM_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_CM_ARGS, 0)
}

func (s *EventCommentMultilineContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventCommentMultilineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventCommentMultilineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventCommentMultilineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventCommentMultiline(s)
	}
}

func (s *EventCommentMultilineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventCommentMultiline(s)
	}
}

func (p *CEEventParser) EventCommentMultiline() (localctx IEventCommentMultilineContext) {
	this := p
	_ = this

	localctx = NewEventCommentMultilineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 118, CEEventParserRULE_eventCommentMultiline)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(605)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_CM:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(600)
			p.Match(CEEventParserEVENT_CM)
		}

	case CEEventParserEVENT_CM_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(601)
			p.Match(CEEventParserEVENT_CM_ARGS)
		}
		p.SetState(603)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CEEventParserSTRING {
			{
				p.SetState(602)
				p.Match(CEEventParserSTRING)
			}

		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventCommentSingleLineContext is an interface to support dynamic dispatch.
type IEventCommentSingleLineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventCommentSingleLineContext differentiates from other interfaces.
	IsEventCommentSingleLineContext()
}

type EventCommentSingleLineContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventCommentSingleLineContext() *EventCommentSingleLineContext {
	var p = new(EventCommentSingleLineContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventCommentSingleLine
	return p
}

func (*EventCommentSingleLineContext) IsEventCommentSingleLineContext() {}

func NewEventCommentSingleLineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventCommentSingleLineContext {
	var p = new(EventCommentSingleLineContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventCommentSingleLine

	return p
}

func (s *EventCommentSingleLineContext) GetParser() antlr.Parser { return s.parser }

func (s *EventCommentSingleLineContext) EVENT_CS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_CS, 0)
}

func (s *EventCommentSingleLineContext) EVENT_CS_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_CS_ARGS, 0)
}

func (s *EventCommentSingleLineContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventCommentSingleLineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventCommentSingleLineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventCommentSingleLineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventCommentSingleLine(s)
	}
}

func (s *EventCommentSingleLineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventCommentSingleLine(s)
	}
}

func (p *CEEventParser) EventCommentSingleLine() (localctx IEventCommentSingleLineContext) {
	this := p
	_ = this

	localctx = NewEventCommentSingleLineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 120, CEEventParserRULE_eventCommentSingleLine)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(612)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_CS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(607)
			p.Match(CEEventParserEVENT_CS)
		}

	case CEEventParserEVENT_CS_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(608)
			p.Match(CEEventParserEVENT_CS_ARGS)
		}
		p.SetState(610)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CEEventParserSTRING {
			{
				p.SetState(609)
				p.Match(CEEventParserSTRING)
			}

		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventCustomBinaryContext is an interface to support dynamic dispatch.
type IEventCustomBinaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventCustomBinaryContext differentiates from other interfaces.
	IsEventCustomBinaryContext()
}

type EventCustomBinaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventCustomBinaryContext() *EventCustomBinaryContext {
	var p = new(EventCustomBinaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventCustomBinary
	return p
}

func (*EventCustomBinaryContext) IsEventCustomBinaryContext() {}

func NewEventCustomBinaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventCustomBinaryContext {
	var p = new(EventCustomBinaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventCustomBinary

	return p
}

func (s *EventCustomBinaryContext) GetParser() antlr.Parser { return s.parser }

func (s *EventCustomBinaryContext) EVENT_CB() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_CB, 0)
}

func (s *EventCustomBinaryContext) CUSTOM_BINARY_TYPE() antlr.TerminalNode {
	return s.GetToken(CEEventParserCUSTOM_BINARY_TYPE, 0)
}

func (s *EventCustomBinaryContext) AllBYTE() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserBYTE)
}

func (s *EventCustomBinaryContext) BYTE(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserBYTE, i)
}

func (s *EventCustomBinaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventCustomBinaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventCustomBinaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventCustomBinary(s)
	}
}

func (s *EventCustomBinaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventCustomBinary(s)
	}
}

func (p *CEEventParser) EventCustomBinary() (localctx IEventCustomBinaryContext) {
	this := p
	_ = this

	localctx = NewEventCustomBinaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 122, CEEventParserRULE_eventCustomBinary)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(614)
		p.Match(CEEventParserEVENT_CB)
	}
	{
		p.SetState(615)
		p.Match(CEEventParserCUSTOM_BINARY_TYPE)
	}
	p.SetState(619)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CEEventParserBYTE {
		{
			p.SetState(616)
			p.Match(CEEventParserBYTE)
		}

		p.SetState(621)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventCustomTextContext is an interface to support dynamic dispatch.
type IEventCustomTextContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventCustomTextContext differentiates from other interfaces.
	IsEventCustomTextContext()
}

type EventCustomTextContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventCustomTextContext() *EventCustomTextContext {
	var p = new(EventCustomTextContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventCustomText
	return p
}

func (*EventCustomTextContext) IsEventCustomTextContext() {}

func NewEventCustomTextContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventCustomTextContext {
	var p = new(EventCustomTextContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventCustomText

	return p
}

func (s *EventCustomTextContext) GetParser() antlr.Parser { return s.parser }

func (s *EventCustomTextContext) EVENT_CT() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_CT, 0)
}

func (s *EventCustomTextContext) CUSTOM_TEXT_TYPE() antlr.TerminalNode {
	return s.GetToken(CEEventParserCUSTOM_TEXT_TYPE, 0)
}

func (s *EventCustomTextContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventCustomTextContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventCustomTextContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventCustomTextContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventCustomText(s)
	}
}

func (s *EventCustomTextContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventCustomText(s)
	}
}

func (p *CEEventParser) EventCustomText() (localctx IEventCustomTextContext) {
	this := p
	_ = this

	localctx = NewEventCustomTextContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, CEEventParserRULE_eventCustomText)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(622)
		p.Match(CEEventParserEVENT_CT)
	}
	{
		p.SetState(623)
		p.Match(CEEventParserCUSTOM_TEXT_TYPE)
	}
	p.SetState(625)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CEEventParserSTRING {
		{
			p.SetState(624)
			p.Match(CEEventParserSTRING)
		}

	}

	return localctx
}

// IEventEdgeContext is an interface to support dynamic dispatch.
type IEventEdgeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventEdgeContext differentiates from other interfaces.
	IsEventEdgeContext()
}

type EventEdgeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventEdgeContext() *EventEdgeContext {
	var p = new(EventEdgeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventEdge
	return p
}

func (*EventEdgeContext) IsEventEdgeContext() {}

func NewEventEdgeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventEdgeContext {
	var p = new(EventEdgeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventEdge

	return p
}

func (s *EventEdgeContext) GetParser() antlr.Parser { return s.parser }

func (s *EventEdgeContext) EVENT_EDGE() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_EDGE, 0)
}

func (s *EventEdgeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventEdgeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventEdgeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventEdge(s)
	}
}

func (s *EventEdgeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventEdge(s)
	}
}

func (p *CEEventParser) EventEdge() (localctx IEventEdgeContext) {
	this := p
	_ = this

	localctx = NewEventEdgeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, CEEventParserRULE_eventEdge)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(627)
		p.Match(CEEventParserEVENT_EDGE)
	}

	return localctx
}

// IEventEndContainerContext is an interface to support dynamic dispatch.
type IEventEndContainerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventEndContainerContext differentiates from other interfaces.
	IsEventEndContainerContext()
}

type EventEndContainerContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventEndContainerContext() *EventEndContainerContext {
	var p = new(EventEndContainerContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventEndContainer
	return p
}

func (*EventEndContainerContext) IsEventEndContainerContext() {}

func NewEventEndContainerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventEndContainerContext {
	var p = new(EventEndContainerContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventEndContainer

	return p
}

func (s *EventEndContainerContext) GetParser() antlr.Parser { return s.parser }

func (s *EventEndContainerContext) EVENT_E() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_E, 0)
}

func (s *EventEndContainerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventEndContainerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventEndContainerContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventEndContainer(s)
	}
}

func (s *EventEndContainerContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventEndContainer(s)
	}
}

func (p *CEEventParser) EventEndContainer() (localctx IEventEndContainerContext) {
	this := p
	_ = this

	localctx = NewEventEndContainerContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 128, CEEventParserRULE_eventEndContainer)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(629)
		p.Match(CEEventParserEVENT_E)
	}

	return localctx
}

// IEventListContext is an interface to support dynamic dispatch.
type IEventListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventListContext differentiates from other interfaces.
	IsEventListContext()
}

type EventListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventListContext() *EventListContext {
	var p = new(EventListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventList
	return p
}

func (*EventListContext) IsEventListContext() {}

func NewEventListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventListContext {
	var p = new(EventListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventList

	return p
}

func (s *EventListContext) GetParser() antlr.Parser { return s.parser }

func (s *EventListContext) EVENT_L() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_L, 0)
}

func (s *EventListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventList(s)
	}
}

func (s *EventListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventList(s)
	}
}

func (p *CEEventParser) EventList() (localctx IEventListContext) {
	this := p
	_ = this

	localctx = NewEventListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 130, CEEventParserRULE_eventList)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(631)
		p.Match(CEEventParserEVENT_L)
	}

	return localctx
}

// IEventMapContext is an interface to support dynamic dispatch.
type IEventMapContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventMapContext differentiates from other interfaces.
	IsEventMapContext()
}

type EventMapContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventMapContext() *EventMapContext {
	var p = new(EventMapContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventMap
	return p
}

func (*EventMapContext) IsEventMapContext() {}

func NewEventMapContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventMapContext {
	var p = new(EventMapContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventMap

	return p
}

func (s *EventMapContext) GetParser() antlr.Parser { return s.parser }

func (s *EventMapContext) EVENT_M() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_M, 0)
}

func (s *EventMapContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventMapContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventMapContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventMap(s)
	}
}

func (s *EventMapContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventMap(s)
	}
}

func (p *CEEventParser) EventMap() (localctx IEventMapContext) {
	this := p
	_ = this

	localctx = NewEventMapContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 132, CEEventParserRULE_eventMap)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(633)
		p.Match(CEEventParserEVENT_M)
	}

	return localctx
}

// IEventMarkerContext is an interface to support dynamic dispatch.
type IEventMarkerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventMarkerContext differentiates from other interfaces.
	IsEventMarkerContext()
}

type EventMarkerContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventMarkerContext() *EventMarkerContext {
	var p = new(EventMarkerContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventMarker
	return p
}

func (*EventMarkerContext) IsEventMarkerContext() {}

func NewEventMarkerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventMarkerContext {
	var p = new(EventMarkerContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventMarker

	return p
}

func (s *EventMarkerContext) GetParser() antlr.Parser { return s.parser }

func (s *EventMarkerContext) EVENT_MARK() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_MARK, 0)
}

func (s *EventMarkerContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventMarkerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventMarkerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventMarkerContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventMarker(s)
	}
}

func (s *EventMarkerContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventMarker(s)
	}
}

func (p *CEEventParser) EventMarker() (localctx IEventMarkerContext) {
	this := p
	_ = this

	localctx = NewEventMarkerContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 134, CEEventParserRULE_eventMarker)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(635)
		p.Match(CEEventParserEVENT_MARK)
	}
	{
		p.SetState(636)
		p.Match(CEEventParserSTRING)
	}

	return localctx
}

// IEventMediaContext is an interface to support dynamic dispatch.
type IEventMediaContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventMediaContext differentiates from other interfaces.
	IsEventMediaContext()
}

type EventMediaContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventMediaContext() *EventMediaContext {
	var p = new(EventMediaContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventMedia
	return p
}

func (*EventMediaContext) IsEventMediaContext() {}

func NewEventMediaContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventMediaContext {
	var p = new(EventMediaContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventMedia

	return p
}

func (s *EventMediaContext) GetParser() antlr.Parser { return s.parser }

func (s *EventMediaContext) EVENT_MEDIA() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_MEDIA, 0)
}

func (s *EventMediaContext) MEDIA_TYPE() antlr.TerminalNode {
	return s.GetToken(CEEventParserMEDIA_TYPE, 0)
}

func (s *EventMediaContext) AllBYTE() []antlr.TerminalNode {
	return s.GetTokens(CEEventParserBYTE)
}

func (s *EventMediaContext) BYTE(i int) antlr.TerminalNode {
	return s.GetToken(CEEventParserBYTE, i)
}

func (s *EventMediaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventMediaContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventMediaContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventMedia(s)
	}
}

func (s *EventMediaContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventMedia(s)
	}
}

func (p *CEEventParser) EventMedia() (localctx IEventMediaContext) {
	this := p
	_ = this

	localctx = NewEventMediaContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 136, CEEventParserRULE_eventMedia)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(638)
		p.Match(CEEventParserEVENT_MEDIA)
	}
	{
		p.SetState(639)
		p.Match(CEEventParserMEDIA_TYPE)
	}
	p.SetState(643)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CEEventParserBYTE {
		{
			p.SetState(640)
			p.Match(CEEventParserBYTE)
		}

		p.SetState(645)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEventNodeContext is an interface to support dynamic dispatch.
type IEventNodeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventNodeContext differentiates from other interfaces.
	IsEventNodeContext()
}

type EventNodeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventNodeContext() *EventNodeContext {
	var p = new(EventNodeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventNode
	return p
}

func (*EventNodeContext) IsEventNodeContext() {}

func NewEventNodeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventNodeContext {
	var p = new(EventNodeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventNode

	return p
}

func (s *EventNodeContext) GetParser() antlr.Parser { return s.parser }

func (s *EventNodeContext) EVENT_NODE() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_NODE, 0)
}

func (s *EventNodeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventNodeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventNodeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventNode(s)
	}
}

func (s *EventNodeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventNode(s)
	}
}

func (p *CEEventParser) EventNode() (localctx IEventNodeContext) {
	this := p
	_ = this

	localctx = NewEventNodeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 138, CEEventParserRULE_eventNode)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(646)
		p.Match(CEEventParserEVENT_NODE)
	}

	return localctx
}

// IEventNullContext is an interface to support dynamic dispatch.
type IEventNullContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventNullContext differentiates from other interfaces.
	IsEventNullContext()
}

type EventNullContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventNullContext() *EventNullContext {
	var p = new(EventNullContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventNull
	return p
}

func (*EventNullContext) IsEventNullContext() {}

func NewEventNullContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventNullContext {
	var p = new(EventNullContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventNull

	return p
}

func (s *EventNullContext) GetParser() antlr.Parser { return s.parser }

func (s *EventNullContext) EVENT_NULL() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_NULL, 0)
}

func (s *EventNullContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventNullContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventNullContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventNull(s)
	}
}

func (s *EventNullContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventNull(s)
	}
}

func (p *CEEventParser) EventNull() (localctx IEventNullContext) {
	this := p
	_ = this

	localctx = NewEventNullContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 140, CEEventParserRULE_eventNull)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(648)
		p.Match(CEEventParserEVENT_NULL)
	}

	return localctx
}

// IEventNumberContext is an interface to support dynamic dispatch.
type IEventNumberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventNumberContext differentiates from other interfaces.
	IsEventNumberContext()
}

type EventNumberContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventNumberContext() *EventNumberContext {
	var p = new(EventNumberContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventNumber
	return p
}

func (*EventNumberContext) IsEventNumberContext() {}

func NewEventNumberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventNumberContext {
	var p = new(EventNumberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventNumber

	return p
}

func (s *EventNumberContext) GetParser() antlr.Parser { return s.parser }

func (s *EventNumberContext) EVENT_N() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_N, 0)
}

func (s *EventNumberContext) INT_BIN() antlr.TerminalNode {
	return s.GetToken(CEEventParserINT_BIN, 0)
}

func (s *EventNumberContext) INT_OCT() antlr.TerminalNode {
	return s.GetToken(CEEventParserINT_OCT, 0)
}

func (s *EventNumberContext) INT_DEC() antlr.TerminalNode {
	return s.GetToken(CEEventParserINT_DEC, 0)
}

func (s *EventNumberContext) INT_HEX() antlr.TerminalNode {
	return s.GetToken(CEEventParserINT_HEX, 0)
}

func (s *EventNumberContext) FLOAT_DEC() antlr.TerminalNode {
	return s.GetToken(CEEventParserFLOAT_DEC, 0)
}

func (s *EventNumberContext) FLOAT_HEX() antlr.TerminalNode {
	return s.GetToken(CEEventParserFLOAT_HEX, 0)
}

func (s *EventNumberContext) FLOAT_INF() antlr.TerminalNode {
	return s.GetToken(CEEventParserFLOAT_INF, 0)
}

func (s *EventNumberContext) FLOAT_NAN() antlr.TerminalNode {
	return s.GetToken(CEEventParserFLOAT_NAN, 0)
}

func (s *EventNumberContext) FLOAT_SNAN() antlr.TerminalNode {
	return s.GetToken(CEEventParserFLOAT_SNAN, 0)
}

func (s *EventNumberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventNumberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventNumberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventNumber(s)
	}
}

func (s *EventNumberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventNumber(s)
	}
}

func (p *CEEventParser) EventNumber() (localctx IEventNumberContext) {
	this := p
	_ = this

	localctx = NewEventNumberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 142, CEEventParserRULE_eventNumber)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(650)
		p.Match(CEEventParserEVENT_N)
	}
	{
		p.SetState(651)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-105)&-(0x1f+1)) == 0 && ((1<<uint((_la-105)))&((1<<(CEEventParserFLOAT_NAN-105))|(1<<(CEEventParserFLOAT_SNAN-105))|(1<<(CEEventParserFLOAT_INF-105))|(1<<(CEEventParserFLOAT_DEC-105))|(1<<(CEEventParserFLOAT_HEX-105))|(1<<(CEEventParserINT_BIN-105))|(1<<(CEEventParserINT_OCT-105))|(1<<(CEEventParserINT_DEC-105))|(1<<(CEEventParserINT_HEX-105)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IEventPadContext is an interface to support dynamic dispatch.
type IEventPadContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventPadContext differentiates from other interfaces.
	IsEventPadContext()
}

type EventPadContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventPadContext() *EventPadContext {
	var p = new(EventPadContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventPad
	return p
}

func (*EventPadContext) IsEventPadContext() {}

func NewEventPadContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventPadContext {
	var p = new(EventPadContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventPad

	return p
}

func (s *EventPadContext) GetParser() antlr.Parser { return s.parser }

func (s *EventPadContext) EVENT_PAD() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_PAD, 0)
}

func (s *EventPadContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventPadContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventPadContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventPad(s)
	}
}

func (s *EventPadContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventPad(s)
	}
}

func (p *CEEventParser) EventPad() (localctx IEventPadContext) {
	this := p
	_ = this

	localctx = NewEventPadContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 144, CEEventParserRULE_eventPad)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(653)
		p.Match(CEEventParserEVENT_PAD)
	}

	return localctx
}

// IEventLocalReferenceContext is an interface to support dynamic dispatch.
type IEventLocalReferenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventLocalReferenceContext differentiates from other interfaces.
	IsEventLocalReferenceContext()
}

type EventLocalReferenceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventLocalReferenceContext() *EventLocalReferenceContext {
	var p = new(EventLocalReferenceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventLocalReference
	return p
}

func (*EventLocalReferenceContext) IsEventLocalReferenceContext() {}

func NewEventLocalReferenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventLocalReferenceContext {
	var p = new(EventLocalReferenceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventLocalReference

	return p
}

func (s *EventLocalReferenceContext) GetParser() antlr.Parser { return s.parser }

func (s *EventLocalReferenceContext) EVENT_REFL() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_REFL, 0)
}

func (s *EventLocalReferenceContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventLocalReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventLocalReferenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventLocalReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventLocalReference(s)
	}
}

func (s *EventLocalReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventLocalReference(s)
	}
}

func (p *CEEventParser) EventLocalReference() (localctx IEventLocalReferenceContext) {
	this := p
	_ = this

	localctx = NewEventLocalReferenceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 146, CEEventParserRULE_eventLocalReference)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(655)
		p.Match(CEEventParserEVENT_REFL)
	}
	{
		p.SetState(656)
		p.Match(CEEventParserSTRING)
	}

	return localctx
}

// IEventRemoteReferenceContext is an interface to support dynamic dispatch.
type IEventRemoteReferenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventRemoteReferenceContext differentiates from other interfaces.
	IsEventRemoteReferenceContext()
}

type EventRemoteReferenceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventRemoteReferenceContext() *EventRemoteReferenceContext {
	var p = new(EventRemoteReferenceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventRemoteReference
	return p
}

func (*EventRemoteReferenceContext) IsEventRemoteReferenceContext() {}

func NewEventRemoteReferenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventRemoteReferenceContext {
	var p = new(EventRemoteReferenceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventRemoteReference

	return p
}

func (s *EventRemoteReferenceContext) GetParser() antlr.Parser { return s.parser }

func (s *EventRemoteReferenceContext) EVENT_REFR() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_REFR, 0)
}

func (s *EventRemoteReferenceContext) EVENT_REFR_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_REFR_ARGS, 0)
}

func (s *EventRemoteReferenceContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventRemoteReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventRemoteReferenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventRemoteReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventRemoteReference(s)
	}
}

func (s *EventRemoteReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventRemoteReference(s)
	}
}

func (p *CEEventParser) EventRemoteReference() (localctx IEventRemoteReferenceContext) {
	this := p
	_ = this

	localctx = NewEventRemoteReferenceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 148, CEEventParserRULE_eventRemoteReference)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(663)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_REFR:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(658)
			p.Match(CEEventParserEVENT_REFR)
		}

	case CEEventParserEVENT_REFR_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(659)
			p.Match(CEEventParserEVENT_REFR_ARGS)
		}
		p.SetState(661)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CEEventParserSTRING {
			{
				p.SetState(660)
				p.Match(CEEventParserSTRING)
			}

		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventResourceIdContext is an interface to support dynamic dispatch.
type IEventResourceIdContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventResourceIdContext differentiates from other interfaces.
	IsEventResourceIdContext()
}

type EventResourceIdContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventResourceIdContext() *EventResourceIdContext {
	var p = new(EventResourceIdContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventResourceId
	return p
}

func (*EventResourceIdContext) IsEventResourceIdContext() {}

func NewEventResourceIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventResourceIdContext {
	var p = new(EventResourceIdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventResourceId

	return p
}

func (s *EventResourceIdContext) GetParser() antlr.Parser { return s.parser }

func (s *EventResourceIdContext) EVENT_RID() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_RID, 0)
}

func (s *EventResourceIdContext) EVENT_RID_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_RID_ARGS, 0)
}

func (s *EventResourceIdContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventResourceIdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventResourceIdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventResourceIdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventResourceId(s)
	}
}

func (s *EventResourceIdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventResourceId(s)
	}
}

func (p *CEEventParser) EventResourceId() (localctx IEventResourceIdContext) {
	this := p
	_ = this

	localctx = NewEventResourceIdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 150, CEEventParserRULE_eventResourceId)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(670)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_RID:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(665)
			p.Match(CEEventParserEVENT_RID)
		}

	case CEEventParserEVENT_RID_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(666)
			p.Match(CEEventParserEVENT_RID_ARGS)
		}
		p.SetState(668)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CEEventParserSTRING {
			{
				p.SetState(667)
				p.Match(CEEventParserSTRING)
			}

		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventStringContext is an interface to support dynamic dispatch.
type IEventStringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventStringContext differentiates from other interfaces.
	IsEventStringContext()
}

type EventStringContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventStringContext() *EventStringContext {
	var p = new(EventStringContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventString
	return p
}

func (*EventStringContext) IsEventStringContext() {}

func NewEventStringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventStringContext {
	var p = new(EventStringContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventString

	return p
}

func (s *EventStringContext) GetParser() antlr.Parser { return s.parser }

func (s *EventStringContext) EVENT_S() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_S, 0)
}

func (s *EventStringContext) EVENT_S_ARGS() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_S_ARGS, 0)
}

func (s *EventStringContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventStringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventStringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventStringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventString(s)
	}
}

func (s *EventStringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventString(s)
	}
}

func (p *CEEventParser) EventString() (localctx IEventStringContext) {
	this := p
	_ = this

	localctx = NewEventStringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 152, CEEventParserRULE_eventString)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(677)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CEEventParserEVENT_S:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(672)
			p.Match(CEEventParserEVENT_S)
		}

	case CEEventParserEVENT_S_ARGS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(673)
			p.Match(CEEventParserEVENT_S_ARGS)
		}
		p.SetState(675)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CEEventParserSTRING {
			{
				p.SetState(674)
				p.Match(CEEventParserSTRING)
			}

		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventStructInstanceContext is an interface to support dynamic dispatch.
type IEventStructInstanceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventStructInstanceContext differentiates from other interfaces.
	IsEventStructInstanceContext()
}

type EventStructInstanceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventStructInstanceContext() *EventStructInstanceContext {
	var p = new(EventStructInstanceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventStructInstance
	return p
}

func (*EventStructInstanceContext) IsEventStructInstanceContext() {}

func NewEventStructInstanceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventStructInstanceContext {
	var p = new(EventStructInstanceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventStructInstance

	return p
}

func (s *EventStructInstanceContext) GetParser() antlr.Parser { return s.parser }

func (s *EventStructInstanceContext) EVENT_SI() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_SI, 0)
}

func (s *EventStructInstanceContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventStructInstanceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventStructInstanceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventStructInstanceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventStructInstance(s)
	}
}

func (s *EventStructInstanceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventStructInstance(s)
	}
}

func (p *CEEventParser) EventStructInstance() (localctx IEventStructInstanceContext) {
	this := p
	_ = this

	localctx = NewEventStructInstanceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 154, CEEventParserRULE_eventStructInstance)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(679)
		p.Match(CEEventParserEVENT_SI)
	}
	{
		p.SetState(680)
		p.Match(CEEventParserSTRING)
	}

	return localctx
}

// IEventStructTemplateContext is an interface to support dynamic dispatch.
type IEventStructTemplateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventStructTemplateContext differentiates from other interfaces.
	IsEventStructTemplateContext()
}

type EventStructTemplateContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventStructTemplateContext() *EventStructTemplateContext {
	var p = new(EventStructTemplateContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventStructTemplate
	return p
}

func (*EventStructTemplateContext) IsEventStructTemplateContext() {}

func NewEventStructTemplateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventStructTemplateContext {
	var p = new(EventStructTemplateContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventStructTemplate

	return p
}

func (s *EventStructTemplateContext) GetParser() antlr.Parser { return s.parser }

func (s *EventStructTemplateContext) EVENT_ST() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_ST, 0)
}

func (s *EventStructTemplateContext) STRING() antlr.TerminalNode {
	return s.GetToken(CEEventParserSTRING, 0)
}

func (s *EventStructTemplateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventStructTemplateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventStructTemplateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventStructTemplate(s)
	}
}

func (s *EventStructTemplateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventStructTemplate(s)
	}
}

func (p *CEEventParser) EventStructTemplate() (localctx IEventStructTemplateContext) {
	this := p
	_ = this

	localctx = NewEventStructTemplateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 156, CEEventParserRULE_eventStructTemplate)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(682)
		p.Match(CEEventParserEVENT_ST)
	}
	{
		p.SetState(683)
		p.Match(CEEventParserSTRING)
	}

	return localctx
}

// IEventTimeContext is an interface to support dynamic dispatch.
type IEventTimeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventTimeContext differentiates from other interfaces.
	IsEventTimeContext()
}

type EventTimeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventTimeContext() *EventTimeContext {
	var p = new(EventTimeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventTime
	return p
}

func (*EventTimeContext) IsEventTimeContext() {}

func NewEventTimeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventTimeContext {
	var p = new(EventTimeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventTime

	return p
}

func (s *EventTimeContext) GetParser() antlr.Parser { return s.parser }

func (s *EventTimeContext) EVENT_T() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_T, 0)
}

func (s *EventTimeContext) DATE() antlr.TerminalNode {
	return s.GetToken(CEEventParserDATE, 0)
}

func (s *EventTimeContext) TIME() antlr.TerminalNode {
	return s.GetToken(CEEventParserTIME, 0)
}

func (s *EventTimeContext) DATETIME() antlr.TerminalNode {
	return s.GetToken(CEEventParserDATETIME, 0)
}

func (s *EventTimeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventTimeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventTimeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventTime(s)
	}
}

func (s *EventTimeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventTime(s)
	}
}

func (p *CEEventParser) EventTime() (localctx IEventTimeContext) {
	this := p
	_ = this

	localctx = NewEventTimeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 158, CEEventParserRULE_eventTime)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(685)
		p.Match(CEEventParserEVENT_T)
	}
	{
		p.SetState(686)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-141)&-(0x1f+1)) == 0 && ((1<<uint((_la-141)))&((1<<(CEEventParserTIME-141))|(1<<(CEEventParserDATE-141))|(1<<(CEEventParserDATETIME-141)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IEventUIDContext is an interface to support dynamic dispatch.
type IEventUIDContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventUIDContext differentiates from other interfaces.
	IsEventUIDContext()
}

type EventUIDContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventUIDContext() *EventUIDContext {
	var p = new(EventUIDContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventUID
	return p
}

func (*EventUIDContext) IsEventUIDContext() {}

func NewEventUIDContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventUIDContext {
	var p = new(EventUIDContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventUID

	return p
}

func (s *EventUIDContext) GetParser() antlr.Parser { return s.parser }

func (s *EventUIDContext) EVENT_UID() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_UID, 0)
}

func (s *EventUIDContext) UID() antlr.TerminalNode {
	return s.GetToken(CEEventParserUID, 0)
}

func (s *EventUIDContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventUIDContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventUIDContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventUID(s)
	}
}

func (s *EventUIDContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventUID(s)
	}
}

func (p *CEEventParser) EventUID() (localctx IEventUIDContext) {
	this := p
	_ = this

	localctx = NewEventUIDContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 160, CEEventParserRULE_eventUID)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(688)
		p.Match(CEEventParserEVENT_UID)
	}
	{
		p.SetState(689)
		p.Match(CEEventParserUID)
	}

	return localctx
}

// IEventVersionContext is an interface to support dynamic dispatch.
type IEventVersionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventVersionContext differentiates from other interfaces.
	IsEventVersionContext()
}

type EventVersionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventVersionContext() *EventVersionContext {
	var p = new(EventVersionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CEEventParserRULE_eventVersion
	return p
}

func (*EventVersionContext) IsEventVersionContext() {}

func NewEventVersionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventVersionContext {
	var p = new(EventVersionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CEEventParserRULE_eventVersion

	return p
}

func (s *EventVersionContext) GetParser() antlr.Parser { return s.parser }

func (s *EventVersionContext) EVENT_V() antlr.TerminalNode {
	return s.GetToken(CEEventParserEVENT_V, 0)
}

func (s *EventVersionContext) VALUE_UINT_DEC() antlr.TerminalNode {
	return s.GetToken(CEEventParserVALUE_UINT_DEC, 0)
}

func (s *EventVersionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventVersionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventVersionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.EnterEventVersion(s)
	}
}

func (s *EventVersionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CEEventParserListener); ok {
		listenerT.ExitEventVersion(s)
	}
}

func (p *CEEventParser) EventVersion() (localctx IEventVersionContext) {
	this := p
	_ = this

	localctx = NewEventVersionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 162, CEEventParserRULE_eventVersion)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(691)
		p.Match(CEEventParserEVENT_V)
	}
	{
		p.SetState(692)
		p.Match(CEEventParserVALUE_UINT_DEC)
	}

	return localctx
}
