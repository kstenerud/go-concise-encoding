Presentation: go-concise-encoding
=================================

## Agenda

- Overview of data encoding and where Concise Encoding fits
  - Brief overview of Concise Encoding
- Reference implementation (`go-concise-encoding`)
  - Overall architecture
  - Code examples
- Light code view of CBE codec
- Personal takeaways


## Data Encoding in Today's World

Two main popular encoding philosophies:

**Record oriented**:
- Generally record based
- Maps directly (or close to directly) to the actual memory layout of the data (at least in C++)
- Explicitly defined data
- Usually involves code generation and data definition files
- Focus on codec speed and mmap based data loading
- Optimized for big data
- Somewhat wasteful in data transmission as a tradeoff
- Examples Protobufs, flatbuffers

**Element oriented**:
- Supports ad-hoc data
- Schema optional
- Many opportunities for data compression
- Generally, these formats are poorly defined, resulting in implementation differences (security issue)
- Generally, these lack many types, resulting in stringification (CPU intensive)
- They tend to be text formats
- Examples: JSON, XML, CBOR, HTML

### Issues in existing formats:
- Lack of types
- Have to choose between text (CPU and bandwidth wasteful) and binary (human-hostile)
- Deprecations and changes to the format are problematic


## Concise Encoding
- Element-oriented
- Representable in binary and text, and transferrable between text/binary with no data loss
- Tightly specified (for better security)
- Native support for all common types, containers and structs (no more stringifying)
- Future-proof (versioning)
- Low energy and size impact
- No definition files or code generation steps; just import the module and go.

https://concise-encoding.org/


### Closest competitor: Amazon Ion
- Big endian
- Less types
- No typed arrays
- No time zones
- No BC dates
- No leap seconds
- No bfloat16
- List must be prefixed with length
- Map keys must be strings


## Reference implementation

- Written in go
- Designed to be efficient, but modular
- SAX-style event model (push)
- Out-of-the-box support for almost all go types, structs, containers

Example: https://play.golang.org/p/6_pD6CQVLuN


### Architecture

- **Events**: The contract that all major components operate by (`DataEventReceiver`).
- **Iterators**: Consume go objects to produce data events that describe them.
- **Builders**: Consume data events to produce go objects.
- **Codecs**: Marshal and unmarshal events to/from CTE or CBE documents.
- **Rules**: Enforce proper structure and content in Concise Encoding documents.

```golang
type DataEventReceiver interface {
	// Must be called before any other event.
	OnBeginDocument()
	// Must be called last of all. No other events may be sent after this call.
	OnEndDocument()

	OnVersion(version uint64)
	...
	OnInt(value int64)
	...
	OnList()
	OnMap()
	...
	OnEnd()
	...
	OnArray(arrayType ArrayType, elementCount uint64, data []uint8)
	OnArrayBegin(arrayType ArrayType)
	OnArrayChunk(length uint64, moreChunksFollow bool)
	OnArrayData(data []byte)
	...
}
```

#### Iterators

```golang
// Iterator interface. All iterators follow this signature.
type IteratorFunction func(context *Context, value reflect.Value)
```

**Integer**:
```golang
func iterateInt(context *Context, v reflect.Value) {
	context.EventReceiver.OnInt(v.Int())
}
```

**Interface**:
```golang
func iterateInterface(context *Context, v reflect.Value) {
	if common.IsNil(v) {
		context.NotifyNil()
	} else {
		elem := v.Elem()
		iterate := context.GetIteratorForType(elem.Type())
		iterate(context, elem)
	}
}
```

**Pointer**:
```golang
func newPointerIterator(ctx *Context, pointerType reflect.Type) IteratorFunction {
	iterate := ctx.GetIteratorForType(pointerType.Elem())

	return func(context *Context, v reflect.Value) {
		if common.IsNil(v) {
			context.NotifyNil()
			return
		}
		if context.TryAddReference(v) {
			return
		}
		iterate(context, v.Elem())
	}
}
```

**Map**:
```golang
func newMapIterator(ctx *Context, mapType reflect.Type) IteratorFunction {
	iterateKey := ctx.GetIteratorForType(mapType.Key())
	iterateValue := ctx.GetIteratorForType(mapType.Elem())

	return func(context *Context, v reflect.Value) {
		if common.IsNil(v) {
			context.NotifyNil()
			return
		}
		if context.TryAddReference(v) {
			return
		}

		context.EventReceiver.OnMap()
		iter := common.MapRange(v)
		for iter.Next() {
			iterateKey(context, iter.Key())
			iterateValue(context, iter.Value())
		}
		context.EventReceiver.OnEnd()
	}
}
```

