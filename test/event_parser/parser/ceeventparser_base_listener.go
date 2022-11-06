// Code generated from /home/karl/Projects/go-concise-encoding/codegen/test/CEEventParser.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // CEEventParser

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseCEEventParserListener is a complete listener for a parse tree produced by CEEventParser.
type BaseCEEventParserListener struct{}

var _ CEEventParserListener = &BaseCEEventParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseCEEventParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseCEEventParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseCEEventParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseCEEventParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BaseCEEventParserListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseCEEventParserListener) ExitStart(ctx *StartContext) {}

// EnterEvent is called when production event is entered.
func (s *BaseCEEventParserListener) EnterEvent(ctx *EventContext) {}

// ExitEvent is called when production event is exited.
func (s *BaseCEEventParserListener) ExitEvent(ctx *EventContext) {}

// EnterEventArrayBits is called when production eventArrayBits is entered.
func (s *BaseCEEventParserListener) EnterEventArrayBits(ctx *EventArrayBitsContext) {}

// ExitEventArrayBits is called when production eventArrayBits is exited.
func (s *BaseCEEventParserListener) ExitEventArrayBits(ctx *EventArrayBitsContext) {}

// EnterEventArrayChunkLast is called when production eventArrayChunkLast is entered.
func (s *BaseCEEventParserListener) EnterEventArrayChunkLast(ctx *EventArrayChunkLastContext) {}

// ExitEventArrayChunkLast is called when production eventArrayChunkLast is exited.
func (s *BaseCEEventParserListener) ExitEventArrayChunkLast(ctx *EventArrayChunkLastContext) {}

// EnterEventArrayChunkMore is called when production eventArrayChunkMore is entered.
func (s *BaseCEEventParserListener) EnterEventArrayChunkMore(ctx *EventArrayChunkMoreContext) {}

// ExitEventArrayChunkMore is called when production eventArrayChunkMore is exited.
func (s *BaseCEEventParserListener) ExitEventArrayChunkMore(ctx *EventArrayChunkMoreContext) {}

// EnterEventArrayDataBits is called when production eventArrayDataBits is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataBits(ctx *EventArrayDataBitsContext) {}

// ExitEventArrayDataBits is called when production eventArrayDataBits is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataBits(ctx *EventArrayDataBitsContext) {}

// EnterEventArrayDataFloat16 is called when production eventArrayDataFloat16 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataFloat16(ctx *EventArrayDataFloat16Context) {}

// ExitEventArrayDataFloat16 is called when production eventArrayDataFloat16 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataFloat16(ctx *EventArrayDataFloat16Context) {}

// EnterEventArrayDataFloat32 is called when production eventArrayDataFloat32 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataFloat32(ctx *EventArrayDataFloat32Context) {}

// ExitEventArrayDataFloat32 is called when production eventArrayDataFloat32 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataFloat32(ctx *EventArrayDataFloat32Context) {}

// EnterEventArrayDataFloat64 is called when production eventArrayDataFloat64 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataFloat64(ctx *EventArrayDataFloat64Context) {}

// ExitEventArrayDataFloat64 is called when production eventArrayDataFloat64 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataFloat64(ctx *EventArrayDataFloat64Context) {}

// EnterEventArrayDataInt16 is called when production eventArrayDataInt16 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataInt16(ctx *EventArrayDataInt16Context) {}

// ExitEventArrayDataInt16 is called when production eventArrayDataInt16 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataInt16(ctx *EventArrayDataInt16Context) {}

// EnterEventArrayDataInt32 is called when production eventArrayDataInt32 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataInt32(ctx *EventArrayDataInt32Context) {}

// ExitEventArrayDataInt32 is called when production eventArrayDataInt32 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataInt32(ctx *EventArrayDataInt32Context) {}

// EnterEventArrayDataInt64 is called when production eventArrayDataInt64 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataInt64(ctx *EventArrayDataInt64Context) {}

// ExitEventArrayDataInt64 is called when production eventArrayDataInt64 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataInt64(ctx *EventArrayDataInt64Context) {}

// EnterEventArrayDataInt8 is called when production eventArrayDataInt8 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataInt8(ctx *EventArrayDataInt8Context) {}

// ExitEventArrayDataInt8 is called when production eventArrayDataInt8 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataInt8(ctx *EventArrayDataInt8Context) {}

// EnterEventArrayDataText is called when production eventArrayDataText is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataText(ctx *EventArrayDataTextContext) {}

