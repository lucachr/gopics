/*
A GoPics' user, with convenience methods.

Copyright (c) 2015, Luca Chiricozzi. All rights reserved.
Released under the MIT License.
http://opensource.org/licenses/MIT
*/
package main

import (
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/lucachr/gopics/reutils"
	"github.com/ungerik/go-gravatar"
)

// An user of the image board
type User struct {
	Name     string `redis:"name"`
	Email    string `redis:"email"`
	Password []byte `redis:"password"`
	PicURL   string `redis:"pic_url"`
	Posts    []Post `redis:"-"`
}

// validate is a convenience method for validating user data.
func (usr *User) validate(conn redis.Conn) error {
	if !reutils.MatchName(usr.Name) {
		return ErrValidation("Your username is invalid!")
	}

	for _, name := range invalidUser {
		if strings.Contains(usr.Name, name) {
			return ErrValidation("You cannot choose that name!")
		}
	}

	reg, err := redisGetUser(conn, usr.Name)
	if err != nil && err != redis.ErrNil {
		return err
	}
	if reg != nil {
		return ErrValidation("A user with the same name already exist!")
	}

	if !reutils.MatchEmail(usr.Email) {
		return ErrValidation("Your email is invalid!")
	}

	if len(usr.Password) < 8 {
		return ErrValidation("Your password is too short!")
	}

	return nil
}

// save is a convenience method for adding a new user.
func (usr *User) save(conn redis.Conn) error {
	usr.PicURL = gravatar.Url(usr.Email)

	_, err := conn.Do("HMSET", redisFlat(userTag+usr.Name, usr)...)
	return err
}
