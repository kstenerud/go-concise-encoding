Architecture
============

This document describes the architecture of this project. As this is the reference implementation of [Concise Encoding](https://concise-encoding.org/), I've done by best to keep it as readable as possible.

The code is separated into five primary sections:

- [Events](#events): Contract that all major components operate by.
- [Iterators](#iterators): Iterators iterate through a go object to produce data events.
- [Builders](#builders): Builders interpret data events to produce go objects.
- [Codecs](#codecs): Codecs encode/decode data events to/from CTE or CBE documents.
- [Rules](#rules): Rules enforce proper structure and content in Concise Encoding documents.

The other secondary sections are: 

- [Code Generation](#code-generation)
- [Debug Helpers](#debug-helpers)
- [Test Helpers](#test-helpers)
- [Conversions](#conversions)
- [Options](#options)



Code Organization
-----------------

| Directory                  | Description                                                     |
| -------------------------- | --------------------------------------------------------------- |
| [buffer](buffer)           | Basic data buffer code used by the codecs                       |
| [builder](builder)         | [Builders](#builders)                                           |
| [cbe](cbe)                 | [CBE codec](https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md) |
| [ce](ce)                   | Top-level API                                                   |
| [codegen](codegen)         | Code generator source (generates all `generated-code.go` files) |
| [conversions](conversions) | Data type converters used by builders and codecs                |
| [cte](cte)                 | [CTE codec](https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md) |
| [debug](debug)             | Tools to help with debugging                                    |
| [events](events)           | [Events](#events), and a "null" event receiver                  |
| [internal](internal)       | Various tools used internally by the library                    |
| [iterator](iterator)       | [Iterators](#iterators)                                         |
| [options](options)         | Configuration for all high level APIs                           |
| [rules](rules)             | [Rules](#rules)                                                 |
| [test](test)               | Test helper code                                                |
| [version](version)         | The currently supported Concise Encoding version                |



Primary Sections
----------------

### Events

Events form the backbone of the entire library. The major components either consume or produce data events. This architecture makes it easy to mix & match components to produce whatever software design you want.

See: [data_events.go](events/data_events.go)


### Iterators

Iterators inspect go objects to produce data events. They support all primitives, as well as arrays, slices, maps, pointers, and structs, and can handle recursive pointers. All iterators follow a common interface defined in [iterator.go](iterator/iterator.go).

Iterators are accessed via an iterator session, which caches iterators so that already examined go structs and primitives don't need to be regenerated on every call (examining structs via reflection is slow). The iterators themselves are functions, and the cache itself stores generator functions that generate the iterator functions.

The [root iterator](iterator/iterator_root.go) acts as a top-level iterator, and coordinates iteration by constructing more specialized iterators depending on the object it's tasked with iterating over.


### Builders

Builders are the opposite of iterators, ingesting data events to produce go objects. Builders follow a common builder interface defined in [builder.go](builder/builder.go).

The root builder adapts data events to builder commands, which are then farmed out to more specialized builders to build go objects.

Builders are accessed via a builder session, which like the iterator session caches builders (due to the slowness of reflection). The cache itself stores builder generator functions.

The [reference filler](builder/reference_filler.go) maintains a list of outstanding [markers](https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#marker) and [references](https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#reference), filling in referenced data as it becomes available.


### Codecs

The codecs encode and decode data events into Concise Encoding documents. The [CBE](cbe) codec is relatively clean and fast due to the simplicity of the format, whereas the [CTE](cte) codec can get a bit hairy in some places. I've done my best to keep the code sane and followable, with varying levels of success.


### Rules

The [rules](rules/rules.go) class is structured as a data event receiver, and is designed to sit in between a codec and an iterator or builder to make sure the events contain valid data and happen in a valid order.



Secondary Sections
------------------

### Code Generation

The [codegen](codegen) directory contains all of the code to generate the more tedious parts of the library. To use it, simply run `go build` inside the [codegen](codegen) directory and then run `./codegen`. It will create/replace files in various places called `generated-do-not-edit.go`. To generate the Unicode character handling code, you'll also need the file `ucd.all.flat.xml` from https://www.unicode.org/Public/UCD/latest/ucdxml/ucd.all.flat.zip

### Debug Helpers

The [debug](debug) directory contains code to help with debugging. Internally, errors are handled via panics, which are then wrapped in error objects at the library boundary. However, in some cases having a stack trace leading up to the error can be very useful, which is where `PassThroughPanics` comes in handy.

### Test Helpers

The [test](test) directory contains code to help with writing and debugging tests.

- `PassThroughPanics()` allows you to selectively turn on [panic pass-through](#debug) for a single test without disrupting the other tests.
- Various constructors for common data types used in the tests.
- Panic/no panic assertions.
- Data generators.
- Data event constructors and generators for quickly building events.
- Event printer that sits in the middle of an event receiver chain and prints out the events passing through.
- Event receiver that converts events into objects representing those events, and an event driver that turns those objects back into events. Much of the code can be accessed locally via `testhelpers_test.go` files in the major subsections.

### Conversions

The [conversions](conversions) directory contains common type conversion functions used by various subsections. It's kept publically accessible so that user-defined codec code can make use of them.

### Options

All configuration options and their defaults are defined in the [options](options) directory.
