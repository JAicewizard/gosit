// +build cgo_bench

package posit

import (
	"testing"
)

func TestAddFuz(t *testing.T) {
	t.Parallel()
	testAddFuz(t)
}

func TestMulFuz(t *testing.T) {
	t.Parallel()
	testMulFuz(t)
}
func TestDivFuz(t *testing.T) {
	t.Parallel()
	testDivFuz(t)
}
func TestSqrtFuz(t *testing.T) {
	t.Parallel()
	testSqrtFuz(t)
}
