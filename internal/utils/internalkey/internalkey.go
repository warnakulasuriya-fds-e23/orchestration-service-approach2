package internalkey

import (
	"math/rand"
	"sync"
)

type internalAPIKey struct {
	key string
}

var (
	instance *internalAPIKey
	once     sync.Once
)

func GetInternalAPIKey() string {
	once.Do(func() {
		//randomly generate a key
		const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		b := make([]byte, 15)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		randomKey := string(b)
		instance = &internalAPIKey{
			key: randomKey,
		}
	})
	return instance.key
}
