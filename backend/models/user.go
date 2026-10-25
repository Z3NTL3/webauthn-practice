package models

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

var (
	_ webauthn.User = User{}
)

type User struct {
	ID          []byte
	Email       string
	Credentials []webauthn.Credential
	webauthn.User
}

func NewUser(email string) *User {
	id, _ := uuid.NewRandom()
	return &User{ID: id[:], Email: email}
}

func (u User) WebAuthnID() []byte {
	return u.ID
}

func (u User) WebAuthnName() string {
	return u.Email
}

func (u User) WebAuthnDisplayName() string {
	return u.Email
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}
