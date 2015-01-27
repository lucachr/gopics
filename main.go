/*
An tutorial image sharing platform built with Go, Redis, UIkit and good intentions.

Copyright (c) 2015, Luca Chiricozzi. All rights reserved.
Released under the MIT License.
http://opensource.org/licenses/MIT
*/
package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Create a new Redis pool
	flag.Parse()
	pool = newPool(*redisServer)

	http.Handle("/", appHandler(handleRoot))
	http.Handle("/register", appHandler(handleRegister))
	http.Handle("/registration", redisHandler(handleRegistration))
	http.Handle("/login", redisHandler(handleLogin))
	http.HandleFunc("/logout", handleLogout)
	http.Handle("/post", redisHandler(handlePost))

	media := http.FileServer(http.Dir(filepath.Join(mediaPath...)))
	http.Handle("/media/", http.StripPrefix("/media/", media))

	static := http.FileServer(http.Dir(filepath.Join(staticPath...)))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
