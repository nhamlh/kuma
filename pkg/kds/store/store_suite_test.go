package store_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSyncResourceStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SyncResourceStore Suite")
}