// ExitEventArrayDataText is called when production eventArrayDataText is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataText(ctx *EventArrayDataTextContext) {}

// EnterEventArrayDataUID is called when production eventArrayDataUID is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataUID(ctx *EventArrayDataUIDContext) {}

// ExitEventArrayDataUID is called when production eventArrayDataUID is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataUID(ctx *EventArrayDataUIDContext) {}

// EnterEventArrayDataUint16 is called when production eventArrayDataUint16 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataUint16(ctx *EventArrayDataUint16Context) {}

// ExitEventArrayDataUint16 is called when production eventArrayDataUint16 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataUint16(ctx *EventArrayDataUint16Context) {}

// EnterEventArrayDataUint16X is called when production eventArrayDataUint16X is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataUint16X(ctx *EventArrayDataUint16XContext) {}

// ExitEventArrayDataUint16X is called when production eventArrayDataUint16X is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataUint16X(ctx *EventArrayDataUint16XContext) {}

// EnterEventArrayDataUint32 is called when production eventArrayDataUint32 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataUint32(ctx *EventArrayDataUint32Context) {}

// ExitEventArrayDataUint32 is called when production eventArrayDataUint32 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataUint32(ctx *EventArrayDataUint32Context) {}

// EnterEventArrayDataUint32X is called when production eventArrayDataUint32X is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataUint32X(ctx *EventArrayDataUint32XContext) {}

// ExitEventArrayDataUint32X is called when production eventArrayDataUint32X is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataUint32X(ctx *EventArrayDataUint32XContext) {}

// EnterEventArrayDataUint64 is called when production eventArrayDataUint64 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataUint64(ctx *EventArrayDataUint64Context) {}

// ExitEventArrayDataUint64 is called when production eventArrayDataUint64 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataUint64(ctx *EventArrayDataUint64Context) {}

// EnterEventArrayDataUint64X is called when production eventArrayDataUint64X is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataUint64X(ctx *EventArrayDataUint64XContext) {}

// ExitEventArrayDataUint64X is called when production eventArrayDataUint64X is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataUint64X(ctx *EventArrayDataUint64XContext) {}

// EnterEventArrayDataUint8 is called when production eventArrayDataUint8 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataUint8(ctx *EventArrayDataUint8Context) {}

// ExitEventArrayDataUint8 is called when production eventArrayDataUint8 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataUint8(ctx *EventArrayDataUint8Context) {}

// EnterEventArrayDataUint8X is called when production eventArrayDataUint8X is entered.
func (s *BaseCEEventParserListener) EnterEventArrayDataUint8X(ctx *EventArrayDataUint8XContext) {}

// ExitEventArrayDataUint8X is called when production eventArrayDataUint8X is exited.
func (s *BaseCEEventParserListener) ExitEventArrayDataUint8X(ctx *EventArrayDataUint8XContext) {}

// EnterEventArrayFloat16 is called when production eventArrayFloat16 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayFloat16(ctx *EventArrayFloat16Context) {}

// ExitEventArrayFloat16 is called when production eventArrayFloat16 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayFloat16(ctx *EventArrayFloat16Context) {}

// EnterEventArrayFloat32 is called when production eventArrayFloat32 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayFloat32(ctx *EventArrayFloat32Context) {}

// ExitEventArrayFloat32 is called when production eventArrayFloat32 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayFloat32(ctx *EventArrayFloat32Context) {}

// EnterEventArrayFloat64 is called when production eventArrayFloat64 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayFloat64(ctx *EventArrayFloat64Context) {}

// ExitEventArrayFloat64 is called when production eventArrayFloat64 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayFloat64(ctx *EventArrayFloat64Context) {}

// EnterEventArrayInt16 is called when production eventArrayInt16 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayInt16(ctx *EventArrayInt16Context) {}

// ExitEventArrayInt16 is called when production eventArrayInt16 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayInt16(ctx *EventArrayInt16Context) {}

// EnterEventArrayInt32 is called when production eventArrayInt32 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayInt32(ctx *EventArrayInt32Context) {}

// ExitEventArrayInt32 is called when production eventArrayInt32 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayInt32(ctx *EventArrayInt32Context) {}

// EnterEventArrayInt64 is called when production eventArrayInt64 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayInt64(ctx *EventArrayInt64Context) {}

// ExitEventArrayInt64 is called when production eventArrayInt64 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayInt64(ctx *EventArrayInt64Context) {}

