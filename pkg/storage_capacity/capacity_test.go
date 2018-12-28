package storage_capacity

import (
	"testing"
)

func TestCapacity_String(t *testing.T) {
	tests := []struct {
		c   Capacity
		str string
	}{
		{
			0,
			"0",
		},
		{
			7,
			"7b",
		},
		{
			8,
			"1B",
		},
		{
			2 * Bit,
			"2b",
		},
		{
			9 * Bit,
			"1.1B",
		},
		{
			18 * Byte,
			"18B",
		},
		{
			1024 * Byte,
			"1KB",
		},
		{
			1023 * Byte,
			"1023B",
		},
		{
			2048 * Byte,
			"2KB",
		},
		{
			1537 * Byte,
			"1.5KB",
		},
		{
			1537 * Megabyte,
			"1.5GB",
		},
	}
	for _, test := range tests {
		if test.str != test.c.String() {
			t.Errorf("expected %s, actual %s", test.str, test.c.String())
		}
	}

}
