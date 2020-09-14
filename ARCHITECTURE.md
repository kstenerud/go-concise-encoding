Architecture
============

This document describes the architecture of this project. As this is the reference implementation of [Concise Encoding](https://concise-encoding.org/), I've done by best to keep it as readable as possible.

The code is separated into five primary areas:

- [Events](#events): Contract that all major components operate by.
- [Iterators](#iterators): Iterators iterate through a go object to produce data events.
- [Builders](#builders): Builders interpret data events to produce go objects.
- [Codecs](#codecs): Codecs encode/decode data events to/from CTE or CBE documents.
- [Rules](#rules): Rules enforce proper structure and content in Concise Encoding documents.



Events
------

Events form the backbone of the entire library. The major components either consume or produce data events. This architecture makes it easy to mix & match components to produce whatever software design you want.

See: [data_events.go](events/data_events.go)



Iterators
---------

Iterators inspect go objects to produce data events. They support all primitives, as well as arrays, slices, maps, pointers, and structs, and can handle recursive pointers. All iterators follow a common interface defined in [iterator.go](iterator/iterator.go).

Iterators are accessed via an iterator session, which caches iterators so that already examined go structs and primitives don't need to be regenerated on every call (examining structs via reflection is slow).

The [root iterator](iterator/iterator_root.go) acts as a top-level iterator, and coordinates iteration by constructing more specialized iterators depending on the object it's tasked with iterating over.



Builders
--------

Builders are the opposite of iterators, ingesting data events to produce go objects. Builders follow a common builder interface defined in [builder.go](builder/builder.go).

Like the root iterator, the [root builder](builder/builder_root.go) acts as a top-level builder, building of complex go objects by coordinating more specialized builders. The root builder adapts data events to builder commands.

Builders are accessed via a builder session, which like the iterator session caches builders (once again due to the slowness of reflection).

The [reference filler](builder/reference_filler.go) maintains a list of outstanding [markers](https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#marker) and [references](https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#reference), filling in referenced data as it becomes available.



Codecs
------

The codecs encode and decode data events into Concise Encoding documents. The [CBE](cbe) codec is relatively clean and fast due to the simplicity of the format, whereas the [CTE](cte) codec can get a bit hairy in some places. I've done my best to keep the code sane and followable, with varying levels of success.



Rules
-----

The rules class is structured as a data event receiver, and is designed to sit in between a codec and an iterator or builder to make sure they're behaving properly. Since it's a data event receiver with passthrough, it can be placed anywhere, allowing for versatile Concise Encoding software designs.



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
