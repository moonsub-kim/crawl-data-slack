package logger

// Field definition
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates field
func NewField(key string, value interface{}) Field {
	return Field{key, value}
}

// NewErrorField creates field containing error
func NewErrorField(err error) Field {
	return Field{"error", err}
}
