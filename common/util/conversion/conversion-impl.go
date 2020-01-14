package conversion

import (
	"bytes"
	"encoding/json"
)

const (
	empty = ""
	tab   = "\t"
)

// Client manage all convert function
type Client struct{}

// Stringify returns a string representation
func (c *Client) Stringify(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return empty, err
	}
	return string(b), nil
}

// Structify returns the original representation
func (c *Client) Structify(data string, value interface{}) error {
	return json.Unmarshal([]byte(data), value)
}

// PrettyJSON returns a pretty json string
func (c *Client) PrettyJSON(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(empty, tab)

	err := encoder.Encode(data)
	if err != nil {
		return empty, err
	}
	return buffer.String(), nil
}