// EnterEventArrayInt8 is called when production eventArrayInt8 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayInt8(ctx *EventArrayInt8Context) {}

// ExitEventArrayInt8 is called when production eventArrayInt8 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayInt8(ctx *EventArrayInt8Context) {}

// EnterEventArrayUID is called when production eventArrayUID is entered.
func (s *BaseCEEventParserListener) EnterEventArrayUID(ctx *EventArrayUIDContext) {}

// ExitEventArrayUID is called when production eventArrayUID is exited.
func (s *BaseCEEventParserListener) ExitEventArrayUID(ctx *EventArrayUIDContext) {}

// EnterEventArrayUint16 is called when production eventArrayUint16 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayUint16(ctx *EventArrayUint16Context) {}

// ExitEventArrayUint16 is called when production eventArrayUint16 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayUint16(ctx *EventArrayUint16Context) {}

// EnterEventArrayUint16X is called when production eventArrayUint16X is entered.
func (s *BaseCEEventParserListener) EnterEventArrayUint16X(ctx *EventArrayUint16XContext) {}

// ExitEventArrayUint16X is called when production eventArrayUint16X is exited.
func (s *BaseCEEventParserListener) ExitEventArrayUint16X(ctx *EventArrayUint16XContext) {}

// EnterEventArrayUint32 is called when production eventArrayUint32 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayUint32(ctx *EventArrayUint32Context) {}

// ExitEventArrayUint32 is called when production eventArrayUint32 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayUint32(ctx *EventArrayUint32Context) {}

// EnterEventArrayUint32X is called when production eventArrayUint32X is entered.
func (s *BaseCEEventParserListener) EnterEventArrayUint32X(ctx *EventArrayUint32XContext) {}

// ExitEventArrayUint32X is called when production eventArrayUint32X is exited.
func (s *BaseCEEventParserListener) ExitEventArrayUint32X(ctx *EventArrayUint32XContext) {}

// EnterEventArrayUint64 is called when production eventArrayUint64 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayUint64(ctx *EventArrayUint64Context) {}

// ExitEventArrayUint64 is called when production eventArrayUint64 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayUint64(ctx *EventArrayUint64Context) {}

// EnterEventArrayUint64X is called when production eventArrayUint64X is entered.
func (s *BaseCEEventParserListener) EnterEventArrayUint64X(ctx *EventArrayUint64XContext) {}

// ExitEventArrayUint64X is called when production eventArrayUint64X is exited.
func (s *BaseCEEventParserListener) ExitEventArrayUint64X(ctx *EventArrayUint64XContext) {}

// EnterEventArrayUint8 is called when production eventArrayUint8 is entered.
func (s *BaseCEEventParserListener) EnterEventArrayUint8(ctx *EventArrayUint8Context) {}

// ExitEventArrayUint8 is called when production eventArrayUint8 is exited.
func (s *BaseCEEventParserListener) ExitEventArrayUint8(ctx *EventArrayUint8Context) {}

// EnterEventArrayUint8X is called when production eventArrayUint8X is entered.
func (s *BaseCEEventParserListener) EnterEventArrayUint8X(ctx *EventArrayUint8XContext) {}

// ExitEventArrayUint8X is called when production eventArrayUint8X is exited.
func (s *BaseCEEventParserListener) ExitEventArrayUint8X(ctx *EventArrayUint8XContext) {}

// EnterEventBeginArrayBits is called when production eventBeginArrayBits is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayBits(ctx *EventBeginArrayBitsContext) {}

// ExitEventBeginArrayBits is called when production eventBeginArrayBits is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayBits(ctx *EventBeginArrayBitsContext) {}

// EnterEventBeginArrayFloat16 is called when production eventBeginArrayFloat16 is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayFloat16(ctx *EventBeginArrayFloat16Context) {}

// ExitEventBeginArrayFloat16 is called when production eventBeginArrayFloat16 is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayFloat16(ctx *EventBeginArrayFloat16Context) {}

// EnterEventBeginArrayFloat32 is called when production eventBeginArrayFloat32 is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayFloat32(ctx *EventBeginArrayFloat32Context) {}

// ExitEventBeginArrayFloat32 is called when production eventBeginArrayFloat32 is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayFloat32(ctx *EventBeginArrayFloat32Context) {}

// EnterEventBeginArrayFloat64 is called when production eventBeginArrayFloat64 is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayFloat64(ctx *EventBeginArrayFloat64Context) {}

