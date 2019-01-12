package bytesize

import (
	"fmt"
)

const (
	_          = iota             // ignore first value
	KB float64 = 1 << (10 * iota) // 1*2^10=1024
	MB                            // 1*2^20
	GB                            // 1*2^30
	TB                            // 1*2^40
	PB                            // 1*2^50
	EB                            // 1*2^60
)

func ByteSize(i uint64) string {
	b := float64(i)
	switch {
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%dB", i)
}