**Struct**:
```golang
func newStructIterator(ctx *Context, structType reflect.Type) IteratorFunction {
	fields := make([]structField, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		reflectField := structType.Field(i)
		if common.IsFieldExported(reflectField.Name) {
			field := structField{
				Name:  reflectField.Name,
				Type:  reflectField.Type,
				Index: i,
			}
			field.applyTags(reflectField.Tag.Get("ce"))
			if ctx.LowercaseStructFieldNames {
				field.Name = common.ASCIIToLower(field.Name)
			}

			if !field.Omit {
				field.Iterate = ctx.GetIteratorForType(field.Type)
				fields = append(fields, field)
			}
		}
	}

	return func(context *Context, v reflect.Value) {
		context.EventReceiver.OnMap()

		for _, field := range fields {
			context.EventReceiver.OnStringlikeArray(events.ArrayTypeString, field.Name)
			field.Iterate(context, v.Field(field.Index))
		}

		context.EventReceiver.OnEnd()
	}
}
```

#### Builders

```golang
// ObjectBuilder responds to external events to progressively build an object.
type Builder interface {
	// External data and structure events
	BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value
	...
	BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value
	...
	// Signals that a new source container has begun.
	// This gets triggered from a data event.
	BuildInitiateList(ctx *Context)
	BuildInitiateMap(ctx *Context)
	...
	// Signals that the source container is finished
	// This gets triggered from a data event.
	BuildEndContainer(ctx *Context)

	// Tells this builder to create a new container to receive the source container's objects.
	// This gets called by the parent builder.
	BuildBeginListContents(ctx *Context)
	BuildBeginMapContents(ctx *Context)
	BuildBeginMarkupContents(ctx *Context, name []byte)

	// Notify that a child builder has finished building a container.
	// This gets triggered from the child builder when the container has ended and the builder unstacked.
	NotifyChildContainerFinished(ctx *Context, container reflect.Value)
}

type BuilderGenerator func(ctx *Context) Builder
type BuilderGeneratorGetter func(reflect.Type) BuilderGenerator
```

**BoolBuilder**:
```golang
type boolBuilder struct{}

var globalBoolBuilder = &boolBuilder{}

func generateBoolBuilder(ctx *Context) Builder { return globalBoolBuilder }

func (_this *boolBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	dst.SetBool(value)
	return dst
}
```

**SliceBuilder**:
```golang
type sliceBuilder struct {
	dstType       reflect.Type
	elemGenerator BuilderGenerator
	ppContainer   **reflect.Value
}

func newSliceBuilderGenerator(getBuilderGeneratorForType BuilderGeneratorGetter, dstType reflect.Type) BuilderGenerator {
	builderGenerator := getBuilderGeneratorForType(dstType.Elem())

	return func(ctx *Context) Builder {
		builder := &sliceBuilder{
			dstType:       dstType,
			elemGenerator: builderGenerator,
		}
		return builder
	}
}

func (_this *sliceBuilder) newElem() reflect.Value {
	return reflect.New(_this.dstType.Elem()).Elem()
}

func (_this *sliceBuilder) storeValue(value reflect.Value) {
	**_this.ppContainer = reflect.Append(**_this.ppContainer, value)
}

func (_this *sliceBuilder) BuildBeginListContents(ctx *Context) {
	ctx.StackBuilder(_this)
	_this.reset()
}

func (_this *sliceBuilder) BuildEndContainer(ctx *Context) {
	object := **_this.ppContainer
	ctx.UnstackBuilderAndNotifyChildFinished(object)
}

func (_this *sliceBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.storeValue(value)
}

func (_this *sliceBuilder) BuildFromBool(ctx *Context, value bool, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.elemGenerator(ctx).BuildFromBool(ctx, value, object)
	_this.storeValue(object)
	return object
}

func (_this *sliceBuilder) BuildInitiateList(ctx *Context) {
	_this.elemGenerator(ctx).BuildBeginListContents(ctx)
}
```

**MapBuilder**:
```golang
const (
	kvBuilderKey   = 0
	kvBuilderValue = 1
)

type mapBuilder struct {
	mapType         reflect.Type
	kvTypes         [2]reflect.Type
	kvGenerators    [2]BuilderGenerator
	container       reflect.Value
	key             reflect.Value
	builderIndex    int
	nextGenerator   BuilderGenerator
	nextStoreMethod func(*mapBuilder, reflect.Value)
}

func newMapBuilderGenerator(getBuilderGeneratorForType BuilderGeneratorGetter, mapType reflect.Type) BuilderGenerator {
	kvTypes := [2]reflect.Type{mapType.Key(), mapType.Elem()}
	kvGenerators := [2]BuilderGenerator{getBuilderGeneratorForType(kvTypes[0]), getBuilderGeneratorForType(kvTypes[1])}

	return func(ctx *Context) Builder {
		builder := &mapBuilder{
			mapType:      mapType,
			kvTypes:      kvTypes,
			kvGenerators: kvGenerators,
		}
		return builder
	}
}

func (_this *mapBuilder) storeKey(value reflect.Value) {
	_this.key = value
}

func (_this *mapBuilder) storeValue(value reflect.Value) {
	_this.container.SetMapIndex(_this.key, value)
}

var mapBuilderKVStoreMethods = []func(*mapBuilder, reflect.Value){
	(*mapBuilder).storeKey,
	(*mapBuilder).storeValue,
}

func (_this *mapBuilder) store(value reflect.Value) {
	_this.nextStoreMethod(_this, value)
	_this.swapKeyValue()
}

func (_this *mapBuilder) swapKeyValue() {
	_this.builderIndex = (_this.builderIndex + 1) & 1
	_this.nextGenerator = _this.kvGenerators[_this.builderIndex]
	_this.nextStoreMethod = mapBuilderKVStoreMethods[_this.builderIndex]
}

func (_this *mapBuilder) newElem() reflect.Value {
	return reflect.New(_this.kvTypes[_this.builderIndex]).Elem()
}

func (_this *mapBuilder) BuildBeginMapContents(ctx *Context) {
	ctx.StackBuilder(_this)
	_this.reset()
}

func (_this *mapBuilder) BuildEndContainer(ctx *Context) {
	object := _this.container
	ctx.UnstackBuilderAndNotifyChildFinished(object)
}

func (_this *mapBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	_this.store(value)
}

func (_this *mapBuilder) BuildFromBool(ctx *Context, value bool, _ reflect.Value) reflect.Value {
	object := _this.newElem()
	_this.nextGenerator(ctx).BuildFromBool(ctx, value, object)
	_this.store(object)
	return object
}

func (_this *mapBuilder) BuildInitiateList(ctx *Context) {
	_this.nextGenerator(ctx).BuildBeginListContents(ctx)
}

func (_this *mapBuilder) BuildInitiateMap(ctx *Context) {
	_this.nextGenerator(ctx).BuildBeginMapContents(ctx)
}
```

