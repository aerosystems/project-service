package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// GenerateToken returns a unique token
func GenerateToken() string {
	rand.Seed(time.Now().UnixNano())
	hash, err := bcrypt.GenerateFromPassword([]byte(strconv.Itoa(rand.Int())), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Example rand int: ", strconv.Itoa(rand.Int()))
	fmt.Println("Hash to store:", string(hash))

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}
