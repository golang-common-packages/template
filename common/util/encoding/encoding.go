package encoding

// Storage interface for encoding package
type Storage interface {
	Base32Encode(data []byte) string
	Base32Decode(data string) ([]byte, error)
	Base64Encode(data []byte) string
	Base64Decode(data string) ([]byte, error)
}
