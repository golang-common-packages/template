package conversion

// Storage interface for conversion package
type Storage interface {
	Stringify(data interface{}) (string, error)
	Structify(data string, value interface{}) error
	PrettyJSON(data interface{}) (string, error)
}
