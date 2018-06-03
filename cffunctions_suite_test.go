package cffunctions_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCffunctions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cffunctions Suite")
}
