package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateUUID() string {
	timestamp := fmt.Sprintf("%16x", time.Now().UnixNano())
	random := fmt.Sprintf("%16x", rand.Uint64())
	uuidStr := fmt.Sprint(timestamp, random)
	uuid := fmt.Sprintf("%s-%s-%s-%s-%s", uuidStr[:8], uuidStr[8:12], uuidStr[12:16], uuidStr[16:20], uuidStr[20:])
	return uuid
}
