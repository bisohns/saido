package inspector

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
)

var byteMap = []string{"B", "KB", "MB", "GB", "TB"}

// ByteSize : helps parse string into individual byte values
type ByteSize struct {
	value float64
}

func (b *ByteSize) format(unit string) float64 {
	return b.value / math.Pow(1024, index(byteMap, unit))
}

// Initialize a NewByteSize
func NewByteSize(byteCount string, unit string) *ByteSize {
	if byteCount == `-` {
		byteCount = "0"
	}
	byteCountInt, err := strconv.ParseFloat(byteCount, 64)
	byteCountInt = byteCountInt * math.Pow(1024, index(byteMap, unit))
	if err != nil {
		log.Fatal(err)
	}
	return &ByteSize{
		value: byteCountInt,
	}
}

func index(arr []string, str string) float64 {
	for index, a := range arr {
		// Match GB to GB and GiB
		if a == str || fmt.Sprintf(
			`%si%s`, string(a[:len(a)-1]), string(a[len(a)-1])) == str {
			return float64(index)
		}
	}
	// defaults to byte if not found in array
	return 0
}
