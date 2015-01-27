/*
A GoPics' post.

Copyright (c) 2015, Luca Chiricozzi. All rights reserved.
Released under the MIT License.
http://opensource.org/licenses/MIT
*/
package main

// An user's post
type Post struct {
	AuthorName   string `redis:"author_name"`
	AuthorPicURL string `redis:"author_pic_url"`
	Name         string `redis:"name"`
	Text         string `redis:"text"`
	Time         string `redis:"time"`
}
