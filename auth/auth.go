/*
Authentication utilities for GoPics.

Copyright (c) 2015, Luca Chiricozzi. All rights reserved.
Released under the MIT License.
http://opensource.org/licenses/MIT
*/
package auth

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

const (
	defaultAge = 60 * 60 * 24 * 7 // One week
	authCookie = "AUTH"           // The name of the auth cookie
)

// A Keyring is used to encrypt and decrypt auth cookies.
type Keyring struct {
	*securecookie.SecureCookie
}

// NewKeyring creates a new Keyring with the given hashKey and blockKey.
// The hashKey is required, it has a length of 32 or 64 bytes.
// The blockKey is optional, it has a length of 16, 24 or 32 bytes.
func NewKeyring(hashKey, blockKey []byte) Keyring {
	return Keyring{securecookie.New(hashKey, blockKey)}
}

// newCookie creates a new authentication cookie
func newCookie(value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     authCookie,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
	}

}

// SetCookie creates a new authentication cookie.
func SetCookie(w http.ResponseWriter, k Keyring, value string) error {

	encoded, err := k.Encode(authCookie, value)
	if err != nil {
		return err
	}

	http.SetCookie(w, newCookie(encoded, defaultAge))
	return nil
}

// GetCookie reads the authentication cookie from a request and returns the
// authentication credentials in the cookie
func GetCookie(r *http.Request, k Keyring) (value string, err error) {
	c, err := r.Cookie(authCookie)
	if err != nil {
		return
	}

	err = k.Decode(authCookie, c.Value, &value)
	if err != nil {
		return
	}

	return
}

// DelCookie removes the authentication cookie.
func DelCookie(w http.ResponseWriter) {
	http.SetCookie(w, newCookie("", -1))
}
