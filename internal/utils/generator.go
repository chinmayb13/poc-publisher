package utils

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func GetRandomStrings(count int) []string {
	var randStr []string
	for i := 0; i < count; i++ {
		randStr = append(randStr, uuid.New().String())
	}
	return randStr
}

func GetRandomSequence(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
