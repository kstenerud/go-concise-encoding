// Package concise_encoding provides codecs and marshalers for Concise Binary
// Encoding and Concise Text Encoding.
//
// If all you're interested in is (de)serializing to go objects, the marshaler
// API is sufficient. The codecs provide more control over the process, and
// can handle more data types (such as comments and metadata). The builders and
// iterators are the lowest level API, providing maximum control but higher
// complexity.
//
// The primary architecture design is one of filtered message pipelines,
// consisting of data events (which report what kind of data is encountered),
// and builder directives (which direct the parts of a complex data structure
// that is to be built). All software pieces are designed around this principle
// to keep the components interchangeable.
//
// High Level API:
// - Marshalers: (de)serializes to/from go objects.
//
// Medium Level API:
// - Encoder: Accepts data events and generates a CBE or CTE encoded document.
// - Decoder: Decodes a CBE or CTE document, generating data events.
//
// Low Level API:
// - Iterator: Iterates through an object, generating data events.
// - Rules: Passes through data events, applying rules to ensure the events
//          conform to the concise encoding specifications.
// - Root Iterator: Iterates through a top-level object to generate data events.
// - Builder: Responds to builder directives to build complex objects.
// - Root Builder: Translates data events to builder directives.
package concise_encoding
