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

type NotFoundHandler interface {
	Handle(interface{}, func(http.ResponseWriter, *http.Request))
}

type GorillaNotFoundHandler struct {
}

func (g *GorillaNotFoundHandler) Handle(router interface{}, handleFunc func(http.ResponseWriter, *http.Request)) {

	var gorillaRouter *mux.Router
	gorillaRouter = router.(*mux.Router)
	gorillaRouter.NotFoundHandler = http.HandlerFunc(handleFunc)

}

func NewGorillaNotFoundHandler() *GorillaNotFoundHandler {

	handler := &GorillaNotFoundHandler{}
	return handler

}

type HttpServer struct {
	port             string
	address          string
	errTemplate      *template.Template
	notFoundTemplate *template.Template
	Router           HttpRouter
	notFoundHandler  NotFoundHandler
}

type HttpRouter interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request)) *mux.Route
	Handle(string, http.Handler) *mux.Route
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func NewHttpServer(a string, p string) *HttpServer {
	router := mux.NewRouter()
	gorillaHandler := NewGorillaNotFoundHandler()
	s := &HttpServer{Router: router, address: a, port: p, notFoundHandler: gorillaHandler}
	return s
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
		s.Router.HandleFunc(fmt.Sprintf("%s/%s", context, r.Path), s.errorHandler(r))
	}
}

func (s *HttpServer) Start() error {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", s.Router)
	fmt.Println("egi")
	s.notFoundHandler.Handle(s.Router, s.NotFound)
	//http.NotFoundHandler = http.HandlerFunc(s.NotFound)
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.address, s.port), nil)
}
