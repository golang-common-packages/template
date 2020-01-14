package condition

// Storage interface for condition package
type Storage interface {
	IfThen(condition bool, a interface{}) interface{}
	IfThenElse(condition bool, a interface{}, b interface{}) interface{}
	DefaultIfNil(value interface{}, defaultValue interface{}) interface{}
	FirstNonNil(values ...interface{}) interface{}
}
