package jwt

import (
	"container/ring"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math"
	"sync"
	"time"
)

type KeyStore struct {
	keys     *ring.Ring
	keysLock *sync.RWMutex
	issuer   string
	tokenTTL time.Duration
	keyTTL   *time.Duration
	ticker   *time.Ticker
}

const (
	DefaultKeyBitSize = 2048
)

func generateNewKey(ks *KeyStore, bitSize int) error {
	fmt.Println("generating a new key")
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return err
	}

	ks.keysLock.Lock()

	ks.keys = ks.keys.Next()
	ks.keys.Value = key

	ks.keysLock.Unlock()
	return nil
}

func NewSingleKeyStore(issuer string, tokenTTL time.Duration) (*KeyStore, error) {
	return NewKeyStore(issuer, tokenTTL, nil, DefaultKeyBitSize)
}

func NewKeyStore(issuer string, tokenTTL time.Duration, keyTTL *time.Duration, keyBitSize int) (*KeyStore, error) {
	var ks *KeyStore
	if keyTTL == nil {
		// single key, no rotation
		ks = &KeyStore{
			keys:     ring.New(1),
			keysLock: &sync.RWMutex{},
			issuer:   issuer,
			tokenTTL: tokenTTL,
		}
	} else {
		var numKeys int
		if tokenTTL.Milliseconds() < keyTTL.Milliseconds() {
			numKeys = 2
		} else {
			numKeys = int(math.Ceil(float64(tokenTTL.Milliseconds()) / float64(keyTTL.Milliseconds())))
		}

		ks = &KeyStore{
			keys:     ring.New(numKeys),
			keysLock: &sync.RWMutex{},
			issuer:   issuer,
			tokenTTL: tokenTTL,
			keyTTL:   keyTTL,
			ticker:   time.NewTicker(*keyTTL),
		}
	}

	if err := generateNewKey(ks, keyBitSize); err != nil {
		return nil, err
	}

	if ks.ticker != nil {
		go func() {
			for range ks.ticker.C {
				if err := generateNewKey(ks, keyBitSize); err != nil {
					fmt.Printf("failed to generate a new key: %s", err)
				}
			}
		}()
	}

	return ks, nil
}

func (ks *KeyStore) currentKey() *rsa.PrivateKey {
	ks.keysLock.RLock()

	curr := ks.keys.Value.(*rsa.PrivateKey)

	ks.keysLock.RUnlock()
	return curr
}

func (ks *KeyStore) allKeys() []*rsa.PrivateKey {
	ks.keysLock.RLock()

	keys := make([]*rsa.PrivateKey, 0)
	ks.keys.Do(func(key interface{}) {
		if key != nil {
			keys = append(keys, key.(*rsa.PrivateKey))
		}
	})

	ks.keysLock.RUnlock()
	return keys
}
