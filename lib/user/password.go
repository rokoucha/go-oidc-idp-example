package user

import "golang.org/x/crypto/argon2"

const (
	time    = 2
	memory  = 15 * 1024
	threads = 1
	keyLen  = 32
)

func hash(salt []byte, password string) []byte {
	return argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
}
