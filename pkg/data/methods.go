package data

func (tb TarBuf) Read(p []byte) (n int, err error) {
	return tb.Buffer.Read(p)
}

func (tb TarBuf) Close() error {
	return nil
}

func (di DataInfo) OpenTar() (*TarBuf, error) {
	return di.TarBuf, di.Err
}
