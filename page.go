/*
A GoPics' page.

Copyright (c) 2015, Luca Chiricozzi. All rights reserved.
Released under the MIT License.
http://opensource.org/licenses/MIT
*/
package main

// A GoPics page.
type Page struct {
	Title      string
	User       *User
	LoggedUser string // Username of the logged user
	ValError   string // Validation error message
}
