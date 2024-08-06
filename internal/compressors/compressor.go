package compressors

import (
	"fmt"
)

type Compressor interface {
	// TODO
	Compress([]byte) ([]byte, error)
	Decompress([]byte) ([]byte, error)
}

var (
	compressorMap = make(map[string]Compressor)
)

func init() {
	compressorMap["none"] = &DummyCompressor{}
}

func New(compressorType string) (Compressor, error) {
	if compressor, ok := compressorMap[compressorType]; ok {
		return compressor, nil
	}

	return nil, fmt.Errorf("unknown compression type %s", compressorType)
}
