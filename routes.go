// Copyright 2015 The httpserver Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.
package httpserver

import (
	"net/http"
)

type Error struct {
	Msg string
}

func NewError(msg string) *Error {
	s := &Error{Msg: msg}
	return s
}

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

type RouteHandler interface {
	GetRoutes() []*Route
}
