package utils

import (
	"github.com/google/uuid"
	"math/rand"
)

func GenerateConfirmationCode() int {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	var code = arr[rand.Int()%len(arr)-1]
	for i := 0; i < 4; i++ {
		code *= 10 + arr[rand.Int()%len(arr)]
	}
	return code
}

func GenerateUUIDString() string {
	return uuid.NewString()
}
