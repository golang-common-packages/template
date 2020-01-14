package manipulation

import (
	"math/rand"
)

// Client manage all manipulation function
type Client struct{}

// Shuffle shuffles the array using a random source
func (c *Client) Shuffle(array []interface{}, source rand.Source) {
	random := rand.New(source)
	for i := len(array) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}
