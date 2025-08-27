package codegen

import (
	"crypto/rand"
	"fmt"
)

type Numeric struct{ digits int }

func NewNumeric(digits int) Numeric { return Numeric{digits: digits} }

func (n Numeric) Generate() (string, error) {
	max := 1
	for i := 0; i < n.digits; i++ {
		max *= 10
	}
	b := make([]byte, n.digits)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	val := 0
	for i := 0; i < n.digits; i++ {
		val = (val*10 + int(b[i])%10) % max
	}
	format := fmt.Sprintf("%%0%dd", n.digits)
	return fmt.Sprintf(format, val), nil
}
