package bconf

// FieldSetStruct is implemented by structs that provide their field set key.
type FieldSetStruct interface {
	FieldSet() string
}
