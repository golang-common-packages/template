package rune

// Storage interface for rune package
type Storage interface {
	IsMark(r rune) bool
}
