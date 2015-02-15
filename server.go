// Copyright 2015 The httpserver Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.
package httpserver

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

const (
	SkipCheckHttpMethod = ""
	GetMethod           = "GET"
	PutMethod           = "PUT"
	PostMethod          = "POST"
	DeleteMethod        = "DELETE"
)

type Route struct {
	Path       string
	Method     string
	HandleFunc http.HandlerFunc
}

func NewRoute(p string, m string, hf http.HandlerFunc) *Route {
	r := &Route{Path: p, Method: m, HandleFunc: hf}
	return r
}

func (r *Route) HasIncorrectHttpMethod(method string) bool {

	if len(r.Method) != 0 && method != r.Method {
		return true
	}
	return false

}

type ServerError struct {
	Msg string
}

func NewError(msg string) *ServerError {
	s := &ServerError{Msg: msg}
	return s
}

type HttpServer struct {
	port             string
	address          string
	errTemplate      *template.Template
	notFoundTemplate *template.Template
	router           *mux.Router
}

func NewHttpServer(a string, p string) *HttpServer {
	r := mux.NewRouter()
	s := &HttpServer{router: r, address: a, port: p}
	return s
}

type RouteHandler interface {
	GetRoutes() []*Route
}

func (s *HttpServer) SetErrTemplate(t *template.Template) {
	s.errTemplate = t
}

func (s *HttpServer) SetNotFoundTemplate(t *template.Template) {
	s.notFoundTemplate = t
}

func (s *HttpServer) errorHandler(route *Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				error := NewError(fmt.Sprintf("\"%v\"", recoverErr))
				w.WriteHeader(500)
				if s.errTemplate != nil {
					s.errTemplate.Execute(w, error)
				} else {
					fmt.Fprintf(w, fmt.Sprintf("\"%v\"", recoverErr))
				}
			}
		}()
		if route.HasIncorrectHttpMethod(r.Method) {
			http.Error(w, "Method not allowed", 405)
			return
		}
		route.HandleFunc(w, r)
	}
}
func (s *HttpServer) NotFound(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(404)
	if s.notFoundTemplate != nil {
		s.notFoundTemplate.Execute(w, nil)
	} else {
		fmt.Fprintf(w, "Not found")
	}

}

func (s *HttpServer) DeployAtBase(h RouteHandler) {
	s.Deploy("", h)
}

func (s *HttpServer) Deploy(context string, h RouteHandler) {
	routes := h.GetRoutes()
	for _, r := range routes {
		s.router.HandleFunc(fmt.Sprintf("%s/%s", context, r.Path), s.errorHandler(r))
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", s.router)
	s.router.NotFoundHandler = http.HandlerFunc(s.NotFound)
}

func (s *HttpServer) Start() error {
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.address, s.port), nil)
}
