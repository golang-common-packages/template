package errhandling

// Storage interface for errhandling package
type Storage interface {
	New(text string) error
	DefaultErrorIfNil(err error, defaultMessage string) string
}
