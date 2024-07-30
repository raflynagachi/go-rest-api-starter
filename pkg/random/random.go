package random

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomInt generates a random integer between min and max
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// RandomID generates a random int64 between 1 to 1mil
func RandomID() int64 {
	return int64(RandomInt(1, 1000000000))
}

// RandomFloat generates a random float between min and max
func RandomFloat(min, max int) float64 {
	return float64(min + rand.Intn(max-min+1))
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(k)]
		sb.WriteByte(ch)
	}

	return sb.String()
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@mail.com", RandomString(6))
}
