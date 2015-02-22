package mocks

import "github.com/regiluze/httpserver"
import "github.com/stretchr/testify/mock"

type RouteHandler struct {
	mock.Mock
}

func (m *RouteHandler) GetRoutes() []*httpserver.Route {
	ret := m.Called()

	r0 := ret.Get(0).([]*httpserver.Route)

	return r0
}
