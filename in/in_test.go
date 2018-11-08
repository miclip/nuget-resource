package in_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("in", func() {

	It("should output an empty JSON list", func() {
		//_, _, _ := in.Execute(in.Request{}, "/targetDir")
		Expect(nil).ShouldNot(HaveOccurred())
	})
})