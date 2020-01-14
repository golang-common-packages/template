package compression

// Storage store function in compression package
type Storage interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}
