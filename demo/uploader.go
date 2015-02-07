// Copyright 2015 The httpserver Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.
package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/regiluze/httpserver"
)

var uploadTemplate = template.Must(template.ParseFiles("html/index.html"))
var ErrorTemplate = template.Must(template.ParseFiles("html/error500.html"))
var NotFoundTemplate = template.Must(template.ParseFiles("html/error404.html"))

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type ImageUploaderHandler struct {
}

func NewImageUploadHandler() *ImageUploaderHandler {

	iuh := &ImageUploaderHandler{}
	return iuh

}

func (iuh *ImageUploaderHandler) upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		uploadTemplate.Execute(w, nil)
		return
	}
	f, _, err := r.FormFile("image")
	check(err)
	defer f.Close()
	t, err := ioutil.TempFile(".", "/static/img/image-")
	check(err)
	defer t.Close()
	_, copyErr := io.Copy(t, f)
	check(copyErr)
	http.Redirect(w, r, "/view?id="+t.Name()[17:], 302)
}

func (iuh *ImageUploaderHandler) view(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, "static/img/image-"+r.FormValue("id"))
}

func (iuh *ImageUploaderHandler) HandleRoutes(errFunc httpserver.ErrHandler) *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/", errFunc(iuh.upload))
	r.HandleFunc("/view", errFunc(iuh.view))
	return r

}
