package support

import (
	"crypto/rand"
	"math/big"
	"time"
)

func SleepRandomBetween(min, max time.Duration) {
	difference := max - min

	randomOffset, err := rand.Int(rand.Reader, big.NewInt(int64(difference)))
	if err != nil {
		panic(err)
	}

	randomDuration := min + time.Duration(randomOffset.Int64())

	time.Sleep(randomDuration * time.Millisecond)
}
