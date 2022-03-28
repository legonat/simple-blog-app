package pkg

import (
	"math/rand"
	"time"
)

func RandomIntFromRange(min, max int) (n int) {
	rand.Seed(time.Now().UnixNano())
	n = rand.Intn(max-min) + min
	return
}