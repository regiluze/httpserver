package mocks

import "github.com/stretchr/testify/mock"

import "net/http"

import "github.com/gorilla/mux"

type HttpRouter struct {
	mock.Mock
}

func (m *HttpRouter) HandleFunc(_a0 string, _a1 func(http.ResponseWriter, *http.Request)) *mux.Route {
	ret := m.Called(_a0, _a1)

	r0 := ret.Get(0).(*mux.Route)

	return r0
}
func (m *HttpRouter) Handle(_a0 string, _a1 http.Handler) *mux.Route {
	ret := m.Called(_a0, _a1)

	r0 := ret.Get(0).(*mux.Route)

	return r0
}
func (m *HttpRouter) ServeHTTP(_a0 http.ResponseWriter, _a1 *http.Request) {
	m.Called(_a0, _a1)
}
