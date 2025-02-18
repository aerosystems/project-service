package entities

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

type Project struct {
	Uuid         uuid.UUID
	CustomerUUID uuid.UUID
	Name         string
	Token        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewProject(customerUuid uuid.UUID, name string) *Project {
	return &Project{
		Uuid:         uuid.New(),
		Token:        generateToken(),
		CustomerUUID: customerUuid,
		Name:         name,
	}
}

func generateToken() string {
	rand.Seed(time.Now().Unix())
	sum := sha256.Sum256([]byte(strconv.Itoa(rand.Int())))
	return fmt.Sprintf("%x", sum)
}
