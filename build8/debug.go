package build8

import (
	"fmt"
	"math"

	"e8vm.io/e8vm/debug8"
	"e8vm.io/e8vm/image"
)

func debugSection(tab *debug8.Table) (*image.Section, error) {
	bs := tab.Marshal()
	if len(bs) > math.MaxInt32-1 {
		return nil, fmt.Errorf("debug section too large")
	}

	return &image.Section{
		Header: &image.Header{
			Type: image.Debug,
			Size: uint32(len(bs)),
		},
		Bytes: bs,
	}, nil
}
