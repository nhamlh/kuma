package vips_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVIPs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VIPs Suite")
}
