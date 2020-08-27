package person

import (
	"fmt"
	"testing"
)

func Test_NewPerson(t *testing.T) {
	p := NewPerson("kiki", 28)
	fmt.Println(p.Get())
	p.Set("jimmy", 26)
	fmt.Println(p.Get())
}
