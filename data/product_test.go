package data

import "testing"

func TestCheckValidate(t *testing.T) {
	p := &Product{
		Name:  "Himanshu",
		Price: 10,
		SKU:   "abc-csa-ssa",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