**StructBuilder**:
```golang
type structBuilder struct {
	dstType                reflect.Type
	generatorDescs         map[string]*structBuilderGeneratorDesc
	nameBuilderGenerator   BuilderGenerator
	ignoreBuilderGenerator BuilderGenerator
	nextBuilderGenerator   BuilderGenerator
	container              reflect.Value
	nextValue              reflect.Value
	nextIsKey              bool
	nextIsIgnored          bool
}

type structBuilderGeneratorDesc struct {
	field            *structBuilderField
	builderGenerator BuilderGenerator
}

func newStructBuilderGenerator(getBuilderGeneratorForType BuilderGeneratorGetter, dstType reflect.Type) BuilderGenerator {
	nameBuilderGenerator := getBuilderGeneratorForType(reflect.TypeOf(""))
	ignoreBuilderGenerator := generateIgnoreBuilder
	generatorDescs := make(map[string]*structBuilderGeneratorDesc)

	for i := 0; i < dstType.NumField(); i++ {
		reflectField := dstType.Field(i)
		if reflectField.PkgPath == "" {
			builderGenerator := getBuilderGeneratorForType(reflectField.Type)
			structField := &structBuilderField{
				Name:  reflectField.Name,
				Index: i,
			}
			structField.applyTags(reflectField.Tag.Get("ce"))
			generatorDescs[structField.Name] = &structBuilderGeneratorDesc{
				field:            structField,
				builderGenerator: builderGenerator,
			}
		}
	}

	// Make lowercase mappings as well in case we later do case-insensitive field name matching
	for _, desc := range generatorDescs {
		lowerName := common.ASCIIToLower(desc.field.Name)
		if _, exists := generatorDescs[lowerName]; !exists {
			generatorDescs[lowerName] = desc
		}
	}

	return func(ctx *Context) Builder {
		builder := &structBuilder{
			dstType:                dstType,
			generatorDescs:         generatorDescs,
			nameBuilderGenerator:   nameBuilderGenerator,
			ignoreBuilderGenerator: ignoreBuilderGenerator,
		}
		return builder
	}
}

func (_this *structBuilder) swapKeyValue() {
	_this.nextIsKey = !_this.nextIsKey
}

func (_this *structBuilder) BuildBeginMapContents(ctx *Context) {
	ctx.StackBuilder(_this)
	_this.reset()
}

func (_this *structBuilder) BuildEndContainer(ctx *Context) {
	object := _this.container
	ctx.UnstackBuilderAndNotifyChildFinished(object)
}

func (_this *structBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	if _this.nextIsIgnored {
		_this.nextIsIgnored = false
		return
	}

	_this.nextValue.Set(value)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromBool(ctx *Context, value bool, _ reflect.Value) reflect.Value {
	_this.nextBuilderGenerator(ctx).BuildFromBool(ctx, value, _this.nextValue)
	object := _this.nextValue
	_this.swapKeyValue()
	return object
}

func (_this *structBuilder) BuildInitiateMap(ctx *Context) {
	_this.nextBuilderGenerator(ctx).BuildBeginMapContents(ctx)
}
```

### CBE Codec

https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#type-field

- [Decoder](cbe/decoder.go)
- [Encoder](cbe/encoder.go)
- [Marshaler](cbe/marshal.go)

## Takeaways

- Memory allocations are your enemy. Nothing slows go down more than allocations.
- Magic allocations are the worst (string to/from byte slice)
- Caching is absolutely essential in larger programs
- Cache locality is king (keep related data close together)
- Go really needs generics (lots of code generation required otherwise)
- Read-only values defined at runtime would be nice (more like C-style const)
- Read-only slices would be nice (would cut down on magic allocations)
