package meterconn

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMeterConn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metered Conn Suite")
}
