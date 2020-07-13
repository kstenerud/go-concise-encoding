// Codecs and marshalers to produce and consume Concise Binary Encoding and
// Concise Text Encoding documents (https://github.com/kstenerud/concise-encoding).
//
// Most people will only care about the codecs or marshalers in the cbe and cte
// packages.
//
// If all you're interested in is (de)serializing to go objects, the marshaler
// API is sufficient. The codecs provide more control over the process, and
// can handle more data types (such as comments and metadata). The event
// handlers, builders and iterators are the lowest level API, providing maximum
// control but higher complexity.
//
// The primary architecture design is one of filtered message pipelines,
// consisting of data events (which report what kind of data is encountered),
// and builder directives (which direct the parts of a complex data structure
// that is to be built). All software components are designed around this
// principle to promote interchangeability.
//
// High Level API
//
// * Marshalers: (de)serializes to/from go objects.
//
// Medium Level API
//
// * Encoder: Accepts data events and generates a CBE or CTE encoded document.
//
// * Decoder: Decodes a CBE or CTE document, generating data events
//
// Low Level API
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
