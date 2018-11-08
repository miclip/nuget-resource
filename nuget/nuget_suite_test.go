package nuget_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNuget(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nuget Suite")
}
