package encoding

import (
	"encoding/base32"
	"encoding/base64"
)

// Client manager all encoding function
type Client struct{}

// Base32Encode base32 encode
func (c *Client) Base32Encode(data []byte) string {
	return base32.StdEncoding.EncodeToString(data)
}

// Base32Decode base32 decode
func (c *Client) Base32Decode(data string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(data)
}

// Base64Encode base64 encode
func (c *Client) Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode base64 decode
func (c *Client) Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
