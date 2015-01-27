/*
Requests handlers for GoPics.

Copyright (c) 2015, Luca Chiricozzi. All rights reserved.
Released under the MIT License.
http://opensource.org/licenses/MIT
*/
package main

import (
	"image"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/garyburd/redigo/redis"
	"github.com/lucachr/gopics/auth"
	"github.com/lucachr/gopics/flash"
	"github.com/nfnt/resize"
	"golang.org/x/crypto/bcrypt"
)

const (

	// Max picture sizes
	maxPicBytes = 2097152 // 2MB
	maxWidth    = 800
	maxHeight   = 600
)

// appHandler is an handler that takes a Page and returns a pointer to an
// appError.
type appHandler func(http.ResponseWriter, *http.Request, *Page) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check for validation error in the form
	msg, err := flash.GetCookie(w, r)
	switch {
	case err == http.ErrNoCookie: // Do nothing
	case err != nil:
		err := &appError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
		httpAppError(w, err)
	}

	p := new(Page)
	p.ValError = msg
	httpAppError(w, fn(w, r, p))
}

// redisHandler is a request handler that needs a connection to Redis and
// returns a pointer to an appError.
type redisHandler func(http.ResponseWriter, *http.Request, redis.Conn) *appError

func (fn redisHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn := pool.Get()
	defer conn.Close()

	httpAppError(w, fn(w, r, conn))
}

// login sets an auth cookie with the given username and redirect the
// user to her home page.
func login(w http.ResponseWriter, r *http.Request, username string) *appError {
	err := auth.SetCookie(w, keyring, username)
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	http.Redirect(w, r, "/"+username, http.StatusSeeOther)
	return nil
}

