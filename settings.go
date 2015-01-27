/*
Settings and configuration for GoPics.

Copyright (c) 2015, Luca Chiricozzi. All rights reserved.
Released under the MIT License.
http://opensource.org/licenses/MIT
*/
package main

import (
	"flag"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/lucachr/gopics/auth"
)

const (
	pageTitle = "GoPics | "

	timeLayout = "Mon 2 Jan 2006 15:04"

	// Redis "tags" for users and posts data.
	userTag      = "user:"
	userTimeline = "timeline:"
	postTag      = "post:"
)

var (
	// Keyring with hash key and block key for authentication cookies.
	keyring = auth.NewKeyring(
		[]byte("afeee502-a34f-11e4-a451-902b34a8c90ffd8c-a34f-11e4-a1d3-902b34a8"),
		[]byte("e12b872e-a34f-11e4-9c70-902b34a8"),
	)

	// Invalid usernames.
	invalidUser = []string{
		"index",
		"register",
		"login",
		"logout",
		"registration",
		"post",
		"media",
		"static",
	}

	pool        *redis.Pool
	redisServer = flag.String("redisServer", redisDefaultAddr, "")

	// A slice with the path of your media directory
	basePath = []string{os.Getenv("GOPATH"), "src", "github.com",
		"lucachr", "gopics"}
	mediaPath     = append(basePath, "media")
	staticPath    = append(basePath, "static")
	templatesPath = append(basePath, "templates")

	templates = buildTemplates(
		"header.html",
		"index.html",
		"register.html",
		"timeline.html",
		"footer.html",
	)
)
