package helpers

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// GenerateToken returns a unique token
func GenerateToken() string {
	rand.Seed(time.Now().Unix())
	sum := sha256.Sum256([]byte(strconv.Itoa(rand.Int())))
	return fmt.Sprintf("%x", sum)
}
