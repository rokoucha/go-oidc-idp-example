package keychain

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"

	"github.com/go-jose/go-jose/v3"
)

type Keychain struct {
	keySet *jose.JSONWebKeySet
}

func New() *Keychain {
	return &Keychain{
		keySet: &jose.JSONWebKeySet{},
	}
}

func (k *Keychain) Get(kid string) *jose.JSONWebKey {
	keys := k.keySet.Key(kid)
	if len(keys) == 0 {
		return nil
	}
	return &keys[0]
}

func (k *Keychain) GetAll() []jose.JSONWebKey {
	return k.keySet.Keys
}

func (k *Keychain) Add(jwk jose.JSONWebKey) {
	k.keySet.Keys = append(k.keySet.Keys, jwk)
}

func (k *Keychain) Create(kid string) (jose.JSONWebKey, error) {
	if k.Get(kid) != nil {
		return jose.JSONWebKey{}, errors.New("key already exists")
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return jose.JSONWebKey{}, err
	}

	jwk := jose.JSONWebKey{
		Key:       privateKey,
		KeyID:     kid,
		Algorithm: string(jose.RS256),
	}

	k.Add(jwk)

	return jwk, nil
}
