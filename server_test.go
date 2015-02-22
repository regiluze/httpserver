package httpserver_test

import (
	"./mocks"
	//	"fmt"
	. "github.com/regiluze/httpserver"
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

const (
	IrrelevantAddress = "0.0.0.0"
	IrrelevantPort    = "80"
)

func buildSingleRoute() []*Route {
	root := NewRoute("", SkipCheckHttpMethod, FakeHandleFunc)
	routes := []*Route{root}
	return routes
}

func buildWrongRoute() []*Route {
	routes := []*Route{}
	return routes
}

func FakeHandleFunc(w http.ResponseWriter, r *http.Request) {}

var _ = Describe("Server", func() {

	var (
		router       *mocks.HttpRouter
		routeHandler *mocks.RouteHandler
		server       *HttpServer
	)

	BeforeEach(func() {
		routeHandler = new(mocks.RouteHandler)
		routeHandler.On("GetRoutes").Return(buildSingleRoute())
		router = new(mocks.HttpRouter)
		router.On("HandleFunc", mock.Anything, mock.Anything, mock.Anything).Return(mux.NewRouter().NewRoute())

		server = NewHttpServer(IrrelevantAddress, IrrelevantPort)
		server.Router = router
	})

	Describe("server init process", func() {
		Context("deploy simple app", func() {
			It("gets single route from handler", func() {

				server.Deploy("/", routeHandler)

				routeHandler.AssertNumberOfCalls(GinkgoT(), "GetRoutes", 1)
			})
			It("adds single handle func to router handlerFunc", func() {

				server.Deploy("/", routeHandler)

				router.AssertNumberOfCalls(GinkgoT(), "HandleFunc", 1)
			})
			It("raises an error with route handler without routes", func() {
				wrongRouteHandler := new(mocks.RouteHandler)
				wrongRouteHandler.On("GetRoutes").Return(buildWrongRoute())
				var err error

				err = server.Deploy("/", wrongRouteHandler)

				Expect(err).To(HaveOccurred())
			})
		})
	})

})
