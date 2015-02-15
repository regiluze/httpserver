package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type Router struct {
	mock.Mock
}

func (m *Router) HandleFunc(context string, fn http.HandlerFunc) {
	m.Called()

}
