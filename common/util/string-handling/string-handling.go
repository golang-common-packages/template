package strhandling

// Storage interface for strhandling package
type Storage interface {
	IsEmpty(text string) bool
	IsNotEmpty(text string) bool
	IsBlank(text string) bool
	IsNotBlank(text string) bool
	Left(text string, size int) string
	Right(text string, size int) string
	Center(text string, size int) string
	Length(text string) int
	Reverse(text string) string
	ByteArrayToString(c []byte) string
}
