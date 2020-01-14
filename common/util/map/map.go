package maptools

// Storage store function in otp package
type Storage interface {
	RemoveKeyFromMap(object interface{}, keys []string) interface{}
}
