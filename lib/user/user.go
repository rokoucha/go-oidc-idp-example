package user

import (
	"crypto/subtle"
	"errors"
	"slices"

	"github.com/oklog/ulid/v2"
)

type UserInfo struct {
	ID       ulid.ULID
	Username string
	password []byte
}

type User struct {
	salt  []byte
	users []UserInfo
}

func New(salt []byte) *User {
	return &User{
		salt:  salt,
		users: []UserInfo{},
	}
}

func (u *User) Get(id ulid.ULID) (UserInfo, bool) {
	i := slices.IndexFunc(u.users, func(u UserInfo) bool {
		return u.ID == id
	})
	if i == -1 {
		return UserInfo{}, false
	}

	return u.users[i], true
}

func (u *User) Register(username, password string) error {
	if slices.IndexFunc(u.users, func(u UserInfo) bool {
		return u.Username == username
	}) != -1 {
		return errors.New("username already exists")
	}

	u.users = append(u.users, UserInfo{
		ID:       ulid.Make(),
		Username: username,
		password: hash(u.salt, password),
	})

	return nil
}

func (u *User) Authenticate(username, password string) (UserInfo, bool) {
	i := slices.IndexFunc(u.users, func(u UserInfo) bool {
		return u.Username == username
	})
	if i == -1 {
		return UserInfo{}, false
	}

	user := u.users[i]

	return user, subtle.ConstantTimeCompare(hash(u.salt, password), user.password) == 1
}
