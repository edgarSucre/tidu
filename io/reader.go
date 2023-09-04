package io

import "fmt"

type BadReader struct{}

func (BadReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("reader error")
}
