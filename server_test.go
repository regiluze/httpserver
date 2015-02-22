package httpserver_test

import (
	"./mocks"
	//	"fmt"
	. "github.com/regiluze/httpserver"
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	//	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

const (
	IrrelevantAddress = "0.0.0.0"
	IrrelevantPort    = "80"
)

type FakeRouteImplementation struct{}

func (fake *FakeRouteImplementation) GetRoutes() []*Route {

	root := NewRoute("", SkipCheckHttpMethod, fake.HandleFunc)
	routes := []*Route{root}
	return routes

}

func (fake *FakeRouteImplementation) HandleFunc(w http.ResponseWriter, r *http.Request) {}

var _ = Describe("Server", func() {

	var (
		router       *mocks.HttpRouter
		routeHandler *FakeRouteImplementation
		server       *HttpServer
	)

	BeforeEach(func() {
		routeHandler = new(FakeRouteImplementation)
		router = new(mocks.HttpRouter)
		router.On("HandleFunc", mock.Anything, mock.Anything, mock.Anything).Return(mux.NewRouter().NewRoute())

		server = NewHttpServer(IrrelevantAddress, IrrelevantPort)
		server.Router = router
	})

	Describe("server init process", func() {
		Context("deploy simple app", func() {
			It("gets single route from handler", func() {

				server.Deploy("/", routeHandler)

				router.AssertNumberOfCalls(GinkgoT(), "HandleFunc", 1)
			})
		})
	})

})
