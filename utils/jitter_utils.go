package foundation

import (
	"math/rand"
	"time"
)

var (
	// seed random number
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// ApplyJitter adds +-25% jitter to the input
func ApplyJitter(input int) int {
	deviation := int(0.25 * float64(input))
	return input - deviation + r.Intn(2*deviation)
}
