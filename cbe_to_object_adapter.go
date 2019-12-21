package cbe

import (
	"fmt"
	"net/url"

	"github.com/kstenerud/go-reconstruct"

	"github.com/kstenerud/go-compact-time"
)

type adhocContainerType int

const (
	adhocContainerTypeNone adhocContainerType = iota
	adhocContainerTypeList
	adhocContainerTypeMap
	adhocContainerTypeBytes
	adhocContainerTypeString
	adhocContainerTypeURI
	adhocContainerTypeIgnored
)

type CBEToObjectAdapter struct {
	callbacks             reconstruct.ObjectIteratorCallbacks
	ignoredContainerDepth int
	containerTypes        []adhocContainerType
	currentArray          []byte
	remainingArrayLength  uint64
	isFinalChunk          bool
}

func NewCBEToObjectAdapter(callbacks reconstruct.ObjectIteratorCallbacks) *CBEToObjectAdapter {
	this := new(CBEToObjectAdapter)
	this.Init(callbacks)
	return this
}

func (this *CBEToObjectAdapter) Init(callbacks reconstruct.ObjectIteratorCallbacks) {
	this.callbacks = callbacks
}

func (this *CBEToObjectAdapter) stackContainer(containerType adhocContainerType) {
	this.containerTypes = append(this.containerTypes, containerType)
}

func (this *CBEToObjectAdapter) unstackContainer() {
	this.containerTypes = this.containerTypes[:len(this.containerTypes)-1]
}

func (this *CBEToObjectAdapter) getCurrentContainerType() adhocContainerType {
	return this.containerTypes[len(this.containerTypes)-1]
}

func (this *CBEToObjectAdapter) OnNil() error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	return this.callbacks.OnNil()
}

func (this *CBEToObjectAdapter) OnBool(value bool) error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	return this.callbacks.OnBool(value)
}

func (this *CBEToObjectAdapter) OnPositiveInt(value uint64) error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	return this.callbacks.OnUint(value)
}

func (this *CBEToObjectAdapter) OnNegativeInt(value uint64) error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	if value > 0x7fffffffffffffff {
		return fmt.Errorf("Value %x cannot fit into an int64", value)
	}
	return this.callbacks.OnInt(-int64(value))
}

func (this *CBEToObjectAdapter) OnFloat(value float64) error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	return this.callbacks.OnFloat(value)
}

func (this *CBEToObjectAdapter) OnTime(time *compact_time.Time) error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	gotime, err := time.AsGoTime()
	if err != nil {
		return err
	}
	return this.callbacks.OnTime(gotime)
}

func (this *CBEToObjectAdapter) OnListBegin() error {
	if this.ignoredContainerDepth > 0 {
		this.ignoredContainerDepth++
		return nil
	}
	this.stackContainer(adhocContainerTypeList)
	return this.callbacks.OnListBegin()
}

func (this *CBEToObjectAdapter) OnMapBegin() error {
	if this.ignoredContainerDepth > 0 {
		this.ignoredContainerDepth++
		return nil
	}
	this.stackContainer(adhocContainerTypeMap)
	return this.callbacks.OnMapBegin()
}

func (this *CBEToObjectAdapter) OnMarkupBegin() error {
	return fmt.Errorf("Adhoc data adapter cannot handle markup")
}

func (this *CBEToObjectAdapter) OnMetadataBegin() error {
	this.ignoredContainerDepth++
	return nil
}

func (this *CBEToObjectAdapter) OnCommentBegin() error {
	this.ignoredContainerDepth++
	return nil
}

func (this *CBEToObjectAdapter) OnContainerEnd() error {
	if this.ignoredContainerDepth > 0 {
		this.ignoredContainerDepth--
		return nil
	}
	containerType := this.getCurrentContainerType()
	this.unstackContainer()
	switch containerType {
	case adhocContainerTypeList:
		return this.callbacks.OnListEnd()
	case adhocContainerTypeMap:
		return this.callbacks.OnMapEnd()
	default:
		panic(fmt.Errorf("BUG: %v: Invalid container type", containerType))
	}
}

func (this *CBEToObjectAdapter) OnMarkerBegin() error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	// TODO
	return nil
}

func (this *CBEToObjectAdapter) OnReferenceBegin() error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	// TODO
	return nil
}

func (this *CBEToObjectAdapter) OnBytesBegin() error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	this.stackContainer(adhocContainerTypeBytes)
	this.currentArray = this.currentArray[:0]
	return nil
}

func (this *CBEToObjectAdapter) OnStringBegin() error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	this.stackContainer(adhocContainerTypeString)
	this.currentArray = this.currentArray[:0]
	return nil
}

func (this *CBEToObjectAdapter) OnURIBegin() error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	this.stackContainer(adhocContainerTypeURI)
	this.currentArray = this.currentArray[:0]
	return nil
}

func (this *CBEToObjectAdapter) endArray() error {
	containerType := this.getCurrentContainerType()
	this.unstackContainer()
	switch containerType {
	case adhocContainerTypeBytes:
		return this.callbacks.OnBytes(this.currentArray)
	case adhocContainerTypeString:
		return this.callbacks.OnString(string(this.currentArray))
	case adhocContainerTypeURI:
		uri, err := url.Parse(string(this.currentArray))
		if err != nil {
			return err
		}
		return this.callbacks.OnURI(uri)
	default:
		panic(fmt.Errorf("BUG: %v: Invalid container type", containerType))
	}
}

func (this *CBEToObjectAdapter) OnArrayChunkBegin(byteCount uint64, isFinalChunk bool) error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	this.remainingArrayLength = byteCount
	this.isFinalChunk = isFinalChunk
	if this.remainingArrayLength == 0 && this.isFinalChunk {
		return this.endArray()
	}
	return nil
}

func (this *CBEToObjectAdapter) OnArrayData(bytes []byte) error {
	if this.ignoredContainerDepth > 0 {
		return nil
	}
	this.currentArray = append(this.currentArray, bytes...)
	this.remainingArrayLength -= uint64(len(bytes))
	if this.remainingArrayLength == 0 && this.isFinalChunk {
		return this.endArray()
	}
	return nil
}

func (this *CBEToObjectAdapter) OnDocumentEnd() error {
	// TODO: Anything?
	return nil
}
