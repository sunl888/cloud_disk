package storage_capacity

import (
	"strconv"
	"strings"
)

type Capacity uint64

const (
	Bit      = 1
	Byte     = Bit << 3
	Kilobyte = Byte << 10
	Megabyte = Kilobyte << 10
	Gigabyte = Megabyte << 10
	Terabyte = Gigabyte << 10
	Petabyte = Terabyte << 10
	Exabyte  = Petabyte << 10
)

var capacityStrMap = map[Capacity]string{
	Bit:      "b",
	Byte:     "B",
	Kilobyte: "KB",
	Megabyte: "MB",
	Gigabyte: "GB",
	Terabyte: "TB",
	Petabyte: "PB",
	Exabyte:  "EB",
}

func (c Capacity) String() string {
	if c == 0 {
		return "0"
	}
	sb := strings.Builder{}
	// cInt64 := uint64(c)

	if c < Byte {

		sb.WriteString(strconv.FormatInt(int64(c), 10))
		sb.WriteString("b")
		return sb.String()
	}

	units := []Capacity{Bit, Byte, Kilobyte, Megabyte, Gigabyte, Terabyte, Petabyte, Exabyte}
	for i := 2; i < len(units); i++ {
		if c < units[i] {
			t := float64(c) / float64(units[i-1])
			sb.WriteString(strings.TrimSuffix(strconv.FormatFloat(t, 'f', 1, 64), ".0"))
			sb.WriteString(capacityStrMap[units[i-1]])
			break
		}
	}
	return sb.String()
}