// ExitEventBeginArrayFloat64 is called when production eventBeginArrayFloat64 is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayFloat64(ctx *EventBeginArrayFloat64Context) {}

// EnterEventBeginArrayInt16 is called when production eventBeginArrayInt16 is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayInt16(ctx *EventBeginArrayInt16Context) {}

// ExitEventBeginArrayInt16 is called when production eventBeginArrayInt16 is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayInt16(ctx *EventBeginArrayInt16Context) {}

// EnterEventBeginArrayInt32 is called when production eventBeginArrayInt32 is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayInt32(ctx *EventBeginArrayInt32Context) {}

// ExitEventBeginArrayInt32 is called when production eventBeginArrayInt32 is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayInt32(ctx *EventBeginArrayInt32Context) {}

// EnterEventBeginArrayInt64 is called when production eventBeginArrayInt64 is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayInt64(ctx *EventBeginArrayInt64Context) {}

// ExitEventBeginArrayInt64 is called when production eventBeginArrayInt64 is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayInt64(ctx *EventBeginArrayInt64Context) {}

// EnterEventBeginArrayInt8 is called when production eventBeginArrayInt8 is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayInt8(ctx *EventBeginArrayInt8Context) {}

// ExitEventBeginArrayInt8 is called when production eventBeginArrayInt8 is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayInt8(ctx *EventBeginArrayInt8Context) {}

// EnterEventBeginArrayUID is called when production eventBeginArrayUID is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayUID(ctx *EventBeginArrayUIDContext) {}

// ExitEventBeginArrayUID is called when production eventBeginArrayUID is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayUID(ctx *EventBeginArrayUIDContext) {}

// EnterEventBeginArrayUint16 is called when production eventBeginArrayUint16 is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayUint16(ctx *EventBeginArrayUint16Context) {}

// ExitEventBeginArrayUint16 is called when production eventBeginArrayUint16 is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayUint16(ctx *EventBeginArrayUint16Context) {}

// EnterEventBeginArrayUint32 is called when production eventBeginArrayUint32 is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayUint32(ctx *EventBeginArrayUint32Context) {}

// ExitEventBeginArrayUint32 is called when production eventBeginArrayUint32 is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayUint32(ctx *EventBeginArrayUint32Context) {}

// EnterEventBeginArrayUint64 is called when production eventBeginArrayUint64 is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayUint64(ctx *EventBeginArrayUint64Context) {}

// ExitEventBeginArrayUint64 is called when production eventBeginArrayUint64 is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayUint64(ctx *EventBeginArrayUint64Context) {}

// EnterEventBeginArrayUint8 is called when production eventBeginArrayUint8 is entered.
func (s *BaseCEEventParserListener) EnterEventBeginArrayUint8(ctx *EventBeginArrayUint8Context) {}

// ExitEventBeginArrayUint8 is called when production eventBeginArrayUint8 is exited.
func (s *BaseCEEventParserListener) ExitEventBeginArrayUint8(ctx *EventBeginArrayUint8Context) {}

// EnterEventBeginCustomBinary is called when production eventBeginCustomBinary is entered.
func (s *BaseCEEventParserListener) EnterEventBeginCustomBinary(ctx *EventBeginCustomBinaryContext) {}

// ExitEventBeginCustomBinary is called when production eventBeginCustomBinary is exited.
func (s *BaseCEEventParserListener) ExitEventBeginCustomBinary(ctx *EventBeginCustomBinaryContext) {}

// EnterEventBeginCustomText is called when production eventBeginCustomText is entered.
func (s *BaseCEEventParserListener) EnterEventBeginCustomText(ctx *EventBeginCustomTextContext) {}

// ExitEventBeginCustomText is called when production eventBeginCustomText is exited.
func (s *BaseCEEventParserListener) ExitEventBeginCustomText(ctx *EventBeginCustomTextContext) {}

// EnterEventBeginMedia is called when production eventBeginMedia is entered.
func (s *BaseCEEventParserListener) EnterEventBeginMedia(ctx *EventBeginMediaContext) {}

// ExitEventBeginMedia is called when production eventBeginMedia is exited.
func (s *BaseCEEventParserListener) ExitEventBeginMedia(ctx *EventBeginMediaContext) {}

// EnterEventBeginRemoteReference is called when production eventBeginRemoteReference is entered.
func (s *BaseCEEventParserListener) EnterEventBeginRemoteReference(ctx *EventBeginRemoteReferenceContext) {
}

