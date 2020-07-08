package ce

import (
	"math/big"
	"time"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// DataEventReceiver receives data events (int, string, etc) and performs
// actions based on those events. Generally, this is used to drive complex
// object builders, and also the encoders.
type DataEventReceiver interface {
	OnVersion(version uint64)
	OnPadding(count int)
	OnNil()
	OnBool(value bool)
	OnTrue()
	OnFalse()
	OnPositiveInt(value uint64)
	OnNegativeInt(value uint64)
	OnInt(value int64)
	OnBigInt(value *big.Int)
	OnFloat(value float64)
	OnBigFloat(value *big.Float)
	OnDecimalFloat(value compact_float.DFloat)
	OnBigDecimalFloat(value *apd.Decimal)
	OnComplex(value complex128)
	OnNan(signaling bool)
	OnUUID(value []byte)
	OnTime(value time.Time)
	OnCompactTime(value *compact_time.Time)
	OnBytes(value []byte)
	OnString(value string)
	OnURI(value string)
	OnCustom(value []byte)
	OnBytesBegin()
	OnStringBegin()
	OnURIBegin()
	OnCustomBegin()
	OnArrayChunk(length uint64, isFinalChunk bool)
	OnArrayData(data []byte)
	OnList()
	OnMap()
	OnMarkup()
	OnMetadata()
	OnComment()
	OnEnd()
	OnMarker()
	OnReference()
	OnEndDocument()
}
