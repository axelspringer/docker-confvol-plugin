package driver_test

import "testing"
import "github.com/onsi/gomega"

func TestBooks(t *testing.T) {
	gomega.RegisterTestingT(t)

	f := farm.New([]string{"Cow", "Horse"})
	Expect(f.HasCow()).To(BeTrue(), "Farm should have cow")
}
