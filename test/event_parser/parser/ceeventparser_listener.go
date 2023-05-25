// Code generated from /home/karl/Projects/go-concise-encoding/codegen/test/CEEventParser.g4 by ANTLR 4.12.0. DO NOT EDIT.

package parser // CEEventParser

import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// CEEventParserListener is a complete listener for a parse tree produced by CEEventParser.
type CEEventParserListener interface {
	antlr.ParseTreeListener

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterEvent is called when entering the event production.
	EnterEvent(c *EventContext)

	// EnterEventArrayBits is called when entering the eventArrayBits production.
	EnterEventArrayBits(c *EventArrayBitsContext)

	// EnterEventArrayChunkLast is called when entering the eventArrayChunkLast production.
	EnterEventArrayChunkLast(c *EventArrayChunkLastContext)

	// EnterEventArrayChunkMore is called when entering the eventArrayChunkMore production.
	EnterEventArrayChunkMore(c *EventArrayChunkMoreContext)

	// EnterEventArrayDataBits is called when entering the eventArrayDataBits production.
	EnterEventArrayDataBits(c *EventArrayDataBitsContext)

	// EnterEventArrayDataFloat16 is called when entering the eventArrayDataFloat16 production.
	EnterEventArrayDataFloat16(c *EventArrayDataFloat16Context)

	// EnterEventArrayDataFloat32 is called when entering the eventArrayDataFloat32 production.
	EnterEventArrayDataFloat32(c *EventArrayDataFloat32Context)

	// EnterEventArrayDataFloat64 is called when entering the eventArrayDataFloat64 production.
	EnterEventArrayDataFloat64(c *EventArrayDataFloat64Context)

	// EnterEventArrayDataInt16 is called when entering the eventArrayDataInt16 production.
	EnterEventArrayDataInt16(c *EventArrayDataInt16Context)

	// EnterEventArrayDataInt32 is called when entering the eventArrayDataInt32 production.
	EnterEventArrayDataInt32(c *EventArrayDataInt32Context)

	// EnterEventArrayDataInt64 is called when entering the eventArrayDataInt64 production.
	EnterEventArrayDataInt64(c *EventArrayDataInt64Context)

	// EnterEventArrayDataInt8 is called when entering the eventArrayDataInt8 production.
	EnterEventArrayDataInt8(c *EventArrayDataInt8Context)

	// EnterEventArrayDataText is called when entering the eventArrayDataText production.
	EnterEventArrayDataText(c *EventArrayDataTextContext)

	// EnterEventArrayDataUID is called when entering the eventArrayDataUID production.
	EnterEventArrayDataUID(c *EventArrayDataUIDContext)

	// EnterEventArrayDataUint16 is called when entering the eventArrayDataUint16 production.
	EnterEventArrayDataUint16(c *EventArrayDataUint16Context)

	// EnterEventArrayDataUint16X is called when entering the eventArrayDataUint16X production.
	EnterEventArrayDataUint16X(c *EventArrayDataUint16XContext)

	// EnterEventArrayDataUint32 is called when entering the eventArrayDataUint32 production.
	EnterEventArrayDataUint32(c *EventArrayDataUint32Context)

	// EnterEventArrayDataUint32X is called when entering the eventArrayDataUint32X production.
	EnterEventArrayDataUint32X(c *EventArrayDataUint32XContext)

	// EnterEventArrayDataUint64 is called when entering the eventArrayDataUint64 production.
	EnterEventArrayDataUint64(c *EventArrayDataUint64Context)

	// EnterEventArrayDataUint64X is called when entering the eventArrayDataUint64X production.
	EnterEventArrayDataUint64X(c *EventArrayDataUint64XContext)

	// EnterEventArrayDataUint8 is called when entering the eventArrayDataUint8 production.
	EnterEventArrayDataUint8(c *EventArrayDataUint8Context)

	// EnterEventArrayDataUint8X is called when entering the eventArrayDataUint8X production.
	EnterEventArrayDataUint8X(c *EventArrayDataUint8XContext)

	// EnterEventArrayFloat16 is called when entering the eventArrayFloat16 production.
	EnterEventArrayFloat16(c *EventArrayFloat16Context)

	// EnterEventArrayFloat32 is called when entering the eventArrayFloat32 production.
	EnterEventArrayFloat32(c *EventArrayFloat32Context)

	// EnterEventArrayFloat64 is called when entering the eventArrayFloat64 production.
	EnterEventArrayFloat64(c *EventArrayFloat64Context)

	// EnterEventArrayInt16 is called when entering the eventArrayInt16 production.
	EnterEventArrayInt16(c *EventArrayInt16Context)

	// EnterEventArrayInt32 is called when entering the eventArrayInt32 production.
	EnterEventArrayInt32(c *EventArrayInt32Context)

	// EnterEventArrayInt64 is called when entering the eventArrayInt64 production.
	EnterEventArrayInt64(c *EventArrayInt64Context)

	// EnterEventArrayInt8 is called when entering the eventArrayInt8 production.
	EnterEventArrayInt8(c *EventArrayInt8Context)

	// EnterEventArrayUID is called when entering the eventArrayUID production.
	EnterEventArrayUID(c *EventArrayUIDContext)

	// EnterEventArrayUint16 is called when entering the eventArrayUint16 production.
	EnterEventArrayUint16(c *EventArrayUint16Context)

	// EnterEventArrayUint16X is called when entering the eventArrayUint16X production.
	EnterEventArrayUint16X(c *EventArrayUint16XContext)

	// EnterEventArrayUint32 is called when entering the eventArrayUint32 production.
	EnterEventArrayUint32(c *EventArrayUint32Context)

	// EnterEventArrayUint32X is called when entering the eventArrayUint32X production.
	EnterEventArrayUint32X(c *EventArrayUint32XContext)

	// EnterEventArrayUint64 is called when entering the eventArrayUint64 production.
	EnterEventArrayUint64(c *EventArrayUint64Context)

	// EnterEventArrayUint64X is called when entering the eventArrayUint64X production.
	EnterEventArrayUint64X(c *EventArrayUint64XContext)

	// EnterEventArrayUint8 is called when entering the eventArrayUint8 production.
	EnterEventArrayUint8(c *EventArrayUint8Context)

	// EnterEventArrayUint8X is called when entering the eventArrayUint8X production.
	EnterEventArrayUint8X(c *EventArrayUint8XContext)

	// EnterEventBeginArrayBits is called when entering the eventBeginArrayBits production.
	EnterEventBeginArrayBits(c *EventBeginArrayBitsContext)

	// EnterEventBeginArrayFloat16 is called when entering the eventBeginArrayFloat16 production.
	EnterEventBeginArrayFloat16(c *EventBeginArrayFloat16Context)

	// EnterEventBeginArrayFloat32 is called when entering the eventBeginArrayFloat32 production.
	EnterEventBeginArrayFloat32(c *EventBeginArrayFloat32Context)

	// EnterEventBeginArrayFloat64 is called when entering the eventBeginArrayFloat64 production.
	EnterEventBeginArrayFloat64(c *EventBeginArrayFloat64Context)

	// EnterEventBeginArrayInt16 is called when entering the eventBeginArrayInt16 production.
	EnterEventBeginArrayInt16(c *EventBeginArrayInt16Context)

	// EnterEventBeginArrayInt32 is called when entering the eventBeginArrayInt32 production.
	EnterEventBeginArrayInt32(c *EventBeginArrayInt32Context)

	// EnterEventBeginArrayInt64 is called when entering the eventBeginArrayInt64 production.
	EnterEventBeginArrayInt64(c *EventBeginArrayInt64Context)

	// EnterEventBeginArrayInt8 is called when entering the eventBeginArrayInt8 production.
	EnterEventBeginArrayInt8(c *EventBeginArrayInt8Context)

	// EnterEventBeginArrayUID is called when entering the eventBeginArrayUID production.
	EnterEventBeginArrayUID(c *EventBeginArrayUIDContext)

	// EnterEventBeginArrayUint16 is called when entering the eventBeginArrayUint16 production.
	EnterEventBeginArrayUint16(c *EventBeginArrayUint16Context)

	// EnterEventBeginArrayUint32 is called when entering the eventBeginArrayUint32 production.
	EnterEventBeginArrayUint32(c *EventBeginArrayUint32Context)

	// EnterEventBeginArrayUint64 is called when entering the eventBeginArrayUint64 production.
	EnterEventBeginArrayUint64(c *EventBeginArrayUint64Context)

	// EnterEventBeginArrayUint8 is called when entering the eventBeginArrayUint8 production.
	EnterEventBeginArrayUint8(c *EventBeginArrayUint8Context)

	// EnterEventBeginCustomBinary is called when entering the eventBeginCustomBinary production.
	EnterEventBeginCustomBinary(c *EventBeginCustomBinaryContext)

	// EnterEventBeginCustomText is called when entering the eventBeginCustomText production.
	EnterEventBeginCustomText(c *EventBeginCustomTextContext)

	// EnterEventBeginMedia is called when entering the eventBeginMedia production.
	EnterEventBeginMedia(c *EventBeginMediaContext)

	// EnterEventBeginRemoteReference is called when entering the eventBeginRemoteReference production.
	EnterEventBeginRemoteReference(c *EventBeginRemoteReferenceContext)

	// EnterEventBeginResourceId is called when entering the eventBeginResourceId production.
	EnterEventBeginResourceId(c *EventBeginResourceIdContext)

	// EnterEventBeginString is called when entering the eventBeginString production.
	EnterEventBeginString(c *EventBeginStringContext)

	// EnterEventBoolean is called when entering the eventBoolean production.
	EnterEventBoolean(c *EventBooleanContext)

	// EnterEventCommentMultiline is called when entering the eventCommentMultiline production.
	EnterEventCommentMultiline(c *EventCommentMultilineContext)

	// EnterEventCommentSingleLine is called when entering the eventCommentSingleLine production.
	EnterEventCommentSingleLine(c *EventCommentSingleLineContext)

	// EnterEventCustomBinary is called when entering the eventCustomBinary production.
	EnterEventCustomBinary(c *EventCustomBinaryContext)

	// EnterEventCustomText is called when entering the eventCustomText production.
	EnterEventCustomText(c *EventCustomTextContext)

	// EnterEventEdge is called when entering the eventEdge production.
	EnterEventEdge(c *EventEdgeContext)

	// EnterEventEndContainer is called when entering the eventEndContainer production.
	EnterEventEndContainer(c *EventEndContainerContext)

	// EnterEventList is called when entering the eventList production.
	EnterEventList(c *EventListContext)

	// EnterEventMap is called when entering the eventMap production.
	EnterEventMap(c *EventMapContext)

	// EnterEventMarker is called when entering the eventMarker production.
	EnterEventMarker(c *EventMarkerContext)

	// EnterEventMedia is called when entering the eventMedia production.
	EnterEventMedia(c *EventMediaContext)

	// EnterEventNode is called when entering the eventNode production.
	EnterEventNode(c *EventNodeContext)

	// EnterEventNull is called when entering the eventNull production.
	EnterEventNull(c *EventNullContext)

	// EnterEventNumber is called when entering the eventNumber production.
	EnterEventNumber(c *EventNumberContext)

	// EnterEventPad is called when entering the eventPad production.
	EnterEventPad(c *EventPadContext)

	// EnterEventLocalReference is called when entering the eventLocalReference production.
	EnterEventLocalReference(c *EventLocalReferenceContext)

	// EnterEventRemoteReference is called when entering the eventRemoteReference production.
	EnterEventRemoteReference(c *EventRemoteReferenceContext)

	// EnterEventResourceId is called when entering the eventResourceId production.
	EnterEventResourceId(c *EventResourceIdContext)

	// EnterEventString is called when entering the eventString production.
	EnterEventString(c *EventStringContext)

	// EnterEventStructInstance is called when entering the eventStructInstance production.
	EnterEventStructInstance(c *EventStructInstanceContext)

	// EnterEventStructTemplate is called when entering the eventStructTemplate production.
	EnterEventStructTemplate(c *EventStructTemplateContext)

	// EnterEventTime is called when entering the eventTime production.
	EnterEventTime(c *EventTimeContext)

	// EnterEventUID is called when entering the eventUID production.
	EnterEventUID(c *EventUIDContext)

	// EnterEventVersion is called when entering the eventVersion production.
	EnterEventVersion(c *EventVersionContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitEvent is called when exiting the event production.
	ExitEvent(c *EventContext)

	// ExitEventArrayBits is called when exiting the eventArrayBits production.
	ExitEventArrayBits(c *EventArrayBitsContext)

	// ExitEventArrayChunkLast is called when exiting the eventArrayChunkLast production.
	ExitEventArrayChunkLast(c *EventArrayChunkLastContext)

	// ExitEventArrayChunkMore is called when exiting the eventArrayChunkMore production.
	ExitEventArrayChunkMore(c *EventArrayChunkMoreContext)

	// ExitEventArrayDataBits is called when exiting the eventArrayDataBits production.
	ExitEventArrayDataBits(c *EventArrayDataBitsContext)

	// ExitEventArrayDataFloat16 is called when exiting the eventArrayDataFloat16 production.
	ExitEventArrayDataFloat16(c *EventArrayDataFloat16Context)

	// ExitEventArrayDataFloat32 is called when exiting the eventArrayDataFloat32 production.
	ExitEventArrayDataFloat32(c *EventArrayDataFloat32Context)

	// ExitEventArrayDataFloat64 is called when exiting the eventArrayDataFloat64 production.
	ExitEventArrayDataFloat64(c *EventArrayDataFloat64Context)

	// ExitEventArrayDataInt16 is called when exiting the eventArrayDataInt16 production.
	ExitEventArrayDataInt16(c *EventArrayDataInt16Context)

	// ExitEventArrayDataInt32 is called when exiting the eventArrayDataInt32 production.
	ExitEventArrayDataInt32(c *EventArrayDataInt32Context)

	// ExitEventArrayDataInt64 is called when exiting the eventArrayDataInt64 production.
	ExitEventArrayDataInt64(c *EventArrayDataInt64Context)

	// ExitEventArrayDataInt8 is called when exiting the eventArrayDataInt8 production.
	ExitEventArrayDataInt8(c *EventArrayDataInt8Context)

	// ExitEventArrayDataText is called when exiting the eventArrayDataText production.
	ExitEventArrayDataText(c *EventArrayDataTextContext)

	// ExitEventArrayDataUID is called when exiting the eventArrayDataUID production.
	ExitEventArrayDataUID(c *EventArrayDataUIDContext)

	// ExitEventArrayDataUint16 is called when exiting the eventArrayDataUint16 production.
	ExitEventArrayDataUint16(c *EventArrayDataUint16Context)

	// ExitEventArrayDataUint16X is called when exiting the eventArrayDataUint16X production.
	ExitEventArrayDataUint16X(c *EventArrayDataUint16XContext)

	// ExitEventArrayDataUint32 is called when exiting the eventArrayDataUint32 production.
	ExitEventArrayDataUint32(c *EventArrayDataUint32Context)

	// ExitEventArrayDataUint32X is called when exiting the eventArrayDataUint32X production.
	ExitEventArrayDataUint32X(c *EventArrayDataUint32XContext)

	// ExitEventArrayDataUint64 is called when exiting the eventArrayDataUint64 production.
	ExitEventArrayDataUint64(c *EventArrayDataUint64Context)

	// ExitEventArrayDataUint64X is called when exiting the eventArrayDataUint64X production.
	ExitEventArrayDataUint64X(c *EventArrayDataUint64XContext)

	// ExitEventArrayDataUint8 is called when exiting the eventArrayDataUint8 production.
	ExitEventArrayDataUint8(c *EventArrayDataUint8Context)

	// ExitEventArrayDataUint8X is called when exiting the eventArrayDataUint8X production.
	ExitEventArrayDataUint8X(c *EventArrayDataUint8XContext)

	// ExitEventArrayFloat16 is called when exiting the eventArrayFloat16 production.
	ExitEventArrayFloat16(c *EventArrayFloat16Context)

	// ExitEventArrayFloat32 is called when exiting the eventArrayFloat32 production.
	ExitEventArrayFloat32(c *EventArrayFloat32Context)

	// ExitEventArrayFloat64 is called when exiting the eventArrayFloat64 production.
	ExitEventArrayFloat64(c *EventArrayFloat64Context)

	// ExitEventArrayInt16 is called when exiting the eventArrayInt16 production.
	ExitEventArrayInt16(c *EventArrayInt16Context)

	// ExitEventArrayInt32 is called when exiting the eventArrayInt32 production.
	ExitEventArrayInt32(c *EventArrayInt32Context)

	// ExitEventArrayInt64 is called when exiting the eventArrayInt64 production.
	ExitEventArrayInt64(c *EventArrayInt64Context)

	// ExitEventArrayInt8 is called when exiting the eventArrayInt8 production.
	ExitEventArrayInt8(c *EventArrayInt8Context)

	// ExitEventArrayUID is called when exiting the eventArrayUID production.
	ExitEventArrayUID(c *EventArrayUIDContext)

	// ExitEventArrayUint16 is called when exiting the eventArrayUint16 production.
	ExitEventArrayUint16(c *EventArrayUint16Context)

	// ExitEventArrayUint16X is called when exiting the eventArrayUint16X production.
	ExitEventArrayUint16X(c *EventArrayUint16XContext)

	// ExitEventArrayUint32 is called when exiting the eventArrayUint32 production.
	ExitEventArrayUint32(c *EventArrayUint32Context)

	// ExitEventArrayUint32X is called when exiting the eventArrayUint32X production.
	ExitEventArrayUint32X(c *EventArrayUint32XContext)

	// ExitEventArrayUint64 is called when exiting the eventArrayUint64 production.
	ExitEventArrayUint64(c *EventArrayUint64Context)

	// ExitEventArrayUint64X is called when exiting the eventArrayUint64X production.
	ExitEventArrayUint64X(c *EventArrayUint64XContext)

	// ExitEventArrayUint8 is called when exiting the eventArrayUint8 production.
	ExitEventArrayUint8(c *EventArrayUint8Context)

	// ExitEventArrayUint8X is called when exiting the eventArrayUint8X production.
	ExitEventArrayUint8X(c *EventArrayUint8XContext)

	// ExitEventBeginArrayBits is called when exiting the eventBeginArrayBits production.
	ExitEventBeginArrayBits(c *EventBeginArrayBitsContext)

	// ExitEventBeginArrayFloat16 is called when exiting the eventBeginArrayFloat16 production.
	ExitEventBeginArrayFloat16(c *EventBeginArrayFloat16Context)

	// ExitEventBeginArrayFloat32 is called when exiting the eventBeginArrayFloat32 production.
	ExitEventBeginArrayFloat32(c *EventBeginArrayFloat32Context)

	// ExitEventBeginArrayFloat64 is called when exiting the eventBeginArrayFloat64 production.
	ExitEventBeginArrayFloat64(c *EventBeginArrayFloat64Context)

	// ExitEventBeginArrayInt16 is called when exiting the eventBeginArrayInt16 production.
	ExitEventBeginArrayInt16(c *EventBeginArrayInt16Context)

	// ExitEventBeginArrayInt32 is called when exiting the eventBeginArrayInt32 production.
	ExitEventBeginArrayInt32(c *EventBeginArrayInt32Context)

	// ExitEventBeginArrayInt64 is called when exiting the eventBeginArrayInt64 production.
	ExitEventBeginArrayInt64(c *EventBeginArrayInt64Context)

	// ExitEventBeginArrayInt8 is called when exiting the eventBeginArrayInt8 production.
	ExitEventBeginArrayInt8(c *EventBeginArrayInt8Context)

	// ExitEventBeginArrayUID is called when exiting the eventBeginArrayUID production.
	ExitEventBeginArrayUID(c *EventBeginArrayUIDContext)

	// ExitEventBeginArrayUint16 is called when exiting the eventBeginArrayUint16 production.
	ExitEventBeginArrayUint16(c *EventBeginArrayUint16Context)

	// ExitEventBeginArrayUint32 is called when exiting the eventBeginArrayUint32 production.
	ExitEventBeginArrayUint32(c *EventBeginArrayUint32Context)

	// ExitEventBeginArrayUint64 is called when exiting the eventBeginArrayUint64 production.
	ExitEventBeginArrayUint64(c *EventBeginArrayUint64Context)

	// ExitEventBeginArrayUint8 is called when exiting the eventBeginArrayUint8 production.
	ExitEventBeginArrayUint8(c *EventBeginArrayUint8Context)

	// ExitEventBeginCustomBinary is called when exiting the eventBeginCustomBinary production.
	ExitEventBeginCustomBinary(c *EventBeginCustomBinaryContext)

	// ExitEventBeginCustomText is called when exiting the eventBeginCustomText production.
	ExitEventBeginCustomText(c *EventBeginCustomTextContext)

	// ExitEventBeginMedia is called when exiting the eventBeginMedia production.
	ExitEventBeginMedia(c *EventBeginMediaContext)

	// ExitEventBeginRemoteReference is called when exiting the eventBeginRemoteReference production.
	ExitEventBeginRemoteReference(c *EventBeginRemoteReferenceContext)

	// ExitEventBeginResourceId is called when exiting the eventBeginResourceId production.
	ExitEventBeginResourceId(c *EventBeginResourceIdContext)

	// ExitEventBeginString is called when exiting the eventBeginString production.
	ExitEventBeginString(c *EventBeginStringContext)

	// ExitEventBoolean is called when exiting the eventBoolean production.
	ExitEventBoolean(c *EventBooleanContext)

	// ExitEventCommentMultiline is called when exiting the eventCommentMultiline production.
	ExitEventCommentMultiline(c *EventCommentMultilineContext)

	// ExitEventCommentSingleLine is called when exiting the eventCommentSingleLine production.
	ExitEventCommentSingleLine(c *EventCommentSingleLineContext)

	// ExitEventCustomBinary is called when exiting the eventCustomBinary production.
	ExitEventCustomBinary(c *EventCustomBinaryContext)

	// ExitEventCustomText is called when exiting the eventCustomText production.
	ExitEventCustomText(c *EventCustomTextContext)

	// ExitEventEdge is called when exiting the eventEdge production.
	ExitEventEdge(c *EventEdgeContext)

	// ExitEventEndContainer is called when exiting the eventEndContainer production.
	ExitEventEndContainer(c *EventEndContainerContext)

	// ExitEventList is called when exiting the eventList production.
	ExitEventList(c *EventListContext)

	// ExitEventMap is called when exiting the eventMap production.
	ExitEventMap(c *EventMapContext)

	// ExitEventMarker is called when exiting the eventMarker production.
	ExitEventMarker(c *EventMarkerContext)

	// ExitEventMedia is called when exiting the eventMedia production.
	ExitEventMedia(c *EventMediaContext)

	// ExitEventNode is called when exiting the eventNode production.
	ExitEventNode(c *EventNodeContext)

	// ExitEventNull is called when exiting the eventNull production.
	ExitEventNull(c *EventNullContext)

	// ExitEventNumber is called when exiting the eventNumber production.
	ExitEventNumber(c *EventNumberContext)

	// ExitEventPad is called when exiting the eventPad production.
	ExitEventPad(c *EventPadContext)

	// ExitEventLocalReference is called when exiting the eventLocalReference production.
	ExitEventLocalReference(c *EventLocalReferenceContext)

	// ExitEventRemoteReference is called when exiting the eventRemoteReference production.
	ExitEventRemoteReference(c *EventRemoteReferenceContext)

	// ExitEventResourceId is called when exiting the eventResourceId production.
	ExitEventResourceId(c *EventResourceIdContext)

	// ExitEventString is called when exiting the eventString production.
	ExitEventString(c *EventStringContext)

	// ExitEventStructInstance is called when exiting the eventStructInstance production.
	ExitEventStructInstance(c *EventStructInstanceContext)

	// ExitEventStructTemplate is called when exiting the eventStructTemplate production.
	ExitEventStructTemplate(c *EventStructTemplateContext)

	// ExitEventTime is called when exiting the eventTime production.
	ExitEventTime(c *EventTimeContext)

	// ExitEventUID is called when exiting the eventUID production.
	ExitEventUID(c *EventUIDContext)

	// ExitEventVersion is called when exiting the eventVersion production.
	ExitEventVersion(c *EventVersionContext)
}
