package utils

import (
	"fmt"
	"math/rand"
)

//NewID returns a random string that can be used as an ID. this implemetation is for tests purposes only
func NewID() string {
	return fmt.Sprintf("%v", rand.Int())
}
