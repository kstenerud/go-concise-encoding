// Package concise_encoding implements a concise encoding codec and marshaler.
// https://concise-encoding.org/
//
// Most people will only care about the highest level APIs for mashaling in
// package ce.
//
// If all you're interested in is (de)serializing to go objects, the marshaler
// API is sufficient. The codecs provide more control over the process, and
// can handle more data types (such as comments). The event
// handlers, builders and iterators are the lowest level API, providing maximum
// control but the highest complexity.
//
// The primary architecture design is one of filtered message pipelines,
// consisting of data events (which report what kind of data is encountered),
// and builder directives (which direct the parts of a complex data structure
// that is to be built). All software components are designed around this
// principle to promote interchangeability.
//
//
// High Level API (package ce)
//
// * Marshalers: (de)serializes to/from go objects.
//
//
// Medium Level API (package ce)
//
// * Encoder: Accepts data events and generates a CBE or CTE encoded document.
//
// * Decoder: Decodes a CBE or CTE document, generating data events
//
//
// Low Level API (packages builder, events, iterator, rules)
//
// * Iterator: Iterates through an object, generating data events.
//
// * DataEventReceiver: Receives data events and acts upon them.
//
// * Builder: DataEventReceiver that builds objects in response to events.
//
// * Rules: DataEventReceiver that validates events, ensuring their contents
// and order match a valid CBE/CTE document.
package concise_encoding