// handleRegistration handles the registration of a new user to GoPics.
func handleRegistration(w http.ResponseWriter, r *http.Request,
	conn redis.Conn) *appError {
	// Create a new user an get user's detail from the form
	usr := new(User)
	usr.Name = r.FormValue("name")
	usr.Email = r.FormValue("email")
	usr.Password = []byte(r.FormValue("password"))

	// Validate the user credentials.
	err := usr.validate(conn)
	switch err.(type) {
	case ErrValidation:
		return setFlashAndRedirect(w, r, "/register", err.Error())
	case nil: // Do nothing
	default:
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// Hash the user password
	pass, err := bcrypt.GenerateFromPassword(usr.Password,
		bcrypt.DefaultCost)
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	usr.Password = pass

	// All right, register the new user.
	if err = usr.save(conn); err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return login(w, r, usr.Name)
}

// handleLogin manages the login of the users
func handleLogin(w http.ResponseWriter, r *http.Request,
	conn redis.Conn) *appError {

	// Get user credential from Redis
	usr, err := redisGetUser(conn, r.FormValue("name"))
	switch {
	case err == redis.ErrNil:
		// The user does not exist.
		return setFlashAndRedirect(w, r, "/",
			"Invalid username or password.")
	case err != nil:
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// Check if the submitted password and the user's one match.
	err = bcrypt.CompareHashAndPassword(usr.Password,
		[]byte(r.FormValue("password")))
	switch {
	case err == bcrypt.ErrMismatchedHashAndPassword:
		// The passwords don't match
		return setFlashAndRedirect(w, r, "/",
			"Invalid username or password.")
	case err != nil:
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// The credentials are ok, go on with login.
	return login(w, r, usr.Name)
}

// handleLogut removes the auth cookie and redirect to
// the root URL.
func handleLogout(w http.ResponseWriter, r *http.Request) {
	auth.DelCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handleRoot manages the app root url, if an index page is required,
// it calls the handleIndex function, otherwise, it calls the
// handleTimeline function handling the path after the root URL
// as an username.
func handleRoot(w http.ResponseWriter, r *http.Request, p *Page) *appError {
	path := r.URL.Path[len("/"):]
	if path == "" || strings.Contains(path, "index") {
		return handleIndex(w, r, p)
	}

	return handleTimeline(w, r, p, path)
}

// handleIndex manages the index page
func handleIndex(w http.ResponseWriter, r *http.Request, p *Page) *appError {
	// Check if an user is logged
	username, err := auth.GetCookie(r, keyring)
	switch {
	case err == http.ErrNoCookie: //Do nothing
	case err != nil:
		return &appError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	case username != "":
		// If an user is found, redirect she to her home
		http.Redirect(w, r, "/"+username, http.StatusSeeOther)
		return nil
	}

	// Render the index
	p.Title = pageTitle + "Welcome!"
	return renderTemplate(w, "index", p)
}

// handleRegister manages the sign up page
func handleRegister(w http.ResponseWriter, r *http.Request, p *Page) *appError {
	// Display the page
	p.Title = pageTitle + "Sign Up"
	return renderTemplate(w, "register", p)
}

// handleTimeLine manages the users' timelines
func handleTimeline(w http.ResponseWriter, r *http.Request, p *Page,
	username string) *appError {
	conn := pool.Get()
	defer conn.Close()

	// Get user's data from Redis
	usr, err := redisGetUser(conn, username)
	switch {
	case err == redis.ErrNil:
		http.NotFound(w, r)
		return nil
	case err != nil:
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// Create the timeline of the user.
	usr.Posts, err = redisGetPosts(conn, userTimeline+usr.Name)
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// If an user is logged, get her name.
	logName, err := auth.GetCookie(r, keyring)
	if err != nil && err != http.ErrNoCookie {
		return &appError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	// Set the page data and display it
	p.Title = pageTitle + username
	p.User = usr
	p.LoggedUser = logName

	return renderTemplate(w, "timeline", p)
}

// handlePost manages posts submission.
func handlePost(w http.ResponseWriter, r *http.Request,
	conn redis.Conn) *appError {
	var img image.Image

	// Get the username from the auth cookie
	username, err := auth.GetCookie(r, keyring)
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	// Check the content lenght
	switch {
	case r.ContentLength == -1:
		return &appError{
			Err:  ErrInvalidLength,
			Code: http.StatusLengthRequired,
		}
	case r.ContentLength > maxPicBytes:
		return &appError{
			Err:  ErrInvalidLength,
			Code: http.StatusRequestEntityTooLarge,
		}
	}

	// Read only the first maxPicBytes of the request's body
	r.Body = http.MaxBytesReader(w, r.Body, maxPicBytes)

	// Try to get the content of the form picture field
	f, _, err := r.FormFile("picture")
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	// Try to decode the content of f as an image
	src, _, err := image.Decode(f)
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusUnsupportedMediaType,
		}
	}

	// Get the ratio d of the image
	bound := src.Bounds()
	x, y := bound.Max.X, bound.Max.Y
	d := float32(x) / float32(y)

	// Check the image sizes
	if x > maxWidth || y > maxHeight {
		if x > y {
			img = resize.Resize(uint(maxWidth), uint(1/d*maxWidth),
				src, resize.Lanczos3)

		} else {
			img = resize.Resize(uint(d*maxHeight), uint(maxHeight),
				src, resize.Lanczos3)
		}
	} else {
		img = resize.Resize(uint(x), uint(y), src, resize.Lanczos3)

	}

	// The image name is generated as an uuid
	picName := uuid.New() + ".jpeg"
	path := append(mediaPath, picName)

	// Create a new file
	dst, err := os.Create(filepath.Join(path...))
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	defer dst.Close()

	// Write the image in the file
	err = jpeg.Encode(dst, img, nil)
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// Build the post
	p := new(Post)
	p.Name = picName
	p.Text = r.FormValue("text")
	p.Time = time.Now().Format(timeLayout)

	// Get the author data from Redis
	usr, err := redisGetUser(conn, username)
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// Add the author data
	p.AuthorName = usr.Name
	p.AuthorPicURL = usr.PicURL

	// Create the post and add it to the user timeline
	conn.Send("MULTI")
	conn.Send("HMSET", redisFlat(postTag+p.Name, p)...)
	conn.Send("ZADD", userTimeline+usr.Name, unixTimeNow(), p.Name)
	_, err = conn.Do("EXEC")
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// All right, redirect to the home.
	http.Redirect(w, r, "/"+usr.Name, http.StatusSeeOther)
	return nil
}
