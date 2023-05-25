// Code generated from /home/karl/Projects/go-concise-encoding/codegen/cte/CTEParser.g4 by ANTLR 4.12.0. DO NOT EDIT.

package parser // CTEParser

import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// CTEParserListener is a complete listener for a parse tree produced by CTEParser.
type CTEParserListener interface {
	antlr.ParseTreeListener

	// EnterCte is called when entering the cte production.
	EnterCte(c *CteContext)

	// EnterVersion is called when entering the version production.
	EnterVersion(c *VersionContext)

	// EnterRecordTypes is called when entering the recordTypes production.
	EnterRecordTypes(c *RecordTypesContext)

	// EnterValue is called when entering the value production.
	EnterValue(c *ValueContext)

	// EnterSeparator is called when entering the separator production.
	EnterSeparator(c *SeparatorContext)

	// EnterCommentLine is called when entering the commentLine production.
	EnterCommentLine(c *CommentLineContext)

	// EnterCommentBlock is called when entering the commentBlock production.
	EnterCommentBlock(c *CommentBlockContext)

	// EnterValueNull is called when entering the valueNull production.
	EnterValueNull(c *ValueNullContext)

	// EnterValueUid is called when entering the valueUid production.
	EnterValueUid(c *ValueUidContext)

	// EnterValueBool is called when entering the valueBool production.
	EnterValueBool(c *ValueBoolContext)

	// EnterValueInt is called when entering the valueInt production.
	EnterValueInt(c *ValueIntContext)

	// EnterValueFloat is called when entering the valueFloat production.
	EnterValueFloat(c *ValueFloatContext)

	// EnterValueInf is called when entering the valueInf production.
	EnterValueInf(c *ValueInfContext)

	// EnterValueNinf is called when entering the valueNinf production.
	EnterValueNinf(c *ValueNinfContext)

	// EnterValueNan is called when entering the valueNan production.
	EnterValueNan(c *ValueNanContext)

	// EnterValueSnan is called when entering the valueSnan production.
	EnterValueSnan(c *ValueSnanContext)

	// EnterValueDate is called when entering the valueDate production.
	EnterValueDate(c *ValueDateContext)

	// EnterValueTime is called when entering the valueTime production.
	EnterValueTime(c *ValueTimeContext)

	// EnterValueString is called when entering the valueString production.
	EnterValueString(c *ValueStringContext)

	// EnterStringContents is called when entering the stringContents production.
	EnterStringContents(c *StringContentsContext)

	// EnterStringEscape is called when entering the stringEscape production.
	EnterStringEscape(c *StringEscapeContext)

	// EnterVerbatimSequence is called when entering the verbatimSequence production.
	EnterVerbatimSequence(c *VerbatimSequenceContext)

	// EnterVerbatimContents is called when entering the verbatimContents production.
	EnterVerbatimContents(c *VerbatimContentsContext)

	// EnterCodepointSequence is called when entering the codepointSequence production.
	EnterCodepointSequence(c *CodepointSequenceContext)

	// EnterCodepointContents is called when entering the codepointContents production.
	EnterCodepointContents(c *CodepointContentsContext)

	// EnterEscapeChar is called when entering the escapeChar production.
	EnterEscapeChar(c *EscapeCharContext)

	// EnterCustomText is called when entering the customText production.
	EnterCustomText(c *CustomTextContext)

	// EnterCustomEscape is called when entering the customEscape production.
	EnterCustomEscape(c *CustomEscapeContext)

	// EnterCustomBinary is called when entering the customBinary production.
	EnterCustomBinary(c *CustomBinaryContext)

	// EnterCustomType is called when entering the customType production.
	EnterCustomType(c *CustomTypeContext)

	// EnterMediaText is called when entering the mediaText production.
	EnterMediaText(c *MediaTextContext)

	// EnterMediaEscape is called when entering the mediaEscape production.
	EnterMediaEscape(c *MediaEscapeContext)

	// EnterMediaBinary is called when entering the mediaBinary production.
	EnterMediaBinary(c *MediaBinaryContext)

	// EnterMediaType is called when entering the mediaType production.
	EnterMediaType(c *MediaTypeContext)

	// EnterValueRid is called when entering the valueRid production.
	EnterValueRid(c *ValueRidContext)

	// EnterValueRemoteRef is called when entering the valueRemoteRef production.
	EnterValueRemoteRef(c *ValueRemoteRefContext)

	// EnterMarkerID is called when entering the markerID production.
	EnterMarkerID(c *MarkerIDContext)

	// EnterMarker is called when entering the marker production.
	EnterMarker(c *MarkerContext)

	// EnterReference is called when entering the reference production.
	EnterReference(c *ReferenceContext)

	// EnterContainerMap is called when entering the containerMap production.
	EnterContainerMap(c *ContainerMapContext)

	// EnterContainerList is called when entering the containerList production.
	EnterContainerList(c *ContainerListContext)

	// EnterContainerRecordType is called when entering the containerRecordType production.
	EnterContainerRecordType(c *ContainerRecordTypeContext)

	// EnterContainerRecord is called when entering the containerRecord production.
	EnterContainerRecord(c *ContainerRecordContext)

	// EnterContainerNode is called when entering the containerNode production.
	EnterContainerNode(c *ContainerNodeContext)

	// EnterContainerEdge is called when entering the containerEdge production.
	EnterContainerEdge(c *ContainerEdgeContext)

	// EnterKvPair is called when entering the kvPair production.
	EnterKvPair(c *KvPairContext)

	// EnterRecordTypeBegin is called when entering the recordTypeBegin production.
	EnterRecordTypeBegin(c *RecordTypeBeginContext)

	// EnterRecordBegin is called when entering the recordBegin production.
	EnterRecordBegin(c *RecordBeginContext)

	// EnterArrayElemInt is called when entering the arrayElemInt production.
	EnterArrayElemInt(c *ArrayElemIntContext)

	// EnterArrayElemIntB is called when entering the arrayElemIntB production.
	EnterArrayElemIntB(c *ArrayElemIntBContext)

	// EnterArrayElemIntO is called when entering the arrayElemIntO production.
	EnterArrayElemIntO(c *ArrayElemIntOContext)

	// EnterArrayElemIntX is called when entering the arrayElemIntX production.
	EnterArrayElemIntX(c *ArrayElemIntXContext)

	// EnterArrayElemUint is called when entering the arrayElemUint production.
	EnterArrayElemUint(c *ArrayElemUintContext)

	// EnterArrayElemUintB is called when entering the arrayElemUintB production.
	EnterArrayElemUintB(c *ArrayElemUintBContext)

	// EnterArrayElemUintO is called when entering the arrayElemUintO production.
	EnterArrayElemUintO(c *ArrayElemUintOContext)

	// EnterArrayElemUintX is called when entering the arrayElemUintX production.
	EnterArrayElemUintX(c *ArrayElemUintXContext)

	// EnterArrayElemFloat is called when entering the arrayElemFloat production.
	EnterArrayElemFloat(c *ArrayElemFloatContext)

	// EnterArrayElemFloatX is called when entering the arrayElemFloatX production.
	EnterArrayElemFloatX(c *ArrayElemFloatXContext)

	// EnterArrayElemNan is called when entering the arrayElemNan production.
	EnterArrayElemNan(c *ArrayElemNanContext)

	// EnterArrayElemSnan is called when entering the arrayElemSnan production.
	EnterArrayElemSnan(c *ArrayElemSnanContext)

	// EnterArrayElemInf is called when entering the arrayElemInf production.
	EnterArrayElemInf(c *ArrayElemInfContext)

	// EnterArrayElemNinf is called when entering the arrayElemNinf production.
	EnterArrayElemNinf(c *ArrayElemNinfContext)

	// EnterArrayElemUid is called when entering the arrayElemUid production.
	EnterArrayElemUid(c *ArrayElemUidContext)

	// EnterArrayElemBits is called when entering the arrayElemBits production.
	EnterArrayElemBits(c *ArrayElemBitsContext)

	// EnterArrayElemByteX is called when entering the arrayElemByteX production.
	EnterArrayElemByteX(c *ArrayElemByteXContext)

	// EnterArrayI8 is called when entering the arrayI8 production.
	EnterArrayI8(c *ArrayI8Context)

	// EnterArrayI16 is called when entering the arrayI16 production.
	EnterArrayI16(c *ArrayI16Context)

	// EnterArrayI32 is called when entering the arrayI32 production.
	EnterArrayI32(c *ArrayI32Context)

	// EnterArrayI64 is called when entering the arrayI64 production.
	EnterArrayI64(c *ArrayI64Context)

	// EnterArrayU8 is called when entering the arrayU8 production.
	EnterArrayU8(c *ArrayU8Context)

	// EnterArrayU16 is called when entering the arrayU16 production.
	EnterArrayU16(c *ArrayU16Context)

	// EnterArrayU32 is called when entering the arrayU32 production.
	EnterArrayU32(c *ArrayU32Context)

	// EnterArrayU64 is called when entering the arrayU64 production.
	EnterArrayU64(c *ArrayU64Context)

	// EnterArrayF16 is called when entering the arrayF16 production.
	EnterArrayF16(c *ArrayF16Context)

	// EnterArrayF32 is called when entering the arrayF32 production.
	EnterArrayF32(c *ArrayF32Context)

	// EnterArrayF64 is called when entering the arrayF64 production.
	EnterArrayF64(c *ArrayF64Context)

	// EnterArrayUid is called when entering the arrayUid production.
	EnterArrayUid(c *ArrayUidContext)

	// EnterArrayBit is called when entering the arrayBit production.
	EnterArrayBit(c *ArrayBitContext)

	// ExitCte is called when exiting the cte production.
	ExitCte(c *CteContext)

	// ExitVersion is called when exiting the version production.
	ExitVersion(c *VersionContext)

	// ExitRecordTypes is called when exiting the recordTypes production.
	ExitRecordTypes(c *RecordTypesContext)

	// ExitValue is called when exiting the value production.
	ExitValue(c *ValueContext)

	// ExitSeparator is called when exiting the separator production.
	ExitSeparator(c *SeparatorContext)

	// ExitCommentLine is called when exiting the commentLine production.
	ExitCommentLine(c *CommentLineContext)

	// ExitCommentBlock is called when exiting the commentBlock production.
	ExitCommentBlock(c *CommentBlockContext)

	// ExitValueNull is called when exiting the valueNull production.
	ExitValueNull(c *ValueNullContext)

	// ExitValueUid is called when exiting the valueUid production.
	ExitValueUid(c *ValueUidContext)

	// ExitValueBool is called when exiting the valueBool production.
	ExitValueBool(c *ValueBoolContext)

	// ExitValueInt is called when exiting the valueInt production.
	ExitValueInt(c *ValueIntContext)

	// ExitValueFloat is called when exiting the valueFloat production.
	ExitValueFloat(c *ValueFloatContext)

	// ExitValueInf is called when exiting the valueInf production.
	ExitValueInf(c *ValueInfContext)

	// ExitValueNinf is called when exiting the valueNinf production.
	ExitValueNinf(c *ValueNinfContext)

	// ExitValueNan is called when exiting the valueNan production.
	ExitValueNan(c *ValueNanContext)

	// ExitValueSnan is called when exiting the valueSnan production.
	ExitValueSnan(c *ValueSnanContext)

	// ExitValueDate is called when exiting the valueDate production.
	ExitValueDate(c *ValueDateContext)

	// ExitValueTime is called when exiting the valueTime production.
	ExitValueTime(c *ValueTimeContext)

	// ExitValueString is called when exiting the valueString production.
	ExitValueString(c *ValueStringContext)

	// ExitStringContents is called when exiting the stringContents production.
	ExitStringContents(c *StringContentsContext)

	// ExitStringEscape is called when exiting the stringEscape production.
	ExitStringEscape(c *StringEscapeContext)

	// ExitVerbatimSequence is called when exiting the verbatimSequence production.
	ExitVerbatimSequence(c *VerbatimSequenceContext)

	// ExitVerbatimContents is called when exiting the verbatimContents production.
	ExitVerbatimContents(c *VerbatimContentsContext)

	// ExitCodepointSequence is called when exiting the codepointSequence production.
	ExitCodepointSequence(c *CodepointSequenceContext)

	// ExitCodepointContents is called when exiting the codepointContents production.
	ExitCodepointContents(c *CodepointContentsContext)

	// ExitEscapeChar is called when exiting the escapeChar production.
	ExitEscapeChar(c *EscapeCharContext)

	// ExitCustomText is called when exiting the customText production.
	ExitCustomText(c *CustomTextContext)

	// ExitCustomEscape is called when exiting the customEscape production.
	ExitCustomEscape(c *CustomEscapeContext)

	// ExitCustomBinary is called when exiting the customBinary production.
	ExitCustomBinary(c *CustomBinaryContext)

	// ExitCustomType is called when exiting the customType production.
	ExitCustomType(c *CustomTypeContext)

	// ExitMediaText is called when exiting the mediaText production.
	ExitMediaText(c *MediaTextContext)

	// ExitMediaEscape is called when exiting the mediaEscape production.
	ExitMediaEscape(c *MediaEscapeContext)

	// ExitMediaBinary is called when exiting the mediaBinary production.
	ExitMediaBinary(c *MediaBinaryContext)

	// ExitMediaType is called when exiting the mediaType production.
	ExitMediaType(c *MediaTypeContext)

	// ExitValueRid is called when exiting the valueRid production.
	ExitValueRid(c *ValueRidContext)

	// ExitValueRemoteRef is called when exiting the valueRemoteRef production.
	ExitValueRemoteRef(c *ValueRemoteRefContext)

	// ExitMarkerID is called when exiting the markerID production.
	ExitMarkerID(c *MarkerIDContext)

	// ExitMarker is called when exiting the marker production.
	ExitMarker(c *MarkerContext)

	// ExitReference is called when exiting the reference production.
	ExitReference(c *ReferenceContext)

	// ExitContainerMap is called when exiting the containerMap production.
	ExitContainerMap(c *ContainerMapContext)

	// ExitContainerList is called when exiting the containerList production.
	ExitContainerList(c *ContainerListContext)

	// ExitContainerRecordType is called when exiting the containerRecordType production.
	ExitContainerRecordType(c *ContainerRecordTypeContext)

	// ExitContainerRecord is called when exiting the containerRecord production.
	ExitContainerRecord(c *ContainerRecordContext)

	// ExitContainerNode is called when exiting the containerNode production.
	ExitContainerNode(c *ContainerNodeContext)

	// ExitContainerEdge is called when exiting the containerEdge production.
	ExitContainerEdge(c *ContainerEdgeContext)

	// ExitKvPair is called when exiting the kvPair production.
	ExitKvPair(c *KvPairContext)

	// ExitRecordTypeBegin is called when exiting the recordTypeBegin production.
	ExitRecordTypeBegin(c *RecordTypeBeginContext)

	// ExitRecordBegin is called when exiting the recordBegin production.
	ExitRecordBegin(c *RecordBeginContext)

	// ExitArrayElemInt is called when exiting the arrayElemInt production.
	ExitArrayElemInt(c *ArrayElemIntContext)

	// ExitArrayElemIntB is called when exiting the arrayElemIntB production.
	ExitArrayElemIntB(c *ArrayElemIntBContext)

	// ExitArrayElemIntO is called when exiting the arrayElemIntO production.
	ExitArrayElemIntO(c *ArrayElemIntOContext)

	// ExitArrayElemIntX is called when exiting the arrayElemIntX production.
	ExitArrayElemIntX(c *ArrayElemIntXContext)

	// ExitArrayElemUint is called when exiting the arrayElemUint production.
	ExitArrayElemUint(c *ArrayElemUintContext)

	// ExitArrayElemUintB is called when exiting the arrayElemUintB production.
	ExitArrayElemUintB(c *ArrayElemUintBContext)

	// ExitArrayElemUintO is called when exiting the arrayElemUintO production.
	ExitArrayElemUintO(c *ArrayElemUintOContext)

	// ExitArrayElemUintX is called when exiting the arrayElemUintX production.
	ExitArrayElemUintX(c *ArrayElemUintXContext)

	// ExitArrayElemFloat is called when exiting the arrayElemFloat production.
	ExitArrayElemFloat(c *ArrayElemFloatContext)

	// ExitArrayElemFloatX is called when exiting the arrayElemFloatX production.
	ExitArrayElemFloatX(c *ArrayElemFloatXContext)

	// ExitArrayElemNan is called when exiting the arrayElemNan production.
	ExitArrayElemNan(c *ArrayElemNanContext)

	// ExitArrayElemSnan is called when exiting the arrayElemSnan production.
	ExitArrayElemSnan(c *ArrayElemSnanContext)

	// ExitArrayElemInf is called when exiting the arrayElemInf production.
	ExitArrayElemInf(c *ArrayElemInfContext)

	// ExitArrayElemNinf is called when exiting the arrayElemNinf production.
	ExitArrayElemNinf(c *ArrayElemNinfContext)

	// ExitArrayElemUid is called when exiting the arrayElemUid production.
	ExitArrayElemUid(c *ArrayElemUidContext)

	// ExitArrayElemBits is called when exiting the arrayElemBits production.
	ExitArrayElemBits(c *ArrayElemBitsContext)

	// ExitArrayElemByteX is called when exiting the arrayElemByteX production.
	ExitArrayElemByteX(c *ArrayElemByteXContext)

	// ExitArrayI8 is called when exiting the arrayI8 production.
	ExitArrayI8(c *ArrayI8Context)

	// ExitArrayI16 is called when exiting the arrayI16 production.
	ExitArrayI16(c *ArrayI16Context)

	// ExitArrayI32 is called when exiting the arrayI32 production.
	ExitArrayI32(c *ArrayI32Context)

	// ExitArrayI64 is called when exiting the arrayI64 production.
	ExitArrayI64(c *ArrayI64Context)

	// ExitArrayU8 is called when exiting the arrayU8 production.
	ExitArrayU8(c *ArrayU8Context)

	// ExitArrayU16 is called when exiting the arrayU16 production.
	ExitArrayU16(c *ArrayU16Context)

	// ExitArrayU32 is called when exiting the arrayU32 production.
	ExitArrayU32(c *ArrayU32Context)

	// ExitArrayU64 is called when exiting the arrayU64 production.
	ExitArrayU64(c *ArrayU64Context)

	// ExitArrayF16 is called when exiting the arrayF16 production.
	ExitArrayF16(c *ArrayF16Context)

	// ExitArrayF32 is called when exiting the arrayF32 production.
	ExitArrayF32(c *ArrayF32Context)

	// ExitArrayF64 is called when exiting the arrayF64 production.
	ExitArrayF64(c *ArrayF64Context)

	// ExitArrayUid is called when exiting the arrayUid production.
	ExitArrayUid(c *ArrayUidContext)

	// ExitArrayBit is called when exiting the arrayBit production.
	ExitArrayBit(c *ArrayBitContext)
}
