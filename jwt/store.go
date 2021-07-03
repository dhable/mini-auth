package jwt

import (
	"container/ring"
	"crypto/rand"
	"crypto/rsa"
	"sync"
	"time"
)

type KeyStore struct {
	keys     *ring.Ring
	keysLock *sync.RWMutex
	issuer   string
	tokenTTL time.Duration
}

const (
	DefaultKeyBitSize = 2048
)

func generateNewKey(ks *KeyStore, bitSize int) error {
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

func SingleKeyWithBitSize(issuer string, tokenTTL time.Duration, keyBitSize int) (*KeyStore, error) {
	ks := &KeyStore{
		keys:     ring.New(1),
		keysLock: &sync.RWMutex{},
		issuer:   issuer,
		tokenTTL: tokenTTL,
	}

	if err := generateNewKey(ks, keyBitSize); err != nil {
		return nil, err
	}

	return ks, nil
}

func SingleKey(issuer string, tokenTTL time.Duration) (*KeyStore, error) {
	return SingleKeyWithBitSize(issuer, tokenTTL, DefaultKeyBitSize)
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
