// Copyright 2015 The httpserver Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.
package main

import (
	"fmt"
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
	context string
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
	fmt.Println("egi>>", iuh.context)
	http.Redirect(w, r, fmt.Sprintf("%s/view/?id=", iuh.context)+t.Name()[17:], 302)
}

func (iuh *ImageUploaderHandler) view(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, "static/img/image-"+r.FormValue("id"))
}

func (iuh *ImageUploaderHandler) HandleRoutes(context string, r *mux.Router, errFunc httpserver.ErrHandler) *mux.Router {
	iuh.context = "/" + context
	r.HandleFunc(fmt.Sprintf("/%s/", context), errFunc(iuh.upload))
	r.HandleFunc(fmt.Sprintf("/%s/view/", context), errFunc(iuh.view))
	return r
}
