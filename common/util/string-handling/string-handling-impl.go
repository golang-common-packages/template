package strhandling

import (
	"bytes"
	"strings"

	"github.com/shomali11/util/xrunes"
)

const (
	space = " "
)

// Client manage all strhandling function
type Client struct{}

// IsEmpty returns true if the string is empty
func (c *Client) IsEmpty(text string) bool {
	return len(text) == 0
}

// IsNotEmpty returns true if the string is not empty
func (c *Client) IsNotEmpty(text string) bool {
	return !IsEmpty(text)
}

// IsBlank returns true if the string is blank (all whitespace)
func (c *Client) IsBlank(text string) bool {
	return len(strings.TrimSpace(text)) == 0
}

// IsNotBlank returns true if the string is not blank
func (c *Client) IsNotBlank(text string) bool {
	return !IsBlank(text)
}

// Left justifies the text to the left
func (c *Client) Left(text string, size int) string {
	spaces := size - Length(text)
	if spaces <= 0 {
		return text
	}

	var buffer bytes.Buffer
	buffer.WriteString(text)

	for i := 0; i < spaces; i++ {
		buffer.WriteString(space)
	}
	return buffer.String()
}

// Right justifies the text to the right
func (c *Client) Right(text string, size int) string {
	spaces := size - Length(text)
	if spaces <= 0 {
		return text
	}

	var buffer bytes.Buffer
	for i := 0; i < spaces; i++ {
		buffer.WriteString(space)
	}

	buffer.WriteString(text)
	return buffer.String()
}

// Center justifies the text in the center
func (c *Client) Center(text string, size int) string {
	left := Right(text, (Length(text)+size)/2)
	return Left(left, size)
}

// Length counts the input while respecting UTF8 encoding and combined characters
func (c *Client) Length(text string) int {
	textRunes := []rune(text)
	textRunesLength := len(textRunes)

	sum, i, j := 0, 0, 0
	for i < textRunesLength && j < textRunesLength {
		j = i + 1
		for j < textRunesLength && xrunes.IsMark(textRunes[j]) {
			j++
		}
		sum++
		i = j
	}
	return sum
}

// Reverse reverses the input while respecting UTF8 encoding and combined characters
func (c *Client) Reverse(text string) string {
	textRunes := []rune(text)
	textRunesLength := len(textRunes)
	if textRunesLength <= 1 {
		return text
	}

	i, j := 0, 0
	for i < textRunesLength && j < textRunesLength {
		j = i + 1
		for j < textRunesLength && xrunes.IsMark(textRunes[j]) {
			j++
		}

		if xrunes.IsMark(textRunes[j-1]) {
			// Reverses Combined Characters
			reverse(textRunes[i:j], j-i)
		}

		i = j
	}

	// Reverses the entire array
	reverse(textRunes, textRunesLength)

	return string(textRunes)
}

// ByteArrayToString convert byte array to string
func (c *Client) ByteArrayToString(s []byte) string {
	n := -1
	for i, b := range s {
		if b == 0 {
			break
		}
		n = i
	}
	return string(s[:n+1])
}

func reverse(runes []rune, length int) {
	for i, j := 0, length-1; i < length/2; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
}