// ExitEventBeginRemoteReference is called when production eventBeginRemoteReference is exited.
func (s *BaseCEEventParserListener) ExitEventBeginRemoteReference(ctx *EventBeginRemoteReferenceContext) {
}

// EnterEventBeginResourceId is called when production eventBeginResourceId is entered.
func (s *BaseCEEventParserListener) EnterEventBeginResourceId(ctx *EventBeginResourceIdContext) {}

// ExitEventBeginResourceId is called when production eventBeginResourceId is exited.
func (s *BaseCEEventParserListener) ExitEventBeginResourceId(ctx *EventBeginResourceIdContext) {}

// EnterEventBeginString is called when production eventBeginString is entered.
func (s *BaseCEEventParserListener) EnterEventBeginString(ctx *EventBeginStringContext) {}

// ExitEventBeginString is called when production eventBeginString is exited.
func (s *BaseCEEventParserListener) ExitEventBeginString(ctx *EventBeginStringContext) {}

// EnterEventBoolean is called when production eventBoolean is entered.
func (s *BaseCEEventParserListener) EnterEventBoolean(ctx *EventBooleanContext) {}

// ExitEventBoolean is called when production eventBoolean is exited.
func (s *BaseCEEventParserListener) ExitEventBoolean(ctx *EventBooleanContext) {}

// EnterEventCommentMultiline is called when production eventCommentMultiline is entered.
func (s *BaseCEEventParserListener) EnterEventCommentMultiline(ctx *EventCommentMultilineContext) {}

// ExitEventCommentMultiline is called when production eventCommentMultiline is exited.
func (s *BaseCEEventParserListener) ExitEventCommentMultiline(ctx *EventCommentMultilineContext) {}

// EnterEventCommentSingleLine is called when production eventCommentSingleLine is entered.
func (s *BaseCEEventParserListener) EnterEventCommentSingleLine(ctx *EventCommentSingleLineContext) {}

// ExitEventCommentSingleLine is called when production eventCommentSingleLine is exited.
func (s *BaseCEEventParserListener) ExitEventCommentSingleLine(ctx *EventCommentSingleLineContext) {}

// EnterEventCustomBinary is called when production eventCustomBinary is entered.
func (s *BaseCEEventParserListener) EnterEventCustomBinary(ctx *EventCustomBinaryContext) {}

// ExitEventCustomBinary is called when production eventCustomBinary is exited.
func (s *BaseCEEventParserListener) ExitEventCustomBinary(ctx *EventCustomBinaryContext) {}

// EnterEventCustomText is called when production eventCustomText is entered.
func (s *BaseCEEventParserListener) EnterEventCustomText(ctx *EventCustomTextContext) {}

// ExitEventCustomText is called when production eventCustomText is exited.
func (s *BaseCEEventParserListener) ExitEventCustomText(ctx *EventCustomTextContext) {}

// EnterEventEdge is called when production eventEdge is entered.
func (s *BaseCEEventParserListener) EnterEventEdge(ctx *EventEdgeContext) {}

// ExitEventEdge is called when production eventEdge is exited.
func (s *BaseCEEventParserListener) ExitEventEdge(ctx *EventEdgeContext) {}

// EnterEventEndContainer is called when production eventEndContainer is entered.
func (s *BaseCEEventParserListener) EnterEventEndContainer(ctx *EventEndContainerContext) {}

// ExitEventEndContainer is called when production eventEndContainer is exited.
func (s *BaseCEEventParserListener) ExitEventEndContainer(ctx *EventEndContainerContext) {}

// EnterEventList is called when production eventList is entered.
func (s *BaseCEEventParserListener) EnterEventList(ctx *EventListContext) {}

// ExitEventList is called when production eventList is exited.
func (s *BaseCEEventParserListener) ExitEventList(ctx *EventListContext) {}

// EnterEventMap is called when production eventMap is entered.
func (s *BaseCEEventParserListener) EnterEventMap(ctx *EventMapContext) {}

// ExitEventMap is called when production eventMap is exited.
func (s *BaseCEEventParserListener) ExitEventMap(ctx *EventMapContext) {}

// EnterEventMarker is called when production eventMarker is entered.
func (s *BaseCEEventParserListener) EnterEventMarker(ctx *EventMarkerContext) {}

// ExitEventMarker is called when production eventMarker is exited.
func (s *BaseCEEventParserListener) ExitEventMarker(ctx *EventMarkerContext) {}

// EnterEventMedia is called when production eventMedia is entered.
func (s *BaseCEEventParserListener) EnterEventMedia(ctx *EventMediaContext) {}

