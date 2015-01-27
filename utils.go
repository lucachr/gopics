/*
Generic utility functions for GoPics.

Copyright (c) 2015, Luca Chiricozzi. All rights reserved.
Released under the MIT License.
http://opensource.org/licenses/MIT
*/

package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/lucachr/gopics/flash"
)

// renderTemplate executes a template with the data contained in the given
// page.
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) *appError {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return nil
}

// httpAppError send an error reponse to the user if
// is not nil.
func httpAppError(w http.ResponseWriter, ae *appError) {
	if ae != nil {
		http.Error(w, ae.Error(), ae.Code)
	}
}

// unixTimeNow returns the current Unix time.
func unixTimeNow() int64 {
	return time.Now().Unix()
}

// buildFilePath builds the absolute path to file, given the directory
// tree as a slice of strings.
func buildFilePath(path []string, file string) string {
	return filepath.Join(append(path, file)...)
}

// buildTemplates parses the given templates.
func buildTemplates(filenames ...string) *template.Template {
	// Build a slice of paths
	fs := []string{}
	for _, fn := range filenames {
		f := buildFilePath(templatesPath, fn)
		fs = append(fs, f)
	}

	return template.Must(template.ParseFiles(fs...))
}

// setFlashAndRedirect sets a flash cookie with a value of msg and redirects
// the user to the given URL.
func setFlashAndRedirect(w http.ResponseWriter, r *http.Request,
	url, msg string) *appError {
	err := flash.SetCookie(w, msg)
	if err != nil {
		return &appError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
	return nil
}
