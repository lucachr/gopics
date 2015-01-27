/*
Flash cookies for GoPics.

Copyright (c) 2015, Luca Chiricozzi. All rights reserved.
Released under the MIT License.
http://opensource.org/licenses/MIT
*/
package flash

import "net/http"

const cookieName = "FLASH"

// SetCookie set a  flash cookie.
func SetCookie(w http.ResponseWriter, value string) error {
	c := &http.Cookie{Name: cookieName, Value: value, HttpOnly: true}
	http.SetCookie(w, c)
	return nil
}

// GetCookie get a flash cookie and delete it.
func GetCookie(w http.ResponseWriter, r *http.Request) (value string,
	err error) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return
	}
	value = c.Value

	// Remove the flash cookie
	c = &http.Cookie{Name: cookieName, Value: "", MaxAge: -1}
	http.SetCookie(w, c)
	return
}