// ExitEventMedia is called when production eventMedia is exited.
func (s *BaseCEEventParserListener) ExitEventMedia(ctx *EventMediaContext) {}

// EnterEventNode is called when production eventNode is entered.
func (s *BaseCEEventParserListener) EnterEventNode(ctx *EventNodeContext) {}

// ExitEventNode is called when production eventNode is exited.
func (s *BaseCEEventParserListener) ExitEventNode(ctx *EventNodeContext) {}

// EnterEventNull is called when production eventNull is entered.
func (s *BaseCEEventParserListener) EnterEventNull(ctx *EventNullContext) {}

// ExitEventNull is called when production eventNull is exited.
func (s *BaseCEEventParserListener) ExitEventNull(ctx *EventNullContext) {}

// EnterEventNumber is called when production eventNumber is entered.
func (s *BaseCEEventParserListener) EnterEventNumber(ctx *EventNumberContext) {}

// ExitEventNumber is called when production eventNumber is exited.
func (s *BaseCEEventParserListener) ExitEventNumber(ctx *EventNumberContext) {}

// EnterEventPad is called when production eventPad is entered.
func (s *BaseCEEventParserListener) EnterEventPad(ctx *EventPadContext) {}

// ExitEventPad is called when production eventPad is exited.
func (s *BaseCEEventParserListener) ExitEventPad(ctx *EventPadContext) {}

// EnterEventLocalReference is called when production eventLocalReference is entered.
func (s *BaseCEEventParserListener) EnterEventLocalReference(ctx *EventLocalReferenceContext) {}

// ExitEventLocalReference is called when production eventLocalReference is exited.
func (s *BaseCEEventParserListener) ExitEventLocalReference(ctx *EventLocalReferenceContext) {}

// EnterEventRemoteReference is called when production eventRemoteReference is entered.
func (s *BaseCEEventParserListener) EnterEventRemoteReference(ctx *EventRemoteReferenceContext) {}

// ExitEventRemoteReference is called when production eventRemoteReference is exited.
func (s *BaseCEEventParserListener) ExitEventRemoteReference(ctx *EventRemoteReferenceContext) {}

// EnterEventResourceId is called when production eventResourceId is entered.
func (s *BaseCEEventParserListener) EnterEventResourceId(ctx *EventResourceIdContext) {}

// ExitEventResourceId is called when production eventResourceId is exited.
func (s *BaseCEEventParserListener) ExitEventResourceId(ctx *EventResourceIdContext) {}

// EnterEventString is called when production eventString is entered.
func (s *BaseCEEventParserListener) EnterEventString(ctx *EventStringContext) {}

// ExitEventString is called when production eventString is exited.
func (s *BaseCEEventParserListener) ExitEventString(ctx *EventStringContext) {}

// EnterEventStructInstance is called when production eventStructInstance is entered.
func (s *BaseCEEventParserListener) EnterEventStructInstance(ctx *EventStructInstanceContext) {}

// ExitEventStructInstance is called when production eventStructInstance is exited.
func (s *BaseCEEventParserListener) ExitEventStructInstance(ctx *EventStructInstanceContext) {}

// EnterEventStructTemplate is called when production eventStructTemplate is entered.
func (s *BaseCEEventParserListener) EnterEventStructTemplate(ctx *EventStructTemplateContext) {}

// ExitEventStructTemplate is called when production eventStructTemplate is exited.
func (s *BaseCEEventParserListener) ExitEventStructTemplate(ctx *EventStructTemplateContext) {}

// EnterEventTime is called when production eventTime is entered.
func (s *BaseCEEventParserListener) EnterEventTime(ctx *EventTimeContext) {}

// ExitEventTime is called when production eventTime is exited.
func (s *BaseCEEventParserListener) ExitEventTime(ctx *EventTimeContext) {}

// EnterEventUID is called when production eventUID is entered.
func (s *BaseCEEventParserListener) EnterEventUID(ctx *EventUIDContext) {}

// ExitEventUID is called when production eventUID is exited.
func (s *BaseCEEventParserListener) ExitEventUID(ctx *EventUIDContext) {}

// EnterEventVersion is called when production eventVersion is entered.
func (s *BaseCEEventParserListener) EnterEventVersion(ctx *EventVersionContext) {}

// ExitEventVersion is called when production eventVersion is exited.
func (s *BaseCEEventParserListener) ExitEventVersion(ctx *EventVersionContext) {}
