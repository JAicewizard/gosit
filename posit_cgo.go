// +build cgo_bench

package posit

//#cgo CFLAGS: -I./positrs
//#cgo LDFLAGS: -L./positrs -lpositrs
//#include "./positrs/positrs.h"
import "C"
import (
	"crypto/rand"
	"fmt"
	fuzz "github.com/google/gofuzz"
	"testing"
)

func testAddFuz(t *testing.T) {
	data := make([]byte, 1000)
	var a uint32
	var b uint32
	for i := 0; true; i++ {
		rand.Read(data)
		fuzz.NewFromGoFuzz(data).Fuzz(&a)
		rand.Read(data)
		fuzz.NewFromGoFuzz(data).Fuzz(&b)
		if i%100 == 0 {
			fmt.Printf("add:%d\n", i)
		}
		exp := posit{
			num: uint32(C.positadd(C.uint(a), C.uint(b))),
			es:  2,
		}
		pa := posit{
			num: a,
			es:  2,
		}
		pb := posit{
			num: b,
			es:  2,
		}
		t.Logf("a:%#032b", a)
		t.Logf("b:%#032b", b)
		t.Logf("exp:%#032b", exp.num)
		res := AddPositSameES(pa, pb)
		t.Logf("got:%#032b", res.num)
		if res != exp {
			t.Fatal("got the wrong number:", Getfloat(res), "expected:", Getfloat(exp))
		}
	}
}

func testMulFuz(t *testing.T) {
	data := make([]byte, 1000)
	var a uint32
	var b uint32
	for i := 0; true; i++ {
		rand.Read(data)
		fuzz.NewFromGoFuzz(data).Fuzz(&a)
		rand.Read(data)
		fuzz.NewFromGoFuzz(data).Fuzz(&b)
		if i%100 == 0 {
			fmt.Printf("mul:%d\n", i)
		}
		exp := posit{
			num: uint32(C.positmul(C.uint(a), C.uint(b))),
			es:  2,
		}
		pa := posit{
			num: a,
			es:  2,
		}
		pb := posit{
			num: b,
			es:  2,
		}
		t.Logf("a:%#032b", a)
		t.Logf("b:%#032b", b)
		t.Logf("exp:%#032b", exp.num)
		res := MulPositSameES(pa, pb)
		t.Logf("got:%#032b", res.num)
		if res != exp {
			t.Fatal("got the wrong number:", Getfloat(res), "expected:", Getfloat(exp))
		}
	}
}

func testDivFuz(t *testing.T) {
	data := make([]byte, 1000)
	var a uint32
	var b uint32
	for i := 0; true; i++ {
		rand.Read(data)
		fuzz.NewFromGoFuzz(data).Fuzz(&a)
		rand.Read(data)
		fuzz.NewFromGoFuzz(data).Fuzz(&b)
		if i%100 == 0 {
			fmt.Printf("div:%d\n", i)
		}
		exp := posit{
			num: uint32(C.positdiv(C.uint(a), C.uint(b))),
			es:  2,
		}
		pa := posit{
			num: a,
			es:  2,
		}
		pb := posit{
			num: b,
			es:  2,
		}
		t.Logf("a:%#032b", a)
		t.Logf("b:%#032b", b)
		t.Logf("exp:%#032b", exp.num)
		res := DivPositSameES(pa, pb)
		t.Logf("got:%#032b", res.num)
		if res != exp {
			t.Fatal("got the wrong number:", Getfloat(res), "expected:", Getfloat(exp))
		}
	}
}

func testSqrtFuz(t *testing.T) {
	data := make([]byte, 1000)
	var a uint32
	for i := 0; true; i++ {
		rand.Read(data)
		fuzz.NewFromGoFuzz(data).Fuzz(&a)
		if i%100 == 0 {
			fmt.Printf("sqrt:%d\n", i)
		}
		exp := posit{
			num: uint32(C.positsqrt(C.uint(a))),
			es:  2,
		}
		pa := posit{
			num: a,
			es:  2,
		}
		res := SqrtPosit(pa)
		if res != exp && (a&(1<<31)) == 0 {
			t.Logf("a:%#032b", a)
			t.Logf("exp:%#032b", exp.num)
			t.Logf("got:%#032b", res.num)
			t.Fatal("got the wrong number:", Getfloat(res), "expected:", Getfloat(exp))
		}
	}
}
