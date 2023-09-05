package document

// Document represents a single document in the search engine
type Document struct {
	ID     string
	Fields map[string]Field
}

// NewDocument creates a new document instance
func NewDocument(id string) *Document {
	return &Document{
		ID:     id,
		Fields: make(map[string]Field),
	}
}

// AddField adds a field to the document
func (d *Document) AddField(name string, value interface{}) {
	d.Fields[name] = Field{
		Type:  TextField,
		Value: value,
	}
}
