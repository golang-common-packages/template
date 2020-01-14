package maptools

import (
	"github.com/fatih/structs"
)

// Client manage all mapping action
type Client struct{}

// RemoveKeyFromMap function will remove "keys" from map
func (c *Client) RemoveKeyFromMap(object interface{}, keys []string) interface{} {
	objectMap := structs.Map(object)
	for keyMap := range objectMap {
		for keyTruct := range keys {
			if keyMap == keys[keyTruct] {
				delete(objectMap, keyMap)
			}
		}
	}

	return objectMap
}
