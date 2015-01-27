/*
Utilities for common regexp patterns

Copyright 2014 Luca Chiricozzi. All rights reserved.
Released under the MIT License.
http://opensource.org/licenses/MIT
*/
package reutils

import "regexp"

// MatchName checks whether a string is a valid name.
func MatchName(s string) bool {
	exp := regexp.MustCompile("^[\\pL\\pN !?'.-]+$")
	return exp.MatchString(s)
}

// MatchEmail checks whether a string is a valid email.
func MatchEmail(s string) bool {
	exp := regexp.MustCompile("^[a-zA-Z0-9+&*-]+(?:\\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,7}$")
	return exp.MatchString(s)
}
