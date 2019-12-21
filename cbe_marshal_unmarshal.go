package cbe

import (
	"github.com/kstenerud/go-cbe/rules"
	"github.com/kstenerud/go-reconstruct"
)

func MarshalCBE(object interface{}) (document []byte, err error) {
	var encoder CBEEncoder
	encoder.Init(InlineContainerTypeNone, nil, rules.DefaultLimits())
	var adapter ObjectToCBEAdapter
	adapter.Init(&encoder)
	var iterator reconstruct.ObjectIterator
	iterator.Init(&adapter)
	if err = iterator.Iterate(object); err != nil {
		return
	}
	if err = encoder.End(); err != nil {
		return
	}
	document = encoder.EncodedBytes()
	return
}

func UnmarshalCBE(cbeDocument []byte, dst interface{}) (err error) {
	var builder reconstruct.AdhocBuilder
	var adapter CBEToObjectAdapter
	adapter.Init(&builder)
	var decoder CBEDecoder
	decoder.Init(InlineContainerTypeNone, rules.DefaultLimits(), &adapter)
	if err = decoder.Decode(cbeDocument); err != nil {
		return
	}
	if err = decoder.EndDocument(); err != nil {
		return
	}

	err = reconstruct.Reconstruct(builder.GetObject(), dst)

	return
}
