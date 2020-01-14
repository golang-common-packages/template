package rune

import (
	"unicode"
)

// Client manage all rune function
type Client struct{}

// IsMark determines whether the rune is a marker
func (c *Client) IsMark(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Mc, r)
}
