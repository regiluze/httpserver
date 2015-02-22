package httpserver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHttpserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Httpserver Suite")
}
