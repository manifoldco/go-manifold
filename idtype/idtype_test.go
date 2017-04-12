package idtype

import "testing"

func TestUpperAndLower(t *testing.T) {

	if Product.Upper() != 0x00 {
		t.Error("Decoded upper byte was incorrect")
	}

	if Product.Lower() != 0xC9 {
		t.Error("Decoded lower byte was incorrect")
	}
}

func TestDecode(t *testing.T) {
	o := Decode(0x00, 0xC9)

	if o != Product {
		t.Error("Decoded type did not match Product")
	}
}
