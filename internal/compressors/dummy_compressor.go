package compressors

type DummyCompressor struct{}

func (c *DummyCompressor) Compress(d []byte) ([]byte, error) {
	return d, nil
}

func (c *DummyCompressor) Decompress(d []byte) ([]byte, error) {
	return d, nil
}
