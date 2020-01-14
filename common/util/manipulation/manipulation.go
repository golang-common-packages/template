package manipulation

import (
	"math/rand"
)

// Storage interface for manipulation package
type Storage interface {
	Shuffle(array []interface{}, source rand.Source)
}
