package document

// Field represents a document's field with its type and value
type Field struct {
	Type  FieldType
	Value interface{}
}

// FieldType represents the type of a field
type FieldType int

const (
	// TextField represents a text field type
	TextField FieldType = iota
	// NumericField represents a numeric field type
	NumericField
)

// NewField creates a new field with given type and value
func NewField(t FieldType, value interface{}) Field {
	return Field{
		Type:  t,
		Value: value,
	}
}

// GetValue returns the value of the field
func (f *Field) GetValue() interface{} {
	return f.Value
}
