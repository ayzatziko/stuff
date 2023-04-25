package genericx

import (
	"fmt"
	"testing"
)

type A struct {
	A string
}

type B struct {
	B int
}

type C struct {
	C bool
}

type AorB interface {
	A | B
}

func AorBPrint[T AorB](aorb T) {
	fmt.Println(aorb)
}

func TestAorBFunc(t *testing.T) {
	a := A{"asdf"}
	AorBPrint(a)

	b := B{12}
	AorBPrint(b)

	// algebraic types do not work in go
	// c := C{true}
	// AorBPrint(c)
}
